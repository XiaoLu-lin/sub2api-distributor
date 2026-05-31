package distributor

import "encoding/json"

// Summary is the distributor dashboard snapshot returned to the frontend.
type Summary struct {
	TotalEarned         float64 `json:"total_earned"`
	FrozenAmount        float64 `json:"frozen_amount"`
	InternalTransferred float64 `json:"internal_transferred_amount"`
	PayingAmount        float64 `json:"paying_amount"`
	PaidAmount          float64 `json:"paid_amount"`
	WithdrawableAmount  float64 `json:"withdrawable_amount"`
	TotalEarnedText         string `json:"-"`
	FrozenAmountText        string `json:"-"`
	InternalTransferredText string `json:"-"`
	PayingAmountText        string `json:"-"`
	PaidAmountText          string `json:"-"`
	WithdrawableAmountText  string `json:"-"`
}

// Profile stores distributor enablement and settlement information.
type Profile struct {
	UserID                 int64  `json:"user_id"`
	Status                 string `json:"status"`
	DisplayName            string `json:"display_name"`
	SettlementChannel      string `json:"settlement_channel"`
	SettlementAccountName  string `json:"settlement_account_name"`
	SettlementAccountNo    string `json:"settlement_account_no"`
	SettlementAccountExtra string `json:"settlement_account_extra"`
	Notes                  string `json:"notes"`
}

// UserLookupItem is the operator-facing projection used when enabling or editing distributors.
type UserLookupItem struct {
	ID            int64  `json:"id"`
	Email         string `json:"email"`
	Username      string `json:"username"`
	Role          string `json:"role"`
	Status        string `json:"status"`
	IsDistributor bool   `json:"is_distributor"`
}

// Invitee describes a user attributed to a distributor invitation relationship.
type Invitee struct {
	UserID      int64   `json:"user_id"`
	Email       string  `json:"email"`
	Username    string  `json:"username"`
	TotalRebate float64 `json:"total_rebate"`
	CreatedAt   string  `json:"created_at"`
}

// InviteMeta exposes the current distributor invite code from the main system.
type InviteMeta struct {
	AffCode string `json:"aff_code"`
}

// RebateRecord is a normalized rebate ledger item shown in the portal UI.
type RebateRecord struct {
	LedgerID      int64   `json:"ledger_id"`
	Amount        float64 `json:"amount"`
	SourceUserID  int64   `json:"source_user_id"`
	SourceEmail   string  `json:"source_email"`
	SourceOrderID *int64  `json:"source_order_id,omitempty"`
	CreatedAt     string  `json:"created_at"`
}

// WithdrawalRequest is the shared summary model for distributor withdrawal rows.
type WithdrawalRequest struct {
	ID                         int64   `json:"id"`
	RequestNo                  string  `json:"request_no"`
	UserID                     int64   `json:"user_id"`
	Amount                     float64 `json:"amount"`
	Status                     string  `json:"status"`
	ApplicantRemark            string  `json:"applicant_remark"`
	PaidChannel                string  `json:"paid_channel"`
	PaidReferenceNo            string  `json:"paid_reference_no"`
	PaidRemark                 string  `json:"paid_remark"`
	PaidAt                     *string `json:"paid_at"`
	CreatedAt                  string  `json:"created_at"`
	SnapshotWithdrawableBefore float64 `json:"snapshot_withdrawable_before"`
	SnapshotWithdrawableAfter  float64 `json:"snapshot_withdrawable_after"`
}

// WithdrawalDetail extends a request row with its audit event history.
type WithdrawalDetail struct {
	WithdrawalRequest
	Events []WithdrawalEvent `json:"events"`
}

// WithdrawalEvent records an auditable status transition for a withdrawal request.
type WithdrawalEvent struct {
	ID             int64  `json:"id"`
	RequestID      int64  `json:"request_id"`
	Action         string `json:"action"`
	OperatorUserID *int64 `json:"operator_user_id,omitempty"`
	Detail         string `json:"detail"`
	CreatedAt      string `json:"created_at"`
}

// MarshalJSON stabilizes decimal-style fields so frontend and acceptance
// outputs stay human-readable instead of exposing floating-point noise.
func (s Summary) MarshalJSON() ([]byte, error) {
	type summaryJSON struct {
		TotalEarned         float64 `json:"total_earned"`
		FrozenAmount        float64 `json:"frozen_amount"`
		InternalTransferred float64 `json:"internal_transferred_amount"`
		PayingAmount        float64 `json:"paying_amount"`
		PaidAmount          float64 `json:"paid_amount"`
		WithdrawableAmount  float64 `json:"withdrawable_amount"`
	}

	return json.Marshal(summaryJSON{
		TotalEarned:         NormalizeDecimalFloat64(s.TotalEarned),
		FrozenAmount:        NormalizeDecimalFloat64(s.FrozenAmount),
		InternalTransferred: NormalizeDecimalFloat64(s.InternalTransferred),
		PayingAmount:        NormalizeDecimalFloat64(s.PayingAmount),
		PaidAmount:          NormalizeDecimalFloat64(s.PaidAmount),
		WithdrawableAmount:  NormalizeDecimalFloat64(s.WithdrawableAmount),
	})
}

// MarshalJSON keeps withdrawal rows consistent with NUMERIC(20,8) storage when
// they are returned through API responses.
func (w WithdrawalRequest) MarshalJSON() ([]byte, error) {
	type withdrawalJSON struct {
		ID                         int64   `json:"id"`
		RequestNo                  string  `json:"request_no"`
		UserID                     int64   `json:"user_id"`
		Amount                     float64 `json:"amount"`
		Status                     string  `json:"status"`
		ApplicantRemark            string  `json:"applicant_remark"`
		PaidChannel                string  `json:"paid_channel"`
		PaidReferenceNo            string  `json:"paid_reference_no"`
		PaidRemark                 string  `json:"paid_remark"`
		PaidAt                     *string `json:"paid_at"`
		CreatedAt                  string  `json:"created_at"`
		SnapshotWithdrawableBefore float64 `json:"snapshot_withdrawable_before"`
		SnapshotWithdrawableAfter  float64 `json:"snapshot_withdrawable_after"`
	}

	return json.Marshal(withdrawalJSON{
		ID:                         w.ID,
		RequestNo:                  w.RequestNo,
		UserID:                     w.UserID,
		Amount:                     NormalizeDecimalFloat64(w.Amount),
		Status:                     w.Status,
		ApplicantRemark:            w.ApplicantRemark,
		PaidChannel:                w.PaidChannel,
		PaidReferenceNo:            w.PaidReferenceNo,
		PaidRemark:                 w.PaidRemark,
		PaidAt:                     w.PaidAt,
		CreatedAt:                  w.CreatedAt,
		SnapshotWithdrawableBefore: NormalizeDecimalFloat64(w.SnapshotWithdrawableBefore),
		SnapshotWithdrawableAfter:  NormalizeDecimalFloat64(w.SnapshotWithdrawableAfter),
	})
}

// MarshalJSON preserves the existing withdrawal row shape while appending the
// audit timeline used by the operator detail dialog.
func (w WithdrawalDetail) MarshalJSON() ([]byte, error) {
	type withdrawalDetailJSON struct {
		ID                         int64             `json:"id"`
		RequestNo                  string            `json:"request_no"`
		UserID                     int64             `json:"user_id"`
		Amount                     float64           `json:"amount"`
		Status                     string            `json:"status"`
		ApplicantRemark            string            `json:"applicant_remark"`
		PaidChannel                string            `json:"paid_channel"`
		PaidReferenceNo            string            `json:"paid_reference_no"`
		PaidRemark                 string            `json:"paid_remark"`
		PaidAt                     *string           `json:"paid_at"`
		CreatedAt                  string            `json:"created_at"`
		SnapshotWithdrawableBefore float64           `json:"snapshot_withdrawable_before"`
		SnapshotWithdrawableAfter  float64           `json:"snapshot_withdrawable_after"`
		Events                     []WithdrawalEvent `json:"events"`
	}

	return json.Marshal(withdrawalDetailJSON{
		ID:                         w.ID,
		RequestNo:                  w.RequestNo,
		UserID:                     w.UserID,
		Amount:                     NormalizeDecimalFloat64(w.Amount),
		Status:                     w.Status,
		ApplicantRemark:            w.ApplicantRemark,
		PaidChannel:                w.PaidChannel,
		PaidReferenceNo:            w.PaidReferenceNo,
		PaidRemark:                 w.PaidRemark,
		PaidAt:                     w.PaidAt,
		CreatedAt:                  w.CreatedAt,
		SnapshotWithdrawableBefore: NormalizeDecimalFloat64(w.SnapshotWithdrawableBefore),
		SnapshotWithdrawableAfter:  NormalizeDecimalFloat64(w.SnapshotWithdrawableAfter),
		Events:                     w.Events,
	})
}
