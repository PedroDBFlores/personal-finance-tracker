package investment

import (
	"context"
	"log/slog"

	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/usecases"
)

type DeleteInvestmentUseCase struct {
	deleter usecases.InvestmentDeleter
	logger  *slog.Logger
}

func NewDeleteInvestmentUseCase(deleter usecases.InvestmentDeleter, logger *slog.Logger) DeleteInvestmentUseCase {
	return DeleteInvestmentUseCase{deleter, logger}
}

func (u *DeleteInvestmentUseCase) Execute(ctx context.Context, accountId domain.Id, id domain.Id) error {
	u.logger.DebugContext(ctx, "deleting investment", slog.Any("accountId", accountId), slog.Any("id", id))
	return u.deleter.Delete(ctx, accountId, id)
}
