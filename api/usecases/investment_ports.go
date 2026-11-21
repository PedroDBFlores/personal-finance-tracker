package usecases

import (
	"context"

	"github.com/pedrodbflores/personal-finance-tracker/domain"
)

//counterfeiter:generate . InvestmentCreator
type InvestmentCreator interface {
	Create(ctx context.Context, accountId domain.Id, investment domain.Investment) (*domain.Id, error)
}

//counterfeiter:generate . InvestmentUpdater
type InvestmentUpdater interface {
	Update(ctx context.Context, accountId domain.Id, id domain.Id, update domain.InvestmentPatch) error
}

//counterfeiter:generate . InvestmentGetter
type InvestmentGetter interface {
	Get(ctx context.Context, accountId domain.Id, id domain.Id) (*domain.Investment, error)
}

//counterfeiter:generate . InvestmentLister
type InvestmentLister interface {
	List(ctx context.Context, accountId domain.Id) ([]domain.Investment, error)
}

//counterfeiter:generate . InvestmentDeleter
type InvestmentDeleter interface {
	Delete(ctx context.Context, accountId domain.Id, id domain.Id) error
}
