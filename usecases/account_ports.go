package usecases

import (
	"context"

	"github.com/pedrodbflores/personal-finance-tracker/domain"
)

//counterfeiter:generate . AccountCreator
type AccountCreator interface {
	Create(ctx context.Context, account domain.Account) (*domain.Id, error)
}

//counterfeiter:generate . AccountUpdater
type AccountUpdater interface {
	Update(ctx context.Context, account domain.Account) error
}

//counterfeiter:generate . AccountGetter
type AccountGetter interface {
	Get(ctx context.Context, id domain.Id) (*domain.Account, error)
}

//counterfeiter:generate . AccountLister
type AccountLister interface {
	List(ctx context.Context) ([]domain.Account, error)
}

//counterfeiter:generate . AccountDeleter
type AccountDeleter interface {
	Delete(ctx context.Context, id domain.Id) error
}
