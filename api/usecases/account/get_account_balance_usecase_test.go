package account_test

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/pedrodbflores/personal-finance-tracker/app_errors"
	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/usecases/account"
	"github.com/pedrodbflores/personal-finance-tracker/usecases/usecasesfakes"
	"github.com/stretchr/testify/require"
)

func mustMoney(t *testing.T, s string) domain.Money {
	t.Helper()
	m, err := domain.ParseMoney(s)
	require.NoError(t, err)
	return m
}

func TestGetAccountBalanceUseCase(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx := context.TODO()
	now := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)

	t.Run("returns the balance combining movements and investments", func(t *testing.T) {
		accountGetter := &usecasesfakes.FakeAccountGetter{}
		movementLister := &usecasesfakes.FakeMovementLister{}
		investmentLister := &usecasesfakes.FakeInvestmentLister{}
		txLister := &usecasesfakes.FakeInvestmentTransactionLister{}

		acc, err := domain.NewAccount(1, "Bank A", domain.AccountTypeBank, domain.CurrencyEUR, now); require.NoError(t, err)
		accountGetter.GetReturns(&acc, nil)

		salary, err := domain.NewMovement(1, "Salary", mustMoney(t, "1000.00"), domain.MovementTypeCredit, now)
		require.NoError(t, err)
		groceries, err := domain.NewMovement(1, "Groceries", mustMoney(t, "25.99"), domain.MovementTypeDebit, now)
		require.NoError(t, err)
		movementLister.ListReturns([]domain.Movement{salary, groceries}, nil)

		etf, err := domain.NewInvestment(2, 1, "Vanguard ETF", domain.ETF, domain.CurrencyEUR, now, nil)
		require.NoError(t, err)
		savings, err := domain.NewInvestment(3, 1, "Savings", domain.Savings, domain.CurrencyEUR, now, nil)
		require.NoError(t, err)
		investmentLister.ListReturns([]domain.Investment{etf, savings}, nil)

		buy, err := domain.NewInvestmentTransaction(10, 2, domain.Buy, 0, mustMoney(t, "10.00"), mustMoney(t, "100.00"), now, "Initial Buy")
		require.NoError(t, err)
		deposit, err := domain.NewInvestmentTransaction(11, 3, domain.Deposit, mustMoney(t, "500.00"), 0, 0, now, "")
		require.NoError(t, err)
		withdraw, err := domain.NewInvestmentTransaction(12, 3, domain.Withdraw, mustMoney(t, "100.00"), 0, 0, now, "")
		require.NoError(t, err)
		txLister.ListStub = func(_ context.Context, _ domain.Id, investmentId domain.Id) ([]domain.InvestmentTransaction, error) {
			switch investmentId {
			case 2:
				return []domain.InvestmentTransaction{buy}, nil
			case 3:
				return []domain.InvestmentTransaction{deposit, withdraw}, nil
			}
			return nil, nil
		}

		useCase := account.NewGetAccountBalanceUseCase(accountGetter, movementLister, investmentLister, txLister, logger)

		result, err := useCase.Execute(ctx, 1)

		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, domain.Id(1), result.AccountId)
		require.Equal(t, domain.CurrencyEUR, result.Currency)
		require.Equal(t, mustMoney(t, "974.01"), result.MovementsBalance)
		require.Equal(t, mustMoney(t, "1400.00"), result.InvestmentsBalance)
		require.Equal(t, mustMoney(t, "2374.01"), result.TotalBalance)

		_, gotId := accountGetter.GetArgsForCall(0)
		require.Equal(t, domain.Id(1), gotId)
		_, listAccId := movementLister.ListArgsForCall(0)
		require.Equal(t, domain.Id(1), listAccId)
		_, invAccId := investmentLister.ListArgsForCall(0)
		require.Equal(t, domain.Id(1), invAccId)
		require.Equal(t, 2, txLister.ListCallCount())
	})

	t.Run("returns zero balances when the account has no movements or investments", func(t *testing.T) {
		accountGetter := &usecasesfakes.FakeAccountGetter{}
		movementLister := &usecasesfakes.FakeMovementLister{}
		investmentLister := &usecasesfakes.FakeInvestmentLister{}
		txLister := &usecasesfakes.FakeInvestmentTransactionLister{}

		acc, err := domain.NewAccount(1, "Empty", domain.AccountTypeBank, domain.CurrencyUSD, now); require.NoError(t, err)
		accountGetter.GetReturns(&acc, nil)
		movementLister.ListReturns([]domain.Movement{}, nil)
		investmentLister.ListReturns([]domain.Investment{}, nil)

		useCase := account.NewGetAccountBalanceUseCase(accountGetter, movementLister, investmentLister, txLister, logger)

		result, err := useCase.Execute(ctx, 1)

		require.NoError(t, err)
		require.Equal(t, domain.Money(0), result.MovementsBalance)
		require.Equal(t, domain.Money(0), result.InvestmentsBalance)
		require.Equal(t, domain.Money(0), result.TotalBalance)
		require.Equal(t, 0, txLister.ListCallCount())
	})

	t.Run("returns an error if the account does not exist", func(t *testing.T) {
		accountGetter := &usecasesfakes.FakeAccountGetter{}
		movementLister := &usecasesfakes.FakeMovementLister{}
		investmentLister := &usecasesfakes.FakeInvestmentLister{}
		txLister := &usecasesfakes.FakeInvestmentTransactionLister{}

		accountGetter.GetReturns(nil, app_errors.NewResourceNotFoundError("Account", "9"))

		useCase := account.NewGetAccountBalanceUseCase(accountGetter, movementLister, investmentLister, txLister, logger)

		result, err := useCase.Execute(ctx, 9)
		require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("Account", "9"))
		require.Nil(t, result)
		require.Equal(t, 0, movementLister.ListCallCount())
	})

	t.Run("propagates movement list errors", func(t *testing.T) {
		accountGetter := &usecasesfakes.FakeAccountGetter{}
		movementLister := &usecasesfakes.FakeMovementLister{}
		investmentLister := &usecasesfakes.FakeInvestmentLister{}
		txLister := &usecasesfakes.FakeInvestmentTransactionLister{}

		acc, err := domain.NewAccount(1, "Bank A", domain.AccountTypeBank, domain.CurrencyEUR, now); require.NoError(t, err)
		accountGetter.GetReturns(&acc, nil)
		movementLister.ListReturns(nil, errors.New("db down"))

		useCase := account.NewGetAccountBalanceUseCase(accountGetter, movementLister, investmentLister, txLister, logger)

		result, err := useCase.Execute(ctx, 1)
		require.EqualError(t, err, "db down")
		require.Nil(t, result)
	})
}
