package investment_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/pedrodbflores/personal-finance-tracker/app_errors"
	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/usecases/investment"
	"github.com/pedrodbflores/personal-finance-tracker/usecases/usecasesfakes"
	"github.com/stretchr/testify/require"
)

func TestDeleteInvestmentUseCase(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx := context.TODO()
	accountId := domain.NewId()
	investmentId := uuid.New()

	t.Run("deletes an investment successfully", func(t *testing.T) {
		deleter := usecasesfakes.FakeInvestmentDeleter{}
		deleter.DeleteReturns(nil)
		useCase := investment.NewDeleteInvestmentUseCase(&deleter, logger)

		err := useCase.Execute(ctx, accountId, investmentId)
		require.NoError(t, err)

		ctxArg, calledAccountId, calledId := deleter.DeleteArgsForCall(0)
		require.Equal(t, ctx, ctxArg)
		require.Equal(t, accountId, calledAccountId)
		require.Equal(t, investmentId, calledId)
	})

	t.Run("returns an error if the investment doesn't exist", func(t *testing.T) {
		deleter := usecasesfakes.FakeInvestmentDeleter{}
		deleter.DeleteReturns(app_errors.NewResourceNotFoundError("Investment", investmentId.String()))
		useCase := investment.NewDeleteInvestmentUseCase(&deleter, logger)

		err := useCase.Execute(ctx, accountId, investmentId)
		require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("Investment", investmentId.String()))
	})
}
