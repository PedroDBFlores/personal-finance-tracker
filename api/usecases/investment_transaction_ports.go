package usecases

import (
	"context"

	"github.com/google/uuid"
	"github.com/pedrodbflores/personal-finance-tracker/domain"
)

//counterfeiter:generate . InvestmentTransactionCreator
type InvestmentTransactionCreator interface {
	Create(ctx context.Context, investmentId uuid.UUID, tx domain.InvestmentTransaction) (*uuid.UUID, error)
}

//counterfeiter:generate . InvestmentTransactionUpdater
type InvestmentTransactionUpdater interface {
	Update(ctx context.Context, investmentId, id uuid.UUID, update domain.InvestmentTransactionPatch) error
}

//counterfeiter:generate . InvestmentTransactionGetter
type InvestmentTransactionGetter interface {
	Get(ctx context.Context, investmentId, id uuid.UUID) (*domain.InvestmentTransaction, error)
}

//counterfeiter:generate . InvestmentTransactionLister
type InvestmentTransactionLister interface {
	List(ctx context.Context, investmentId uuid.UUID) ([]domain.InvestmentTransaction, error)
}

//counterfeiter:generate . InvestmentTransactionDeleter
type InvestmentTransactionDeleter interface {
	Delete(ctx context.Context, investmentId, id uuid.UUID) error
}
