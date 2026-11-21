package domain_test

import (
	"testing"
	"time"

	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/stretchr/testify/require"
)

func TestNewInvestmentTransaction(t *testing.T) {
	id := domain.NewId()
	investmentId := domain.NewId()
	date := time.Now()

	t.Run("creates a buy transaction", func(t *testing.T) {
		tx := domain.NewInvestmentTransaction(id, investmentId, domain.Buy, 0, domain.Money(1000), domain.Money(10000), date, "Buy ETF")

		require.Equal(t, id, tx.Id)
		require.Equal(t, investmentId, tx.InvestmentId)
		require.Equal(t, domain.Buy, tx.Type)
		require.Equal(t, domain.Money(0), tx.Amount)
		require.Equal(t, domain.Money(1000), tx.Quantity)
		require.Equal(t, domain.Money(10000), tx.Price)
		require.True(t, tx.Date.Equal(date))
		require.Equal(t, "Buy ETF", tx.Description)
	})

	t.Run("creates a deposit transaction", func(t *testing.T) {
		tx := domain.NewInvestmentTransaction(id, investmentId, domain.Deposit, domain.Money(50000), 0, 0, date, "Deposit to savings")
		require.Equal(t, domain.Deposit, tx.Type)
		require.Equal(t, domain.Money(50000), tx.Amount)
	})
}

func TestParseInvestmentTransactionType(t *testing.T) {
	t.Run("parses valid types", func(t *testing.T) {
		for s, expected := range map[string]domain.InvestmentTransactionType{
			"Buy": domain.Buy, "Sell": domain.Sell, "Deposit": domain.Deposit, "Withdraw": domain.Withdraw,
		} {
			got, err := domain.ParseInvestmentTransactionType(s)
			require.NoError(t, err)
			require.Equal(t, expected, got)
		}
	})

	t.Run("rejects invalid type", func(t *testing.T) {
		_, err := domain.ParseInvestmentTransactionType("Transfer")
		require.Error(t, err)
	})
}
