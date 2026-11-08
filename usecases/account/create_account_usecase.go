package account

import (
	"context"
	"log/slog"

	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/dtos"
	"github.com/pedrodbflores/personal-finance-tracker/usecases"
)

type CreateAccountUseCase struct {
	creator usecases.AccountCreator
	logger  *slog.Logger
}

func NewCreateAccountUseCase(creator usecases.AccountCreator, logger *slog.Logger) CreateAccountUseCase {
	return CreateAccountUseCase{creator, logger}
}

func (u *CreateAccountUseCase) Execute(ctx context.Context, dto dtos.CreateAccountDto) (*domain.Id, error) {
	account := domain.Account{
		Name:      dto.Name,
		Type:      domain.AccountType(dto.Type),
		Currency:  dto.Currency,
		CreatedAt: dto.CreatedAt,
	}
	u.logger.DebugContext(ctx, "calling account creator", slog.Any("account", account))
	return u.creator.Create(ctx, account)
}
