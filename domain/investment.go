package domain

import (
	"math/big"
	"time"

	"github.com/google/uuid"
)

// Investment represents a financial investment (e.g., ETF, Savings, TermDeposit)
type Investment struct {
	// Id represents the unique investment identifier.
	Id uuid.UUID
	// AccountId links the investment to an account.
	AccountId uuid.UUID
	// Name defines the investment name (e.g., "Vanguard S&P 500 ETF", "Ally Bank CD")
	Name string
	// Type defines the type of investment.
	Type InvestmentType
	// CurrentValue represents the current value of the investment (calculated as Quantity * Price).
	// Last two digits are cents.
	CurrentValue big.Int
	// Currency defines the investment's currency (e.g., "USD", "EUR")
	Currency string
	// PurchaseDate defines when the investment was purchased or opened.
	PurchaseDate time.Time
	// MaturityDate defines when a term deposit matures (optional).
	MaturityDate *time.Time
}

// NewInvestment creates a new Investment based on the provided parameters.
// currentValue is calculated from transactions, so it can be zero initially.
func NewInvestment(
	id uuid.UUID,
	accountId uuid.UUID,
	name string,
	investmentType InvestmentType,
	currency string,
	purchaseDate time.Time,
	maturityDate *time.Time,
) Investment {
	return Investment{
		Id:           id,
		AccountId:    accountId,
		Name:         name,
		Type:         investmentType,
		CurrentValue: *big.NewInt(0),
		Currency:     currency,
		PurchaseDate: purchaseDate,
		MaturityDate: maturityDate,
	}
}

// UpdateCurrentValue recalculates the CurrentValue based on the latest transactions.
// For ETFs: sum of (Quantity * Price) for all Buy transactions minus Sell transactions.
// For Savings/TermDeposit: sum of Deposit transactions minus Withdrawal transactions.
func (i *Investment) UpdateCurrentValue(transactions []InvestmentTransaction) {
	value := big.NewInt(0)

	for _, tx := range transactions {
		if tx.InvestmentId != i.Id {
			continue
		}

		// Handle ETF transactions (Buy/Sell)
		if i.Type == ETF {
			if tx.Type == Buy {
				// Add to value: Quantity * Price
				product := new(big.Int).Mul(&tx.Quantity, &tx.Price)
				value.Add(value, product)
			} else if tx.Type == Sell {
				// Subtract from value: Quantity * Price
				product := new(big.Int).Mul(&tx.Quantity, &tx.Price)
				value.Sub(value, product)
			}
		} else {
			// Handle Savings/TermDeposit transactions (Deposit/Withdraw)
			if tx.Type == Deposit {
				value.Add(value, &tx.Amount)
			} else if tx.Type == Withdraw {
				value.Sub(value, &tx.Amount)
			}
		}
	}

	i.CurrentValue = *value
}
