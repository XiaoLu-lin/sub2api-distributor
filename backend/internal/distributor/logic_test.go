package distributor

import (
	"encoding/json"
	"testing"
)

func TestComputeWithdrawableAmount(t *testing.T) {
	tests := []struct {
		name             string
		totalEarned      float64
		frozenAmount     float64
		internalTransfer float64
		payingAmount     float64
		paidAmount       float64
		wantWithdrawable float64
	}{
		{
			name:             "all earned amount is withdrawable",
			totalEarned:      120,
			wantWithdrawable: 120,
		},
		{
			name:             "frozen and settled amounts are excluded",
			totalEarned:      200,
			frozenAmount:     20,
			internalTransfer: 30,
			payingAmount:     40,
			paidAmount:       50,
			wantWithdrawable: 60,
		},
		{
			name:             "result is clamped to zero",
			totalEarned:      100,
			frozenAmount:     30,
			internalTransfer: 50,
			payingAmount:     60,
			wantWithdrawable: 0,
		},
		{
			name:             "paid amount also reduces withdrawable balance",
			totalEarned:      300,
			paidAmount:       120,
			wantWithdrawable: 180,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ComputeWithdrawableAmount(
				tt.totalEarned,
				tt.frozenAmount,
				tt.internalTransfer,
				tt.payingAmount,
				tt.paidAmount,
			)
			if got != tt.wantWithdrawable {
				t.Fatalf("ComputeWithdrawableAmount() = %v, want %v", got, tt.wantWithdrawable)
			}
		})
	}
}

func TestComputeWithdrawableAmountDecimal(t *testing.T) {
	tests := []struct {
		name                 string
		totalEarned          string
		frozenAmount         string
		internalTransferred  string
		payingAmount         string
		paidAmount           string
		wantWithdrawableText string
	}{
		{
			name:                 "preserves decimal precision across multiple subtractions",
			totalEarned:          "0.30",
			frozenAmount:         "0",
			internalTransferred:  "0",
			payingAmount:         "0.10",
			paidAmount:           "0.20",
			wantWithdrawableText: "0",
		},
		{
			name:                 "keeps a cent-accurate result",
			totalEarned:          "42.50",
			frozenAmount:         "0",
			internalTransferred:  "0",
			payingAmount:         "10.10",
			paidAmount:           "8.20",
			wantWithdrawableText: "24.20",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ComputeWithdrawableAmountDecimal(
				tt.totalEarned,
				tt.frozenAmount,
				tt.internalTransferred,
				tt.payingAmount,
				tt.paidAmount,
			)
			if got != tt.wantWithdrawableText {
				t.Fatalf("ComputeWithdrawableAmountDecimal() = %s, want %s", got, tt.wantWithdrawableText)
			}
		})
	}
}

func TestDecimalAmountWithinBalance(t *testing.T) {
	tests := []struct {
		name       string
		amount     string
		balance    string
		wantAccept bool
	}{
		{
			name:       "accepts amount equal to balance",
			amount:     "15.50",
			balance:    "15.50",
			wantAccept: true,
		},
		{
			name:       "rejects amount larger than balance by one cent",
			amount:     "15.51",
			balance:    "15.50",
			wantAccept: false,
		},
		{
			name:       "accepts decimal-safe amount without float noise",
			amount:     "0.30",
			balance:    "0.30",
			wantAccept: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DecimalAmountWithinBalance(tt.amount, tt.balance)
			if got != tt.wantAccept {
				t.Fatalf("DecimalAmountWithinBalance(%s, %s) = %v, want %v", tt.amount, tt.balance, got, tt.wantAccept)
			}
		})
	}
}

func TestNormalizeDecimalFloat64(t *testing.T) {
	tests := []struct {
		name  string
		input float64
		want  string
	}{
		{
			name:  "normalizes recurring float noise to a stable decimal",
			input: 8.399999999999999,
			want:  "8.4",
		},
		{
			name:  "preserves cents for common two-decimal values",
			input: 7.169999999999998,
			want:  "7.17",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload, err := json.Marshal(NormalizeDecimalFloat64(tt.input))
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}
			got := string(payload)
			if got != tt.want {
				t.Fatalf("NormalizeDecimalFloat64(%v) marshals to %s, want %s", tt.input, got, tt.want)
			}
		})
	}
}

func TestCanTransitionWithdrawalStatus(t *testing.T) {
	tests := []struct {
		name      string
		current   WithdrawalStatus
		next      WithdrawalStatus
		wantAllow bool
	}{
		{
			name:      "paying can become paid",
			current:   WithdrawalStatusPaying,
			next:      WithdrawalStatusPaid,
			wantAllow: true,
		},
		{
			name:      "paying can become cancelled",
			current:   WithdrawalStatusPaying,
			next:      WithdrawalStatusCancelled,
			wantAllow: true,
		},
		{
			name:      "paid cannot transition again",
			current:   WithdrawalStatusPaid,
			next:      WithdrawalStatusCancelled,
			wantAllow: false,
		},
		{
			name:      "cancelled cannot transition again",
			current:   WithdrawalStatusCancelled,
			next:      WithdrawalStatusPaid,
			wantAllow: false,
		},
		{
			name:      "paying cannot stay paying",
			current:   WithdrawalStatusPaying,
			next:      WithdrawalStatusPaying,
			wantAllow: false,
		},
		{
			name:      "unknown status cannot transition",
			current:   WithdrawalStatus("draft"),
			next:      WithdrawalStatusPaid,
			wantAllow: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CanTransitionWithdrawalStatus(tt.current, tt.next)
			if got != tt.wantAllow {
				t.Fatalf("CanTransitionWithdrawalStatus(%q, %q) = %v, want %v", tt.current, tt.next, got, tt.wantAllow)
			}
		})
	}
}

func TestNormalizeUserLookupKeyword(t *testing.T) {
	tests := []struct {
		name    string
		keyword string
		want    string
	}{
		{
			name:    "trim surrounding spaces",
			keyword: "  demo@example.com  ",
			want:    "demo@example.com",
		},
		{
			name:    "keep internal spaces",
			keyword: "demo user",
			want:    "demo user",
		},
		{
			name:    "empty stays empty",
			keyword: "   ",
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeUserLookupKeyword(tt.keyword)
			if got != tt.want {
				t.Fatalf("NormalizeUserLookupKeyword(%q) = %q, want %q", tt.keyword, got, tt.want)
			}
		})
	}
}

func TestIsValidDistributorStatus(t *testing.T) {
	tests := []struct {
		name   string
		status string
		want   bool
	}{
		{name: "active is allowed", status: "active", want: true},
		{name: "disabled is allowed", status: "disabled", want: true},
		{name: "unknown is rejected", status: "pending", want: false},
		{name: "empty is rejected", status: "", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidDistributorStatus(tt.status)
			if got != tt.want {
				t.Fatalf("IsValidDistributorStatus(%q) = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}
