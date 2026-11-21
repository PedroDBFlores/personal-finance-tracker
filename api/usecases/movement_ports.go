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
	Update(ctx context.Context, accountId, id domain.Id, update domain.MovementPatch) error
}

//counterfeiter:generate . MovementGetter
type MovementGetter interface {
	Get(ctx context.Context, accountId, id domain.Id) (*domain.Movement, error)
}

//counterfeiter:generate . MovementLister
type MovementLister interface {
	List(ctx context.Context, accountId domain.Id) ([]domain.Movement, error)
}

//counterfeiter:generate . MovementMonthLister
type MovementMonthLister interface {
	ListMonths(ctx context.Context, accountId domain.Id) ([]domain.YearMonth, error)
}

//counterfeiter:generate . MovementByMonthLister
type MovementByMonthLister interface {
	ListByMonth(ctx context.Context, accountId domain.Id, year, month int) ([]domain.Movement, error)
}

//counterfeiter:generate . MovementDeleter
type MovementDeleter interface {
	Delete(ctx context.Context, accountId, id domain.Id) error
}
