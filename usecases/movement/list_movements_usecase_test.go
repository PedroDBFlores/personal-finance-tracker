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

func TestListMovementsUseCase(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx := context.TODO()
	accountId := domain.NewId()

	t.Run("lists movements of an account", func(t *testing.T) {
		lister := usecasesfakes.FakeMovementLister{}
		movements := []domain.Movement{
			func() domain.Movement { m := domain.NewMovement(accountId, "Salary", domain.Money(100000), domain.MovementTypeCredit, time.Now()); m.Id = 1; return m }(),
			func() domain.Movement { m := domain.NewMovement(accountId, "Groceries", domain.Money(2599), domain.MovementTypeDebit, time.Now()); m.Id = 2; return m }(),
		}
		lister.ListReturns(movements, nil)
		useCase := movement.NewListMovementsUseCase(&lister, logger)

		result, err := useCase.Execute(ctx, accountId)
		require.NoError(t, err)
		require.Len(t, result, 2)
		require.Equal(t, dtos.NewMovementDto(movements[0]), result[0])
		require.Equal(t, dtos.NewMovementDto(movements[1]), result[1])

		ctxArg, calledAccountId := lister.ListArgsForCall(0)
		require.Equal(t, ctx, ctxArg)
		require.Equal(t, accountId, calledAccountId)
	})

	t.Run("returns an error if listing fails", func(t *testing.T) {
		lister := usecasesfakes.FakeMovementLister{}
		lister.ListReturns(nil, app_errors.NewResourceNotFoundError("Account", accountId.String()))
		useCase := movement.NewListMovementsUseCase(&lister, logger)

		result, err := useCase.Execute(ctx, accountId)
		require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("Account", accountId.String()))
		require.Empty(t, result)
	})
}
