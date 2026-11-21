package movement

import (
	"context"
	"log/slog"

	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/dtos"
	"github.com/pedrodbflores/personal-finance-tracker/usecases"
)

type CreateMovementUseCase struct {
	creator usecases.MovementCreator
	logger  *slog.Logger
}

func NewCreateMovementUseCase(creator usecases.MovementCreator, logger *slog.Logger) CreateMovementUseCase {
	return CreateMovementUseCase{creator, logger}
}

func (u *CreateMovementUseCase) Execute(ctx context.Context, accountId domain.Id, dto dtos.CreateMovementDto) (*domain.Id, error) {
	movement, err := domain.NewMovement(accountId, dto.Description, dto.Value, dto.Type, dto.Date)
	if err != nil {
		return nil, err
	}
	u.logger.DebugContext(ctx, "calling movement creator", slog.Any("movement", movement))
	return u.creator.Create(ctx, accountId, movement)
}
