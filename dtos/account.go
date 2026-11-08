package dtos

import (
	"time"

	"github.com/pedrodbflores/personal-finance-tracker/domain"
)

type CreateAccountDto struct {
	// Name defines the account name (e.g., "Chase Savings", "Fidelity Brokerage")
	Name string
	// Type defines the type of account.
	Type domain.AccountType
	// Currency defines the account's currency (e.g., "USD", "EUR")
	Currency domain.Currency
	// CreatedAt defines when the account was created.
	CreatedAt time.Time
}

func NewCreateAccountDto(name string, accType domain.AccountType, currency domain.Currency, createdAt time.Time) CreateAccountDto {
	return CreateAccountDto{name, accType, currency, createdAt}
}

type AccountDto struct {
	domain.Account
}

func NewAccountDto(acc domain.Account) AccountDto {
	return AccountDto{acc}
}
