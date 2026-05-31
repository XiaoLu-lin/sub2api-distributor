package distributor

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
)

const (
	distributorAffiliateCodeLength      = 12
	distributorAffiliateCodeMaxAttempts = 12
)

var distributorAffiliateCodeCharset = []byte("ABCDEFGHJKLMNPQRSTUVWXYZ23456789")

type affiliateIdentityQueryExecer interface {
	QueryRowContext(context.Context, string, ...any) *sql.Row
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}

// ensureAffiliateIdentity makes sure the target user has a source-system
// affiliate identity row before distributor invite features are used.
func (s *Service) ensureAffiliateIdentity(ctx context.Context, userID int64) error {
	return ensureAffiliateIdentityWithExecer(ctx, s.db, userID)
}

func ensureAffiliateIdentityWithExecer(ctx context.Context, execer affiliateIdentityQueryExecer, userID int64) error {
	if userID <= 0 {
		return ErrInvalidProfile
	}

	var existingCode string
	err := execer.QueryRowContext(ctx, `
SELECT aff_code
FROM user_affiliates
WHERE user_id = $1
`, userID).Scan(&existingCode)
	if err == nil {
		return nil
	}
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("query affiliate identity: %w", err)
	}

	for attempt := 0; attempt < distributorAffiliateCodeMaxAttempts; attempt++ {
		code, codeErr := generateDistributorAffiliateCode()
		if codeErr != nil {
			return codeErr
		}

		_, insertErr := execer.ExecContext(ctx, `
INSERT INTO user_affiliates (user_id, aff_code, created_at, updated_at)
VALUES ($1, $2, NOW(), NOW())
ON CONFLICT (user_id) DO NOTHING
`, userID, code)
		if insertErr == nil {
			return nil
		}
		if isAffiliateCodeConflict(insertErr) {
			continue
		}
		return fmt.Errorf("insert affiliate identity: %w", insertErr)
	}

	return fmt.Errorf("generate affiliate identity: exceeded retry limit")
}

func generateDistributorAffiliateCode() (string, error) {
	buf := make([]byte, distributorAffiliateCodeLength)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate affiliate code entropy: %w", err)
	}

	for i, value := range buf {
		buf[i] = distributorAffiliateCodeCharset[int(value)%len(distributorAffiliateCodeCharset)]
	}
	return string(buf), nil
}

func isAffiliateCodeConflict(err error) bool {
	if err == nil {
		return false
	}
	return containsAny(err.Error(),
		"duplicate key value violates unique constraint",
		"user_affiliates_aff_code_key",
	)
}

func containsAny(value string, patterns ...string) bool {
	for _, pattern := range patterns {
		if pattern != "" && value != "" && containsString(value, pattern) {
			return true
		}
	}
	return false
}

func containsString(value string, pattern string) bool {
	if len(pattern) == 0 {
		return true
	}
	if len(value) < len(pattern) {
		return false
	}
	for i := 0; i <= len(value)-len(pattern); i++ {
		if value[i:i+len(pattern)] == pattern {
			return true
		}
	}
	return false
}
