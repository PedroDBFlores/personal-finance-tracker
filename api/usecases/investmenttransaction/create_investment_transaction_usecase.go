package investmenttransaction

import (
	"context"
	"log/slog"

	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/dtos"
	"github.com/pedrodbflores/personal-finance-tracker/usecases"
)

type CreateInvestmentTransactionUseCase struct {
	creator usecases.InvestmentTransactionCreator
	logger  *slog.Logger
}

func NewCreateInvestmentTransactionUseCase(creator usecases.InvestmentTransactionCreator, logger *slog.Logger) CreateInvestmentTransactionUseCase {
	return CreateInvestmentTransactionUseCase{creator, logger}
}

func (u *CreateInvestmentTransactionUseCase) Execute(ctx context.Context, accountId domain.Id, investmentId domain.Id, dto dtos.CreateInvestmentTransactionDto) (*domain.Id, error) {
	tx, err := domain.NewInvestmentTransaction(0, investmentId, dto.Type, dto.Amount, dto.Quantity, dto.Price, dto.Date, dto.Description)
	if err != nil {
		return nil, err
	}
	u.logger.DebugContext(ctx, "calling investment transaction creator", slog.Any("accountId", accountId), slog.Any("transaction", tx))
	return u.creator.Create(ctx, accountId, investmentId, tx)
}
