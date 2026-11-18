package account

import (
	"context"
	"log/slog"

	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/usecases"
)

type DeleteAccountUseCase struct {
	deleter usecases.AccountDeleter
	logger  *slog.Logger
}

func NewDeleteAccountUseCase(deleter usecases.AccountDeleter, logger *slog.Logger) DeleteAccountUseCase {
	return DeleteAccountUseCase{deleter, logger}
}

func (u *DeleteAccountUseCase) Execute(ctx context.Context, id domain.Id) error {
	u.logger.DebugContext(ctx, "calling account deleter", slog.Any("id", id))
	return u.deleter.Delete(ctx, id)
}
