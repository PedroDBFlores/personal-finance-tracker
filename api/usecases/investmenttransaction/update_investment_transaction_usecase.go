package investmenttransaction

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/usecases"
)

type UpdateInvestmentTransactionUseCase struct {
	updater usecases.InvestmentTransactionUpdater
	logger  *slog.Logger
}

func NewUpdateInvestmentTransactionUseCase(updater usecases.InvestmentTransactionUpdater, logger *slog.Logger) UpdateInvestmentTransactionUseCase {
	return UpdateInvestmentTransactionUseCase{updater, logger}
}

func (u *UpdateInvestmentTransactionUseCase) Execute(ctx context.Context, accountId domain.Id, investmentId, id uuid.UUID, update domain.InvestmentTransactionPatch) error {
	u.logger.DebugContext(ctx, "updating investment transaction", slog.Any("accountId", accountId), slog.Any("investmentId", investmentId), slog.Any("id", id), slog.Any("update", update))
	return u.updater.Update(ctx, accountId, investmentId, id, update)
}
