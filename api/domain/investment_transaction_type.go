package domain

// InvestmentTransactionType defines the type of investment transaction (Buy, Sell, Deposit, Withdraw)
type InvestmentTransactionType string

const (
	Buy      InvestmentTransactionType = "Buy"
	Sell     InvestmentTransactionType = "Sell"
	Deposit  InvestmentTransactionType = "Deposit"
	Withdraw InvestmentTransactionType = "Withdraw"
)
