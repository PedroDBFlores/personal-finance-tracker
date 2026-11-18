package dtos

import (
	"time"

	"github.com/pedrodbflores/personal-finance-tracker/domain"
)

type CreateInvestmentDto struct {
	// Name defines the investment name (e.g., "Vanguard S&P 500 ETF")
	Name string
	// Type defines the type of investment.
	Type domain.InvestmentType
	// Currency defines the investment's currency (e.g., "USD", "EUR")
	Currency domain.Currency
	// PurchaseDate defines when the investment was purchased or opened.
	PurchaseDate time.Time
	// MaturityDate defines when a term deposit matures (optional).
	MaturityDate *time.Time
}

func NewCreateInvestmentDto(name string, investmentType domain.InvestmentType, currency domain.Currency, purchaseDate time.Time, maturityDate *time.Time) CreateInvestmentDto {
	return CreateInvestmentDto{name, investmentType, currency, purchaseDate, maturityDate}
}

type InvestmentDto struct {
	domain.Investment
}

func NewInvestmentDto(inv domain.Investment) InvestmentDto {
	return InvestmentDto{inv}
}
