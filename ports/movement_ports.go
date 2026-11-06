package ports

import (
	"uuid"

	"github.com/pedrodbflores/personal-finance-tracker/models"
)

//counterfeiter:generate . MovementCreator
type MovementCreator interface {
	Create(accountId uuid.UUID, movement models.Movement) (uuid.UUID, error)
}

//counterfeiter:generate . MovementUpdater
type MovementUpdater interface {
	Update(accountId uuid.UUID, movement models.Movement) error
}

//counterfeiter:generate . MovementGetter
type MovementGetter interface {
	Get(accountId, id uuid.UUID) (*models.Movement, error)
}

//counterfeiter:generate . MovementLister
type MovementLister interface {
	List(accountId uuid.UUID) ([]models.Movement, error)
}

//counterfeiter:generate . MovementDeleter
type MovementDeleter interface {
	Delete(accountId, id uuid.UUID) error
}
