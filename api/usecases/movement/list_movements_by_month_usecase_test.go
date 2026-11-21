package movement_test

import (
	"context"
	"errors"
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

func TestListMovementsByMonthUseCase(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx := context.TODO()
	now := time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)

	t.Run("returns the movements of the given year and month", func(t *testing.T) {
		accountGetter := &usecasesfakes.FakeAccountGetter{}
		byMonthLister := &usecasesfakes.FakeMovementByMonthLister{}

		acc, err := domain.NewAccount(1, "Bank A", domain.AccountTypeBank, domain.CurrencyEUR, now); require.NoError(t, err)
		accountGetter.GetReturns(&acc, nil)
		mov, err := domain.NewMovement(1, "Salary", domain.Money(100000), domain.MovementTypeCredit, now); require.NoError(t, err)
		mov.Id = 7
		byMonthLister.ListByMonthReturns([]domain.Movement{mov}, nil)

		useCase := movement.NewListMovementsByMonthUseCase(accountGetter, byMonthLister, logger)

		result, err := useCase.Execute(ctx, 1, 2025, 1)

		require.NoError(t, err)
		require.Len(t, result, 1)
		require.Equal(t, dtos.NewMovementDto(mov), result[0])

		_, calledAccId, year, month := byMonthLister.ListByMonthArgsForCall(0)
		require.Equal(t, domain.Id(1), calledAccId)
		require.Equal(t, 2025, year)
		require.Equal(t, 1, month)
	})

	t.Run("returns an error if the account does not exist", func(t *testing.T) {
		accountGetter := &usecasesfakes.FakeAccountGetter{}
		byMonthLister := &usecasesfakes.FakeMovementByMonthLister{}

		accountGetter.GetReturns(nil, app_errors.NewResourceNotFoundError("Account", "9"))

		useCase := movement.NewListMovementsByMonthUseCase(accountGetter, byMonthLister, logger)

		result, err := useCase.Execute(ctx, 9, 2025, 1)
		require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("Account", "9"))
		require.Nil(t, result)
		require.Equal(t, 0, byMonthLister.ListByMonthCallCount())
	})

	t.Run("propagates lister errors", func(t *testing.T) {
		accountGetter := &usecasesfakes.FakeAccountGetter{}
		byMonthLister := &usecasesfakes.FakeMovementByMonthLister{}

		acc, err := domain.NewAccount(1, "Bank A", domain.AccountTypeBank, domain.CurrencyEUR, now); require.NoError(t, err)
		accountGetter.GetReturns(&acc, nil)
		byMonthLister.ListByMonthReturns(nil, errors.New("db down"))

		useCase := movement.NewListMovementsByMonthUseCase(accountGetter, byMonthLister, logger)

		result, err := useCase.Execute(ctx, 1, 2025, 1)
		require.EqualError(t, err, "db down")
		require.Nil(t, result)
	})
}
