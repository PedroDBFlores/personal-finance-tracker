package account_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/pedrodbflores/personal-finance-tracker/app_errors"
	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/usecases/account"
	"github.com/pedrodbflores/personal-finance-tracker/usecases/usecasesfakes"
	"github.com/stretchr/testify/require"
)

func TestUpdateAccountUseCase(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx := context.TODO()
	update := domain.AccountPatch{
		Name: ptr("New Name"),
		Type: ptr(domain.AccountTypeBrokerage),
	}

	t.Run("forwards the update to the updater", func(t *testing.T) {
		accountUpdater := usecasesfakes.FakeAccountUpdater{}
		accountUpdater.UpdateReturns(nil)
		useCase := account.NewUpdateAccountUseCase(&accountUpdater, logger)

		err := useCase.Execute(ctx, 3, update)
		require.NoError(t, err)

		_, _, patchArg := accountUpdater.UpdateArgsForCall(0)
		require.Equal(t, update.Name, patchArg.Name)
		require.Equal(t, update.Type, patchArg.Type)
	})

	t.Run("propagates updater errors", func(t *testing.T) {
		accountUpdater := usecasesfakes.FakeAccountUpdater{}
		accountUpdater.UpdateReturns(app_errors.NewResourceNotFoundError("Account", "3"))
		useCase := account.NewUpdateAccountUseCase(&accountUpdater, logger)

		err := useCase.Execute(ctx, 3, update)
		require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("Account", "3"))
	})
}

func ptr[T any](v T) *T { return &v }
