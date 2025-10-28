package services

import (
	"testing"

	"github.com/fraidev/hedix-wallet/models"
)

func TestNewWallet(t *testing.T) {
	wallet := NewWallet()

	if wallet == nil {
		t.Fatal("Expected non-nil wallet")
	}

	if wallet.GetBalance(models.BTC) != 0.0 {
		t.Error("Expected BTC balance to be 0.0 for new wallet")
	}
	if wallet.GetBalance(models.ETH) != 0.0 {
		t.Error("Expected ETH balance to be 0.0 for new wallet")
	}
	if wallet.GetBalance(models.USD) != 0.0 {
		t.Error("Expected USD balance to be 0.0 for new wallet")
	}
}

func TestWallet_DepositSuccess(t *testing.T) {
	wallet := NewWallet()

	tx := models.Transaction{
		Type:   models.Deposit,
		Asset:  models.BTC,
		Amount: 2.5,
	}

	err := wallet.ProcessTransaction(tx)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	balance := wallet.GetBalance(models.BTC)
	if balance != 2.5 {
		t.Errorf("Expected balance 2.5, got %f", balance)
	}
}

func TestWallet_WithdrawSuccess(t *testing.T) {
	wallet := NewWallet()

	// First deposit
	wallet.ProcessTransaction(models.Transaction{
		Type:   models.Deposit,
		Asset:  models.BTC,
		Amount: 5.0,
	})

	// Then withdraw
	tx := models.Transaction{
		Type:   models.Withdraw,
		Asset:  models.BTC,
		Amount: 2.0,
	}

	err := wallet.ProcessTransaction(tx)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	balance := wallet.GetBalance(models.BTC)
	expected := 3.0
	if balance != expected {
		t.Errorf("Expected balance %f, got %f", expected, balance)
	}
}

func TestWallet_WithdrawInsufficientFunds(t *testing.T) {
	wallet := NewWallet()

	// Deposit small amount
	wallet.ProcessTransaction(models.Transaction{
		Type:   models.Deposit,
		Asset:  models.BTC,
		Amount: 1.0,
	})

	// Try to withdraw more than available
	tx := models.Transaction{
		Type:   models.Withdraw,
		Asset:  models.BTC,
		Amount: 2.0,
	}

	err := wallet.ProcessTransaction(tx)
	if err == nil {
		t.Fatal("Expected error for insufficient funds")
	}

	// Balance should remain unchanged
	balance := wallet.GetBalance(models.BTC)
	if balance != 1.0 {
		t.Errorf("Expected balance 1.0 (unchanged), got %f", balance)
	}

	// Transaction should not be recorded
	history := wallet.GetTransactionHistory()
	if len(history) != 1 {
		t.Errorf("Expected 1 transaction in history (only deposit), got %d", len(history))
	}
}

func TestWallet_WithdrawFromEmptyWallet(t *testing.T) {
	wallet := NewWallet()

	tx := models.Transaction{
		Type:   models.Withdraw,
		Asset:  models.BTC,
		Amount: 1.0,
	}

	err := wallet.ProcessTransaction(tx)
	if err == nil {
		t.Fatal("Expected error for withdrawing from empty wallet")
	}

	balance := wallet.GetBalance(models.BTC)
	if balance != 0.0 {
		t.Errorf("Expected balance 0.0, got %f", balance)
	}
}

func TestWallet_MultipleAssets(t *testing.T) {
	wallet := NewWallet()

	// Deposit to different assets
	wallet.ProcessTransaction(models.Transaction{Type: models.Deposit, Asset: models.BTC, Amount: 1.0})
	wallet.ProcessTransaction(models.Transaction{Type: models.Deposit, Asset: models.ETH, Amount: 2.0})
	wallet.ProcessTransaction(models.Transaction{Type: models.Deposit, Asset: models.USD, Amount: 100.0})

	// Verify balances are independent
	if wallet.GetBalance(models.BTC) != 1.0 {
		t.Errorf("Expected BTC balance 1.0, got %f", wallet.GetBalance(models.BTC))
	}
	if wallet.GetBalance(models.ETH) != 2.0 {
		t.Errorf("Expected ETH balance 2.0, got %f", wallet.GetBalance(models.ETH))
	}
	if wallet.GetBalance(models.USD) != 100.0 {
		t.Errorf("Expected USD balance 100.0, got %f", wallet.GetBalance(models.USD))
	}
}

func TestWallet_GetAllBalances(t *testing.T) {
	wallet := NewWallet()

	wallet.ProcessTransaction(models.Transaction{Type: models.Deposit, Asset: models.BTC, Amount: 1.5})
	wallet.ProcessTransaction(models.Transaction{Type: models.Deposit, Asset: models.ETH, Amount: 3.0})
	wallet.ProcessTransaction(models.Transaction{Type: models.Deposit, Asset: models.USD, Amount: 250.0})

	balances := wallet.GetAllBalances()

	if balances[models.BTC] != 1.5 {
		t.Errorf("Expected BTC balance 1.5, got %f", balances[models.BTC])
	}
	if balances[models.ETH] != 3.0 {
		t.Errorf("Expected ETH balance 3.0, got %f", balances[models.ETH])
	}
	if balances[models.USD] != 250.0 {
		t.Errorf("Expected USD balance 250.0, got %f", balances[models.USD])
	}
}

func TestWallet_TransactionHistory(t *testing.T) {
	wallet := NewWallet()

	// Add successful transactions
	wallet.ProcessTransaction(models.Transaction{Type: models.Deposit, Asset: models.BTC, Amount: 1.0})
	wallet.ProcessTransaction(models.Transaction{Type: models.Deposit, Asset: models.ETH, Amount: 2.0})

	// Try failed transaction
	wallet.ProcessTransaction(models.Transaction{Type: models.Withdraw, Asset: models.BTC, Amount: 5.0})

	history := wallet.GetTransactionHistory()

	// Only successful transactions should be in history
	if len(history) != 2 {
		t.Errorf("Expected 2 transactions in history, got %d", len(history))
	}
}

func TestWallet_String(t *testing.T) {
	wallet := NewWallet()

	wallet.ProcessTransaction(models.Transaction{Type: models.Deposit, Asset: models.BTC, Amount: 1.5})
	wallet.ProcessTransaction(models.Transaction{Type: models.Deposit, Asset: models.ETH, Amount: 2.0})
	wallet.ProcessTransaction(models.Transaction{Type: models.Deposit, Asset: models.USD, Amount: 100.50})

	result := wallet.String()
	expected := "BTC: 1.50000000 | ETH: 2.00000000 | USD: 100.50"

	if result != expected {
		t.Errorf("Expected string: %s, got: %s", expected, result)
	}
}

func TestWallet_ComplexScenario(t *testing.T) {
	wallet := NewWallet()

	// Multiple deposits and withdrawals
	transactions := []struct {
		tx      models.Transaction
		wantErr bool
	}{
		{models.Transaction{Type: models.Deposit, Asset: models.BTC, Amount: 10.0}, false},
		{models.Transaction{Type: models.Withdraw, Asset: models.BTC, Amount: 3.0}, false},
		{models.Transaction{Type: models.Deposit, Asset: models.BTC, Amount: 5.0}, false},
		{models.Transaction{Type: models.Withdraw, Asset: models.BTC, Amount: 2.0}, false},
		{models.Transaction{Type: models.Withdraw, Asset: models.BTC, Amount: 20.0}, true}, // Should fail
	}

	for i, tt := range transactions {
		err := wallet.ProcessTransaction(tt.tx)
		if (err != nil) != tt.wantErr {
			t.Errorf("Transaction %d: expected error=%v, got error=%v", i, tt.wantErr, err != nil)
		}
	}

	// Expected: 10 - 3 + 5 - 2 = 10
	expectedBalance := 10.0
	actualBalance := wallet.GetBalance(models.BTC)
	if actualBalance != expectedBalance {
		t.Errorf("Expected final balance %f, got %f", expectedBalance, actualBalance)
	}
}
