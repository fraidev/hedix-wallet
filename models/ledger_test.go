package models

import "testing"

func TestNewLedger(t *testing.T) {
	ledger := NewLedger()

	if ledger == nil {
		t.Fatal("Expected non-nil ledger")
	}

	if len(ledger.GetTransactions()) != 0 {
		t.Errorf("Expected empty ledger, got %d transactions", len(ledger.GetTransactions()))
	}
}

func TestLedger_AddTransaction(t *testing.T) {
	ledger := NewLedger()

	tx := Transaction{
		Type:   Deposit,
		Asset:  BTC,
		Amount: 1.5,
	}

	ledger.AddTransaction(tx)

	transactions := ledger.GetTransactions()
	if len(transactions) != 1 {
		t.Fatalf("Expected 1 transaction, got %d", len(transactions))
	}

	if transactions[0].Type != Deposit {
		t.Errorf("Expected Deposit, got %s", transactions[0].Type)
	}
	if transactions[0].Asset != BTC {
		t.Errorf("Expected BTC, got %s", transactions[0].Asset)
	}
	if transactions[0].Amount != 1.5 {
		t.Errorf("Expected 1.5, got %f", transactions[0].Amount)
	}
}

func TestLedger_CalculateBalance_SingleDeposit(t *testing.T) {
	ledger := NewLedger()

	ledger.AddTransaction(Transaction{
		Type:   Deposit,
		Asset:  BTC,
		Amount: 2.5,
	})

	balance := ledger.CalculateBalance(BTC)
	if balance != 2.5 {
		t.Errorf("Expected balance 2.5, got %f", balance)
	}
}

func TestLedger_CalculateBalance_MultipleTransactions(t *testing.T) {
	ledger := NewLedger()

	ledger.AddTransaction(Transaction{Type: Deposit, Asset: BTC, Amount: 3.0})
	ledger.AddTransaction(Transaction{Type: Deposit, Asset: BTC, Amount: 2.0})
	ledger.AddTransaction(Transaction{Type: Withdraw, Asset: BTC, Amount: 1.5})

	balance := ledger.CalculateBalance(BTC)
	expected := 3.5 // 3.0 + 2.0 - 1.5
	if balance != expected {
		t.Errorf("Expected balance %f, got %f", expected, balance)
	}
}

func TestLedger_CalculateBalance_MultipleAssets(t *testing.T) {
	ledger := NewLedger()

	ledger.AddTransaction(Transaction{Type: Deposit, Asset: BTC, Amount: 1.0})
	ledger.AddTransaction(Transaction{Type: Deposit, Asset: ETH, Amount: 2.0})
	ledger.AddTransaction(Transaction{Type: Deposit, Asset: USD, Amount: 100.0})

	btcBalance := ledger.CalculateBalance(BTC)
	ethBalance := ledger.CalculateBalance(ETH)
	usdBalance := ledger.CalculateBalance(USD)

	if btcBalance != 1.0 {
		t.Errorf("Expected BTC balance 1.0, got %f", btcBalance)
	}
	if ethBalance != 2.0 {
		t.Errorf("Expected ETH balance 2.0, got %f", ethBalance)
	}
	if usdBalance != 100.0 {
		t.Errorf("Expected USD balance 100.0, got %f", usdBalance)
	}
}

func TestLedger_CalculateAllBalances(t *testing.T) {
	ledger := NewLedger()

	ledger.AddTransaction(Transaction{Type: Deposit, Asset: BTC, Amount: 1.0})
	ledger.AddTransaction(Transaction{Type: Deposit, Asset: ETH, Amount: 2.0})
	ledger.AddTransaction(Transaction{Type: Deposit, Asset: USD, Amount: 100.0})

	balances := ledger.CalculateAllBalances()

	if balances[BTC] != 1.0 {
		t.Errorf("Expected BTC balance 1.0, got %f", balances[BTC])
	}
	if balances[ETH] != 2.0 {
		t.Errorf("Expected ETH balance 2.0, got %f", balances[ETH])
	}
	if balances[USD] != 100.0 {
		t.Errorf("Expected USD balance 100.0, got %f", balances[USD])
	}
}

func TestLedger_EmptyBalance(t *testing.T) {
	ledger := NewLedger()

	balance := ledger.CalculateBalance(BTC)
	if balance != 0.0 {
		t.Errorf("Expected balance 0.0 for empty ledger, got %f", balance)
	}
}
