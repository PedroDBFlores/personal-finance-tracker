package domain

import (
	"fmt"
	"uuid"
)

type Id uuid.UUID

func (i Id) String() string {
	return uuid.UUID(i).String()
}

func NewId() Id {
	return Id(uuid.NewV7())
}

type Currency string

const (
	CurrencyUSD Currency = "USD"
	CurrencyEUR Currency = "EUR"
	CurrencyGBP Currency = "GBP"
)

func ParseCurrency(s string) Currency {
	switch s {
	case "USD", "EUR", "GBP":
		return Currency(s)
	default:
		panic(fmt.Sprintf("unsupported currency %s", s))
	}
}

// AccountType defines the type of account (Bank, Brokerage, CreditCard, etc.)
type AccountType string

const (
	AccountTypeBank       AccountType = "Bank"
	AccountTypeBrokerage  AccountType = "Brokerage"
	AccountTypeCreditCard AccountType = "CreditCard"
)

func ParseAccountType(s string) AccountType {
	switch s {
	case "Bank":
		return AccountTypeBank
	case "Brokerage":
		return AccountTypeBrokerage
	case "CreditCard":
		return AccountTypeCreditCard
	default:
		panic(fmt.Sprintf("unsupported account type %s", s))
	}
}
