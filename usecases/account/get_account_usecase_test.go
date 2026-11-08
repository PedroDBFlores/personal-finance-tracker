package account_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/pedrodbflores/personal-finance-tracker/app_errors"
	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/dtos"
	"github.com/pedrodbflores/personal-finance-tracker/usecases/account"
	"github.com/pedrodbflores/personal-finance-tracker/usecases/usecasesfakes"
	"github.com/stretchr/testify/require"
)

func TestGetAccountUseCase(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx := context.TODO()

	t.Run("gets an account by it's id", func(t *testing.T) {
		accountGetter := usecasesfakes.FakeAccountGetter{}
		acc := domain.NewAccount(domain.NewId(), "Bank Bank", domain.AccountTypeBank, domain.CurrencyEUR, time.Now())
		dto := dtos.NewAccountDto(acc)
		accountGetter.GetReturns(&acc, nil)
		useCase := account.NewGetAccountUseCase(&accountGetter, logger)

		result, err := useCase.Execute(ctx, acc.Id)

		require.NoError(t, err)
		require.Equal(t, *result, dto)

		calledCtx, calledId := accountGetter.GetArgsForCall(0)
		require.Equal(t, ctx, calledCtx)
		require.Equal(t, acc.Id, calledId)
	})

	t.Run("returns an error if getting the account fails", func(t *testing.T) {
		accountGetter := usecasesfakes.FakeAccountGetter{}
		accountGetter.GetReturns(nil, app_errors.NewResourceNotFoundError("Account", "123"))
		useCase := account.NewGetAccountUseCase(&accountGetter, logger)

		result, err := useCase.Execute(ctx, domain.NewId())
		require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("Account", "123"))
		require.Nil(t, result)
	})
}
