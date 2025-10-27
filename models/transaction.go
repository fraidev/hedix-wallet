package models

import (
	"fmt"
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
type Transaction struct {
	Type   TransactionType
	Asset  Asset
	Amount float64
}

// ParseTransaction parses a transaction from a string input
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

	amount, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return Transaction{}, fmt.Errorf("invalid amount: %s", err)
	}

	if amount < 0 {
		return Transaction{}, fmt.Errorf("amount must be positive")
	}

	return Transaction{
		Type:   txType,
		Asset:  asset,
		Amount: amount,
	}, nil
}
