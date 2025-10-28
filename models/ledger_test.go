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
		Amount: 150000000, // 1.5 BTC in satoshis
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
	if transactions[0].Amount != 150000000 {
		t.Errorf("Expected 150000000, got %d", transactions[0].Amount)
	}
}

func TestLedger_CalculateBalance_SingleDeposit(t *testing.T) {
	ledger := NewLedger()

	ledger.AddTransaction(Transaction{
		Type:   Deposit,
		Asset:  BTC,
		Amount: 250000000, // 2.5 BTC
	})

	balance := ledger.CalculateBalance(BTC)
	expected := int64(250000000)
	if balance != expected {
		t.Errorf("Expected balance %d, got %d", expected, balance)
	}
}

func TestLedger_CalculateBalance_MultipleTransactions(t *testing.T) {
	ledger := NewLedger()

	ledger.AddTransaction(Transaction{Type: Deposit, Asset: BTC, Amount: 300000000})  // 3.0 BTC
	ledger.AddTransaction(Transaction{Type: Deposit, Asset: BTC, Amount: 200000000})  // 2.0 BTC
	ledger.AddTransaction(Transaction{Type: Withdraw, Asset: BTC, Amount: 150000000}) // 1.5 BTC

	balance := ledger.CalculateBalance(BTC)
	expected := int64(350000000) // 3.0 + 2.0 - 1.5 = 3.5 BTC
	if balance != expected {
		t.Errorf("Expected balance %d, got %d", expected, balance)
	}
}

func TestLedger_CalculateBalance_MultipleAssets(t *testing.T) {
	ledger := NewLedger()

	ledger.AddTransaction(Transaction{Type: Deposit, Asset: BTC, Amount: 100000000})                 // 1.0 BTC
	ledger.AddTransaction(Transaction{Type: Deposit, Asset: ETH, Amount: 2000000000000000000})       // 2.0 ETH
	ledger.AddTransaction(Transaction{Type: Deposit, Asset: USD, Amount: 10000})                     // 100.00 USD

	btcBalance := ledger.CalculateBalance(BTC)
	ethBalance := ledger.CalculateBalance(ETH)
	usdBalance := ledger.CalculateBalance(USD)

	if btcBalance != 100000000 {
		t.Errorf("Expected BTC balance 100000000, got %d", btcBalance)
	}
	if ethBalance != 2000000000000000000 {
		t.Errorf("Expected ETH balance 2000000000000000000, got %d", ethBalance)
	}
	if usdBalance != 10000 {
		t.Errorf("Expected USD balance 10000, got %d", usdBalance)
	}
}

func TestLedger_CalculateAllBalances(t *testing.T) {
	ledger := NewLedger()

	ledger.AddTransaction(Transaction{Type: Deposit, Asset: BTC, Amount: 100000000})           // 1.0 BTC
	ledger.AddTransaction(Transaction{Type: Deposit, Asset: ETH, Amount: 2000000000000000000}) // 2.0 ETH
	ledger.AddTransaction(Transaction{Type: Deposit, Asset: USD, Amount: 10000})               // 100.00 USD

	balances := ledger.CalculateAllBalances()

	if balances[BTC] != 100000000 {
		t.Errorf("Expected BTC balance 100000000, got %d", balances[BTC])
	}
	if balances[ETH] != 2000000000000000000 {
		t.Errorf("Expected ETH balance 2000000000000000000, got %d", balances[ETH])
	}
	if balances[USD] != 10000 {
		t.Errorf("Expected USD balance 10000, got %d", balances[USD])
	}
}

func TestLedger_EmptyBalance(t *testing.T) {
	ledger := NewLedger()

	balance := ledger.CalculateBalance(BTC)
	if balance != 0 {
		t.Errorf("Expected balance 0 for empty ledger, got %d", balance)
	}
}
