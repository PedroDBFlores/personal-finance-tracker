package movement

import (
	"context"
	"log/slog"

	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/dtos"
	"github.com/pedrodbflores/personal-finance-tracker/usecases"
)

type ListMovementsUseCase struct {
	lister usecases.MovementLister
	logger *slog.Logger
}

func NewListMovementsUseCase(lister usecases.MovementLister, logger *slog.Logger) ListMovementsUseCase {
	return ListMovementsUseCase{lister, logger}
}

func (u *ListMovementsUseCase) Execute(ctx context.Context, accountId domain.Id) ([]dtos.MovementDto, error) {
	result, err := u.lister.List(ctx, accountId)
	if err != nil {
		return []dtos.MovementDto{}, err
	}
	var listMovementDtos []dtos.MovementDto
	for _, movement := range result {
		listMovementDtos = append(listMovementDtos, dtos.NewMovementDto(movement))
	}
	return listMovementDtos, nil
}
