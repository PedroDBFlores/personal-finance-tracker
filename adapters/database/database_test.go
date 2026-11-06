package database_test

import (
	"log/slog"
	"testing"

	"github.com/dgraph-io/badger/v4"
	"github.com/pedrodbflores/personal-finance-tracker/adapters/database"
	"github.com/pedrodbflores/personal-finance-tracker/adapters/database/logger"
)

func getInMemoryDatabase(t *testing.T, _logger *slog.Logger) (*badger.DB, error) {
	t.Helper()
	opts := badger.DefaultOptions("").
		WithInMemory(true).
		WithLogger(logger.NewBadgerSlog(_logger))
	return database.GetDatabase(&opts)
}
