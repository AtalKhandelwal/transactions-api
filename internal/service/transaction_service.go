package service

import "math"

func NormalizeAmount(operationTypeID int, amount float64) float64 {
	if operationTypeID == 4 {
		return math.Abs(amount)
	}
	return -math.Abs(amount)
}
