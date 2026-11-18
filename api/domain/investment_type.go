package domain

import "fmt"

// InvestmentType defines the type of investment (ETF, Savings, or TermDeposit)
type InvestmentType string

const (
	ETF         InvestmentType = "ETF"
	Savings     InvestmentType = "Savings"
	TermDeposit InvestmentType = "TermDeposit"
)

// ParseInvestmentType parses a string into an InvestmentType.
func ParseInvestmentType(s string) (InvestmentType, error) {
	switch s {
	case "ETF":
		return ETF, nil
	case "Savings":
		return Savings, nil
	case "TermDeposit":
		return TermDeposit, nil
	default:
		return "", fmt.Errorf("unsupported investment type %s", s)
	}
}
