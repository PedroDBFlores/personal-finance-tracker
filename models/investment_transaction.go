package models

import (
	"math/big"
	"time"

	"github.com/google/uuid"
)

// InvestmentTransaction represents a transaction for an investment (e.g., Buy, Sell, Deposit, Withdraw)
type InvestmentTransaction struct {
	// Id represents the unique transaction identifier.
	Id uuid.UUID
	// InvestmentId links the transaction to an investment.
	InvestmentId uuid.UUID
	// Type defines the type of transaction.
	Type InvestmentTransactionType
	// Amount represents the monetary value of the transaction (for Deposit/Withdraw).
	// Last two digits are cents.
	Amount big.Int
	// Quantity represents the number of shares or units (for ETFs: shares; for Savings: optional).
	Quantity big.Int
	// Price represents the price per share/unit (for ETFs).
	// Last two digits are cents.
	Price big.Int
	// Date defines the transaction's date.
	Date time.Time
	// Description provides additional details (optional).
	Description string
}

// NewInvestmentTransaction creates a new InvestmentTransaction based on the provided parameters.
// For Deposit/Withdraw, use Amount and set Quantity/Price to 0.
// For Buy/Sell, use Quantity and Price, and set Amount to Quantity * Price (or leave it for calculation).
func NewInvestmentTransaction(
	id uuid.UUID,
	investmentId uuid.UUID,
	txType InvestmentTransactionType,
	amount big.Int,
	quantity big.Int,
	price big.Int,
	date time.Time,
	description string,
) InvestmentTransaction {
	return InvestmentTransaction{
		Id:            id,
		InvestmentId:  investmentId,
		Type:          txType,
		Amount:        amount,
		Quantity:      quantity,
		Price:         price,
		Date:          date,
		Description:   description,
	}
}
