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
	// 1.5 BTC = 150000000 satoshis
	expected := int64(150000000)
	if tx.Amount != expected {
		t.Errorf("Expected amount %d satoshis, got: %d", expected, tx.Amount)
	}
}

func TestParseTransaction_ValidWithdraw(t *testing.T) {
	input := "WITHDRAW USD 100.50"
	tx, err := ParseTransaction(input)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if tx.Type != Withdraw {
		t.Errorf("Expected type Withdraw, got: %s", tx.Type)
	}
	if tx.Asset != USD {
		t.Errorf("Expected asset USD, got: %s", tx.Asset)
	}
	// 100.50 USD = 10050 cents
	expected := int64(10050)
	if tx.Amount != expected {
		t.Errorf("Expected amount %d cents, got: %d", expected, tx.Amount)
	}
}

func TestParseTransaction_CaseInsensitive(t *testing.T) {
	input := "deposit eth 2.5"
	tx, err := ParseTransaction(input)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if tx.Type != Deposit {
		t.Errorf("Expected type Deposit, got: %s", tx.Type)
	}
	if tx.Asset != ETH {
		t.Errorf("Expected asset ETH, got: %s", tx.Asset)
	}
	// 2.5 ETH = 2.5 * 10^18 wei
	expected := int64(2500000000000000000)
	if tx.Amount != expected {
		t.Errorf("Expected amount %d wei, got: %d", expected, tx.Amount)
	}
}

func TestParseTransaction_SmallAmount(t *testing.T) {
	input := "DEPOSIT BTC 0.00000001"
	tx, err := ParseTransaction(input)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// 0.00000001 BTC = 1 satoshi
	expected := int64(1)
	if tx.Amount != expected {
		t.Errorf("Expected amount %d satoshi, got: %d", expected, tx.Amount)
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

func TestTransaction_FormatAmount(t *testing.T) {
	testCases := []struct {
		name     string
		tx       Transaction
		expected string
	}{
		{
			name: "BTC",
			tx: Transaction{
				Type:   Deposit,
				Asset:  BTC,
				Amount: 150000000, // 1.5 BTC
			},
			expected: "1.50000000",
		},
		{
			name: "USD",
			tx: Transaction{
				Type:   Deposit,
				Asset:  USD,
				Amount: 10050, // 100.50 USD
			},
			expected: "100.50",
		},
		{
			name: "ETH",
			tx: Transaction{
				Type:   Deposit,
				Asset:  ETH,
				Amount: 2500000000000000000, // 2.5 ETH
			},
			expected: "2.500000000000000000",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.tx.FormatAmount()
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}
