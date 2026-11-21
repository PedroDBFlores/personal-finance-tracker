package domain

import (
	"strings"
	"time"
)

// Investment represents a financial investment (e.g., ETF, Savings, TermDeposit)
type Investment struct {
	// Id represents the unique investment identifier.
	Id Id
	// AccountId links the investment to an account.
	AccountId Id
	// Name defines the investment name (e.g., "Vanguard S&P 500 ETF", "Ally Bank CD")
	Name string
	// Type defines the type of investment.
	Type InvestmentType
	// CurrentValue represents the current value of the investment (calculated from its transactions).
	// Last two digits are cents.
	CurrentValue Money
	// Currency defines the investment's currency (e.g., "USD", "EUR")
	Currency Currency
	// PurchaseDate defines when the investment was purchased or opened.
	PurchaseDate time.Time
	// MaturityDate defines when a term deposit matures (optional).
	MaturityDate *time.Time
}

// NewInvestment creates a new Investment based on the provided parameters.
// CurrentValue is calculated from transactions, so it starts at zero.
// It returns a ValidationError if the name is empty, the investment type or
// currency is not supported, or the maturity date is before the purchase date.
func NewInvestment(
	id Id,
	accountId Id,
	name string,
	investmentType InvestmentType,
	currency Currency,
	purchaseDate time.Time,
	maturityDate *time.Time,
) (Investment, error) {
	if strings.TrimSpace(name) == "" {
		return Investment{}, newValidationError("investment name is required")
	}
	switch investmentType {
	case ETF, Savings, TermDeposit:
	default:
		return Investment{}, newValidationError("unsupported investment type %q", investmentType)
	}
	switch currency {
	case CurrencyUSD, CurrencyEUR, CurrencyGBP:
	default:
		return Investment{}, newValidationError("unsupported currency %q", currency)
	}
	if maturityDate != nil && maturityDate.Before(purchaseDate) {
		return Investment{}, newValidationError("maturity date must not be before the purchase date")
	}
	return Investment{
		Id:           id,
		AccountId:    accountId,
		Name:         name,
		Type:         investmentType,
		CurrentValue: 0,
		Currency:     currency,
		PurchaseDate: purchaseDate,
		MaturityDate: maturityDate,
	}, nil
}

// UpdateCurrentValue recalculates the CurrentValue based on the given transactions.
// For ETFs: sum of (Quantity * Price) for Buy transactions minus Sell transactions.
// For Savings/TermDeposit: sum of Deposit amounts minus Withdrawal amounts.
func (i *Investment) UpdateCurrentValue(transactions []InvestmentTransaction) {
	value := Money(0)

	for _, tx := range transactions {
		if tx.InvestmentId != i.Id {
			continue
		}

		if i.Type == ETF {
			// Quantity is in hundredths of a share and Price in cents,
			// so the product is Quantity * Price / 100 cents.
			product := tx.Price * tx.Quantity / 100
			if tx.Type == Buy {
				value = value.Add(product)
			} else if tx.Type == Sell {
				value = value.Subtract(product)
			}
		} else {
			if tx.Type == Deposit {
				value = value.Add(tx.Amount)
			} else if tx.Type == Withdraw {
				value = value.Subtract(tx.Amount)
			}
		}
	}

	i.CurrentValue = value
}

// InvestmentPatch represents the fields to modify on an Investment. Nil fields are left unchanged.
type InvestmentPatch struct {
	// Name replaces the investment name when set.
	Name *string
	// Type replaces the investment type when set.
	Type *InvestmentType
	// Currency replaces the investment currency when set.
	Currency *Currency
	// PurchaseDate replaces the purchase date when set.
	PurchaseDate *time.Time
	// MaturityDate replaces the maturity date when set; a non-nil pointer to nil clears it.
	MaturityDate **time.Time
	// CurrentValue replaces the current value when set.
	CurrentValue *Money
}
