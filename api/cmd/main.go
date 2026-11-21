package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/pedrodbflores/personal-finance-tracker/adapters/database"
	"github.com/pedrodbflores/personal-finance-tracker/adapters/rest"
	"github.com/pedrodbflores/personal-finance-tracker/usecases/account"
	"github.com/pedrodbflores/personal-finance-tracker/usecases/investment"
	"github.com/pedrodbflores/personal-finance-tracker/usecases/investmenttransaction"
	"github.com/pedrodbflores/personal-finance-tracker/usecases/movement"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("tracker.db"), &gorm.Config{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open database: %v\n", err)
		os.Exit(1)
	}
	if err := db.AutoMigrate(&database.Account{}, &database.Movement{}, &database.Investment{}, &database.InvestmentTransaction{}); err != nil {
		fmt.Fprintf(os.Stderr, "failed to migrate database: %v\n", err)
		os.Exit(1)
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	accountRepo := database.NewAccountRepository(db)
	movementRepo := database.NewMovementRepository(db)
	investmentRepo := database.NewInvestmentRepository(db)
	investmentTransactionRepo := database.NewInvestmentTransactionRepository(db)

	createAccount := account.NewCreateAccountUseCase(&accountRepo, logger)
	getAccount := account.NewGetAccountUseCase(&accountRepo, logger)
	listAccounts := account.NewListAccountsUseCase(&accountRepo, logger)
	updateAccount := account.NewUpdateAccountUseCase(&accountRepo, logger)
	deleteAccount := account.NewDeleteAccountUseCase(&accountRepo, logger)
	getAccountBalance := account.NewGetAccountBalanceUseCase(&accountRepo, &movementRepo, &investmentRepo, &investmentTransactionRepo, logger)

	createMovement := movement.NewCreateMovementUseCase(&movementRepo, logger)
	getMovement := movement.NewGetMovementUseCase(&movementRepo, logger)
	listMovements := movement.NewListMovementsUseCase(&movementRepo, logger)
	updateMovement := movement.NewUpdateMovementUseCase(&movementRepo, logger)
	deleteMovement := movement.NewDeleteMovementUseCase(&movementRepo, logger)

	createInvestment := investment.NewCreateInvestmentUseCase(&investmentRepo, logger)
	getInvestment := investment.NewGetInvestmentUseCase(&investmentRepo, &investmentTransactionRepo, logger)
	listInvestments := investment.NewListInvestmentsUseCase(&investmentRepo, &investmentTransactionRepo, logger)
	updateInvestment := investment.NewUpdateInvestmentUseCase(&investmentRepo, logger)
	deleteInvestment := investment.NewDeleteInvestmentUseCase(&investmentRepo, logger)

	createInvestmentTransaction := investmenttransaction.NewCreateInvestmentTransactionUseCase(&investmentTransactionRepo, logger)
	getInvestmentTransaction := investmenttransaction.NewGetInvestmentTransactionUseCase(&investmentTransactionRepo, logger)
	listInvestmentTransactions := investmenttransaction.NewListInvestmentTransactionsUseCase(&investmentTransactionRepo, logger)
	updateInvestmentTransaction := investmenttransaction.NewUpdateInvestmentTransactionUseCase(&investmentTransactionRepo, logger)
	deleteInvestmentTransaction := investmenttransaction.NewDeleteInvestmentTransactionUseCase(&investmentTransactionRepo, logger)

	accountHandler := rest.NewAccountHandler(&createAccount, &getAccount, &listAccounts, &updateAccount, &deleteAccount, &getAccountBalance, logger)
	movementHandler := rest.NewMovementHandler(&createMovement, &getMovement, &listMovements, &updateMovement, &deleteMovement, logger)
	investmentHandler := rest.NewInvestmentHandler(&createInvestment, &getInvestment, &listInvestments, &updateInvestment, &deleteInvestment, logger)
	investmentTransactionHandler := rest.NewInvestmentTransactionHandler(&createInvestmentTransaction, &getInvestmentTransaction, &listInvestmentTransactions, &updateInvestmentTransaction, &deleteInvestmentTransaction, logger)

	addr := ":8080"
	fmt.Fprintf(os.Stderr, "listening on %s\n", addr)
	if err := http.ListenAndServe(addr, rest.NewRouter(accountHandler, movementHandler, investmentHandler, investmentTransactionHandler)); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}
