package database_test

import (
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"
	"uuid"

	"github.com/dgraph-io/badger/v4"
	"github.com/pedrodbflores/personal-finance-tracker/adapters/database"
	"github.com/pedrodbflores/personal-finance-tracker/models"
	"github.com/stretchr/testify/require"
)

func findMovementInDb(t *testing.T, db *badger.DB, accountId, id uuid.UUID) *map[string]any {
	t.Helper()
	key := fmt.Appendf(nil, "account:%s:movement:%s", accountId, id)
	out, err := database.GetByKey[*map[string]any](db, key)
	if err != nil {
		t.Fatalf("failed to get movement by key: %v", err)
	}
	return out
}

func TestMovementRepository(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	baseAccount := models.Account{
		Name:      "Test Account",
		Type:      models.AccountTypeBank,
		Currency:  "USD",
		CreatedAt: time.Now(),
	}

	t.Run("it creates a movement in the database", func(t *testing.T) {
		db, err := getInMemoryDatabase(t, logger)
		require.NoError(t, err)
		accRepo := database.NewAccountRepository(db)
		accountId, err := accRepo.Create(baseAccount)
		require.NoError(t, err)
		repo := database.NewMovementRepository(db)
		movement := models.NewMovement(
			uuid.New(),
			"Test Movement",
			int64(1000),
			models.MovementTypeDebit,
			time.Now(),
		)
		id, err := repo.Create(accountId, movement)
		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil(), id)

		mov := findMovementInDb(t, db, accountId, id)
		require.NotNil(t, mov)
		require.Equal(t, accountId, uuid.MustParse((*mov)["accountId"].(string)))
		require.Equal(t, string(models.MovementTypeDebit), (*mov)["type"])
		require.Equal(t, 1000, (*mov)["amount"])
		require.Equal(t, "USD", (*mov)["currency"])
		dateParsed, err := time.Parse(time.RFC3339, (*mov)["date"].(string))
		require.NoError(t, err)
		require.WithinDuration(t, movement.Date, dateParsed, 0)
	})
}
