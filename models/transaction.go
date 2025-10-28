package models

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// TransactionType represents the type of transaction
type TransactionType string

const (
	Deposit  TransactionType = "DEPOSIT"
	Withdraw TransactionType = "WITHDRAW"
)

// Transaction represents a single wallet transaction
// Amount is stored as the smallest unit (satoshis, wei, cents)
type Transaction struct {
	Type   TransactionType
	Asset  Asset
	Amount int64 // Smallest unit: satoshis for BTC, wei for ETH, cents for USD
}

// ParseTransaction parses a transaction from a string input
// Input amount is in the main unit (BTC, ETH, USD) and is converted to smallest unit
func ParseTransaction(input string) (Transaction, error) {
	parts := strings.Fields(input)
	if len(parts) != 3 {
		return Transaction{}, fmt.Errorf("invalid format. Expected: <TYPE> <ASSET> <AMOUNT>")
	}

	txType := TransactionType(strings.ToUpper(parts[0]))
	if txType != Deposit && txType != Withdraw {
		return Transaction{}, fmt.Errorf("invalid transaction type. Must be DEPOSIT or WITHDRAW")
	}

	asset := Asset(strings.ToUpper(parts[1]))
	if asset != BTC && asset != ETH && asset != USD {
		return Transaction{}, fmt.Errorf("invalid asset. Must be BTC, ETH, or USD")
	}

	// Parse as float to handle decimal input
	amountFloat, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return Transaction{}, fmt.Errorf("invalid amount: %s", err)
	}

	if amountFloat < 0 {
		return Transaction{}, fmt.Errorf("amount must be positive")
	}

	// Convert to smallest unit based on asset decimals
	decimals := asset.GetDecimals()
	multiplier := math.Pow(10, float64(decimals))
	amountSmallestUnit := int64(amountFloat * multiplier)

	return Transaction{
		Type:   txType,
		Asset:  asset,
		Amount: amountSmallestUnit,
	}, nil
}

// FormatAmount formats the amount from smallest unit to human-readable string
func (t Transaction) FormatAmount() string {
	decimals := t.Asset.GetDecimals()
	divisor := math.Pow(10, float64(decimals))
	amount := float64(t.Amount) / divisor
	return fmt.Sprintf("%.*f", decimals, amount)
}
