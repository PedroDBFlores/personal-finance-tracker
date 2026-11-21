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

func TestListMovementMonthsUseCase(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx := context.TODO()
	now := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)

	t.Run("returns the year/month pairs of the account's movements", func(t *testing.T) {
		accountGetter := &usecasesfakes.FakeAccountGetter{}
		monthLister := &usecasesfakes.FakeMovementMonthLister{}

		acc := domain.NewAccount(1, "Bank A", domain.AccountTypeBank, domain.CurrencyEUR, now)
		accountGetter.GetReturns(&acc, nil)
		months := []domain.YearMonth{{Year: 2025, Month: 2}, {Year: 2025, Month: 1}}
		monthLister.ListMonthsReturns(months, nil)

		useCase := movement.NewListMovementMonthsUseCase(accountGetter, monthLister, logger)

		result, err := useCase.Execute(ctx, 1)

		require.NoError(t, err)
		require.Len(t, result, 2)
		require.Equal(t, dtos.NewMovementMonthDto(domain.YearMonth{Year: 2025, Month: 2}), result[0])
		require.Equal(t, dtos.NewMovementMonthDto(domain.YearMonth{Year: 2025, Month: 1}), result[1])

		_, calledId := accountGetter.GetArgsForCall(0)
		require.Equal(t, domain.Id(1), calledId)
		_, listAccId := monthLister.ListMonthsArgsForCall(0)
		require.Equal(t, domain.Id(1), listAccId)
	})

	t.Run("returns an empty list when the account has no movements", func(t *testing.T) {
		accountGetter := &usecasesfakes.FakeAccountGetter{}
		monthLister := &usecasesfakes.FakeMovementMonthLister{}

		acc := domain.NewAccount(1, "Bank A", domain.AccountTypeBank, domain.CurrencyEUR, now)
		accountGetter.GetReturns(&acc, nil)
		monthLister.ListMonthsReturns([]domain.YearMonth{}, nil)

		useCase := movement.NewListMovementMonthsUseCase(accountGetter, monthLister, logger)

		result, err := useCase.Execute(ctx, 1)

		require.NoError(t, err)
		require.Empty(t, result)
	})

	t.Run("returns an error if the account does not exist", func(t *testing.T) {
		accountGetter := &usecasesfakes.FakeAccountGetter{}
		monthLister := &usecasesfakes.FakeMovementMonthLister{}

		accountGetter.GetReturns(nil, app_errors.NewResourceNotFoundError("Account", "9"))

		useCase := movement.NewListMovementMonthsUseCase(accountGetter, monthLister, logger)

		result, err := useCase.Execute(ctx, 9)
		require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("Account", "9"))
		require.Nil(t, result)
		require.Equal(t, 0, monthLister.ListMonthsCallCount())
	})

	t.Run("propagates lister errors", func(t *testing.T) {
		accountGetter := &usecasesfakes.FakeAccountGetter{}
		monthLister := &usecasesfakes.FakeMovementMonthLister{}

		acc := domain.NewAccount(1, "Bank A", domain.AccountTypeBank, domain.CurrencyEUR, now)
		accountGetter.GetReturns(&acc, nil)
		monthLister.ListMonthsReturns(nil, errors.New("db down"))

		useCase := movement.NewListMovementMonthsUseCase(accountGetter, monthLister, logger)

		result, err := useCase.Execute(ctx, 1)
		require.EqualError(t, err, "db down")
		require.Nil(t, result)
	})
}
