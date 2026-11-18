package account

import (
	"context"
	"log/slog"

	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/dtos"
	"github.com/pedrodbflores/personal-finance-tracker/usecases"
)

type UpdateAccountUseCase struct {
	updater usecases.AccountUpdater
	logger  *slog.Logger
}

func NewUpdateAccountUseCase(updater usecases.AccountUpdater, logger *slog.Logger) UpdateAccountUseCase {
	return UpdateAccountUseCase{updater, logger}
}

func (u *UpdateAccountUseCase) Execute(ctx context.Context, dto dtos.UpdateAccountDto) error {
	account := domain.Account{
		Id:        dto.Id,
		Name:      dto.Name,
		Type:      domain.AccountType(dto.Type),
		Currency:  dto.Currency,
		CreatedAt: dto.CreatedAt,
	}
	u.logger.DebugContext(ctx, "calling account updater", slog.Any("account", account))
	return u.updater.Update(ctx, account)
}
