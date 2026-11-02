package ports

import (
	"uuid"

	"github.com/pedrodbflores/personal-finance-tracker/models"
)

//go:generate go tool counterfeiter -generate

//counterfeiter:generate . AccountCreator
type AccountCreator interface {
	Create(account models.Account) (uuid.UUID, error)
}

//counterfeiter:generate . AccountUpdater
type AccountUpdater interface {
	Update(account models.Account) error
}

//counterfeiter:generate . AccountGetter
type AccountGetter interface {
	Get(id uuid.UUID) (*models.Account, error)
}

//counterfeiter:generate . AccountLister
type AccountLister interface {
	List() ([]models.Account, error)
}

//counterfeiter:generate . AccountDeleter
type AccountDeleter interface {
	Delete(id uuid.UUID) error
}
