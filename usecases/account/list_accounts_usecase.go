package account

import (
	"context"
	"log/slog"

	"github.com/pedrodbflores/personal-finance-tracker/dtos"
	"github.com/pedrodbflores/personal-finance-tracker/usecases"
)

type ListAccountsUseCase struct {
	lister usecases.AccountLister
	logger *slog.Logger
}

func NewListAccountsUseCase(lister usecases.AccountLister, logger *slog.Logger) ListAccountsUseCase {
	return ListAccountsUseCase{lister, logger}
}

func (u *ListAccountsUseCase) Execute(ctx context.Context) ([]dtos.AccountDto, error) {
	result, err := u.lister.List(ctx)
	var listAccountDtos []dtos.AccountDto
	if err != nil {
		return listAccountDtos, err
	}
	for _, account := range result {
		listAccountDtos = append(listAccountDtos, dtos.NewAccountDto(account))
	}
	return listAccountDtos, nil
}
