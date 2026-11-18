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

func TestUpdateAccountUseCase(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx := context.TODO()
	dto := dtos.NewUpdateAccountDto(
		domain.NewId(),
		"Bank A",
		domain.AccountTypeBank,
		domain.CurrencyEUR,
		time.Now(),
	)

	t.Run("updates an account successfully", func(t *testing.T) {
		accountUpdater := usecasesfakes.FakeAccountUpdater{}
		accountUpdater.UpdateReturns(nil)
		useCase := account.NewUpdateAccountUseCase(&accountUpdater, logger)

		err := useCase.Execute(ctx, dto)
		require.NoError(t, err)

		ctxArg, modelAccount := accountUpdater.UpdateArgsForCall(0)
		require.Equal(t, ctx, ctxArg)
		require.Equal(t, domain.Account{
			Id:        dto.Id,
			Name:      dto.Name,
			Type:      domain.AccountType(dto.Type),
			Currency:  dto.Currency,
			CreatedAt: dto.CreatedAt,
		}, modelAccount)
	})

	t.Run("it returns an error if it fails to update an account", func(t *testing.T) {
		accountUpdater := usecasesfakes.FakeAccountUpdater{}
		accountUpdater.UpdateReturns(app_errors.NewResourceNotFoundError("Account", dto.Id.String()))
		useCase := account.NewUpdateAccountUseCase(&accountUpdater, logger)

		err := useCase.Execute(ctx, dto)
		require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("Account", dto.Id.String()))
	})
}
