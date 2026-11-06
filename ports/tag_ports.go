package ports

import "github.com/pedrodbflores/personal-finance-tracker/domain"

//counterfeiter:generate . TagCreator
type TagCreator interface {
	Create(tag domain.Tag) (domain.Id, error)
}

//counterfeiter:generate . TagUpdater
type TagUpdater interface {
	Update(tag domain.Tag) error
}

//counterfeiter:generate . TagGetter
type TagGetter interface {
	Get(id domain.Id) (*domain.Tag, error)
}

//counterfeiter:generate . TagLister
type TagLister interface {
	List() ([]domain.Tag, error)
}

//counterfeiter:generate . TagDeleter
type TagDeleter interface {
	Delete(id domain.Id) error
}
