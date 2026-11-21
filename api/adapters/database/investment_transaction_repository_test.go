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

func TestInvestmentTransactionRepository(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx := context.TODO()
	baseAccount := domain.Account{
		Name:      "Test Account",
		Type:      domain.AccountTypeBrokerage,
		Currency:  domain.ParseCurrency("USD"),
		CreatedAt: time.Now(),
	}

	setup := func(t *testing.T) (domain.Id, domain.Id, database.InvestmentRepository, database.InvestmentTransactionRepository) {
		t.Helper()
		db, err := getInMemoryDatabase(t, logger)
		require.NoError(t, err)
		accRepo := database.NewAccountRepository(db)
		accountId, err := accRepo.Create(ctx, baseAccount)
		require.NoError(t, err)
		invRepo := database.NewInvestmentRepository(db)
		inv, err := domain.NewInvestment(0, *accountId, "Vanguard ETF", domain.ETF, domain.CurrencyUSD, time.Now(), nil)
		require.NoError(t, err)
		investmentId, err := invRepo.Create(ctx, *accountId, inv)
		require.NoError(t, err)
		return *accountId, *investmentId, invRepo, database.NewInvestmentTransactionRepository(db)
	}

	newTx := func(investmentId domain.Id) domain.InvestmentTransaction {
		tx, err := domain.NewInvestmentTransaction(0, investmentId, domain.Buy, 0, domain.Money(1000), domain.Money(10000), time.Now(), "Initial Buy")
		if err != nil {
			panic(err)
		}
		return tx
	}

	t.Run("it creates a transaction in the database", func(t *testing.T) {
		accountId, investmentId, _, repo := setup(t)

		id, err := repo.Create(ctx, accountId, investmentId, newTx(investmentId))
		require.NoError(t, err)
		require.NotEqual(t, domain.Id(0), *id)

		tx, err := repo.Get(ctx, accountId, investmentId, *id)
		require.NoError(t, err)
		require.NotNil(t, tx)
		require.Equal(t, *id, tx.Id)
		require.Equal(t, investmentId, tx.InvestmentId)
		require.Equal(t, domain.Buy, tx.Type)
		require.Equal(t, domain.Money(1000), tx.Quantity)
		require.Equal(t, domain.Money(10000), tx.Price)
		require.Equal(t, "Initial Buy", tx.Description)
	})

	t.Run("create returns ResourceNotFoundError if the investment belongs to another account", func(t *testing.T) {
		_, investmentId, _, repo := setup(t)

		id, err := repo.Create(ctx, domain.Id(999), investmentId, newTx(investmentId))
		require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("Investment", investmentId.String()))
		require.Nil(t, id)
	})

	t.Run("get returns ResourceNotFoundError if the transaction doesn't exist", func(t *testing.T) {
		accountId, investmentId, _, repo := setup(t)

		id := domain.NewId()
		tx, err := repo.Get(ctx, accountId, investmentId, id)
		require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("InvestmentTransaction", id.String()))
		require.Nil(t, tx)
	})

	t.Run("get returns ResourceNotFoundError if the transaction belongs to another investment", func(t *testing.T) {
		accountId, investmentId, _, repo := setup(t)
		id, err := repo.Create(ctx, accountId, investmentId, newTx(investmentId))
		require.NoError(t, err)

		tx, err := repo.Get(ctx, accountId, domain.NewId(), *id)
		require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("InvestmentTransaction", id.String()))
		require.Nil(t, tx)
	})

	t.Run("get returns ResourceNotFoundError if the investment belongs to another account", func(t *testing.T) {
		accountId, investmentId, _, repo := setup(t)
		id, err := repo.Create(ctx, accountId, investmentId, newTx(investmentId))
		require.NoError(t, err)

		tx, err := repo.Get(ctx, domain.Id(999), investmentId, *id)
		require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("InvestmentTransaction", id.String()))
		require.Nil(t, tx)
	})

	t.Run("update applies only the given fields", func(t *testing.T) {
		accountId, investmentId, _, repo := setup(t)
		id, err := repo.Create(ctx, accountId, investmentId, newTx(investmentId))
		require.NoError(t, err)

		newDescription := "Updated"
		newPrice := domain.Money(12000)
		err = repo.Update(ctx, accountId, investmentId, *id, domain.InvestmentTransactionPatch{Description: &newDescription, Price: &newPrice})
		require.NoError(t, err)

		tx, err := repo.Get(ctx, accountId, investmentId, *id)
		require.NoError(t, err)
		require.Equal(t, "Updated", tx.Description)
		require.Equal(t, domain.Money(12000), tx.Price)
		require.Equal(t, domain.Buy, tx.Type)
		require.Equal(t, domain.Money(1000), tx.Quantity)
	})

	t.Run("update returns ResourceNotFoundError if the transaction doesn't exist", func(t *testing.T) {
		accountId, investmentId, _, repo := setup(t)

		id := domain.NewId()
		newPrice := domain.Money(12000)
		err := repo.Update(ctx, accountId, investmentId, id, domain.InvestmentTransactionPatch{Price: &newPrice})
		require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("InvestmentTransaction", id.String()))
	})

	t.Run("update returns ResourceNotFoundError if the investment belongs to another account", func(t *testing.T) {
		accountId, investmentId, _, repo := setup(t)
		id, err := repo.Create(ctx, accountId, investmentId, newTx(investmentId))
		require.NoError(t, err)

		newPrice := domain.Money(12000)
		err = repo.Update(ctx, domain.Id(999), investmentId, *id, domain.InvestmentTransactionPatch{Price: &newPrice})
		require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("InvestmentTransaction", id.String()))
	})

	t.Run("list returns only the transactions of the given investment", func(t *testing.T) {
		accountId, investmentId, _, repo := setup(t)
		firstId, err := repo.Create(ctx, accountId, investmentId, newTx(investmentId))
		require.NoError(t, err)
		sell, err := domain.NewInvestmentTransaction(0, investmentId, domain.Sell, 0, domain.Money(200), domain.Money(12000), time.Now(), "Sell")
		require.NoError(t, err)
		secondId, err := repo.Create(ctx, accountId, investmentId, sell)
		require.NoError(t, err)

		txs, err := repo.List(ctx, accountId, investmentId)
		require.NoError(t, err)
		require.Len(t, txs, 2)
		require.Equal(t, *firstId, txs[0].Id)
		require.Equal(t, *secondId, txs[1].Id)

		other, err := repo.List(ctx, accountId, domain.NewId())
		require.NoError(t, err)
		require.Empty(t, other)
	})

	t.Run("list returns empty if the investment belongs to another account", func(t *testing.T) {
		accountId, investmentId, _, repo := setup(t)
		_, err := repo.Create(ctx, accountId, investmentId, newTx(investmentId))
		require.NoError(t, err)

		txs, err := repo.List(ctx, domain.Id(999), investmentId)
		require.NoError(t, err)
		require.Empty(t, txs)
	})

	t.Run("delete removes the transaction", func(t *testing.T) {
		accountId, investmentId, _, repo := setup(t)
		id, err := repo.Create(ctx, accountId, investmentId, newTx(investmentId))
		require.NoError(t, err)

		err = repo.Delete(ctx, accountId, investmentId, *id)
		require.NoError(t, err)

		tx, err := repo.Get(ctx, accountId, investmentId, *id)
		require.Error(t, err)
		require.Nil(t, tx)
	})

	t.Run("delete returns ResourceNotFoundError if the transaction doesn't exist", func(t *testing.T) {
		accountId, investmentId, _, repo := setup(t)

		id := domain.NewId()
		err := repo.Delete(ctx, accountId, investmentId, id)
		require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("InvestmentTransaction", id.String()))
	})

	t.Run("delete returns ResourceNotFoundError if the investment belongs to another account", func(t *testing.T) {
		accountId, investmentId, _, repo := setup(t)
		id, err := repo.Create(ctx, accountId, investmentId, newTx(investmentId))
		require.NoError(t, err)

		err = repo.Delete(ctx, domain.Id(999), investmentId, *id)
		require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("InvestmentTransaction", id.String()))

		// the transaction is untouched
		tx, err := repo.Get(ctx, accountId, investmentId, *id)
		require.NoError(t, err)
		require.NotNil(t, tx)
	})
}
