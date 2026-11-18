package movement_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/pedrodbflores/personal-finance-tracker/app_errors"
	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/dtos"
	"github.com/pedrodbflores/personal-finance-tracker/usecases/movement"
	"github.com/pedrodbflores/personal-finance-tracker/usecases/usecasesfakes"
	"github.com/stretchr/testify/require"
)

func TestGetMovementUseCase(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx := context.TODO()
	accountId := domain.NewId()

	t.Run("gets a movement by id", func(t *testing.T) {
		getter := usecasesfakes.FakeMovementGetter{}
		mov := domain.NewMovement(accountId, "Groceries", domain.Money(2599), domain.MovementTypeDebit, time.Now())
		mov.Id = domain.NewId()
		getter.GetReturns(&mov, nil)
		useCase := movement.NewGetMovementUseCase(&getter, logger)

		result, err := useCase.Execute(ctx, accountId, mov.Id)
		require.NoError(t, err)
		require.Equal(t, dtos.NewMovementDto(mov), *result)

		ctxArg, calledAccountId, calledId := getter.GetArgsForCall(0)
		require.Equal(t, ctx, ctxArg)
		require.Equal(t, accountId, calledAccountId)
		require.Equal(t, mov.Id, calledId)
	})

	t.Run("returns an error if getting the movement fails", func(t *testing.T) {
		getter := usecasesfakes.FakeMovementGetter{}
		getter.GetReturns(nil, app_errors.NewResourceNotFoundError("Movement", "123"))
		useCase := movement.NewGetMovementUseCase(&getter, logger)

		result, err := useCase.Execute(ctx, accountId, domain.NewId())
		require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("Movement", "123"))
		require.Nil(t, result)
	})
}
