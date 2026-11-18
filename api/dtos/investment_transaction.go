package dtos

import (
	"time"

	"github.com/pedrodbflores/personal-finance-tracker/domain"
)

type CreateInvestmentTransactionDto struct {
	// Type defines the type of transaction (Buy, Sell, Deposit, Withdraw).
	Type domain.InvestmentTransactionType
	// Amount is the monetary value of the transaction (for Deposit/Withdraw).
	Amount domain.Money
	// Quantity is the number of shares/units in hundredths (for Buy/Sell).
	Quantity domain.Money
	// Price is the price per share/unit (for Buy/Sell).
	Price domain.Money
	// Date defines the transaction's date.
	Date time.Time
	// Description provides additional details (optional).
	Description string
}

func NewCreateInvestmentTransactionDto(txType domain.InvestmentTransactionType, amount, quantity, price domain.Money, date time.Time, description string) CreateInvestmentTransactionDto {
	return CreateInvestmentTransactionDto{txType, amount, quantity, price, date, description}
}

type InvestmentTransactionDto struct {
	domain.InvestmentTransaction
}

func NewInvestmentTransactionDto(tx domain.InvestmentTransaction) InvestmentTransactionDto {
	return InvestmentTransactionDto{tx}
}
