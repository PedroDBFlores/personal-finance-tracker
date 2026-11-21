package usecases

import (
	"context"

	"github.com/pedrodbflores/personal-finance-tracker/domain"
)

// Investment-transaction ports are scoped by the full parent chain
// account -> investment -> transaction, so every method takes the account id.

//counterfeiter:generate . InvestmentTransactionCreator
type InvestmentTransactionCreator interface {
	Create(ctx context.Context, accountId domain.Id, investmentId domain.Id, tx domain.InvestmentTransaction) (*domain.Id, error)
}

//counterfeiter:generate . InvestmentTransactionUpdater
type InvestmentTransactionUpdater interface {
	Update(ctx context.Context, accountId domain.Id, investmentId, id domain.Id, update domain.InvestmentTransactionPatch) error
}

//counterfeiter:generate . InvestmentTransactionGetter
type InvestmentTransactionGetter interface {
	Get(ctx context.Context, accountId domain.Id, investmentId, id domain.Id) (*domain.InvestmentTransaction, error)
}

//counterfeiter:generate . InvestmentTransactionLister
type InvestmentTransactionLister interface {
	List(ctx context.Context, accountId domain.Id, investmentId domain.Id) ([]domain.InvestmentTransaction, error)
}

//counterfeiter:generate . InvestmentTransactionDeleter
type InvestmentTransactionDeleter interface {
	Delete(ctx context.Context, accountId domain.Id, investmentId, id domain.Id) error
}
