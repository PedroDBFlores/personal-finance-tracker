package database_test

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"
	"uuid"

	"github.com/dgraph-io/badger/v4"
	"github.com/pedrodbflores/personal-finance-tracker/adapters/database"
	"github.com/pedrodbflores/personal-finance-tracker/adapters/database/logger"
	"github.com/pedrodbflores/personal-finance-tracker/app_errors"
	"github.com/pedrodbflores/personal-finance-tracker/models"
	"github.com/stretchr/testify/require"
)

func getInMemoryDatabase(t *testing.T, _logger *slog.Logger) (*badger.DB, error) {
	t.Helper()
	opts := badger.DefaultOptions("").
		WithInMemory(true).
		WithLogger(logger.NewBadgerSlog(_logger))
	return database.GetDatabase(&opts)
}

func findAccountInDb(t *testing.T, db *badger.DB, id uuid.UUID) *database.Account {
	t.Helper()
	var out *database.Account
	key := fmt.Appendf(nil, "account:%s", id)
	_ = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return nil
		}
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &out)
		})
	})
	return out
}

func TestAccountRepository(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	t.Run("it creates a new account", func(t *testing.T) {
		db, err := getInMemoryDatabase(t, logger)
		require.NoError(t, err)
		repo := database.NewAccountRepository(db)
		now := time.Now()
		account := models.Account{
			Name:      "Random Bank",
			Type:      models.Bank,
			Currency:  "EUR",
			CreatedAt: now,
		}

		id, err := repo.Create(account)
		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil(), id)

		acc := findAccountInDb(t, db, id)
		require.NotNil(t, acc)
		require.Equal(t, id, acc.Id)
		require.Equal(t, "Random Bank", acc.Name)
		require.Equal(t, string(models.Bank), acc.Type)
		require.Equal(t, "EUR", acc.Currency)
		require.WithinDuration(t, now, acc.CreatedAt, 0)
	})

	t.Run("it gets an account by id", func(t *testing.T) {
		db, err := getInMemoryDatabase(t, logger)
		require.NoError(t, err)
		repo := database.NewAccountRepository(db)
		now := time.Now()
		account := models.Account{
			Name:      "Random Bank",
			Type:      models.Bank,
			Currency:  "EUR",
			CreatedAt: now,
		}
		id, err := repo.Create(account)
		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil(), id)

		acc, err := repo.Get(id)

		require.NoError(t, err)
		require.Equal(t, id, acc.Id)
		require.Equal(t, "Random Bank", acc.Name)
		require.Equal(t, models.Bank, acc.Type)
		require.Equal(t, "EUR", acc.Currency)
		require.WithinDuration(t, now, acc.CreatedAt, 0)
	})

	t.Run("returns ResourceNotFoundError if the account doesn't exist", func(t *testing.T) {
		db, err := getInMemoryDatabase(t, logger)
		require.NoError(t, err)
		repo := database.NewAccountRepository(db)

		id := uuid.NewV7()
		acc, err := repo.Get(id)
		require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("account", id.String()))
		require.Nil(t, acc)
	})

	t.Run("it updates an account", func(t *testing.T) {
		db, err := getInMemoryDatabase(t, logger)
		require.NoError(t, err)
		repo := database.NewAccountRepository(db)
		now := time.Now()
		account := models.Account{
			Name:      "Random Bank",
			Type:      models.Bank,
			Currency:  "EUR",
			CreatedAt: now,
		}

		id, err := repo.Create(account)
		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil(), id)

		acc, err := repo.Get(id)
		require.NoError(t, err)
		acc.Name = "Right Bank"
		err = repo.Update(*acc)
		require.NoError(t, err)

		updatedAcc, err := repo.Get(id)
		require.NoError(t, err)
		require.Equal(t, "Right Bank", updatedAcc.Name)
	})

	t.Run("lists all accounts", func(t *testing.T) {
		db, err := getInMemoryDatabase(t, logger)
		require.NoError(t, err)
		repo := database.NewAccountRepository(db)
		now := time.Now()
		accountOne := models.Account{
			Name:      "Random Bank",
			Type:      models.Bank,
			Currency:  "EUR",
			CreatedAt: now,
		}
		accountTwo := models.Account{
			Name:      "A fine brokerage",
			Type:      models.Brokerage,
			Currency:  "USD",
			CreatedAt: now,
		}

		repo.Create(accountOne)
		repo.Create(accountTwo)

		accounts, err := repo.List()

		require.NoError(t, err)
		require.Len(t, accounts, 2)
		require.Equal(t, "Random Bank", accounts[0].Name)
		require.Equal(t, models.Bank, accounts[0].Type)
		require.Equal(t, "EUR", accounts[0].Currency)
		require.WithinDuration(t, now, accounts[0].CreatedAt, 0)
		require.Equal(t, "A fine brokerage", accounts[1].Name)
		require.Equal(t, models.Brokerage, accounts[1].Type)
		require.Equal(t, "USD", accounts[1].Currency)
		require.WithinDuration(t, now, accounts[1].CreatedAt, 0)
	})

	t.Run("deletes an account by id", func(t *testing.T) {
		db, err := getInMemoryDatabase(t, logger)
		require.NoError(t, err)
		repo := database.NewAccountRepository(db)
		now := time.Now()
		account := models.Account{
			Name:      "Random Bank",
			Type:      models.Bank,
			Currency:  "EUR",
			CreatedAt: now,
		}
		id, err := repo.Create(account)
		require.NoError(t, err)

		err = repo.Delete(id)
		require.NoError(t, err)

		acc, err := repo.Get(id)
		require.Error(t, err)
		require.Nil(t, acc)
	})
}
