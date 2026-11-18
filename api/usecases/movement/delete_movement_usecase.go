package movement

import (
	"context"
	"log/slog"

	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/usecases"
)

type DeleteMovementUseCase struct {
	deleter usecases.MovementDeleter
	logger  *slog.Logger
}

func NewDeleteMovementUseCase(deleter usecases.MovementDeleter, logger *slog.Logger) DeleteMovementUseCase {
	return DeleteMovementUseCase{deleter, logger}
}

func (u *DeleteMovementUseCase) Execute(ctx context.Context, accountId, id domain.Id) error {
	u.logger.DebugContext(ctx, "deleting movement", slog.Any("accountId", accountId), slog.Any("id", id))
	return u.deleter.Delete(ctx, accountId, id)
}
