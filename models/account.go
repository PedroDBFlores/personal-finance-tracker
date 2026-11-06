package models

import (
	"time"

	"uuid"
)

// AccountType defines the type of account (Bank, Brokerage, CreditCard, etc.)
type AccountType string

const (
	AccountTypeBank       AccountType = "Bank"
	AccountTypeBrokerage  AccountType = "Brokerage"
	AccountTypeCreditCard AccountType = "CreditCard"
)

// Account represents a financial account (e.g., savings, brokerage, credit card)
type Account struct {
	// Id represents the unique account identifier.
	Id uuid.UUID
	// Name defines the account name (e.g., "Chase Savings", "Fidelity Brokerage")
	Name string
	// Type defines the type of account.
	Type AccountType
	// Currency defines the account's currency (e.g., "USD", "EUR")
	Currency string
	// CreatedAt defines when the account was created.
	CreatedAt time.Time
}

// NewAccount creates a new Account based on the provided parameters
func NewAccount(id uuid.UUID, name string, accountType AccountType, currency string, createdAt time.Time) Account {
	return Account{
		Id:        id,
		Name:      name,
		Type:      accountType,
		Currency:  currency,
		CreatedAt: createdAt,
	}
}
