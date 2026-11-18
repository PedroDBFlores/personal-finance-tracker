package movement

import (
	"context"
	"log/slog"

	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/usecases"
)

type UpdateMovementUseCase struct {
	updater usecases.MovementUpdater
	logger  *slog.Logger
}

func NewUpdateMovementUseCase(updater usecases.MovementUpdater, logger *slog.Logger) UpdateMovementUseCase {
	return UpdateMovementUseCase{updater, logger}
}

func (u *UpdateMovementUseCase) Execute(ctx context.Context, accountId, id domain.Id, update domain.MovementPatch) error {
	u.logger.DebugContext(ctx, "updating movement", slog.Any("accountId", accountId), slog.Any("id", id), slog.Any("update", update))
	return u.updater.Update(ctx, accountId, id, update)
}
