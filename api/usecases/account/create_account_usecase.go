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
	account, err := domain.NewAccount(0, dto.Name, dto.Type, dto.Currency, dto.CreatedAt)
	if err != nil {
		return nil, err
	}
	u.logger.DebugContext(ctx, "calling account creator", slog.Any("account", account))
	return u.creator.Create(ctx, account)
}
