package models

import "testing"

func TestParseTransaction_ValidDeposit(t *testing.T) {
	input := "DEPOSIT BTC 1.5"
	tx, err := ParseTransaction(input)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if tx.Type != Deposit {
		t.Errorf("Expected type Deposit, got: %s", tx.Type)
	}
	if tx.Asset != BTC {
		t.Errorf("Expected asset BTC, got: %s", tx.Asset)
	}
	if tx.Amount != 1.5 {
		t.Errorf("Expected amount 1.5, got: %f", tx.Amount)
	}
}

func TestParseTransaction_ValidWithdraw(t *testing.T) {
	input := "WITHDRAW ETH 2.0"
	tx, err := ParseTransaction(input)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if tx.Type != Withdraw {
		t.Errorf("Expected type Withdraw, got: %s", tx.Type)
	}
	if tx.Asset != ETH {
		t.Errorf("Expected asset ETH, got: %s", tx.Asset)
	}
	if tx.Amount != 2.0 {
		t.Errorf("Expected amount 2.0, got: %f", tx.Amount)
	}
}

func TestParseTransaction_CaseInsensitive(t *testing.T) {
	input := "deposit usd 100.50"
	tx, err := ParseTransaction(input)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if tx.Type != Deposit {
		t.Errorf("Expected type Deposit, got: %s", tx.Type)
	}
	if tx.Asset != USD {
		t.Errorf("Expected asset USD, got: %s", tx.Asset)
	}
	if tx.Amount != 100.50 {
		t.Errorf("Expected amount 100.50, got: %f", tx.Amount)
	}
}

func TestParseTransaction_InvalidFormat(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"Empty", ""},
		{"TooFewArgs", "DEPOSIT BTC"},
		{"TooManyArgs", "DEPOSIT BTC 1.5 EXTRA"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseTransaction(tc.input)
			if err == nil {
				t.Errorf("Expected error for input: %s", tc.input)
			}
		})
	}
}

func TestParseTransaction_InvalidType(t *testing.T) {
	input := "TRANSFER BTC 1.5"
	_, err := ParseTransaction(input)

	if err == nil {
		t.Error("Expected error for invalid transaction type")
	}
}

func TestParseTransaction_InvalidAsset(t *testing.T) {
	input := "DEPOSIT XRP 1.5"
	_, err := ParseTransaction(input)

	if err == nil {
		t.Error("Expected error for invalid asset type")
	}
}

func TestParseTransaction_InvalidAmount(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"NotANumber", "DEPOSIT BTC abc"},
		{"NegativeAmount", "DEPOSIT BTC -1.5"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseTransaction(tc.input)
			if err == nil {
				t.Errorf("Expected error for input: %s", tc.input)
			}
		})
	}
}
