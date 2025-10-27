package services

import (
	"fmt"

	"github.com/fraidev/hedix-wallet/models"
)

// Wallet represents an in-memory wallet with ledger-based storage
type Wallet struct {
	ledger *models.Ledger
}

// NewWallet creates a new wallet
func NewWallet() *Wallet {
	return &Wallet{
		ledger: models.NewLedger(),
	}
}

// ProcessTransaction processes a transaction attempt in the ledger
// It validates the transaction based on current balance and records the result
func (w *Wallet) ProcessTransaction(tx models.Transaction) error {
	// Calculate current balance for the asset
	currentBalance := w.ledger.CalculateBalance(tx.Asset)

	var err error

	switch tx.Type {
	case models.Deposit:
		// Deposits always succeed
		err = nil
	case models.Withdraw:
		// Withdrawals only succeed if there are sufficient funds
		if currentBalance < tx.Amount {
			err = fmt.Errorf("insufficient funds for withdrawal: requested %.2f, available %.2f", tx.Amount, currentBalance)
		}
	default:
		err = fmt.Errorf("unknown transaction type: %s", tx.Type)
	}

	// Record the transaction in the ledger
	if err == nil {
		w.ledger.AddTransaction(tx)
	}

	return err
}

// GetBalance returns the current balance for a specific asset
func (w *Wallet) GetBalance(asset models.Asset) float64 {
	return w.ledger.CalculateBalance(asset)
}

// GetAllBalances returns balances for all assets
func (w *Wallet) GetAllBalances() map[models.Asset]float64 {
	return w.ledger.CalculateAllBalances()
}

// GetLedger returns the underlying ledger (for testing/debugging)
func (w *Wallet) GetLedger() *models.Ledger {
	return w.ledger
}

// GetTransactionHistory returns all ledger entries
func (w *Wallet) GetTransactionHistory() []models.Transaction {
	return w.ledger.GetTransactions()
}

// String returns a string representation of the wallet balances
func (w *Wallet) String() string {
	balances := w.GetAllBalances()
	return fmt.Sprintf("BTC: %.8f | ETH: %.8f | USD: %.2f",
		balances[models.BTC], balances[models.ETH], balances[models.USD])
}
