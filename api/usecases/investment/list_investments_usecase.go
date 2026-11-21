package investment

import (
	"context"
	"log/slog"

	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/dtos"
	"github.com/pedrodbflores/personal-finance-tracker/usecases"
)

type ListInvestmentsUseCase struct {
	lister   usecases.InvestmentLister
	txLister usecases.InvestmentTransactionLister
	logger   *slog.Logger
}

func NewListInvestmentsUseCase(lister usecases.InvestmentLister, txLister usecases.InvestmentTransactionLister, logger *slog.Logger) ListInvestmentsUseCase {
	return ListInvestmentsUseCase{lister, txLister, logger}
}

func (u *ListInvestmentsUseCase) Execute(ctx context.Context, accountId domain.Id) ([]dtos.InvestmentDto, error) {
	investments, err := u.lister.List(ctx, accountId)
	if err != nil {
		return []dtos.InvestmentDto{}, err
	}
	listInvestmentDtos := make([]dtos.InvestmentDto, 0, len(investments))
	for i := range investments {
		u.logger.DebugContext(ctx, "listing transactions to calculate current value", slog.Any("investmentId", investments[i].Id))
		transactions, err := u.txLister.List(ctx, accountId, investments[i].Id)
		if err != nil {
			return []dtos.InvestmentDto{}, err
		}
		investments[i].UpdateCurrentValue(transactions)
		listInvestmentDtos = append(listInvestmentDtos, dtos.NewInvestmentDto(investments[i]))
	}
	return listInvestmentDtos, nil
}
