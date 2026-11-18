package investment_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pedrodbflores/personal-finance-tracker/app_errors"
	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/dtos"
	"github.com/pedrodbflores/personal-finance-tracker/usecases/investment"
	"github.com/pedrodbflores/personal-finance-tracker/usecases/usecasesfakes"
	"github.com/stretchr/testify/require"
)

func TestCreateInvestmentUseCase(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx := context.TODO()
	accountId := domain.NewId()
	dto := dtos.NewCreateInvestmentDto("Vanguard ETF", domain.ETF, domain.CurrencyUSD, time.Now(), nil)

	t.Run("creates an investment successfully", func(t *testing.T) {
		creator := usecasesfakes.FakeInvestmentCreator{}
		randomId := uuid.New()
		creator.CreateReturns(&randomId, nil)
		useCase := investment.NewCreateInvestmentUseCase(&creator, logger)

		id, err := useCase.Execute(ctx, accountId, dto)
		require.NoError(t, err)
		require.Equal(t, randomId, *id)

		ctxArg, calledAccountId, modelInvestment := creator.CreateArgsForCall(0)
		require.Equal(t, ctx, ctxArg)
		require.Equal(t, accountId, calledAccountId)
		require.Equal(t, accountId, modelInvestment.AccountId)
		require.NotEqual(t, uuid.Nil, modelInvestment.Id)
		require.Equal(t, dto.Name, modelInvestment.Name)
		require.Equal(t, dto.Type, modelInvestment.Type)
		require.Equal(t, dto.Currency, modelInvestment.Currency)
		require.Equal(t, domain.Money(0), modelInvestment.CurrentValue)
	})

	t.Run("returns an error if it fails to create an investment", func(t *testing.T) {
		creator := usecasesfakes.FakeInvestmentCreator{}
		creator.CreateReturns(nil, app_errors.NewResourceNotFoundError("Investment", "123"))
		useCase := investment.NewCreateInvestmentUseCase(&creator, logger)

		id, err := useCase.Execute(ctx, accountId, dto)
		require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("Investment", "123"))
		require.Nil(t, id)
	})
}
