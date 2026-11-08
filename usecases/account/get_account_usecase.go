package account

import (
	"context"
	"log/slog"

	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/dtos"
	"github.com/pedrodbflores/personal-finance-tracker/usecases"
)

type GetAccountUseCase struct {
	getter usecases.AccountGetter
	logger *slog.Logger
}

func NewGetAccountUseCase(getter usecases.AccountGetter, logger *slog.Logger) GetAccountUseCase {
	return GetAccountUseCase{getter, logger}
}

func (u *GetAccountUseCase) Execute(ctx context.Context, id domain.Id) (*dtos.AccountDto, error) {
	result, err := u.getter.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	dto := dtos.NewAccountDto(*result)
	return &dto, nil
}
