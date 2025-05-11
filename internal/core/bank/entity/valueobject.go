package entity

import (
	"fmt"
	"regexp"
)

// Currency представляет тип валюты.
type Currency string

const (
	RUB Currency = "RUB"
)

// Validate проверяет, является ли валюта допустимой.
func (c Currency) Validate() error {
	if c != RUB {
		return fmt.Errorf("invalid currency: %s, only RUB is supported", c)
	}
	return nil
}

func (c Currency) String() string {
	return string(c)
}

// AccountNumber представляет номер банковского счета.
type AccountNumber string

// Validate проверяет корректность номера счета (например, должен содержать только цифры и быть длиной 20 символов).
func (a AccountNumber) Validate() error {
	re := regexp.MustCompile(`^\d{20}$`)
	if !re.MatchString(string(a)) {
		return fmt.Errorf("invalid account number: %s, must be 20 digits", a)
	}
	return nil
}

func (a AccountNumber) String() string {
	return string(a)
}
