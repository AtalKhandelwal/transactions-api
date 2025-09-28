package service

import "math"

const (
	OpCashPurchase        = 1
	OpInstallmentPurchase = 2
	OpWithdrawal          = 3
	OpPayment             = 4
)

func NormalizeAmount(operationTypeID int, amount float64) float64 {
	if operationTypeID == OpPayment {
		return math.Abs(amount)
	}
	return -math.Abs(amount)
}
