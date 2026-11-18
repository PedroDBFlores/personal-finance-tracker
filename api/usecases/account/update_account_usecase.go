package account

import (
	"context"
	"log/slog"

	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/usecases"
)

type UpdateAccountUseCase struct {
	updater usecases.AccountUpdater
	logger  *slog.Logger
}

func NewUpdateAccountUseCase(updater usecases.AccountUpdater, logger *slog.Logger) UpdateAccountUseCase {
	return UpdateAccountUseCase{updater, logger}
}

func (u *UpdateAccountUseCase) Execute(ctx context.Context, id domain.Id, update domain.AccountPatch) error {
	u.logger.DebugContext(ctx, "updating account", slog.Any("id", id), slog.Any("update", update))
	return u.updater.Update(ctx, id, update)
}
