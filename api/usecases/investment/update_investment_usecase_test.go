package investment_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/pedrodbflores/personal-finance-tracker/app_errors"
	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/usecases/investment"
	"github.com/pedrodbflores/personal-finance-tracker/usecases/usecasesfakes"
	"github.com/stretchr/testify/require"
)

func TestUpdateInvestmentUseCase(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx := context.TODO()
	accountId := domain.NewId()
	investmentId := domain.NewId()
	name := "Renamed ETF"
	patch := domain.InvestmentPatch{Name: &name}

	t.Run("updates an investment successfully", func(t *testing.T) {
		updater := usecasesfakes.FakeInvestmentUpdater{}
		updater.UpdateReturns(nil)
		useCase := investment.NewUpdateInvestmentUseCase(&updater, logger)

		err := useCase.Execute(ctx, accountId, investmentId, patch)
		require.NoError(t, err)

		ctxArg, calledAccountId, calledId, calledPatch := updater.UpdateArgsForCall(0)
		require.Equal(t, ctx, ctxArg)
		require.Equal(t, accountId, calledAccountId)
		require.Equal(t, investmentId, calledId)
		require.Equal(t, patch, calledPatch)
	})

	t.Run("returns an error if the investment doesn't exist", func(t *testing.T) {
		updater := usecasesfakes.FakeInvestmentUpdater{}
		updater.UpdateReturns(app_errors.NewResourceNotFoundError("Investment", investmentId.String()))
		useCase := investment.NewUpdateInvestmentUseCase(&updater, logger)

		err := useCase.Execute(ctx, accountId, investmentId, patch)
		require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("Investment", investmentId.String()))
	})
}
