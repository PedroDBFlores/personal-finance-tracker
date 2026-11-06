package domain

import (
	"math/big"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewInvestment(t *testing.T) {
	id := uuid.New()
	accountId := uuid.New()
	name := "Test Investment"
	investmentType := ETF
	currency := "USD"
	purchaseDate := time.Now()
	maturityDate := time.Now().AddDate(1, 0, 0) // 1 year from now

	investment := NewInvestment(id, accountId, name, investmentType, currency, purchaseDate, &maturityDate)

	if investment.Id != id {
		t.Errorf("Expected Id %v, got %v", id, investment.Id)
	}
	if investment.AccountId != accountId {
		t.Errorf("Expected AccountId %v, got %v", accountId, investment.AccountId)
	}
	if investment.Name != name {
		t.Errorf("Expected Name %s, got %s", name, investment.Name)
	}
	if investment.Type != investmentType {
		t.Errorf("Expected Type %s, got %s", investmentType, investment.Type)
	}
	if investment.CurrentValue.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("Expected CurrentValue 0, got %v", investment.CurrentValue)
	}
	if investment.Currency != currency {
		t.Errorf("Expected Currency %s, got %s", currency, investment.Currency)
	}
	if !investment.PurchaseDate.Equal(purchaseDate) {
		t.Errorf("Expected PurchaseDate %v, got %v", purchaseDate, investment.PurchaseDate)
	}
	if investment.MaturityDate == nil || !investment.MaturityDate.Equal(maturityDate) {
		t.Errorf("Expected MaturityDate %v, got %v", maturityDate, investment.MaturityDate)
	}
}

func TestNewInvestmentWithoutMaturityDate(t *testing.T) {
	id := uuid.New()
	accountId := uuid.New()
	name := "Test Investment"
	investmentType := Savings
	currency := "USD"
	purchaseDate := time.Now()

	investment := NewInvestment(id, accountId, name, investmentType, currency, purchaseDate, nil)

	if investment.Id != id {
		t.Errorf("Expected Id %v, got %v", id, investment.Id)
	}
	if investment.MaturityDate != nil {
		t.Errorf("Expected MaturityDate nil, got %v", investment.MaturityDate)
	}
}

func TestUpdateCurrentValueForETF(t *testing.T) {
	investmentId := uuid.New()
	accountId := uuid.New()
	investment := NewInvestment(investmentId, accountId, "Vanguard ETF", ETF, "USD", time.Now(), nil)

	// Buy 10 shares at $100 each (value = $1000.00)
	// Quantity is 10 shares (10 * 10000 (price) = 100000 cents = $1000.00)
	buyQuantity := big.NewInt(10) // 10 shares
	buyPrice := big.NewInt(10000) // $100.00 (10000 cents)
	buyTransaction := NewInvestmentTransaction(
		uuid.New(),
		investmentId,
		Buy,
		*big.NewInt(0),
		*buyQuantity,
		*buyPrice,
		time.Now(),
		"Initial Buy",
	)

	// Sell 2 shares at $120 each (value = $240.00)
	sellQuantity := big.NewInt(2)  // 2 shares
	sellPrice := big.NewInt(12000) // $120.00 (12000 cents)
	sellTransaction := NewInvestmentTransaction(
		uuid.New(),
		investmentId,
		Sell,
		*big.NewInt(0),
		*sellQuantity,
		*sellPrice,
		time.Now().AddDate(0, 1, 0),
		"Partial Sell",
	)

	transactions := []InvestmentTransaction{buyTransaction, sellTransaction}
	investment.UpdateCurrentValue(transactions)

	// Expected: (10 * 10000) - (2 * 12000) = 100000 - 24000 = 76000 ($760.00)
	expectedValue := big.NewInt(76000)
	if investment.CurrentValue.Cmp(expectedValue) != 0 {
		t.Errorf("Expected CurrentValue %v, got %v", expectedValue, investment.CurrentValue)
	}
}

func TestUpdateCurrentValueForSavings(t *testing.T) {
	investmentId := uuid.New()
	accountId := uuid.New()
	investment := NewInvestment(investmentId, accountId, "Savings Account", Savings, "USD", time.Now(), nil)

	// Deposit $500.00
	depositAmount := big.NewInt(50000)
	depositTransaction := NewInvestmentTransaction(
		uuid.New(),
		investmentId,
		Deposit,
		*depositAmount,
		*big.NewInt(0),
		*big.NewInt(0),
		time.Now(),
		"Deposit",
	)

	// Withdraw $200.00
	withdrawAmount := big.NewInt(20000)
	withdrawTransaction := NewInvestmentTransaction(
		uuid.New(),
		investmentId,
		Withdraw,
		*withdrawAmount,
		*big.NewInt(0),
		*big.NewInt(0),
		time.Now().AddDate(0, 1, 0),
		"Withdraw",
	)

	transactions := []InvestmentTransaction{depositTransaction, withdrawTransaction}
	investment.UpdateCurrentValue(transactions)

	// Expected: 50000 - 20000 = 30000 ($300.00)
	expectedValue := big.NewInt(30000)
	if investment.CurrentValue.Cmp(expectedValue) != 0 {
		t.Errorf("Expected CurrentValue %v, got %v", expectedValue, investment.CurrentValue)
	}
}

func TestUpdateCurrentValueWithUnrelatedTransactions(t *testing.T) {
	investmentId := uuid.New()
	otherInvestmentId := uuid.New()
	accountId := uuid.New()
	investment := NewInvestment(investmentId, accountId, "Test Investment", ETF, "USD", time.Now(), nil)

	// Transaction for a different investment
	otherDepositAmount := big.NewInt(100000)
	otherTransaction := NewInvestmentTransaction(
		uuid.New(),
		otherInvestmentId,
		Deposit,
		*otherDepositAmount,
		*big.NewInt(0),
		*big.NewInt(0),
		time.Now(),
		"Deposit",
	)

	transactions := []InvestmentTransaction{otherTransaction}
	investment.UpdateCurrentValue(transactions)

	// Should remain zero since the transaction is for a different investment
	expectedValue := big.NewInt(0)
	if investment.CurrentValue.Cmp(expectedValue) != 0 {
		t.Errorf("Expected CurrentValue %v, got %v", expectedValue, investment.CurrentValue)
	}
}
