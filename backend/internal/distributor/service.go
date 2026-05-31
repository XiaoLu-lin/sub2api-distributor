package distributor

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"time"
)

// Service encapsulates all distributor profile, rebate summary, invite, and
// withdrawal persistence against the shared PostgreSQL database.
type Service struct {
	db *sql.DB
}

type profileTransactionalStore interface {
	upsertProfile(context.Context, Profile) error
}

// NewService creates a distributor domain service backed by the shared database.
func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

// GetProfile returns the settlement profile for a specific distributor user.
func (s *Service) GetProfile(ctx context.Context, userID int64) (*Profile, error) {
	var profile Profile
	var extra []byte
	err := s.db.QueryRowContext(ctx, `
SELECT user_id, status, display_name, settlement_channel, settlement_account_name,
       settlement_account_no, settlement_account_extra, notes
FROM distributor_profiles
WHERE user_id = $1
`, userID).Scan(
		&profile.UserID,
		&profile.Status,
		&profile.DisplayName,
		&profile.SettlementChannel,
		&profile.SettlementAccountName,
		&profile.SettlementAccountNo,
		&extra,
		&profile.Notes,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProfileNotFound
		}
		return nil, fmt.Errorf("get profile: %w", err)
	}
	profile.SettlementAccountExtra = string(extra)
	return &profile, nil
}

// UpsertProfile creates or updates a distributor profile while preserving the
// one-row-per-user invariant.
func (s *Service) UpsertProfile(ctx context.Context, profile Profile) error {
	if profile.UserID <= 0 {
		return ErrInvalidProfile
	}
	if profile.Status == "" {
		profile.Status = "active"
	}
	if !IsValidDistributorStatus(profile.Status) {
		return ErrInvalidProfile
	}

	extra := profile.SettlementAccountExtra
	if extra == "" {
		extra = "{}"
	}

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin profile upsert tx: %w", err)
	}
	return runProfileUpsertTransaction(
		ctx,
		tx.Commit,
		func() { _ = tx.Rollback() },
		profileStoreFunc(func(ctx context.Context, profile Profile) error {
			_, err := tx.ExecContext(ctx, `
INSERT INTO distributor_profiles (
  user_id, status, display_name, settlement_channel, settlement_account_name,
  settlement_account_no, settlement_account_extra, notes, created_at, updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7::jsonb, $8, NOW(), NOW())
ON CONFLICT (user_id) DO UPDATE
SET status = EXCLUDED.status,
    display_name = EXCLUDED.display_name,
    settlement_channel = EXCLUDED.settlement_channel,
    settlement_account_name = EXCLUDED.settlement_account_name,
    settlement_account_no = EXCLUDED.settlement_account_no,
    settlement_account_extra = EXCLUDED.settlement_account_extra,
    notes = EXCLUDED.notes,
    updated_at = NOW()
`, profile.UserID, profile.Status, profile.DisplayName, profile.SettlementChannel,
				profile.SettlementAccountName, profile.SettlementAccountNo, extra, profile.Notes)
			if err != nil {
				return fmt.Errorf("upsert profile: %w", err)
			}
			return nil
		}),
		func(ctx context.Context, userID int64) error {
			return ensureAffiliateIdentityWithExecer(ctx, tx, userID)
		},
		profile,
	)
}

type profileStoreFunc func(context.Context, Profile) error

func (fn profileStoreFunc) upsertProfile(ctx context.Context, profile Profile) error {
	return fn(ctx, profile)
}

func upsertProfileAtomically(
	ctx context.Context,
	store profileTransactionalStore,
	ensureAffiliateIdentity func(context.Context, int64) error,
	profile Profile,
) error {
	if err := store.upsertProfile(ctx, profile); err != nil {
		return err
	}
	if profile.Status == "active" {
		if err := ensureAffiliateIdentity(ctx, profile.UserID); err != nil {
			return err
		}
	}
	return nil
}

func runProfileUpsertTransaction(
	ctx context.Context,
	commit func() error,
	rollback func(),
	store profileTransactionalStore,
	ensureAffiliateIdentity func(context.Context, int64) error,
	profile Profile,
) error {
	committed := false
	defer func() {
		if !committed {
			rollback()
		}
	}()

	if err := upsertProfileAtomically(ctx, store, ensureAffiliateIdentity, profile); err != nil {
		return err
	}
	if err := commit(); err != nil {
		return fmt.Errorf("commit profile upsert tx: %w", err)
	}
	committed = true
	return nil
}

// LookupUsers lets operators search the main-system user table before enabling
// or editing distributor access.
func (s *Service) LookupUsers(ctx context.Context, keyword string) ([]UserLookupItem, error) {
	normalizedKeyword := NormalizeUserLookupKeyword(keyword)
	if normalizedKeyword == "" {
		return []UserLookupItem{}, nil
	}

	pattern := "%" + normalizedKeyword + "%"
	rows, err := s.db.QueryContext(ctx, `
SELECT
  u.id,
  u.email,
  u.username,
  u.role,
  u.status,
  CASE WHEN dp.user_id IS NULL THEN false ELSE true END AS is_distributor
FROM users u
LEFT JOIN distributor_profiles dp ON dp.user_id = u.id
WHERE u.deleted_at IS NULL
  AND (
    u.email ILIKE $1
    OR u.username ILIKE $1
    OR CAST(u.id AS TEXT) = $2
  )
ORDER BY u.email ASC
LIMIT 20
`, pattern, normalizedKeyword)
	if err != nil {
		return nil, fmt.Errorf("lookup users: %w", err)
	}
	defer rows.Close()

	items := make([]UserLookupItem, 0)
	for rows.Next() {
		var item UserLookupItem
		if err := rows.Scan(
			&item.ID,
			&item.Email,
			&item.Username,
			&item.Role,
			&item.Status,
			&item.IsDistributor,
		); err != nil {
			return nil, fmt.Errorf("scan lookup user: %w", err)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// GetSummary aggregates rebate earnings and withdrawal balances for the distributor dashboard.
func (s *Service) GetSummary(ctx context.Context, userID int64) (*Summary, error) {
	var summary Summary
	err := s.db.QueryRowContext(ctx, `
SELECT
  COALESCE(SUM(CASE WHEN action = 'accrue' THEN amount ELSE 0 END), 0)::text AS total_earned,
  COALESCE(SUM(CASE WHEN action = 'accrue' AND frozen_until IS NOT NULL AND frozen_until > NOW() THEN amount ELSE 0 END), 0)::text AS frozen_amount,
  COALESCE(SUM(CASE WHEN action = 'transfer' THEN amount ELSE 0 END), 0)::text AS internal_transferred_amount
FROM user_affiliate_ledger
WHERE user_id = $1
`, userID).Scan(&summary.TotalEarnedText, &summary.FrozenAmountText, &summary.InternalTransferredText)
	if err != nil {
		return nil, fmt.Errorf("query ledger summary: %w", err)
	}

	err = s.db.QueryRowContext(ctx, `
SELECT
  COALESCE(SUM(CASE WHEN status = 'paying' THEN amount ELSE 0 END), 0)::text AS paying_amount,
  COALESCE(SUM(CASE WHEN status = 'paid' THEN amount ELSE 0 END), 0)::text AS paid_amount
FROM distributor_withdraw_requests
WHERE user_id = $1
`, userID).Scan(&summary.PayingAmountText, &summary.PaidAmountText)
	if err != nil {
		return nil, fmt.Errorf("query withdraw summary: %w", err)
	}

	summary.WithdrawableAmountText = ComputeWithdrawableAmountDecimal(
		summary.TotalEarnedText,
		summary.FrozenAmountText,
		summary.InternalTransferredText,
		summary.PayingAmountText,
		summary.PaidAmountText,
	)
	summary.TotalEarned = DecimalStringToFloat64(summary.TotalEarnedText)
	summary.FrozenAmount = DecimalStringToFloat64(summary.FrozenAmountText)
	summary.InternalTransferred = DecimalStringToFloat64(summary.InternalTransferredText)
	summary.PayingAmount = DecimalStringToFloat64(summary.PayingAmountText)
	summary.PaidAmount = DecimalStringToFloat64(summary.PaidAmountText)
	summary.WithdrawableAmount = DecimalStringToFloat64(summary.WithdrawableAmountText)
	return &summary, nil
}

// ListInvitees returns the users attributed to a distributor and their accrued rebate totals.
func (s *Service) ListInvitees(ctx context.Context, userID int64) ([]Invitee, error) {
	rows, err := s.db.QueryContext(ctx, `
SELECT
  u.id,
  u.email,
  u.username,
  COALESCE(SUM(CASE WHEN l.action = 'accrue' THEN l.amount ELSE 0 END), 0)::double precision AS total_rebate,
  u.created_at::text
FROM user_affiliates ua
JOIN users u ON u.id = ua.user_id
LEFT JOIN user_affiliate_ledger l ON l.source_user_id = u.id AND l.user_id = $1
WHERE ua.inviter_id = $1
GROUP BY u.id, u.email, u.username, u.created_at
ORDER BY u.created_at DESC
`, userID)
	if err != nil {
		return nil, fmt.Errorf("list invitees: %w", err)
	}
	defer rows.Close()

	items := make([]Invitee, 0)
	for rows.Next() {
		var item Invitee
		if err := rows.Scan(&item.UserID, &item.Email, &item.Username, &item.TotalRebate, &item.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan invitee: %w", err)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// GetInviteMeta returns the real invite code currently assigned to the distributor.
func (s *Service) GetInviteMeta(ctx context.Context, userID int64) (*InviteMeta, error) {
	var meta InviteMeta
	err := s.db.QueryRowContext(ctx, `
SELECT aff_code
FROM user_affiliates
WHERE user_id = $1
`, userID).Scan(&meta.AffCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInviteMetaNotFound
		}
		return nil, fmt.Errorf("get invite meta: %w", err)
	}
	return &meta, nil
}

// ListRebates returns accrue ledger records that count toward distributor earnings.
func (s *Service) ListRebates(ctx context.Context, userID int64) ([]RebateRecord, error) {
	rows, err := s.db.QueryContext(ctx, `
SELECT
  l.id,
  l.amount::double precision,
  l.source_user_id,
  COALESCE(u.email, ''),
  l.source_order_id,
  l.created_at::text
FROM user_affiliate_ledger l
LEFT JOIN users u ON u.id = l.source_user_id
WHERE l.user_id = $1
  AND l.action = 'accrue'
ORDER BY l.created_at DESC
`, userID)
	if err != nil {
		return nil, fmt.Errorf("list rebates: %w", err)
	}
	defer rows.Close()

	items := make([]RebateRecord, 0)
	for rows.Next() {
		var item RebateRecord
		if err := rows.Scan(&item.LedgerID, &item.Amount, &item.SourceUserID, &item.SourceEmail, &item.SourceOrderID, &item.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan rebate: %w", err)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// ListWithdrawals returns the current user's withdrawal requests ordered by creation time.
func (s *Service) ListWithdrawals(ctx context.Context, userID int64) ([]WithdrawalRequest, error) {
	rows, err := s.db.QueryContext(ctx, `
SELECT id, request_no, user_id, amount::double precision, status, applicant_remark,
       COALESCE(paid_channel, ''), COALESCE(paid_reference_no, ''), COALESCE(paid_remark, ''),
       CASE WHEN paid_at IS NULL THEN NULL ELSE paid_at::text END,
       created_at::text, snapshot_withdrawable_before::double precision, snapshot_withdrawable_after::double precision
FROM distributor_withdraw_requests
WHERE user_id = $1
ORDER BY created_at DESC
`, userID)
	if err != nil {
		return nil, fmt.Errorf("list withdrawals: %w", err)
	}
	defer rows.Close()

	items := make([]WithdrawalRequest, 0)
	for rows.Next() {
		var item WithdrawalRequest
		if err := rows.Scan(
			&item.ID,
			&item.RequestNo,
			&item.UserID,
			&item.Amount,
			&item.Status,
			&item.ApplicantRemark,
			&item.PaidChannel,
			&item.PaidReferenceNo,
			&item.PaidRemark,
			&item.PaidAt,
			&item.CreatedAt,
			&item.SnapshotWithdrawableBefore,
			&item.SnapshotWithdrawableAfter,
		); err != nil {
			return nil, fmt.Errorf("scan withdrawal: %w", err)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// CreateWithdrawal validates the requested amount, snapshots the pre-request
// balance, and creates both the request row and its initial audit event.
func (s *Service) CreateWithdrawal(ctx context.Context, userID int64, amount float64, remark string) (*WithdrawalRequest, error) {
	if amount <= 0 || math.IsNaN(amount) || math.IsInf(amount, 0) {
		return nil, ErrInvalidWithdrawalAmount
	}

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	summary, err := s.getSummaryTx(ctx, tx, userID)
	if err != nil {
		return nil, err
	}
	amountText := DecimalStringFromFloat64(amount)
	if !DecimalAmountWithinBalance(amountText, summary.WithdrawableAmountText) {
		return nil, ErrWithdrawalAmountTooLarge
	}

	// Request numbers only need to be unique enough for local operator workflows,
	// so a timestamp plus a short nanosecond suffix is sufficient for this MVP.
	requestNo := fmt.Sprintf("DW%s%04d", time.Now().Format("20060102150405"), time.Now().Nanosecond()%10000)
	var requestID int64
	var createdAt string
	withdrawableAfterText := ComputeWithdrawableAmountDecimal(
		summary.TotalEarnedText,
		summary.FrozenAmountText,
		summary.InternalTransferredText,
		AddDecimalStrings(summary.PayingAmountText, amountText),
		summary.PaidAmountText,
	)
	withdrawableAfter := DecimalStringToFloat64(withdrawableAfterText)

	err = tx.QueryRowContext(ctx, `
INSERT INTO distributor_withdraw_requests (
  request_no, user_id, amount, status,
  snapshot_total_earned, snapshot_internal_transferred_amount,
  snapshot_paying_before, snapshot_paid_before,
  snapshot_withdrawable_before, snapshot_withdrawable_after,
  applicant_remark, created_at, updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW())
RETURNING id, created_at::text
`, requestNo, userID, amount, WithdrawalStatusPaying, summary.TotalEarned,
		summary.InternalTransferred, summary.PayingAmount, summary.PaidAmount,
		summary.WithdrawableAmount, withdrawableAfter, remark).
		Scan(&requestID, &createdAt)
	if err != nil {
		return nil, fmt.Errorf("insert withdraw request: %w", err)
	}

	eventDetail, _ := json.Marshal(map[string]any{
		"amount": amount,
		"remark": remark,
	})
	if _, err := tx.ExecContext(ctx, `
INSERT INTO distributor_withdraw_events (request_id, action, operator_user_id, detail, created_at)
VALUES ($1, 'create', $2, $3::jsonb, NOW())
`, requestID, userID, string(eventDetail)); err != nil {
		return nil, fmt.Errorf("insert withdraw event: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit withdraw request: %w", err)
	}

	return &WithdrawalRequest{
		ID:                         requestID,
		RequestNo:                  requestNo,
		UserID:                     userID,
		Amount:                     amount,
		Status:                     string(WithdrawalStatusPaying),
		ApplicantRemark:            remark,
		CreatedAt:                  createdAt,
		SnapshotWithdrawableBefore: summary.WithdrawableAmount,
		SnapshotWithdrawableAfter:  withdrawableAfter,
	}, nil
}

// CancelWithdrawal lets a distributor cancel their own payout request while it is still paying.
func (s *Service) CancelWithdrawal(ctx context.Context, userID int64, requestID int64) error {
	return s.updateWithdrawalStatus(ctx, requestID, userID, userID, WithdrawalStatusCancelled, "", "", "")
}

// CancelWithdrawalAsOperator lets an operator cancel any payout request that is still paying.
func (s *Service) CancelWithdrawalAsOperator(ctx context.Context, operatorUserID int64, requestID int64) error {
	return s.updateWithdrawalStatus(ctx, requestID, 0, operatorUserID, WithdrawalStatusCancelled, "", "", "")
}

// MarkWithdrawalPaid finalizes a paying request after the operator completes the offline payout.
func (s *Service) MarkWithdrawalPaid(ctx context.Context, operatorUserID int64, requestID int64, paidChannel string, paidReferenceNo string, paidRemark string) error {
	return s.updateWithdrawalStatus(ctx, requestID, 0, operatorUserID, WithdrawalStatusPaid, paidChannel, paidReferenceNo, paidRemark)
}

// updateWithdrawalStatus centralizes request ownership checks, status transition
// validation, row locking, and event insertion for payout state changes.
func (s *Service) updateWithdrawalStatus(
	ctx context.Context,
	requestID int64,
	ownerUserID int64,
	operatorUserID int64,
	nextStatus WithdrawalStatus,
	paidChannel string,
	paidReferenceNo string,
	paidRemark string,
) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	var currentStatus string
	var requestUserID int64
	err = tx.QueryRowContext(ctx, `
SELECT user_id, status
FROM distributor_withdraw_requests
WHERE id = $1
FOR UPDATE
`, requestID).Scan(&requestUserID, &currentStatus)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrWithdrawalNotFound
		}
		return fmt.Errorf("query withdraw request: %w", err)
	}
	if ownerUserID > 0 && ownerUserID != requestUserID {
		return ErrWithdrawalNotFound
	}
	if !CanTransitionWithdrawalStatus(WithdrawalStatus(currentStatus), nextStatus) {
		return ErrInvalidWithdrawalTransition
	}

	// Paid requests capture payout metadata, while cancelled requests only need
	// the state transition itself.
	if nextStatus == WithdrawalStatusPaid {
		if _, err := tx.ExecContext(ctx, `
UPDATE distributor_withdraw_requests
SET status = $2,
    paid_at = NOW(),
    paid_channel = $3,
    paid_reference_no = $4,
    paid_remark = $5,
    updated_at = NOW()
WHERE id = $1
`, requestID, nextStatus, paidChannel, paidReferenceNo, paidRemark); err != nil {
			return fmt.Errorf("mark paid: %w", err)
		}
	} else {
		if _, err := tx.ExecContext(ctx, `
UPDATE distributor_withdraw_requests
SET status = $2,
    updated_at = NOW()
WHERE id = $1
`, requestID, nextStatus); err != nil {
			return fmt.Errorf("cancel request: %w", err)
		}
	}

	detail, _ := json.Marshal(map[string]any{
		"paid_channel":      paidChannel,
		"paid_reference_no": paidReferenceNo,
		"paid_remark":       paidRemark,
	})
	action := "cancel"
	if nextStatus == WithdrawalStatusPaid {
		action = "mark_paid"
	}
	if _, err := tx.ExecContext(ctx, `
INSERT INTO distributor_withdraw_events (request_id, action, operator_user_id, detail, created_at)
VALUES ($1, $2, $3, $4::jsonb, NOW())
`, requestID, action, operatorUserID, string(detail)); err != nil {
		return fmt.Errorf("insert event: %w", err)
	}

	return tx.Commit()
}

// ListDistributorProfiles returns all enabled or disabled distributor profiles for the operator console.
func (s *Service) ListDistributorProfiles(ctx context.Context) ([]Profile, error) {
	rows, err := s.db.QueryContext(ctx, `
SELECT user_id, status, display_name, settlement_channel, settlement_account_name,
       settlement_account_no, settlement_account_extra, notes
FROM distributor_profiles
ORDER BY updated_at DESC, user_id DESC
`)
	if err != nil {
		return nil, fmt.Errorf("list distributor profiles: %w", err)
	}
	defer rows.Close()

	items := make([]Profile, 0)
	for rows.Next() {
		var item Profile
		var extra []byte
		if err := rows.Scan(
			&item.UserID,
			&item.Status,
			&item.DisplayName,
			&item.SettlementChannel,
			&item.SettlementAccountName,
			&item.SettlementAccountNo,
			&extra,
			&item.Notes,
		); err != nil {
			return nil, fmt.Errorf("scan distributor profile: %w", err)
		}
		item.SettlementAccountExtra = string(extra)
		items = append(items, item)
	}
	return items, rows.Err()
}

// ListAllWithdrawals returns every withdrawal request for the operator queue view.
func (s *Service) ListAllWithdrawals(ctx context.Context) ([]WithdrawalRequest, error) {
	rows, err := s.db.QueryContext(ctx, `
SELECT id, request_no, user_id, amount::double precision, status, applicant_remark,
       COALESCE(paid_channel, ''), COALESCE(paid_reference_no, ''), COALESCE(paid_remark, ''),
       CASE WHEN paid_at IS NULL THEN NULL ELSE paid_at::text END,
       created_at::text, snapshot_withdrawable_before::double precision, snapshot_withdrawable_after::double precision
FROM distributor_withdraw_requests
ORDER BY created_at DESC
`)
	if err != nil {
		return nil, fmt.Errorf("list all withdrawals: %w", err)
	}
	defer rows.Close()

	items := make([]WithdrawalRequest, 0)
	for rows.Next() {
		var item WithdrawalRequest
		if err := rows.Scan(
			&item.ID,
			&item.RequestNo,
			&item.UserID,
			&item.Amount,
			&item.Status,
			&item.ApplicantRemark,
			&item.PaidChannel,
			&item.PaidReferenceNo,
			&item.PaidRemark,
			&item.PaidAt,
			&item.CreatedAt,
			&item.SnapshotWithdrawableBefore,
			&item.SnapshotWithdrawableAfter,
		); err != nil {
			return nil, fmt.Errorf("scan withdrawal: %w", err)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// GetWithdrawalDetail returns a request row together with its audit history.
func (s *Service) GetWithdrawalDetail(ctx context.Context, requestID int64) (*WithdrawalDetail, error) {
	var item WithdrawalRequest
	err := s.db.QueryRowContext(ctx, `
SELECT id, request_no, user_id, amount::double precision, status, applicant_remark,
       COALESCE(paid_channel, ''), COALESCE(paid_reference_no, ''), COALESCE(paid_remark, ''),
       CASE WHEN paid_at IS NULL THEN NULL ELSE paid_at::text END,
       created_at::text, snapshot_withdrawable_before::double precision, snapshot_withdrawable_after::double precision
FROM distributor_withdraw_requests
WHERE id = $1
`, requestID).Scan(
		&item.ID,
		&item.RequestNo,
		&item.UserID,
		&item.Amount,
		&item.Status,
		&item.ApplicantRemark,
		&item.PaidChannel,
		&item.PaidReferenceNo,
		&item.PaidRemark,
		&item.PaidAt,
		&item.CreatedAt,
		&item.SnapshotWithdrawableBefore,
		&item.SnapshotWithdrawableAfter,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrWithdrawalNotFound
		}
		return nil, fmt.Errorf("get withdrawal detail: %w", err)
	}

	rows, err := s.db.QueryContext(ctx, `
SELECT id, request_id, action, operator_user_id, detail::text, created_at::text
FROM distributor_withdraw_events
WHERE request_id = $1
ORDER BY created_at DESC, id DESC
`, requestID)
	if err != nil {
		return nil, fmt.Errorf("list withdrawal events: %w", err)
	}
	defer rows.Close()

	events := make([]WithdrawalEvent, 0)
	for rows.Next() {
		var event WithdrawalEvent
		if err := rows.Scan(
			&event.ID,
			&event.RequestID,
			&event.Action,
			&event.OperatorUserID,
			&event.Detail,
			&event.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan withdrawal event: %w", err)
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate withdrawal events: %w", err)
	}

	return &WithdrawalDetail{
		WithdrawalRequest: item,
		Events:            events,
	}, nil
}

// getSummaryTx recomputes balances inside a transaction and locks withdrawal
// rows to avoid concurrent double-withdrawal races.
func (s *Service) getSummaryTx(ctx context.Context, tx *sql.Tx, userID int64) (*Summary, error) {
	var summary Summary
	err := tx.QueryRowContext(ctx, `
SELECT
  COALESCE(SUM(CASE WHEN action = 'accrue' THEN amount ELSE 0 END), 0)::text AS total_earned,
  COALESCE(SUM(CASE WHEN action = 'accrue' AND frozen_until IS NOT NULL AND frozen_until > NOW() THEN amount ELSE 0 END), 0)::text AS frozen_amount,
  COALESCE(SUM(CASE WHEN action = 'transfer' THEN amount ELSE 0 END), 0)::text AS internal_transferred_amount
FROM user_affiliate_ledger
WHERE user_id = $1
`, userID).Scan(&summary.TotalEarnedText, &summary.FrozenAmountText, &summary.InternalTransferredText)
	if err != nil {
		return nil, fmt.Errorf("query ledger summary in tx: %w", err)
	}

	rows, err := tx.QueryContext(ctx, `
SELECT amount::text, status
FROM distributor_withdraw_requests
WHERE user_id = $1
FOR UPDATE
`, userID)
	if err != nil {
		return nil, fmt.Errorf("query withdrawal rows in tx: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var amountText string
		var status string
		if err := rows.Scan(&amountText, &status); err != nil {
			return nil, fmt.Errorf("scan withdrawal row in tx: %w", err)
		}
		switch status {
		case string(WithdrawalStatusPaying):
			summary.PayingAmountText = AddDecimalStrings(summary.PayingAmountText, amountText)
		case string(WithdrawalStatusPaid):
			summary.PaidAmountText = AddDecimalStrings(summary.PaidAmountText, amountText)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate withdrawal rows in tx: %w", err)
	}
	summary.WithdrawableAmountText = ComputeWithdrawableAmountDecimal(
		summary.TotalEarnedText,
		summary.FrozenAmountText,
		summary.InternalTransferredText,
		summary.PayingAmountText,
		summary.PaidAmountText,
	)
	summary.TotalEarned = DecimalStringToFloat64(summary.TotalEarnedText)
	summary.FrozenAmount = DecimalStringToFloat64(summary.FrozenAmountText)
	summary.InternalTransferred = DecimalStringToFloat64(summary.InternalTransferredText)
	summary.PayingAmount = DecimalStringToFloat64(summary.PayingAmountText)
	summary.PaidAmount = DecimalStringToFloat64(summary.PaidAmountText)
	summary.WithdrawableAmount = DecimalStringToFloat64(summary.WithdrawableAmountText)
	return &summary, nil
}

var (
	// ErrProfileNotFound indicates the requested distributor profile row does not exist.
	ErrProfileNotFound = errors.New("未找到分销商资料")
	// ErrInviteMetaNotFound indicates the user has no invite code row in the source affiliate table.
	ErrInviteMetaNotFound = errors.New("未找到邀请码信息")
	// ErrInvalidProfile indicates the submitted distributor profile payload is incomplete or invalid.
	ErrInvalidProfile = errors.New("分销商资料不合法")
	// ErrInvalidWithdrawalAmount rejects zero, negative, NaN, or infinite payout amounts.
	ErrInvalidWithdrawalAmount = errors.New("提现金额不合法")
	// ErrWithdrawalAmountTooLarge rejects requests above the current withdrawable balance.
	ErrWithdrawalAmountTooLarge = errors.New("提现金额超过当前可申请金额")
	// ErrWithdrawalNotFound indicates the payout request row is missing or inaccessible to the caller.
	ErrWithdrawalNotFound = errors.New("未找到提现申请")
	// ErrInvalidWithdrawalTransition rejects paid/cancelled requests from being changed again.
	ErrInvalidWithdrawalTransition = errors.New("当前提现状态不允许执行该操作")
)
