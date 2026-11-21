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

func TestGetInvestmentUseCase(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx := context.TODO()
	accountId := domain.NewId()
	investmentId := domain.NewId()

	sampleInvestment := func() domain.Investment {
		inv, err := domain.NewInvestment(investmentId, accountId, "Vanguard ETF", domain.ETF, domain.CurrencyUSD, time.Now(), nil); require.NoError(t, err)
		return inv
	}

	t.Run("returns the investment with current value calculated from transactions", func(t *testing.T) {
		getter := usecasesfakes.FakeInvestmentGetter{}
		txLister := usecasesfakes.FakeInvestmentTransactionLister{}
		inv := sampleInvestment()
		getter.GetReturns(&inv, nil)
		buy, err := domain.NewInvestmentTransaction(domain.NewId(), investmentId, domain.Buy, 0, domain.Money(1000), domain.Money(10000), time.Now(), "Buy"); require.NoError(t, err)
		txLister.ListReturns([]domain.InvestmentTransaction{buy}, nil)
		useCase := investment.NewGetInvestmentUseCase(&getter, &txLister, logger)

		dto, err := useCase.Execute(ctx, accountId, investmentId)
		require.NoError(t, err)
		require.Equal(t, investmentId, dto.Id)
		require.Equal(t, domain.Money(100000), dto.CurrentValue)

		ctxArg, calledAccountId, calledId := getter.GetArgsForCall(0)
		require.Equal(t, ctx, ctxArg)
		require.Equal(t, accountId, calledAccountId)
		require.Equal(t, investmentId, calledId)
		_, calledListAccountId, calledInvestmentId := txLister.ListArgsForCall(0)
		require.Equal(t, accountId, calledListAccountId)
		require.Equal(t, investmentId, calledInvestmentId)
	})

	t.Run("returns an error if the investment doesn't exist", func(t *testing.T) {
		getter := usecasesfakes.FakeInvestmentGetter{}
		txLister := usecasesfakes.FakeInvestmentTransactionLister{}
		getter.GetReturns(nil, app_errors.NewResourceNotFoundError("Investment", investmentId.String()))
		useCase := investment.NewGetInvestmentUseCase(&getter, &txLister, logger)

		dto, err := useCase.Execute(ctx, accountId, investmentId)
		require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("Investment", investmentId.String()))
		require.Nil(t, dto)
	})

	t.Run("returns an error if listing transactions fails", func(t *testing.T) {
		getter := usecasesfakes.FakeInvestmentGetter{}
		txLister := usecasesfakes.FakeInvestmentTransactionLister{}
		inv := sampleInvestment()
		getter.GetReturns(&inv, nil)
		txLister.ListReturns(nil, app_errors.NewResourceNotFoundError("InvestmentTransaction", "x"))
		useCase := investment.NewGetInvestmentUseCase(&getter, &txLister, logger)

		dto, err := useCase.Execute(ctx, accountId, investmentId)
		require.Error(t, err)
		require.Nil(t, dto)
	})
}
