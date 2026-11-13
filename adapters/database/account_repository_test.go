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

func TestAccountRepository(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx := context.TODO()
	t.Run("it creates a new account", func(t *testing.T) {
		db, err := getInMemoryDatabase(t, logger)
		require.NoError(t, err)
		repo := database.NewAccountRepository(db)
		now := time.Now()
		account := domain.Account{
			Name:      "Random Bank",
			Type:      domain.AccountTypeBank,
			Currency:  domain.CurrencyEUR,
			CreatedAt: now,
		}

		id, err := repo.Create(ctx, account)
		require.NoError(t, err)
		require.Greater(t, *id, uint(0))

		acc, err := repo.Get(ctx, *id)
		require.NoError(t, err)
		require.NotNil(t, acc)
		require.Equal(t, *id, acc.Id)
		require.Equal(t, "Random Bank", acc.Name)
		require.Equal(t, domain.AccountTypeBank, acc.Type)
		require.Equal(t, domain.CurrencyEUR, acc.Currency)
		require.WithinDuration(t, now, acc.CreatedAt, 0)
	})

	t.Run("get", func(t *testing.T) {
		t.Run("it gets an account by id", func(t *testing.T) {
			db, err := getInMemoryDatabase(t, logger)
			require.NoError(t, err)
			repo := database.NewAccountRepository(db)
			now := time.Now()
			account := domain.Account{
				Name:      "Random Bank",
				Type:      domain.AccountTypeBank,
				Currency:  domain.CurrencyEUR,
				CreatedAt: now,
			}
			id, err := repo.Create(ctx, account)
			require.NoError(t, err)
			require.Greater(t, *id, uint(0))

			acc, err := repo.Get(ctx, *id)

			require.NoError(t, err)
			require.Equal(t, *id, acc.Id)
			require.Equal(t, "Random Bank", acc.Name)
			require.Equal(t, domain.AccountTypeBank, acc.Type)
			require.Equal(t, domain.CurrencyEUR, acc.Currency)
			require.WithinDuration(t, now, acc.CreatedAt, 0)
		})

		t.Run("returns ResourceNotFoundError if the account doesn't exist", func(t *testing.T) {
			db, err := getInMemoryDatabase(t, logger)
			require.NoError(t, err)
			repo := database.NewAccountRepository(db)

			id := domain.NewId()
			acc, err := repo.Get(ctx, id)
			require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("Account", id.String()))
			require.Nil(t, acc)
		})
	})

	t.Run("update", func(t *testing.T) {
		t.Run("it updates an account", func(t *testing.T) {
			db, err := getInMemoryDatabase(t, logger)
			require.NoError(t, err)
			repo := database.NewAccountRepository(db)
			now := time.Now()
			account := domain.Account{
				Name:      "Random Bank",
				Type:      domain.AccountTypeBank,
				Currency:  domain.CurrencyEUR,
				CreatedAt: now,
			}

			id, err := repo.Create(ctx, account)
			require.NoError(t, err)
			require.NotEqual(t, domain.Id(0), id)

			acc, err := repo.Get(ctx, *id)
			require.NoError(t, err)
			acc.Name = "Right Bank"
			err = repo.Update(ctx, *acc)
			require.NoError(t, err)

			updatedAcc, err := repo.Get(ctx, *id)
			require.NoError(t, err)
			require.Equal(t, "Right Bank", updatedAcc.Name)
		})

		t.Run("returns ResourceNotFoundError if updating a non-existent account", func(t *testing.T) {
			db, err := getInMemoryDatabase(t, logger)
			require.NoError(t, err)
			repo := database.NewAccountRepository(db)

			account := domain.Account{
				Id:        domain.NewId(),
				Name:      "Non-existent Bank",
				Type:      domain.AccountTypeBank,
				Currency:  domain.CurrencyEUR,
				CreatedAt: time.Now(),
			}
			err = repo.Update(ctx, account)
			require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("Account", account.Id.String()))
		})
	})

	t.Run("lists all accounts", func(t *testing.T) {
		db, err := getInMemoryDatabase(t, logger)
		require.NoError(t, err)
		repo := database.NewAccountRepository(db)
		now := time.Now()
		accountOne := domain.Account{
			Name:      "Random Bank",
			Type:      domain.AccountTypeBank,
			Currency:  domain.CurrencyEUR,
			CreatedAt: now,
		}
		accountTwo := domain.Account{
			Name:      "A fine brokerage",
			Type:      domain.AccountTypeBrokerage,
			Currency:  domain.CurrencyUSD,
			CreatedAt: now,
		}

		repo.Create(ctx, accountOne)
		repo.Create(ctx, accountTwo)

		accounts, err := repo.List(ctx)

		require.NoError(t, err)
		require.Len(t, accounts, 2)
		require.Equal(t, "Random Bank", accounts[0].Name)
		require.Equal(t, domain.AccountTypeBank, accounts[0].Type)
		require.Equal(t, domain.CurrencyEUR, accounts[0].Currency)
		require.WithinDuration(t, now, accounts[0].CreatedAt, 0)
		require.Equal(t, "A fine brokerage", accounts[1].Name)
		require.Equal(t, domain.AccountTypeBrokerage, accounts[1].Type)
		require.Equal(t, domain.CurrencyUSD, accounts[1].Currency)
		require.WithinDuration(t, now, accounts[1].CreatedAt, 0)
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("deletes an account by id", func(t *testing.T) {
			db, err := getInMemoryDatabase(t, logger)
			require.NoError(t, err)
			repo := database.NewAccountRepository(db)
			now := time.Now()
			account := domain.Account{
				Name:      "Random Bank",
				Type:      domain.AccountTypeBank,
				Currency:  domain.CurrencyEUR,
				CreatedAt: now,
			}
			id, err := repo.Create(ctx, account)
			require.NoError(t, err)

			err = repo.Delete(ctx, *id)
			require.NoError(t, err)

			acc, err := repo.Get(ctx, *id)
			require.Error(t, err)
			require.Nil(t, acc)
		})

		t.Run("returns ResourceNotFoundError if the account doesn't exist", func(t *testing.T) {
			db, err := getInMemoryDatabase(t, logger)
			require.NoError(t, err)
			repo := database.NewAccountRepository(db)

			id := domain.NewId()
			err = repo.Delete(ctx, id)

			require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("Account", id.String()))
		})
	})
}
