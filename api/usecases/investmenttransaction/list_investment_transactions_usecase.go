package investmenttransaction

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/dtos"
	"github.com/pedrodbflores/personal-finance-tracker/usecases"
)

type ListInvestmentTransactionsUseCase struct {
	lister usecases.InvestmentTransactionLister
	logger *slog.Logger
}

func NewListInvestmentTransactionsUseCase(lister usecases.InvestmentTransactionLister, logger *slog.Logger) ListInvestmentTransactionsUseCase {
	return ListInvestmentTransactionsUseCase{lister, logger}
}

func (u *ListInvestmentTransactionsUseCase) Execute(ctx context.Context, accountId domain.Id, investmentId uuid.UUID) ([]dtos.InvestmentTransactionDto, error) {
	result, err := u.lister.List(ctx, accountId, investmentId)
	if err != nil {
		return []dtos.InvestmentTransactionDto{}, err
	}
	listDtos := make([]dtos.InvestmentTransactionDto, 0, len(result))
	for _, tx := range result {
		listDtos = append(listDtos, dtos.NewInvestmentTransactionDto(tx))
	}
	return listDtos, nil
}
