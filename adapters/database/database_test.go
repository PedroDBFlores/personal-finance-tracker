package database_test

import (
	"log/slog"
	"testing"

	"github.com/pedrodbflores/personal-finance-tracker/adapters/database"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func getInMemoryDatabase(t *testing.T, _logger *slog.Logger) (*gorm.DB, error) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		// Logger: _logger,
	})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&database.Account{}, &database.Movement{})
	return db, err
}
