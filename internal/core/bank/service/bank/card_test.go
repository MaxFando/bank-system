package bank

import (
	"testing"
)

func TestGenerateCVV(t *testing.T) {
	cvv := generateCVV()
	if len(cvv) != 3 {
		t.Errorf("invalid CVV length, got %d, want 3", len(cvv))
	}
}

func TestGenerateCardNumber(t *testing.T) {
	cardNumber := generateCardNumber(42)
	if len(cardNumber) != 16 {
		t.Errorf("invalid card number length, got %d, want 16", len(cardNumber))
	}
}
