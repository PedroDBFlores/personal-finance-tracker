package movement_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/pedrodbflores/personal-finance-tracker/app_errors"
	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/usecases/movement"
	"github.com/pedrodbflores/personal-finance-tracker/usecases/usecasesfakes"
	"github.com/stretchr/testify/require"
)

func TestUpdateMovementUseCase(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx := context.TODO()
	accountId := domain.NewId()
	movementId := domain.NewId()
	newValue := domain.Money(5000)
	update := domain.MovementPatch{Value: &newValue}

	t.Run("forwards the patch to the updater", func(t *testing.T) {
		updater := usecasesfakes.FakeMovementUpdater{}
		updater.UpdateReturns(nil)
		useCase := movement.NewUpdateMovementUseCase(&updater, logger)

		err := useCase.Execute(ctx, accountId, movementId, update)
		require.NoError(t, err)

		ctxArg, calledAccountId, calledId, patchArg := updater.UpdateArgsForCall(0)
		require.Equal(t, ctx, ctxArg)
		require.Equal(t, accountId, calledAccountId)
		require.Equal(t, movementId, calledId)
		require.Equal(t, update, patchArg)
	})

	t.Run("propagates updater errors", func(t *testing.T) {
		updater := usecasesfakes.FakeMovementUpdater{}
		updater.UpdateReturns(app_errors.NewResourceNotFoundError("Movement", movementId.String()))
		useCase := movement.NewUpdateMovementUseCase(&updater, logger)

		err := useCase.Execute(ctx, accountId, movementId, update)
		require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("Movement", movementId.String()))
	})
}
