package domain

import (
	"testing"
	"time"
)

func TestNewAccount(t *testing.T) {
	id := Id(1)
	name := "Test Account"
	accountType := AccountTypeBank
	currency := ParseCurrency("USD")
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
