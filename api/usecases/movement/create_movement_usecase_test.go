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

func TestCreateMovementUseCase(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx := context.TODO()
	accountId := domain.NewId()
	dto := dtos.NewCreateMovementDto("Salary", domain.Money(100000), domain.MovementTypeCredit, time.Now())

	t.Run("creates a movement successfully", func(t *testing.T) {
		creator := usecasesfakes.FakeMovementCreator{}
		randomId := domain.NewId()
		creator.CreateReturns(&randomId, nil)
		useCase := movement.NewCreateMovementUseCase(&creator, logger)

		id, err := useCase.Execute(ctx, accountId, dto)
		require.NoError(t, err)
		require.Equal(t, randomId, *id)

		ctxArg, calledAccountId, modelMovement := creator.CreateArgsForCall(0)
		require.Equal(t, ctx, ctxArg)
		require.Equal(t, accountId, calledAccountId)
		require.Equal(t, domain.Movement{
			AccountId:   accountId,
			Description: dto.Description,
			Value:       dto.Value,
			Type:        dto.Type,
			Date:        dto.Date,
		}, modelMovement)
	})

	t.Run("returns an error if it fails to create a movement", func(t *testing.T) {
		creator := usecasesfakes.FakeMovementCreator{}
		creator.CreateReturns(nil, app_errors.NewResourceNotFoundError("Movement", "123"))
		useCase := movement.NewCreateMovementUseCase(&creator, logger)

		id, err := useCase.Execute(ctx, accountId, dto)
		require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("Movement", "123"))
		require.Nil(t, id)
	})
}
