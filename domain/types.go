package domain

import (
	"fmt"
	"math/rand/v2"
	"strconv"
)

type Id uint

func (i Id) String() string {
	return strconv.FormatUint(uint64(i), 10)
}

func NewId() Id {
	return Id(rand.Uint())
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
