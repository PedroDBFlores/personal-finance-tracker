package ports

import (
	"uuid"

	"github.com/pedrodbflores/personal-finance-tracker/models"
)

//counterfeiter:generate . TagCreator
type TagCreator interface {
	Create(tag models.Tag) (uuid.UUID, error)
}

//counterfeiter:generate . TagUpdater
type TagUpdater interface {
	Update(tag models.Tag) error
}

//counterfeiter:generate . TagGetter
type TagGetter interface {
	Get(id uuid.UUID) (*models.Tag, error)
}

//counterfeiter:generate . TagLister
type TagLister interface {
	List() ([]models.Tag, error)
}

//counterfeiter:generate . TagDeleter
type TagDeleter interface {
	Delete(id uuid.UUID) error
}
