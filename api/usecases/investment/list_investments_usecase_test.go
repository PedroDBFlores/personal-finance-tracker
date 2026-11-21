package investment_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/pedrodbflores/personal-finance-tracker/app_errors"
	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/usecases/investment"
	"github.com/pedrodbflores/personal-finance-tracker/usecases/usecasesfakes"
	"github.com/stretchr/testify/require"
)

func TestListInvestmentsUseCase(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx := context.TODO()
	accountId := domain.NewId()

	t.Run("lists investments with current values calculated from transactions", func(t *testing.T) {
		lister := usecasesfakes.FakeInvestmentLister{}
		txLister := usecasesfakes.FakeInvestmentTransactionLister{}
		invId := domain.NewId()
		inv, err := domain.NewInvestment(invId, accountId, "Savings", domain.Savings, domain.CurrencyUSD, time.Now(), nil); require.NoError(t, err)
		lister.ListReturns([]domain.Investment{inv}, nil)
		deposit, err := domain.NewInvestmentTransaction(domain.NewId(), invId, domain.Deposit, domain.Money(50000), 0, 0, time.Now(), "Deposit"); require.NoError(t, err)
		txLister.ListReturns([]domain.InvestmentTransaction{deposit}, nil)
		useCase := investment.NewListInvestmentsUseCase(&lister, &txLister, logger)

		dtos, err := useCase.Execute(ctx, accountId)
		require.NoError(t, err)
		require.Len(t, dtos, 1)
		require.Equal(t, "Savings", dtos[0].Name)
		require.Equal(t, domain.Money(50000), dtos[0].CurrentValue)

		ctxArg, calledAccountId := lister.ListArgsForCall(0)
		require.Equal(t, ctx, ctxArg)
		require.Equal(t, accountId, calledAccountId)
	})

	t.Run("returns an empty list when there are no investments", func(t *testing.T) {
		lister := usecasesfakes.FakeInvestmentLister{}
		txLister := usecasesfakes.FakeInvestmentTransactionLister{}
		lister.ListReturns([]domain.Investment{}, nil)
		useCase := investment.NewListInvestmentsUseCase(&lister, &txLister, logger)

		dtos, err := useCase.Execute(ctx, accountId)
		require.NoError(t, err)
		require.Empty(t, dtos)
	})

	t.Run("returns an error if listing fails", func(t *testing.T) {
		lister := usecasesfakes.FakeInvestmentLister{}
		txLister := usecasesfakes.FakeInvestmentTransactionLister{}
		lister.ListReturns(nil, app_errors.NewResourceNotFoundError("Investment", "1"))
		useCase := investment.NewListInvestmentsUseCase(&lister, &txLister, logger)

		dtos, err := useCase.Execute(ctx, accountId)
		require.Error(t, err)
		require.Empty(t, dtos)
	})
}
