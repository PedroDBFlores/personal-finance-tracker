package investment

import (
	"context"
	"log/slog"

	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/usecases"
)

type UpdateInvestmentUseCase struct {
	updater usecases.InvestmentUpdater
	logger  *slog.Logger
}

func NewUpdateInvestmentUseCase(updater usecases.InvestmentUpdater, logger *slog.Logger) UpdateInvestmentUseCase {
	return UpdateInvestmentUseCase{updater, logger}
}

func (u *UpdateInvestmentUseCase) Execute(ctx context.Context, accountId domain.Id, id domain.Id, update domain.InvestmentPatch) error {
	u.logger.DebugContext(ctx, "updating investment", slog.Any("accountId", accountId), slog.Any("id", id), slog.Any("update", update))
	return u.updater.Update(ctx, accountId, id, update)
}
