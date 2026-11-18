package investment

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/dtos"
	"github.com/pedrodbflores/personal-finance-tracker/usecases"
)

type GetInvestmentUseCase struct {
	getter   usecases.InvestmentGetter
	txLister usecases.InvestmentTransactionLister
	logger   *slog.Logger
}

func NewGetInvestmentUseCase(getter usecases.InvestmentGetter, txLister usecases.InvestmentTransactionLister, logger *slog.Logger) GetInvestmentUseCase {
	return GetInvestmentUseCase{getter, txLister, logger}
}

func (u *GetInvestmentUseCase) Execute(ctx context.Context, accountId domain.Id, id uuid.UUID) (*dtos.InvestmentDto, error) {
	investment, err := u.getter.Get(ctx, accountId, id)
	if err != nil {
		return nil, err
	}
	u.logger.DebugContext(ctx, "listing transactions to calculate current value", slog.Any("investmentId", id))
	transactions, err := u.txLister.List(ctx, id)
	if err != nil {
		return nil, err
	}
	investment.UpdateCurrentValue(transactions)
	dto := dtos.NewInvestmentDto(*investment)
	return &dto, nil
}
