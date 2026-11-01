package database_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
	"uuid"

	"github.com/dgraph-io/badger/v4"
	"github.com/pedrodbflores/personal-finance-tracker/adapters/database"
	"github.com/pedrodbflores/personal-finance-tracker/models"
	"github.com/stretchr/testify/require"
)

func getInMemoryDatabase(t *testing.T) (*badger.DB, error) {
	t.Helper()
	opts := badger.DefaultOptions("").WithInMemory(true)
	return database.GetDatabase(&opts)
}

func findAccountInDb(t *testing.T, db *badger.DB, id uuid.UUID) *database.Account {
	t.Helper()
	var out *database.Account
	key := []byte(fmt.Sprintf("account#%s", id))
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
	t.Run("it creates a new account", func(t *testing.T) {
		db, err := getInMemoryDatabase(t)
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
		require.Equal(t, &database.Account{
			Id:        id,
			Name:      "Random Bank",
			Type:      string(models.Bank),
			Currency:  "EUR",
			CreatedAt: now,
		}, acc)
	})
}
