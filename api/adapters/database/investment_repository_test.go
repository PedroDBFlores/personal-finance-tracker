package database_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/pedrodbflores/personal-finance-tracker/adapters/database"
	"github.com/pedrodbflores/personal-finance-tracker/app_errors"
	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/stretchr/testify/require"
)

func TestInvestmentRepository(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx := context.TODO()
	baseAccount := domain.Account{
		Name:      "Test Account",
		Type:      domain.AccountTypeBrokerage,
		Currency:  domain.ParseCurrency("USD"),
		CreatedAt: time.Now(),
	}

	newInvestment := func(accountId domain.Id) domain.Investment {
		return domain.NewInvestment(0, accountId, "Vanguard ETF", domain.ETF, domain.CurrencyUSD, time.Now(), nil)
	}

	t.Run("it creates an investment in the database", func(t *testing.T) {
		db, err := getInMemoryDatabase(t, logger)
		require.NoError(t, err)
		accRepo := database.NewAccountRepository(db)
		accountId, err := accRepo.Create(ctx, baseAccount)
		require.NoError(t, err)
		repo := database.NewInvestmentRepository(db)

		id, err := repo.Create(ctx, *accountId, newInvestment(*accountId))
		require.NoError(t, err)
		require.NotEqual(t, domain.Id(0), *id)

		inv, err := repo.Get(ctx, *accountId, *id)
		require.NoError(t, err)
		require.NotNil(t, inv)
		require.Equal(t, *id, inv.Id)
		require.Equal(t, *accountId, inv.AccountId)
		require.Equal(t, "Vanguard ETF", inv.Name)
		require.Equal(t, domain.ETF, inv.Type)
		require.Equal(t, domain.Money(0), inv.CurrentValue)
		require.Equal(t, domain.CurrencyUSD, inv.Currency)
	})

	t.Run("get returns ResourceNotFoundError if the investment doesn't exist", func(t *testing.T) {
		db, err := getInMemoryDatabase(t, logger)
		require.NoError(t, err)
		repo := database.NewInvestmentRepository(db)

		id := domain.NewId()
		inv, err := repo.Get(ctx, domain.NewId(), id)
		require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("Investment", id.String()))
		require.Nil(t, inv)
	})

	t.Run("get returns ResourceNotFoundError if the investment belongs to another account", func(t *testing.T) {
		db, err := getInMemoryDatabase(t, logger)
		require.NoError(t, err)
		accRepo := database.NewAccountRepository(db)
		accountId, err := accRepo.Create(ctx, baseAccount)
		require.NoError(t, err)
		repo := database.NewInvestmentRepository(db)
		id, err := repo.Create(ctx, *accountId, newInvestment(*accountId))
		require.NoError(t, err)

		inv, err := repo.Get(ctx, domain.NewId(), *id)
		require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("Investment", id.String()))
		require.Nil(t, inv)
	})

	t.Run("update applies only the given fields", func(t *testing.T) {
		db, err := getInMemoryDatabase(t, logger)
		require.NoError(t, err)
		accRepo := database.NewAccountRepository(db)
		accountId, err := accRepo.Create(ctx, baseAccount)
		require.NoError(t, err)
		repo := database.NewInvestmentRepository(db)
		id, err := repo.Create(ctx, *accountId, newInvestment(*accountId))
		require.NoError(t, err)

		newName := "Renamed ETF"
		newValue := domain.Money(76000)
		err = repo.Update(ctx, *accountId, *id, domain.InvestmentPatch{Name: &newName, CurrentValue: &newValue})
		require.NoError(t, err)

		inv, err := repo.Get(ctx, *accountId, *id)
		require.NoError(t, err)
		require.Equal(t, "Renamed ETF", inv.Name)
		require.Equal(t, domain.Money(76000), inv.CurrentValue)
		require.Equal(t, domain.ETF, inv.Type)
	})

	t.Run("update returns ResourceNotFoundError if the investment doesn't exist", func(t *testing.T) {
		db, err := getInMemoryDatabase(t, logger)
		require.NoError(t, err)
		repo := database.NewInvestmentRepository(db)

		id := domain.NewId()
		newName := "Nope"
		err = repo.Update(ctx, domain.NewId(), id, domain.InvestmentPatch{Name: &newName})
		require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("Investment", id.String()))
	})

	t.Run("list returns only the investments of the given account", func(t *testing.T) {
		db, err := getInMemoryDatabase(t, logger)
		require.NoError(t, err)
		accRepo := database.NewAccountRepository(db)
		accountId, err := accRepo.Create(ctx, baseAccount)
		require.NoError(t, err)
		repo := database.NewInvestmentRepository(db)
		firstId, err := repo.Create(ctx, *accountId, newInvestment(*accountId))
		require.NoError(t, err)
		secondId, err := repo.Create(ctx, *accountId, domain.NewInvestment(0, *accountId, "Savings", domain.Savings, domain.CurrencyUSD, time.Now(), nil))
		require.NoError(t, err)

		investments, err := repo.List(ctx, *accountId)
		require.NoError(t, err)
		require.Len(t, investments, 2)
		require.Equal(t, *firstId, investments[0].Id)
		require.Equal(t, *secondId, investments[1].Id)

		other, err := repo.List(ctx, domain.NewId())
		require.NoError(t, err)
		require.Empty(t, other)
	})

	t.Run("delete removes the investment", func(t *testing.T) {
		db, err := getInMemoryDatabase(t, logger)
		require.NoError(t, err)
		accRepo := database.NewAccountRepository(db)
		accountId, err := accRepo.Create(ctx, baseAccount)
		require.NoError(t, err)
		repo := database.NewInvestmentRepository(db)
		id, err := repo.Create(ctx, *accountId, newInvestment(*accountId))
		require.NoError(t, err)

		err = repo.Delete(ctx, *accountId, *id)
		require.NoError(t, err)

		inv, err := repo.Get(ctx, *accountId, *id)
		require.Error(t, err)
		require.Nil(t, inv)
	})

	t.Run("delete returns ResourceNotFoundError if the investment doesn't exist", func(t *testing.T) {
		db, err := getInMemoryDatabase(t, logger)
		require.NoError(t, err)
		repo := database.NewInvestmentRepository(db)

		id := domain.NewId()
		err = repo.Delete(ctx, domain.NewId(), id)
		require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("Investment", id.String()))
	})
}
