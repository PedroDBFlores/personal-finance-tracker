package domain

import (
	"time"
)

// InvestmentTransaction represents a transaction for an investment (e.g., Buy, Sell, Deposit, Withdraw)
type InvestmentTransaction struct {
	// Id represents the unique transaction identifier.
	Id Id
	// InvestmentId links the transaction to an investment.
	InvestmentId Id
	// Type defines the type of transaction.
	Type InvestmentTransactionType
	// Amount represents the monetary value of the transaction (for Deposit/Withdraw).
	// Last two digits are cents.
	Amount Money
	// Quantity represents the number of shares or units (for Buy/Sell), in hundredths
	// of a unit — the last two digits are always the hundredths.
	Quantity Money
	// Price represents the price per share/unit (for Buy/Sell).
	// Last two digits are cents.
	Price Money
	// Date defines the transaction's date.
	Date time.Time
	// Description provides additional details (optional).
	Description string
}

// NewInvestmentTransaction creates a new InvestmentTransaction based on the provided parameters.
// For Deposit/Withdraw, use Amount and set Quantity/Price to 0.
// For Buy/Sell, use Quantity and Price.
// It returns a ValidationError if the transaction type is not supported, an
// amount/quantity/price is negative, Deposit/Withdraw has no amount, or
// Buy/Sell is missing quantity or price.
func NewInvestmentTransaction(
	id Id,
	investmentId Id,
	txType InvestmentTransactionType,
	amount Money,
	quantity Money,
	price Money,
	date time.Time,
	description string,
) (InvestmentTransaction, error) {
	switch txType {
	case Buy, Sell, Deposit, Withdraw:
	default:
		return InvestmentTransaction{}, newValidationError("unsupported transaction type %q", txType)
	}
	if amount < 0 || quantity < 0 || price < 0 {
		return InvestmentTransaction{}, newValidationError("amount, quantity and price must not be negative")
	}
	switch txType {
	case Deposit, Withdraw:
		if amount == 0 {
			return InvestmentTransaction{}, newValidationError("%s requires a non-zero amount", txType)
		}
	case Buy, Sell:
		if quantity == 0 || price == 0 {
			return InvestmentTransaction{}, newValidationError("%s requires a non-zero quantity and price", txType)
		}
	}
	return InvestmentTransaction{
		Id:           id,
		InvestmentId: investmentId,
		Type:         txType,
		Amount:       amount,
		Quantity:     quantity,
		Price:        price,
		Date:         date,
		Description:  description,
	}, nil
}

// InvestmentTransactionPatch represents the fields to modify on an InvestmentTransaction.
// Nil fields are left unchanged.
type InvestmentTransactionPatch struct {
	// Type replaces the transaction type when set.
	Type *InvestmentTransactionType
	// Amount replaces the amount when set.
	Amount *Money
	// Quantity replaces the quantity when set.
	Quantity *Money
	// Price replaces the price when set.
	Price *Money
	// Date replaces the transaction date when set.
	Date *time.Time
	// Description replaces the description when set.
	Description *string
}
