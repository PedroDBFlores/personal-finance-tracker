package usecases

import (
	"github.com/pedrodbflores/personal-finance-tracker/domain"
)

//counterfeiter:generate . MovementCreator
type MovementCreator interface {
	Create(accountId domain.Id, movement domain.Movement) (domain.Id, error)
}

//counterfeiter:generate . MovementUpdater
type MovementUpdater interface {
	Update(accountId domain.Id, movement domain.Movement) error
}

//counterfeiter:generate . MovementGetter
type MovementGetter interface {
	Get(accountId, id domain.Id) (*domain.Movement, error)
}

//counterfeiter:generate . MovementLister
type MovementLister interface {
	List(accountId domain.Id) ([]domain.Movement, error)
}

//counterfeiter:generate . MovementDeleter
type MovementDeleter interface {
	Delete(accountId, id domain.Id) error
}
