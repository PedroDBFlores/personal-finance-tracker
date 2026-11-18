package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/stretchr/testify/require"
)

func TestNewInvestment(t *testing.T) {
	id := uuid.New()
	accountId := domain.NewId()
	purchaseDate := time.Now()
	maturityDate := purchaseDate.AddDate(1, 0, 0)

	t.Run("creates an investment with all fields", func(t *testing.T) {
		inv := domain.NewInvestment(id, accountId, "Test Investment", domain.ETF, domain.CurrencyUSD, purchaseDate, &maturityDate)

		require.Equal(t, id, inv.Id)
		require.Equal(t, accountId, inv.AccountId)
		require.Equal(t, "Test Investment", inv.Name)
		require.Equal(t, domain.ETF, inv.Type)
		require.Equal(t, domain.Money(0), inv.CurrentValue)
		require.Equal(t, domain.CurrencyUSD, inv.Currency)
		require.True(t, inv.PurchaseDate.Equal(purchaseDate))
		require.NotNil(t, inv.MaturityDate)
		require.True(t, inv.MaturityDate.Equal(maturityDate))
	})

	t.Run("creates an investment without maturity date", func(t *testing.T) {
		inv := domain.NewInvestment(id, accountId, "Savings", domain.Savings, domain.CurrencyUSD, purchaseDate, nil)
		require.Nil(t, inv.MaturityDate)
	})
}

func TestUpdateCurrentValue(t *testing.T) {
	investmentId := uuid.New()
	accountId := domain.NewId()
	now := time.Now()

	t.Run("for ETF: buy minus sell", func(t *testing.T) {
		inv := domain.NewInvestment(investmentId, accountId, "Vanguard ETF", domain.ETF, domain.CurrencyUSD, now, nil)

		// Buy 10 shares at $100.00 each -> 10*100*100 = 100000 cents
		buy := domain.NewInvestmentTransaction(uuid.New(), investmentId, domain.Buy, 0, domain.Money(1000), domain.Money(10000), now, "Initial Buy")
		// Sell 2 shares at $120.00 each -> 2*120*100 = 24000 cents
		sell := domain.NewInvestmentTransaction(uuid.New(), investmentId, domain.Sell, 0, domain.Money(200), domain.Money(12000), now, "Partial Sell")

		inv.UpdateCurrentValue([]domain.InvestmentTransaction{buy, sell})
		require.Equal(t, domain.Money(76000), inv.CurrentValue)
	})

	t.Run("for savings: deposit minus withdraw", func(t *testing.T) {
		inv := domain.NewInvestment(investmentId, accountId, "Savings", domain.Savings, domain.CurrencyUSD, now, nil)

		deposit := domain.NewInvestmentTransaction(uuid.New(), investmentId, domain.Deposit, domain.Money(50000), 0, 0, now, "Deposit")
		withdraw := domain.NewInvestmentTransaction(uuid.New(), investmentId, domain.Withdraw, domain.Money(20000), 0, 0, now, "Withdraw")

		inv.UpdateCurrentValue([]domain.InvestmentTransaction{deposit, withdraw})
		require.Equal(t, domain.Money(30000), inv.CurrentValue)
	})

	t.Run("ignores transactions of other investments", func(t *testing.T) {
		inv := domain.NewInvestment(investmentId, accountId, "ETF", domain.ETF, domain.CurrencyUSD, now, nil)

		other := domain.NewInvestmentTransaction(uuid.New(), uuid.New(), domain.Deposit, domain.Money(100000), 0, 0, now, "Other")

		inv.UpdateCurrentValue([]domain.InvestmentTransaction{other})
		require.Equal(t, domain.Money(0), inv.CurrentValue)
	})
}

func TestParseInvestmentType(t *testing.T) {
	t.Run("parses valid types", func(t *testing.T) {
		for s, expected := range map[string]domain.InvestmentType{
			"ETF": domain.ETF, "Savings": domain.Savings, "TermDeposit": domain.TermDeposit,
		} {
			got, err := domain.ParseInvestmentType(s)
			require.NoError(t, err)
			require.Equal(t, expected, got)
		}
	})

	t.Run("rejects invalid type", func(t *testing.T) {
		_, err := domain.ParseInvestmentType("Crypto")
		require.Error(t, err)
	})
}
