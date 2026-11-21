package investmenttransaction

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/usecases"
)

type DeleteInvestmentTransactionUseCase struct {
	deleter usecases.InvestmentTransactionDeleter
	logger  *slog.Logger
}

func NewDeleteInvestmentTransactionUseCase(deleter usecases.InvestmentTransactionDeleter, logger *slog.Logger) DeleteInvestmentTransactionUseCase {
	return DeleteInvestmentTransactionUseCase{deleter, logger}
}

func (u *DeleteInvestmentTransactionUseCase) Execute(ctx context.Context, accountId domain.Id, investmentId, id uuid.UUID) error {
	u.logger.DebugContext(ctx, "deleting investment transaction", slog.Any("accountId", accountId), slog.Any("investmentId", investmentId), slog.Any("id", id))
	return u.deleter.Delete(ctx, accountId, investmentId, id)
}
