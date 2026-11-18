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

func TestDeleteAccountUseCase(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx := context.TODO()
	id := domain.NewId()

	t.Run("deletes an account successfully", func(t *testing.T) {
		accountDeleter := usecasesfakes.FakeAccountDeleter{}
		accountDeleter.DeleteReturns(nil)
		useCase := account.NewDeleteAccountUseCase(&accountDeleter, logger)

		err := useCase.Execute(ctx, id)
		require.NoError(t, err)

		ctxArg, calledId := accountDeleter.DeleteArgsForCall(0)
		require.Equal(t, ctx, ctxArg)
		require.Equal(t, id, calledId)
	})

	t.Run("it returns an error if it fails to delete an account", func(t *testing.T) {
		accountDeleter := usecasesfakes.FakeAccountDeleter{}
		accountDeleter.DeleteReturns(app_errors.NewResourceNotFoundError("Account", id.String()))
		useCase := account.NewDeleteAccountUseCase(&accountDeleter, logger)

		err := useCase.Execute(ctx, id)
		require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("Account", id.String()))
	})
}
