package movement

import (
	"context"
	"log/slog"

	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/dtos"
	"github.com/pedrodbflores/personal-finance-tracker/usecases"
)

type GetMovementUseCase struct {
	getter usecases.MovementGetter
	logger *slog.Logger
}

func NewGetMovementUseCase(getter usecases.MovementGetter, logger *slog.Logger) GetMovementUseCase {
	return GetMovementUseCase{getter, logger}
}

func (u *GetMovementUseCase) Execute(ctx context.Context, accountId, id domain.Id) (*dtos.MovementDto, error) {
	result, err := u.getter.Get(ctx, accountId, id)
	if err != nil {
		return nil, err
	}
	dto := dtos.NewMovementDto(*result)
	return &dto, nil
}
