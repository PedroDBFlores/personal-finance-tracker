package models

import (
	"math/big"
	"time"

	"uuid"
)

// AccountType defines the type of account (Bank, Brokerage, CreditCard, etc.)
type AccountType string

const (
	Bank       AccountType = "Bank"
	Brokerage  AccountType = "Brokerage"
	CreditCard AccountType = "CreditCard"
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

// CalculateBalance calculates the current balance of the account from its linked movements.
// Returns the balance as a big.Int (last two digits are cents).
func (a Account) CalculateBalance(movements []Movement) big.Int {
	balance := big.NewInt(0)
	for _, movement := range movements {
		if movement.AccountId == a.Id {
			switch movement.Type {
			case Credit:
				balance.Add(balance, &movement.Value)
			case Debit:
				balance.Sub(balance, &movement.Value)
			}
		}
	}
	return *balance
}

// CalculateBalanceAt calculates the account balance at a specific point in time.
// Useful for snapshots (e.g., end-of-year balance).
func (a Account) CalculateBalanceAt(movements []Movement, snapshotTime time.Time) big.Int {
	balance := big.NewInt(0)
	for _, movement := range movements {
		if (movement.AccountId == a.Id) && (movement.Date.Before(snapshotTime) || movement.Date.Equal(snapshotTime)) {
			switch movement.Type {
			case Credit:
				balance.Add(balance, &movement.Value)
			case Debit:
				balance.Sub(balance, &movement.Value)
			}
		}
	}
	return *balance
}
