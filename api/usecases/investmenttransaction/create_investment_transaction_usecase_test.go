package investmenttransaction_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/pedrodbflores/personal-finance-tracker/app_errors"
	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/dtos"
	"github.com/pedrodbflores/personal-finance-tracker/usecases/investmenttransaction"
	"github.com/pedrodbflores/personal-finance-tracker/usecases/usecasesfakes"
	"github.com/stretchr/testify/require"
)

func TestCreateInvestmentTransactionUseCase(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx := context.TODO()
	accountId := domain.Id(7)
	investmentId := domain.NewId()
	dto := dtos.NewCreateInvestmentTransactionDto(domain.Buy, 0, domain.Money(1000), domain.Money(10000), time.Now(), "Initial Buy")

	t.Run("creates a transaction successfully", func(t *testing.T) {
		creator := usecasesfakes.FakeInvestmentTransactionCreator{}
		randomId := domain.NewId()
		creator.CreateReturns(&randomId, nil)
		useCase := investmenttransaction.NewCreateInvestmentTransactionUseCase(&creator, logger)

		id, err := useCase.Execute(ctx, accountId, investmentId, dto)
		require.NoError(t, err)
		require.Equal(t, randomId, *id)

		ctxArg, calledAccountId, calledInvestmentId, modelTx := creator.CreateArgsForCall(0)
		require.Equal(t, ctx, ctxArg)
		require.Equal(t, accountId, calledAccountId)
		require.Equal(t, investmentId, calledInvestmentId)
		require.Equal(t, investmentId, modelTx.InvestmentId)
		require.Equal(t, domain.Id(0), modelTx.Id) // id assigned by the repository
		require.Equal(t, dto.Type, modelTx.Type)
		require.Equal(t, dto.Amount, modelTx.Amount)
		require.Equal(t, dto.Quantity, modelTx.Quantity)
		require.Equal(t, dto.Price, modelTx.Price)
		require.Equal(t, dto.Description, modelTx.Description)
	})

	t.Run("returns an error if it fails to create a transaction", func(t *testing.T) {
		creator := usecasesfakes.FakeInvestmentTransactionCreator{}
		creator.CreateReturns(nil, app_errors.NewResourceNotFoundError("InvestmentTransaction", "123"))
		useCase := investmenttransaction.NewCreateInvestmentTransactionUseCase(&creator, logger)

		id, err := useCase.Execute(ctx, accountId, investmentId, dto)
		require.ErrorIs(t, err, app_errors.NewResourceNotFoundError("InvestmentTransaction", "123"))
		require.Nil(t, id)
	})
}
