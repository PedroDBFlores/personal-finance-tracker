package investmenttransaction_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/pedrodbflores/personal-finance-tracker/app_errors"
	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/usecases/investmenttransaction"
	"github.com/pedrodbflores/personal-finance-tracker/usecases/usecasesfakes"
	"github.com/stretchr/testify/require"
)

func TestUpdateInvestmentTransactionUseCase(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx := context.TODO()
	investmentId := uuid.New()
	txId := uuid.New()
	description := "Updated description"
	patch := domain.InvestmentTransactionPatch{Description: &description}

	t.Run("updates a transaction successfully", func(t *testing.T) {
		updater := usecasesfakes.FakeInvestmentTransactionUpdater{}
		updater.UpdateReturns(nil)
		useCase := investmenttransaction.NewUpdateInvestmentTransactionUseCase(&updater, logger)

		err := useCase.Execute(ctx, investmentId, txId, patch)
		require.NoError(t, err)

		ctxArg, calledInvestmentId, calledId, calledPatch := updater.UpdateArgsForCall(0)
		require.Equal(t, ctx, ctxArg)
		require.Equal(t, investmentId, calledInvestmentId)
		require.Equal(t, txId, calledId)
		require.Equal(t, patch, calledPatch)
	})

	t.Run("returns an error if the transaction doesn't exist", func(t *testing.T) {
		updater := usecasesfakes.FakeInvestmentTransactionUpdater{}
		updater.UpdateReturns(app_errors.NewResourceNotFoundError("InvestmentTransaction", txId.String()))
		useCase := investmenttransaction.NewUpdateInvestmentTransactionUseCase(&updater, logger)

		err := useCase.Execute(ctx, investmentId, txId, patch)
		require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("InvestmentTransaction", txId.String()))
	})
}
