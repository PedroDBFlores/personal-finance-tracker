package account

import "github.com/pedrodbflores/personal-finance-tracker/usecases"

type CreateAccountUseCase struct {
	creator usecases.AccountCreator
}

func NewCreateAccountUseCase(creator usecases.AccountCreator) CreateAccountUseCase {
	return CreateAccountUseCase{creator}
}
