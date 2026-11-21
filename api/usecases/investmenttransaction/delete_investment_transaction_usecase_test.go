package investmenttransaction_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/pedrodbflores/personal-finance-tracker/app_errors"
	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/usecases/investmenttransaction"
	"github.com/pedrodbflores/personal-finance-tracker/usecases/usecasesfakes"
	"github.com/stretchr/testify/require"
)

func TestDeleteInvestmentTransactionUseCase(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx := context.TODO()
	accountId := domain.Id(7)
	investmentId := uuid.New()
	txId := uuid.New()

	t.Run("deletes a transaction successfully", func(t *testing.T) {
		deleter := usecasesfakes.FakeInvestmentTransactionDeleter{}
		deleter.DeleteReturns(nil)
		useCase := investmenttransaction.NewDeleteInvestmentTransactionUseCase(&deleter, logger)

		err := useCase.Execute(ctx, accountId, investmentId, txId)
		require.NoError(t, err)

		ctxArg, calledAccountId, calledInvestmentId, calledId := deleter.DeleteArgsForCall(0)
		require.Equal(t, ctx, ctxArg)
		require.Equal(t, accountId, calledAccountId)
		require.Equal(t, investmentId, calledInvestmentId)
		require.Equal(t, txId, calledId)
	})

	t.Run("returns an error if the transaction doesn't exist", func(t *testing.T) {
		deleter := usecasesfakes.FakeInvestmentTransactionDeleter{}
		deleter.DeleteReturns(app_errors.NewResourceNotFoundError("InvestmentTransaction", txId.String()))
		useCase := investmenttransaction.NewDeleteInvestmentTransactionUseCase(&deleter, logger)

		err := useCase.Execute(ctx, accountId, investmentId, txId)
		require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("InvestmentTransaction", txId.String()))
	})
}
