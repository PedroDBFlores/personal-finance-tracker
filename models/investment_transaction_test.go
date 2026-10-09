package models

import (
	"math/big"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewInvestmentTransaction(t *testing.T) {
	id := uuid.New()
	investmentId := uuid.New()
	txType := Buy
	amount := big.NewInt(10000) // $100.00
	quantity := big.NewInt(100)  // 1.00 shares
	price := big.NewInt(10000)   // $100.00 per share
	date := time.Now()
	description := "Test Transaction"

	transaction := NewInvestmentTransaction(id, investmentId, txType, *amount, *quantity, *price, date, description)

	if transaction.Id != id {
		t.Errorf("Expected Id %v, got %v", id, transaction.Id)
	}
	if transaction.InvestmentId != investmentId {
		t.Errorf("Expected InvestmentId %v, got %v", investmentId, transaction.InvestmentId)
	}
	if transaction.Type != txType {
		t.Errorf("Expected Type %s, got %s", txType, transaction.Type)
	}
	if transaction.Amount.Cmp(amount) != 0 {
		t.Errorf("Expected Amount %v, got %v", amount, transaction.Amount)
	}
	if transaction.Quantity.Cmp(quantity) != 0 {
		t.Errorf("Expected Quantity %v, got %v", quantity, transaction.Quantity)
	}
	if transaction.Price.Cmp(price) != 0 {
		t.Errorf("Expected Price %v, got %v", price, transaction.Price)
	}
	if !transaction.Date.Equal(date) {
		t.Errorf("Expected Date %v, got %v", date, transaction.Date)
	}
	if transaction.Description != description {
		t.Errorf("Expected Description %s, got %s", description, transaction.Description)
	}
}

func TestNewInvestmentTransactionWithDeposit(t *testing.T) {
	id := uuid.New()
	investmentId := uuid.New()
	txType := Deposit
	amount := big.NewInt(50000) // $500.00
	quantity := big.NewInt(0)
	price := big.NewInt(0)
	date := time.Now()
	description := "Deposit to savings"

	transaction := NewInvestmentTransaction(id, investmentId, txType, *amount, *quantity, *price, date, description)

	if transaction.Type != txType {
		t.Errorf("Expected Type %s, got %s", txType, transaction.Type)
	}
	if transaction.Amount.Cmp(amount) != 0 {
		t.Errorf("Expected Amount %v, got %v", amount, transaction.Amount)
	}
	if transaction.Quantity.Cmp(quantity) != 0 {
		t.Errorf("Expected Quantity %v, got %v", quantity, transaction.Quantity)
	}
	if transaction.Price.Cmp(price) != 0 {
		t.Errorf("Expected Price %v, got %v", price, transaction.Price)
	}
}

func TestNewInvestmentTransactionWithSell(t *testing.T) {
	id := uuid.New()
	investmentId := uuid.New()
	txType := Sell
	amount := big.NewInt(0)
	quantity := big.NewInt(500) // 5.00 shares
	price := big.NewInt(15000)  // $150.00 per share
	date := time.Now()
	description := "Sell ETF shares"

	transaction := NewInvestmentTransaction(id, investmentId, txType, *amount, *quantity, *price, date, description)

	if transaction.Type != txType {
		t.Errorf("Expected Type %s, got %s", txType, transaction.Type)
	}
	if transaction.Quantity.Cmp(quantity) != 0 {
		t.Errorf("Expected Quantity %v, got %v", quantity, transaction.Quantity)
	}
	if transaction.Price.Cmp(price) != 0 {
		t.Errorf("Expected Price %v, got %v", price, transaction.Price)
	}
}

func TestInvestmentTransactionDefaults(t *testing.T) {
	id := uuid.Nil
	investmentId := uuid.Nil
	txType := InvestmentTransactionType("")
	amount := big.NewInt(0)
	quantity := big.NewInt(0)
	price := big.NewInt(0)
	date := time.Time{}
	description := ""

	transaction := InvestmentTransaction{
		Id:            id,
		InvestmentId:  investmentId,
		Type:          txType,
		Amount:        *amount,
		Quantity:      *quantity,
		Price:         *price,
		Date:          date,
		Description:   description,
	}

	if transaction.Id != id {
		t.Errorf("Expected Id %v, got %v", id, transaction.Id)
	}
	if transaction.InvestmentId != investmentId {
		t.Errorf("Expected InvestmentId %v, got %v", investmentId, transaction.InvestmentId)
	}
	if transaction.Type != txType {
		t.Errorf("Expected Type %s, got %s", txType, transaction.Type)
	}
	if transaction.Amount.Cmp(amount) != 0 {
		t.Errorf("Expected Amount %v, got %v", amount, transaction.Amount)
	}
	if transaction.Quantity.Cmp(quantity) != 0 {
		t.Errorf("Expected Quantity %v, got %v", quantity, transaction.Quantity)
	}
	if transaction.Price.Cmp(price) != 0 {
		t.Errorf("Expected Price %v, got %v", price, transaction.Price)
	}
	if !transaction.Date.Equal(date) {
		t.Errorf("Expected Date %v, got %v", date, transaction.Date)
	}
	if transaction.Description != description {
		t.Errorf("Expected Description %s, got %s", description, transaction.Description)
	}
}
