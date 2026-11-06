package domain

import (
	"time"
)

// Account represents a financial account (e.g., savings, brokerage, credit card)
type Account struct {
	// Id represents the unique account identifier.
	Id Id
	// Name defines the account name (e.g., "Chase Savings", "Fidelity Brokerage")
	Name string
	// Type defines the type of account.
	Type AccountType
	// Currency defines the account's currency (e.g., "USD", "EUR")
	Currency Currency
	// CreatedAt defines when the account was created.
	CreatedAt time.Time
}

// NewAccount creates a new Account based on the provided parameters
func NewAccount(id Id, name string, accountType AccountType, currency Currency, createdAt time.Time) Account {
	return Account{
		Id:        id,
		Name:      name,
		Type:      accountType,
		Currency:  currency,
		CreatedAt: createdAt,
	}
}
