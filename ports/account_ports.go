package ports

import (
	"github.com/pedrodbflores/personal-finance-tracker/domain"
)

//go:generate go tool counterfeiter -generate

//counterfeiter:generate . AccountCreator
type AccountCreator interface {
	Create(account domain.Account) (domain.Id, error)
}

//counterfeiter:generate . AccountUpdater
type AccountUpdater interface {
	Update(account domain.Account) error
}

//counterfeiter:generate . AccountGetter
type AccountGetter interface {
	Get(id domain.Id) (*domain.Account, error)
}

//counterfeiter:generate . AccountLister
type AccountLister interface {
	List() ([]domain.Account, error)
}

//counterfeiter:generate . AccountDeleter
type AccountDeleter interface {
	Delete(id domain.Id) error
}
