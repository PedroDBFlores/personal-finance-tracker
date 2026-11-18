package investment

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/dtos"
	"github.com/pedrodbflores/personal-finance-tracker/usecases"
)

type CreateInvestmentUseCase struct {
	creator usecases.InvestmentCreator
	logger  *slog.Logger
}

func NewCreateInvestmentUseCase(creator usecases.InvestmentCreator, logger *slog.Logger) CreateInvestmentUseCase {
	return CreateInvestmentUseCase{creator, logger}
}

func (u *CreateInvestmentUseCase) Execute(ctx context.Context, accountId domain.Id, dto dtos.CreateInvestmentDto) (*uuid.UUID, error) {
	investment := domain.NewInvestment(uuid.New(), accountId, dto.Name, dto.Type, dto.Currency, dto.PurchaseDate, dto.MaturityDate)
	u.logger.DebugContext(ctx, "calling investment creator", slog.Any("investment", investment))
	return u.creator.Create(ctx, accountId, investment)
}
