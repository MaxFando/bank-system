package entity

import "fmt"

var (
	ErrInsufficientFunds     = fmt.Errorf("insufficient funds")
	ErrDepositNegativeAmount = fmt.Errorf("deposit amount must be positive")

	ErrCreditNotActive      = fmt.Errorf("credit is not active")
	ErrCreditAmountExceeded = fmt.Errorf("credit amount exceeds limit")
)
