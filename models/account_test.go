package models

import (
	"math/big"
	"testing"
	"time"
	"uuid"
)

func TestNewAccount(t *testing.T) {
	id := uuid.New()
	name := "Test Account"
	accountType := Bank
	currency := "USD"
	createdAt := time.Now()

	account := NewAccount(id, name, accountType, currency, createdAt)

	if account.Id != id {
		t.Errorf("Expected Id %v, got %v", id, account.Id)
	}
	if account.Name != name {
		t.Errorf("Expected Name %s, got %s", name, account.Name)
	}
	if account.Type != accountType {
		t.Errorf("Expected Type %s, got %s", accountType, account.Type)
	}
	if account.Currency != currency {
		t.Errorf("Expected Currency %s, got %s", currency, account.Currency)
	}
	if !account.CreatedAt.Equal(createdAt) {
		t.Errorf("Expected CreatedAt %v, got %v", createdAt, account.CreatedAt)
	}
}

func TestCalculateBalance(t *testing.T) {
	accountId := uuid.New()
	account := NewAccount(accountId, "Test Account", Bank, "USD", time.Now())

	// Create movements for the account
	creditValue := big.NewInt(100000) // $1000.00
	debitValue := big.NewInt(50000)   // $500.00
	movements := []Movement{
		{
			Id:          uuid.New(),
			Description: "Salary",
			Value:       *creditValue,
			Type:        Credit,
			Date:        time.Now(),
			AccountId:   accountId,
			Tags:        nil,
		},
		{
			Id:          uuid.New(),
			Description: "Rent",
			Value:       *debitValue,
			Type:        Debit,
			Date:        time.Now(),
			AccountId:   accountId,
			Tags:        nil,
		},
	}

	balance := account.CalculateBalance(movements)
	expectedBalance := big.NewInt(50000) // $1000.00 - $500.00 = $500.00

	if balance.Cmp(expectedBalance) != 0 {
		t.Errorf("Expected balance %v, got %v", expectedBalance, balance)
	}
}

func TestCalculateBalanceAt(t *testing.T) {
	accountId := uuid.New()
	account := NewAccount(accountId, "Test Account", Bank, "USD", time.Now())

	now := time.Now()
	snapshotTime := now.AddDate(0, 1, 0) // 1 month from now

	// Movement before snapshot time
	creditValue := big.NewInt(100000) // $1000.00
	movements := []Movement{
		{
			Id:          uuid.New(),
			Description: "Salary",
			Value:       *creditValue,
			Type:        Credit,
			Date:        now,
			AccountId:   accountId,
			Tags:        nil,
		},
	}

	balance := account.CalculateBalanceAt(movements, snapshotTime)
	expectedBalance := big.NewInt(100000) // $1000.00

	if balance.Cmp(expectedBalance) != 0 {
		t.Errorf("Expected balance %v, got %v", expectedBalance, balance)
	}

	// Movement after snapshot time
	futureMovement := Movement{
		Id:          uuid.New(),
		Description: "Future Salary",
		Value:       *creditValue,
		Type:        Credit,
		Date:        snapshotTime.AddDate(0, 1, 0), // 2 months from now
		AccountId:   accountId,
		Tags:        nil,
	}
	movementsWithFuture := append(movements, futureMovement)

	balanceAfter := account.CalculateBalanceAt(movementsWithFuture, snapshotTime)
	// Should still only include the first movement
	if balanceAfter.Cmp(expectedBalance) != 0 {
		t.Errorf("Expected balance %v, got %v (should ignore future movement)", expectedBalance, balanceAfter)
	}
}

func TestCalculateBalanceWithDifferentAccount(t *testing.T) {
	accountId1 := uuid.New()
	accountId2 := uuid.New()
	account1 := NewAccount(accountId1, "Account 1", Bank, "USD", time.Now())
	account2 := NewAccount(accountId2, "Account 2", Bank, "USD", time.Now())

	// Movement for account 1
	creditValue := big.NewInt(100000) // $1000.00
	movements := []Movement{
		{
			Id:          uuid.New(),
			Description: "Salary",
			Value:       *creditValue,
			Type:        Credit,
			Date:        time.Now(),
			AccountId:   accountId1,
			Tags:        nil,
		},
	}

	// Account 1 should have balance
	balance1 := account1.CalculateBalance(movements)
	expectedBalance1 := big.NewInt(100000)
	if balance1.Cmp(expectedBalance1) != 0 {
		t.Errorf("Expected balance for account1 %v, got %v", expectedBalance1, balance1)
	}

	// Account 2 should have zero balance
	balance2 := account2.CalculateBalance(movements)
	expectedBalance2 := big.NewInt(0)
	if balance2.Cmp(expectedBalance2) != 0 {
		t.Errorf("Expected balance for account2 %v, got %v", expectedBalance2, balance2)
	}
}
