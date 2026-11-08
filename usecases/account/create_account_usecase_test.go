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

func TestCreateAccountUseCase(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx := context.TODO()
	dto := dtos.NewCreateAccountDto(
		"Bank A",
		domain.AccountTypeBank,
		domain.CurrencyEUR,
		time.Now(),
	)

	t.Run("create an account successfully", func(t *testing.T) {
		accountCreator := usecasesfakes.FakeAccountCreator{}
		randomId := domain.NewId()
		accountCreator.CreateReturns(&randomId, nil)
		useCase := account.NewCreateAccountUseCase(&accountCreator, logger)

		id, err := useCase.Execute(ctx, dto)
		require.NoError(t, err)
		require.NotNil(t, id)
		require.Equal(t, *id, randomId)

		ctxArg, modelAccount := accountCreator.CreateArgsForCall(0)
		require.Equal(t, ctx, ctxArg)
		require.Equal(t, domain.Account{
			Name:      dto.Name,
			Type:      domain.AccountType(dto.Type),
			Currency:  dto.Currency,
			CreatedAt: dto.CreatedAt,
		}, modelAccount)
	})

	t.Run("it retuns an error if it fails to create an account", func(t *testing.T) {
		accountCreator := usecasesfakes.FakeAccountCreator{}
		accountCreator.CreateReturns(nil, app_errors.NewResourceNotFoundError("Account", "123"))
		useCase := account.NewCreateAccountUseCase(&accountCreator, logger)

		id, err := useCase.Execute(ctx, dto)
		require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("Account", "123"))
		require.Nil(t, id)
	})
}
