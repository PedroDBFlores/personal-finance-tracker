package database_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
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

	setup := func(t *testing.T) (domain.Id, uuid.UUID, database.InvestmentRepository, database.InvestmentTransactionRepository) {
		t.Helper()
		db, err := getInMemoryDatabase(t, logger)
		require.NoError(t, err)
		accRepo := database.NewAccountRepository(db)
		accountId, err := accRepo.Create(ctx, baseAccount)
		require.NoError(t, err)
		invRepo := database.NewInvestmentRepository(db)
		investmentId, err := invRepo.Create(ctx, *accountId, domain.NewInvestment(uuid.New(), *accountId, "Vanguard ETF", domain.ETF, domain.CurrencyUSD, time.Now(), nil))
		require.NoError(t, err)
		return *accountId, *investmentId, invRepo, database.NewInvestmentTransactionRepository(db)
	}

	newTx := func(investmentId uuid.UUID) domain.InvestmentTransaction {
		return domain.NewInvestmentTransaction(uuid.New(), investmentId, domain.Buy, 0, domain.Money(1000), domain.Money(10000), time.Now(), "Initial Buy")
	}

	t.Run("it creates a transaction in the database", func(t *testing.T) {
		_, investmentId, _, repo := setup(t)

		id, err := repo.Create(ctx, investmentId, newTx(investmentId))
		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil, *id)

		tx, err := repo.Get(ctx, investmentId, *id)
		require.NoError(t, err)
		require.NotNil(t, tx)
		require.Equal(t, *id, tx.Id)
		require.Equal(t, investmentId, tx.InvestmentId)
		require.Equal(t, domain.Buy, tx.Type)
		require.Equal(t, domain.Money(1000), tx.Quantity)
		require.Equal(t, domain.Money(10000), tx.Price)
		require.Equal(t, "Initial Buy", tx.Description)
	})

	t.Run("get returns ResourceNotFoundError if the transaction doesn't exist", func(t *testing.T) {
		_, investmentId, _, repo := setup(t)

		id := uuid.New()
		tx, err := repo.Get(ctx, investmentId, id)
		require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("InvestmentTransaction", id.String()))
		require.Nil(t, tx)
	})

	t.Run("get returns ResourceNotFoundError if the transaction belongs to another investment", func(t *testing.T) {
		_, investmentId, _, repo := setup(t)
		id, err := repo.Create(ctx, investmentId, newTx(investmentId))
		require.NoError(t, err)

		tx, err := repo.Get(ctx, uuid.New(), *id)
		require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("InvestmentTransaction", id.String()))
		require.Nil(t, tx)
	})

	t.Run("update applies only the given fields", func(t *testing.T) {
		_, investmentId, _, repo := setup(t)
		id, err := repo.Create(ctx, investmentId, newTx(investmentId))
		require.NoError(t, err)

		newDescription := "Updated"
		newPrice := domain.Money(12000)
		err = repo.Update(ctx, investmentId, *id, domain.InvestmentTransactionPatch{Description: &newDescription, Price: &newPrice})
		require.NoError(t, err)

		tx, err := repo.Get(ctx, investmentId, *id)
		require.NoError(t, err)
		require.Equal(t, "Updated", tx.Description)
		require.Equal(t, domain.Money(12000), tx.Price)
		require.Equal(t, domain.Buy, tx.Type)
		require.Equal(t, domain.Money(1000), tx.Quantity)
	})

	t.Run("update returns ResourceNotFoundError if the transaction doesn't exist", func(t *testing.T) {
		_, investmentId, _, repo := setup(t)

		id := uuid.New()
		newPrice := domain.Money(12000)
		err := repo.Update(ctx, investmentId, id, domain.InvestmentTransactionPatch{Price: &newPrice})
		require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("InvestmentTransaction", id.String()))
	})

	t.Run("list returns only the transactions of the given investment", func(t *testing.T) {
		_, investmentId, _, repo := setup(t)
		firstId, err := repo.Create(ctx, investmentId, newTx(investmentId))
		require.NoError(t, err)
		secondId, err := repo.Create(ctx, investmentId, domain.NewInvestmentTransaction(uuid.New(), investmentId, domain.Sell, 0, domain.Money(200), domain.Money(12000), time.Now(), "Sell"))
		require.NoError(t, err)

		txs, err := repo.List(ctx, investmentId)
		require.NoError(t, err)
		require.Len(t, txs, 2)
		require.Equal(t, *firstId, txs[0].Id)
		require.Equal(t, *secondId, txs[1].Id)

		other, err := repo.List(ctx, uuid.New())
		require.NoError(t, err)
		require.Empty(t, other)
	})

	t.Run("delete removes the transaction", func(t *testing.T) {
		_, investmentId, _, repo := setup(t)
		id, err := repo.Create(ctx, investmentId, newTx(investmentId))
		require.NoError(t, err)

		err = repo.Delete(ctx, investmentId, *id)
		require.NoError(t, err)

		tx, err := repo.Get(ctx, investmentId, *id)
		require.Error(t, err)
		require.Nil(t, tx)
	})

	t.Run("delete returns ResourceNotFoundError if the transaction doesn't exist", func(t *testing.T) {
		_, investmentId, _, repo := setup(t)

		id := uuid.New()
		err := repo.Delete(ctx, investmentId, id)
		require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("InvestmentTransaction", id.String()))
	})
}
