package investmenttransaction

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
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

func (u *CreateInvestmentTransactionUseCase) Execute(ctx context.Context, accountId domain.Id, investmentId uuid.UUID, dto dtos.CreateInvestmentTransactionDto) (*uuid.UUID, error) {
	tx := domain.NewInvestmentTransaction(uuid.New(), investmentId, dto.Type, dto.Amount, dto.Quantity, dto.Price, dto.Date, dto.Description)
	u.logger.DebugContext(ctx, "calling investment transaction creator", slog.Any("accountId", accountId), slog.Any("transaction", tx))
	return u.creator.Create(ctx, accountId, investmentId, tx)
}
