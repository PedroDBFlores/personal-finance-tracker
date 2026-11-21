package investmenttransaction_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pedrodbflores/personal-finance-tracker/app_errors"
	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/usecases/investmenttransaction"
	"github.com/pedrodbflores/personal-finance-tracker/usecases/usecasesfakes"
	"github.com/stretchr/testify/require"
)

func TestListInvestmentTransactionsUseCase(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx := context.TODO()
	accountId := domain.Id(7)
	investmentId := uuid.New()

	t.Run("lists transactions of an investment", func(t *testing.T) {
		lister := usecasesfakes.FakeInvestmentTransactionLister{}
		tx := domain.NewInvestmentTransaction(uuid.New(), investmentId, domain.Buy, 0, domain.Money(1000), domain.Money(10000), time.Now(), "Buy")
		lister.ListReturns([]domain.InvestmentTransaction{tx}, nil)
		useCase := investmenttransaction.NewListInvestmentTransactionsUseCase(&lister, logger)

		dtos, err := useCase.Execute(ctx, accountId, investmentId)
		require.NoError(t, err)
		require.Len(t, dtos, 1)
		require.Equal(t, domain.Buy, dtos[0].Type)

		ctxArg, calledAccountId, calledInvestmentId := lister.ListArgsForCall(0)
		require.Equal(t, ctx, ctxArg)
		require.Equal(t, accountId, calledAccountId)
		require.Equal(t, investmentId, calledInvestmentId)
	})

	t.Run("returns an error if listing fails", func(t *testing.T) {
		lister := usecasesfakes.FakeInvestmentTransactionLister{}
		lister.ListReturns(nil, app_errors.NewResourceNotFoundError("InvestmentTransaction", "1"))
		useCase := investmenttransaction.NewListInvestmentTransactionsUseCase(&lister, logger)

		dtos, err := useCase.Execute(ctx, accountId, investmentId)
		require.Error(t, err)
		require.Empty(t, dtos)
	})
}
