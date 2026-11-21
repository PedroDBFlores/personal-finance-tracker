package domain

import (
	"errors"
	"testing"
	"time"
)

func TestNewAccount(t *testing.T) {
	id := Id(1)
	name := "Test Account"
	accountType := AccountTypeBank
	currency := ParseCurrency("USD")
	createdAt := time.Now()

	account, err := NewAccount(id, name, accountType, currency, createdAt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

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

func TestNewAccountValidation(t *testing.T) {
	now := time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)

	t.Run("accepts valid input", func(t *testing.T) {
		acc, err := NewAccount(1, "Bank A", AccountTypeBank, CurrencyEUR, now)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if acc.Name != "Bank A" {
			t.Errorf("Expected Name Bank A, got %s", acc.Name)
		}
	})

	t.Run("rejects an empty name", func(t *testing.T) {
		_, err := NewAccount(1, "", AccountTypeBank, CurrencyEUR, now)
		var ve ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("expected ValidationError, got %v", err)
		}
	})

	t.Run("rejects a whitespace-only name", func(t *testing.T) {
		_, err := NewAccount(1, "   ", AccountTypeBank, CurrencyEUR, now)
		var ve ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("expected ValidationError, got %v", err)
		}
	})

	t.Run("rejects an unsupported account type", func(t *testing.T) {
		_, err := NewAccount(1, "Bank A", AccountType("Crypto"), CurrencyEUR, now)
		var ve ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("expected ValidationError, got %v", err)
		}
	})

	t.Run("rejects an unsupported currency", func(t *testing.T) {
		_, err := NewAccount(1, "Bank A", AccountTypeBank, Currency("BTC"), now)
		var ve ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("expected ValidationError, got %v", err)
		}
	})
}
