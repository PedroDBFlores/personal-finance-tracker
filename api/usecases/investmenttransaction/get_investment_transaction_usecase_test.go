package investmenttransaction_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/pedrodbflores/personal-finance-tracker/app_errors"
	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/usecases/investmenttransaction"
	"github.com/pedrodbflores/personal-finance-tracker/usecases/usecasesfakes"
	"github.com/stretchr/testify/require"
)

func TestGetInvestmentTransactionUseCase(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx := context.TODO()
	accountId := domain.Id(7)
	investmentId := domain.NewId()
	txId := domain.NewId()

	t.Run("returns the transaction", func(t *testing.T) {
		getter := usecasesfakes.FakeInvestmentTransactionGetter{}
		tx := domain.NewInvestmentTransaction(txId, investmentId, domain.Deposit, domain.Money(50000), 0, 0, time.Now(), "Deposit")
		getter.GetReturns(&tx, nil)
		useCase := investmenttransaction.NewGetInvestmentTransactionUseCase(&getter, logger)

		dto, err := useCase.Execute(ctx, accountId, investmentId, txId)
		require.NoError(t, err)
		require.Equal(t, txId, dto.Id)
		require.Equal(t, domain.Money(50000), dto.Amount)

		ctxArg, calledAccountId, calledInvestmentId, calledId := getter.GetArgsForCall(0)
		require.Equal(t, ctx, ctxArg)
		require.Equal(t, accountId, calledAccountId)
		require.Equal(t, investmentId, calledInvestmentId)
		require.Equal(t, txId, calledId)
	})

	t.Run("returns an error if the transaction doesn't exist", func(t *testing.T) {
		getter := usecasesfakes.FakeInvestmentTransactionGetter{}
		getter.GetReturns(nil, app_errors.NewResourceNotFoundError("InvestmentTransaction", txId.String()))
		useCase := investmenttransaction.NewGetInvestmentTransactionUseCase(&getter, logger)

		dto, err := useCase.Execute(ctx, accountId, investmentId, txId)
		require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("InvestmentTransaction", txId.String()))
		require.Nil(t, dto)
	})
}
