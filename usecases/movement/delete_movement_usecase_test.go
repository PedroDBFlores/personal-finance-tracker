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

func TestDeleteMovementUseCase(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx := context.TODO()
	accountId := domain.NewId()
	movementId := domain.NewId()

	t.Run("deletes a movement", func(t *testing.T) {
		deleter := usecasesfakes.FakeMovementDeleter{}
		deleter.DeleteReturns(nil)
		useCase := movement.NewDeleteMovementUseCase(&deleter, logger)

		err := useCase.Execute(ctx, accountId, movementId)
		require.NoError(t, err)

		ctxArg, calledAccountId, calledId := deleter.DeleteArgsForCall(0)
		require.Equal(t, ctx, ctxArg)
		require.Equal(t, accountId, calledAccountId)
		require.Equal(t, movementId, calledId)
	})

	t.Run("propagates deleter errors", func(t *testing.T) {
		deleter := usecasesfakes.FakeMovementDeleter{}
		deleter.DeleteReturns(app_errors.NewResourceNotFoundError("Movement", movementId.String()))
		useCase := movement.NewDeleteMovementUseCase(&deleter, logger)

		err := useCase.Execute(ctx, accountId, movementId)
		require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("Movement", movementId.String()))
	})
}
