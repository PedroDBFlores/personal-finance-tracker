package usecases

import (
	"context"

	"github.com/google/uuid"
	"github.com/pedrodbflores/personal-finance-tracker/domain"
)

// Investment-transaction ports are scoped by the full parent chain
// account -> investment -> transaction, so every method takes the account id.

//counterfeiter:generate . InvestmentTransactionCreator
type InvestmentTransactionCreator interface {
	Create(ctx context.Context, accountId domain.Id, investmentId uuid.UUID, tx domain.InvestmentTransaction) (*uuid.UUID, error)
}

//counterfeiter:generate . InvestmentTransactionUpdater
type InvestmentTransactionUpdater interface {
	Update(ctx context.Context, accountId domain.Id, investmentId, id uuid.UUID, update domain.InvestmentTransactionPatch) error
}

//counterfeiter:generate . InvestmentTransactionGetter
type InvestmentTransactionGetter interface {
	Get(ctx context.Context, accountId domain.Id, investmentId, id uuid.UUID) (*domain.InvestmentTransaction, error)
}

//counterfeiter:generate . InvestmentTransactionLister
type InvestmentTransactionLister interface {
	List(ctx context.Context, accountId domain.Id, investmentId uuid.UUID) ([]domain.InvestmentTransaction, error)
}

//counterfeiter:generate . InvestmentTransactionDeleter
type InvestmentTransactionDeleter interface {
	Delete(ctx context.Context, accountId domain.Id, investmentId, id uuid.UUID) error
}
