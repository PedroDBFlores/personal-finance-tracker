package account

import (
	"context"

	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/dtos"
	"github.com/pedrodbflores/personal-finance-tracker/usecases"
)

type CreateAccountUseCase struct {
	creator usecases.AccountCreator
}

func NewCreateAccountUseCase(creator usecases.AccountCreator) CreateAccountUseCase {
	return CreateAccountUseCase{creator}
}

func (u *CreateAccountUseCase) Execute(ctx context.Context, dto dtos.CreateAccountDto) (*domain.Id, error) {
	return nil, nil
}
