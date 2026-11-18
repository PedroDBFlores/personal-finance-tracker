package domain

import "fmt"

// InvestmentTransactionType defines the type of investment transaction (Buy, Sell, Deposit, Withdraw)
type InvestmentTransactionType string

const (
	Buy      InvestmentTransactionType = "Buy"
	Sell     InvestmentTransactionType = "Sell"
	Deposit  InvestmentTransactionType = "Deposit"
	Withdraw InvestmentTransactionType = "Withdraw"
)

// ParseInvestmentTransactionType parses a string into an InvestmentTransactionType.
func ParseInvestmentTransactionType(s string) (InvestmentTransactionType, error) {
	switch s {
	case "Buy":
		return Buy, nil
	case "Sell":
		return Sell, nil
	case "Deposit":
		return Deposit, nil
	case "Withdraw":
		return Withdraw, nil
	default:
		return "", fmt.Errorf("unsupported investment transaction type %s", s)
	}
}
