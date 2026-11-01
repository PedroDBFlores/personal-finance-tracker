package dtos

import "time"

// AccountType defines the type of account (Bank, Brokerage, CreditCard, etc.)
type AccountType string

const (
	Bank       AccountType = "Bank"
	Brokerage  AccountType = "Brokerage"
	CreditCard AccountType = "CreditCard"
)

type CreateAccountDto struct {
	// Name defines the account name (e.g., "Chase Savings", "Fidelity Brokerage")
	Name string
	// Type defines the type of account.
	Type AccountType
	// Currency defines the account's currency (e.g., "USD", "EUR")
	Currency string
	// CreatedAt defines when the account was created.
	CreatedAt time.Time
}
