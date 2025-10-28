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

	if wallet.GetBalance(models.BTC) != 0 {
		t.Error("Expected BTC balance to be 0 for new wallet")
	}
	if wallet.GetBalance(models.ETH) != 0 {
		t.Error("Expected ETH balance to be 0 for new wallet")
	}
	if wallet.GetBalance(models.USD) != 0 {
		t.Error("Expected USD balance to be 0 for new wallet")
	}
}

func TestWallet_DepositSuccess(t *testing.T) {
	wallet := NewWallet()

	tx := models.Transaction{
		Type:   models.Deposit,
		Asset:  models.BTC,
		Amount: 250000000, // 2.5 BTC
	}

	err := wallet.ProcessTransaction(tx)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	balance := wallet.GetBalance(models.BTC)
	expected := int64(250000000)
	if balance != expected {
		t.Errorf("Expected balance %d, got %d", expected, balance)
	}
}

func TestWallet_WithdrawSuccess(t *testing.T) {
	wallet := NewWallet()

	// First deposit
	wallet.ProcessTransaction(models.Transaction{
		Type:   models.Deposit,
		Asset:  models.BTC,
		Amount: 500000000, // 5.0 BTC
	})

	// Then withdraw
	tx := models.Transaction{
		Type:   models.Withdraw,
		Asset:  models.BTC,
		Amount: 200000000, // 2.0 BTC
	}

	err := wallet.ProcessTransaction(tx)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	balance := wallet.GetBalance(models.BTC)
	expected := int64(300000000) // 3.0 BTC
	if balance != expected {
		t.Errorf("Expected balance %d, got %d", expected, balance)
	}
}

func TestWallet_WithdrawInsufficientFunds(t *testing.T) {
	wallet := NewWallet()

	// Deposit small amount
	wallet.ProcessTransaction(models.Transaction{
		Type:   models.Deposit,
		Asset:  models.BTC,
		Amount: 100000000, // 1.0 BTC
	})

	// Try to withdraw more than available
	tx := models.Transaction{
		Type:   models.Withdraw,
		Asset:  models.BTC,
		Amount: 200000000, // 2.0 BTC
	}

	err := wallet.ProcessTransaction(tx)
	if err == nil {
		t.Fatal("Expected error for insufficient funds")
	}

	// Balance should remain unchanged
	balance := wallet.GetBalance(models.BTC)
	expected := int64(100000000)
	if balance != expected {
		t.Errorf("Expected balance %d (unchanged), got %d", expected, balance)
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
		Amount: 100000000, // 1.0 BTC
	}

	err := wallet.ProcessTransaction(tx)
	if err == nil {
		t.Fatal("Expected error for withdrawing from empty wallet")
	}

	balance := wallet.GetBalance(models.BTC)
	if balance != 0 {
		t.Errorf("Expected balance 0, got %d", balance)
	}
}

func TestWallet_MultipleAssets(t *testing.T) {
	wallet := NewWallet()

	// Deposit to different assets
	wallet.ProcessTransaction(models.Transaction{Type: models.Deposit, Asset: models.BTC, Amount: 100000000})           // 1.0 BTC
	wallet.ProcessTransaction(models.Transaction{Type: models.Deposit, Asset: models.ETH, Amount: 2000000000000000000}) // 2.0 ETH
	wallet.ProcessTransaction(models.Transaction{Type: models.Deposit, Asset: models.USD, Amount: 10000})               // 100.00 USD

	// Verify balances are independent
	if wallet.GetBalance(models.BTC) != 100000000 {
		t.Errorf("Expected BTC balance 100000000, got %d", wallet.GetBalance(models.BTC))
	}
	if wallet.GetBalance(models.ETH) != 2000000000000000000 {
		t.Errorf("Expected ETH balance 2000000000000000000, got %d", wallet.GetBalance(models.ETH))
	}
	if wallet.GetBalance(models.USD) != 10000 {
		t.Errorf("Expected USD balance 10000, got %d", wallet.GetBalance(models.USD))
	}
}

func TestWallet_GetAllBalances(t *testing.T) {
	wallet := NewWallet()

	wallet.ProcessTransaction(models.Transaction{Type: models.Deposit, Asset: models.BTC, Amount: 150000000})           // 1.5 BTC
	wallet.ProcessTransaction(models.Transaction{Type: models.Deposit, Asset: models.ETH, Amount: 3000000000000000000}) // 3.0 ETH
	wallet.ProcessTransaction(models.Transaction{Type: models.Deposit, Asset: models.USD, Amount: 25000})               // 250.00 USD

	balances := wallet.GetAllBalances()

	if balances[models.BTC] != 150000000 {
		t.Errorf("Expected BTC balance 150000000, got %d", balances[models.BTC])
	}
	if balances[models.ETH] != 3000000000000000000 {
		t.Errorf("Expected ETH balance 3000000000000000000, got %d", balances[models.ETH])
	}
	if balances[models.USD] != 25000 {
		t.Errorf("Expected USD balance 25000, got %d", balances[models.USD])
	}
}

func TestWallet_TransactionHistory(t *testing.T) {
	wallet := NewWallet()

	// Add successful transactions
	wallet.ProcessTransaction(models.Transaction{Type: models.Deposit, Asset: models.BTC, Amount: 100000000})           // 1.0 BTC
	wallet.ProcessTransaction(models.Transaction{Type: models.Deposit, Asset: models.ETH, Amount: 2000000000000000000}) // 2.0 ETH

	// Try failed transaction
	wallet.ProcessTransaction(models.Transaction{Type: models.Withdraw, Asset: models.BTC, Amount: 500000000}) // 5.0 BTC - should fail

	history := wallet.GetTransactionHistory()

	// Only successful transactions should be in history
	if len(history) != 2 {
		t.Errorf("Expected 2 transactions in history, got %d", len(history))
	}
}

func TestWallet_String(t *testing.T) {
	wallet := NewWallet()

	wallet.ProcessTransaction(models.Transaction{Type: models.Deposit, Asset: models.BTC, Amount: 150000000})           // 1.5 BTC
	wallet.ProcessTransaction(models.Transaction{Type: models.Deposit, Asset: models.ETH, Amount: 2000000000000000000}) // 2.0 ETH
	wallet.ProcessTransaction(models.Transaction{Type: models.Deposit, Asset: models.USD, Amount: 10050})               // 100.50 USD

	result := wallet.String()
	expected := "BTC: 1.50000000 | ETH: 2.000000000000000000 | USD: 100.50"

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
		{models.Transaction{Type: models.Deposit, Asset: models.BTC, Amount: 1000000000}, false},  // 10.0 BTC
		{models.Transaction{Type: models.Withdraw, Asset: models.BTC, Amount: 300000000}, false},  // 3.0 BTC
		{models.Transaction{Type: models.Deposit, Asset: models.BTC, Amount: 500000000}, false},   // 5.0 BTC
		{models.Transaction{Type: models.Withdraw, Asset: models.BTC, Amount: 200000000}, false},  // 2.0 BTC
		{models.Transaction{Type: models.Withdraw, Asset: models.BTC, Amount: 2000000000}, true},  // 20.0 BTC - should fail
	}

	for i, tt := range transactions {
		err := wallet.ProcessTransaction(tt.tx)
		if (err != nil) != tt.wantErr {
			t.Errorf("Transaction %d: expected error=%v, got error=%v", i, tt.wantErr, err != nil)
		}
	}

	// Expected: 10 - 3 + 5 - 2 = 10 BTC = 1000000000 satoshis
	expectedBalance := int64(1000000000)
	actualBalance := wallet.GetBalance(models.BTC)
	if actualBalance != expectedBalance {
		t.Errorf("Expected final balance %d, got %d", expectedBalance, actualBalance)
	}
}
