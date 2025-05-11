package entity

import "fmt"

var (
	ErrInsufficientFunds     = fmt.Errorf("insufficient funds")
	ErrDepositNegativeAmount = fmt.Errorf("deposit amount must be positive")
)
