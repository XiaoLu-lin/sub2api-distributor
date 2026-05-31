package distributor

import (
	"math/big"
	"strconv"
	"strings"
)

// WithdrawalStatus models the linear state machine used by offline payout requests.
type WithdrawalStatus string

const (
	// WithdrawalStatusPaying means the request has been submitted and is waiting for offline payout completion.
	WithdrawalStatusPaying WithdrawalStatus = "paying"
	// WithdrawalStatusPaid means the operator confirmed the offline payout.
	WithdrawalStatusPaid WithdrawalStatus = "paid"
	// WithdrawalStatusCancelled means the request was voided before payout completion.
	WithdrawalStatusCancelled WithdrawalStatus = "cancelled"
)

// ComputeWithdrawableAmount derives the currently claimable rebate balance from
// total earnings after subtracting frozen, transferred, paying, and paid amounts.
func ComputeWithdrawableAmount(
	totalEarned float64,
	frozenAmount float64,
	internalTransferredAmount float64,
	payingAmount float64,
	paidAmount float64,
) float64 {
	withdrawable := totalEarned - frozenAmount - internalTransferredAmount - payingAmount - paidAmount
	if withdrawable < 0 {
		return 0
	}
	return withdrawable
}

// ComputeWithdrawableAmountDecimal performs the same withdrawable-balance
// calculation using exact decimal string arithmetic.
func ComputeWithdrawableAmountDecimal(
	totalEarned string,
	frozenAmount string,
	internalTransferredAmount string,
	payingAmount string,
	paidAmount string,
) string {
	total := mustDecimalRat(totalEarned)
	total.Sub(total, mustDecimalRat(frozenAmount))
	total.Sub(total, mustDecimalRat(internalTransferredAmount))
	total.Sub(total, mustDecimalRat(payingAmount))
	total.Sub(total, mustDecimalRat(paidAmount))
	if total.Sign() < 0 {
		return "0"
	}
	return normalizeDecimalRat(total)
}

// DecimalAmountWithinBalance compares two decimal strings without floating-point
// conversion and returns whether amount is less than or equal to balance.
func DecimalAmountWithinBalance(amount string, balance string) bool {
	return mustDecimalRat(amount).Cmp(mustDecimalRat(balance)) <= 0
}

// AddDecimalStrings returns the exact decimal sum of two decimal strings.
func AddDecimalStrings(left string, right string) string {
	sum := mustDecimalRat(left)
	sum.Add(sum, mustDecimalRat(right))
	return normalizeDecimalRat(sum)
}

// DecimalStringToFloat64 converts a decimal string into a float64 only at the
// response boundary.
func DecimalStringToFloat64(value string) float64 {
	number, _ := mustDecimalRat(value).Float64()
	return number
}

// DecimalStringFromFloat64 converts a request amount into normalized decimal text.
func DecimalStringFromFloat64(value float64) string {
	return normalizeDecimalRat(mustDecimalRat(strconv.FormatFloat(value, 'f', -1, 64)))
}

// NormalizeDecimalFloat64 rounds a float through the same decimal-normalization
// path used for request and response amounts so JSON output stays stable.
func NormalizeDecimalFloat64(value float64) float64 {
	return DecimalStringToFloat64(DecimalStringFromFloat64(value))
}

// CanTransitionWithdrawalStatus validates whether a request can move from its
// current state to the target state.
func CanTransitionWithdrawalStatus(current WithdrawalStatus, next WithdrawalStatus) bool {
	switch current {
	case WithdrawalStatusPaying:
		return next == WithdrawalStatusPaid || next == WithdrawalStatusCancelled
	case WithdrawalStatusPaid, WithdrawalStatusCancelled:
		return false
	default:
		return false
	}
}

// NormalizeUserLookupKeyword trims operator search input before it is used in SQL filters.
func NormalizeUserLookupKeyword(keyword string) string {
	return strings.TrimSpace(keyword)
}

// IsValidDistributorStatus constrains profile status values to the supported enum.
func IsValidDistributorStatus(status string) bool {
	switch strings.TrimSpace(status) {
	case "active", "disabled":
		return true
	default:
		return false
	}
}

func mustDecimalRat(value string) *big.Rat {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return big.NewRat(0, 1)
	}

	rat, ok := new(big.Rat).SetString(trimmed)
	if !ok {
		return big.NewRat(0, 1)
	}
	return rat
}

func normalizeDecimalRat(value *big.Rat) string {
	if value == nil || value.Sign() == 0 {
		return "0"
	}

	text := value.FloatString(8)
	if strings.Contains(text, ".") {
		parts := strings.SplitN(text, ".", 2)
		fraction := strings.TrimRight(parts[1], "0")
		if len(fraction) == 1 {
			return parts[0] + "." + fraction + "0"
		}
	}
	text = strings.TrimRight(text, "0")
	text = strings.TrimRight(text, ".")
	if text == "" || text == "-0" {
		return "0"
	}
	return text
}
