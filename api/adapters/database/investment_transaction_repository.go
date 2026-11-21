package database

import (
	"context"

	"github.com/pedrodbflores/personal-finance-tracker/app_errors"
	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type InvestmentTransactionRepository struct {
	db *gorm.DB
}

func NewInvestmentTransactionRepository(db *gorm.DB) InvestmentTransactionRepository {
	return InvestmentTransactionRepository{db: db}
}

// accountScope restricts a transaction query to investments owned by the given
// account (transactions carry no account_id of their own).
const accountScope = "investment_id IN (SELECT id FROM investments WHERE account_id = ?)"

func (r *InvestmentTransactionRepository) Create(ctx context.Context, accountId domain.Id, investmentId domain.Id, tx domain.InvestmentTransaction) (*domain.Id, error) {
	if _, err := gorm.G[Investment](r.db).Where("id = ? AND account_id = ?", investmentId, accountId).First(ctx); err != nil {
		return nil, app_errors.NewResourceNotFoundError("Investment", investmentId.String())
	}
	tx.InvestmentId = investmentId
	model := domainInvestmentTransactionToDBInvestmentTransaction(tx)
	if err := gorm.G[InvestmentTransaction](r.db).Create(ctx, &model); err != nil {
		return nil, err
	}
	id := domain.Id(model.ID)
	return &id, nil
}

func (r *InvestmentTransactionRepository) Update(ctx context.Context, accountId domain.Id, investmentId, id domain.Id, update domain.InvestmentTransactionPatch) error {
	assignments := map[string]interface{}{}
	if update.Type != nil {
		assignments["type"] = *update.Type
	}
	if update.Amount != nil {
		assignments["amount"] = *update.Amount
	}
	if update.Quantity != nil {
		assignments["quantity"] = *update.Quantity
	}
	if update.Price != nil {
		assignments["price"] = *update.Price
	}
	if update.Date != nil {
		assignments["date"] = *update.Date
	}
	if update.Description != nil {
		assignments["description"] = *update.Description
	}
	if len(assignments) == 0 {
		return nil
	}
	if rowsAffected, err := gorm.G[InvestmentTransaction](r.db).
		Where("id = ? AND investment_id = ? AND "+accountScope, id, investmentId, accountId).
		Set(clause.Assignments(assignments)).Update(ctx); err != nil {
		return err
	} else if rowsAffected == 0 {
		return app_errors.NewResourceNotFoundError("InvestmentTransaction", id.String())
	}
	return nil
}

func (r *InvestmentTransactionRepository) Get(ctx context.Context, accountId domain.Id, investmentId, id domain.Id) (*domain.InvestmentTransaction, error) {
	if tx, err := gorm.G[InvestmentTransaction](r.db).Where("id = ? AND investment_id = ? AND "+accountScope, id, investmentId, accountId).First(ctx); err != nil {
		return nil, app_errors.NewResourceNotFoundError("InvestmentTransaction", id.String())
	} else {
		transaction, err := dbInvestmentTransactionToDomainInvestmentTransaction(tx)
		if err != nil {
			return nil, err
		}
		return &transaction, nil
	}
}

func (r *InvestmentTransactionRepository) List(ctx context.Context, accountId domain.Id, investmentId domain.Id) ([]domain.InvestmentTransaction, error) {
	transactions, err := gorm.G[InvestmentTransaction](r.db).Where("investment_id = ? AND "+accountScope, investmentId, accountId).Find(ctx)
	if err != nil {
		return []domain.InvestmentTransaction{}, err
	}
	modelTransactions := make([]domain.InvestmentTransaction, 0, len(transactions))
	for _, tx := range transactions {
		t, err := dbInvestmentTransactionToDomainInvestmentTransaction(tx)
		if err != nil {
			return nil, err
		}
		modelTransactions = append(modelTransactions, t)
	}
	return modelTransactions, nil
}

func (r *InvestmentTransactionRepository) Delete(ctx context.Context, accountId domain.Id, investmentId, id domain.Id) error {
	if rowsAffected, err := gorm.G[InvestmentTransaction](r.db).Where("id = ? AND investment_id = ? AND "+accountScope, id, investmentId, accountId).Delete(ctx); err != nil {
		return err
	} else if rowsAffected == 0 {
		return app_errors.NewResourceNotFoundError("InvestmentTransaction", id.String())
	}
	return nil
}

func dbInvestmentTransactionToDomainInvestmentTransaction(t InvestmentTransaction) (domain.InvestmentTransaction, error) {
	return domain.NewInvestmentTransaction(
		domain.Id(t.ID),
		domain.Id(t.InvestmentId),
		t.Type,
		t.Amount,
		t.Quantity,
		t.Price,
		t.Date,
		t.Description,
	)
}

func domainInvestmentTransactionToDBInvestmentTransaction(t domain.InvestmentTransaction) InvestmentTransaction {
	return InvestmentTransaction{
		ID:           uint(t.Id),
		InvestmentId: uint(t.InvestmentId),
		Type:         t.Type,
		Amount:       t.Amount,
		Quantity:     t.Quantity,
		Price:        t.Price,
		Date:         t.Date,
		Description:  t.Description,
	}
}
