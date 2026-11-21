package investmenttransaction

import (
	"context"
	"log/slog"

	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/dtos"
	"github.com/pedrodbflores/personal-finance-tracker/usecases"
)

type GetInvestmentTransactionUseCase struct {
	getter usecases.InvestmentTransactionGetter
	logger *slog.Logger
}

func NewGetInvestmentTransactionUseCase(getter usecases.InvestmentTransactionGetter, logger *slog.Logger) GetInvestmentTransactionUseCase {
	return GetInvestmentTransactionUseCase{getter, logger}
}

func (u *GetInvestmentTransactionUseCase) Execute(ctx context.Context, accountId domain.Id, investmentId, id domain.Id) (*dtos.InvestmentTransactionDto, error) {
	result, err := u.getter.Get(ctx, accountId, investmentId, id)
	if err != nil {
		return nil, err
	}
	dto := dtos.NewInvestmentTransactionDto(*result)
	return &dto, nil
}
