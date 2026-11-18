package main

import (
	"fmt"
	"log/slog"
	"os"

	tui "github.com/grindlemire/go-tui"
	"github.com/pedrodbflores/personal-finance-tracker/adapters/database"
	"github.com/pedrodbflores/personal-finance-tracker/usecases/account"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("tracker.db"), &gorm.Config{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open database: %v\n", err)
		os.Exit(1)
	}
	if err := db.AutoMigrate(&database.Account{}, &database.Movement{}); err != nil {
		fmt.Fprintf(os.Stderr, "failed to migrate database: %v\n", err)
		os.Exit(1)
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	accountRepo := database.NewAccountRepository(db)
	listAccounts := account.NewListAccountsUseCase(&accountRepo, logger)

	screen := newAccountsScreen(&listAccounts)
	app, err := tui.NewApp(tui.WithRootComponent(screen))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create app: %v\n", err)
		os.Exit(1)
	}
	defer app.Close()

	screen.load(app.EventQueue())
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error running app: %v\n", err)
		os.Exit(1)
	}
}
