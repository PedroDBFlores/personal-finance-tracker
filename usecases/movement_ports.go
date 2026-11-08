package usecases

import (
	"context"

	"github.com/pedrodbflores/personal-finance-tracker/domain"
)

//counterfeiter:generate . MovementCreator
type MovementCreator interface {
	Create(ctx context.Context, accountId domain.Id, movement domain.Movement) (*domain.Id, error)
}

//counterfeiter:generate . MovementUpdater
type MovementUpdater interface {
	Update(ctx context.Context, accountId domain.Id, movement domain.Movement) error
}

//counterfeiter:generate . MovementGetter
type MovementGetter interface {
	Get(ctx context.Context, accountId, id domain.Id) (*domain.Movement, error)
}

//counterfeiter:generate . MovementLister
type MovementLister interface {
	List(ctx context.Context, accountId domain.Id) ([]domain.Movement, error)
}

//counterfeiter:generate . MovementDeleter
type MovementDeleter interface {
	Delete(ctx context.Context, accountId, id domain.Id) error
}
