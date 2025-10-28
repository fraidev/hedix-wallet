package models

// Ledger represents a transaction ledger that stores all transaction history
type Ledger struct {
	transactions []Transaction
}

// NewLedger creates a new empty ledger
func NewLedger() *Ledger {
	return &Ledger{
		transactions: make([]Transaction, 0),
	}
}

// AddTransaction adds a new transaction entry to the ledger
func (l *Ledger) AddTransaction(tx Transaction) {
	l.transactions = append(l.transactions, tx)
}

// GetTransactions returns all transactions in the ledger
func (l *Ledger) GetTransactions() []Transaction {
	return l.transactions
}

// CalculateBalance calculates the current balance for a specific asset
// by replaying all successful transactions from the ledger
// Returns balance in smallest unit (satoshis, wei, cents)
func (l *Ledger) CalculateBalance(asset Asset) int64 {
	var balance int64 = 0

	for _, transaction := range l.transactions {
		// Only process successful transactions for the requested asset
		if transaction.Asset != asset {
			continue
		}

		switch transaction.Type {
		case Deposit:
			balance += transaction.Amount
		case Withdraw:
			balance -= transaction.Amount
		}
	}

	return balance
}

// CalculateAllBalances calculates balances for all assets
// Returns balances in smallest units
func (l *Ledger) CalculateAllBalances() map[Asset]int64 {
	return map[Asset]int64{
		BTC: l.CalculateBalance(BTC),
		ETH: l.CalculateBalance(ETH),
		USD: l.CalculateBalance(USD),
	}
}
