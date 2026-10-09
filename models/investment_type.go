package models

// InvestmentType defines the type of investment (ETF, Savings, or TermDeposit)
type InvestmentType string

const (
	ETF         InvestmentType = "ETF"
	Savings     InvestmentType = "Savings"
	TermDeposit InvestmentType = "TermDeposit"
)
