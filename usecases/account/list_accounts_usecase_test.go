package account_test

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/dtos"
	"github.com/pedrodbflores/personal-finance-tracker/usecases/account"
	"github.com/pedrodbflores/personal-finance-tracker/usecases/usecasesfakes"
	"github.com/stretchr/testify/require"
)

func TestListAccountsUseCase(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx := context.TODO()

	t.Run("returns a list of accounts", func(t *testing.T) {
		accountLister := usecasesfakes.FakeAccountLister{}
		listAccounts := []domain.Account{
			domain.NewAccount(domain.NewId(), "Bank A", domain.AccountTypeBank, domain.CurrencyEUR, time.Now()),
			domain.NewAccount(domain.NewId(), "Credit Card B", domain.AccountTypeCreditCard, domain.CurrencyEUR, time.Now()),
		}
		var listAccountDtos []dtos.AccountDto
		for _, account := range listAccounts {
			listAccountDtos = append(listAccountDtos, dtos.NewAccountDto(account))
		}
		accountLister.ListReturns(listAccounts, nil)
		useCase := account.NewListAccountsUseCase(&accountLister, logger)

		result, err := useCase.Execute(ctx)

		require.NoError(t, err)
		require.Equal(t, listAccountDtos, result)

		calledCtx := accountLister.ListArgsForCall(0)
		require.Equal(t, ctx, calledCtx)
	})

	t.Run("returns an error if it fails to list accounts", func(t *testing.T) {
		accountLister := usecasesfakes.FakeAccountLister{}
		accountLister.ListReturns([]domain.Account{}, errors.New("oopsie"))
		useCase := account.NewListAccountsUseCase(&accountLister, logger)

		result, err := useCase.Execute(ctx)

		require.Error(t, err, "oopsie")
		require.Zero(t, len(result))
	})
}
