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
		tx, err := domain.NewInvestmentTransaction(id, investmentId, domain.Buy, 0, domain.Money(1000), domain.Money(10000), date, "Buy ETF")
		require.NoError(t, err)

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
		tx, err := domain.NewInvestmentTransaction(id, investmentId, domain.Deposit, domain.Money(50000), 0, 0, date, "Deposit to savings")
		require.NoError(t, err)
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

func TestNewInvestmentTransactionValidation(t *testing.T) {
	now := time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)
	expectValidationError := func(t *testing.T, err error) {
		t.Helper()
		var ve domain.ValidationError
		require.ErrorAs(t, err, &ve)
	}

	t.Run("accepts a valid buy", func(t *testing.T) {
		tx, err := domain.NewInvestmentTransaction(1, 1, domain.Buy, 0, domain.Money(1000), domain.Money(10000), now, "Initial Buy")
		require.NoError(t, err)
		require.Equal(t, domain.Buy, tx.Type)
	})

	t.Run("accepts a valid deposit", func(t *testing.T) {
		_, err := domain.NewInvestmentTransaction(1, 1, domain.Deposit, domain.Money(50000), 0, 0, now, "")
		require.NoError(t, err)
	})

	t.Run("rejects an unsupported transaction type", func(t *testing.T) {
		_, err := domain.NewInvestmentTransaction(1, 1, domain.InvestmentTransactionType("Transfer"), domain.Money(100), 0, 0, now, "")
		expectValidationError(t, err)
	})

	t.Run("rejects a deposit without an amount", func(t *testing.T) {
		_, err := domain.NewInvestmentTransaction(1, 1, domain.Deposit, 0, 0, 0, now, "")
		expectValidationError(t, err)
	})

	t.Run("rejects a buy without a quantity", func(t *testing.T) {
		_, err := domain.NewInvestmentTransaction(1, 1, domain.Buy, 0, 0, domain.Money(10000), now, "")
		expectValidationError(t, err)
	})

	t.Run("rejects a buy without a price", func(t *testing.T) {
		_, err := domain.NewInvestmentTransaction(1, 1, domain.Buy, 0, domain.Money(1000), 0, now, "")
		expectValidationError(t, err)
	})

	t.Run("rejects a negative amount", func(t *testing.T) {
		_, err := domain.NewInvestmentTransaction(1, 1, domain.Withdraw, domain.Money(-100), 0, 0, now, "")
		expectValidationError(t, err)
	})
}
