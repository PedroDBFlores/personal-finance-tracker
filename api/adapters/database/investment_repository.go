package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/pedrodbflores/personal-finance-tracker/app_errors"
	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type InvestmentRepository struct {
	db *gorm.DB
}

func NewInvestmentRepository(db *gorm.DB) InvestmentRepository {
	return InvestmentRepository{db: db}
}

func (r *InvestmentRepository) Create(ctx context.Context, accountId domain.Id, investment domain.Investment) (*uuid.UUID, error) {
	investment.AccountId = accountId
	if investment.Id == uuid.Nil {
		investment.Id = uuid.New()
	}
	inv := domainInvestmentToDBInvestment(investment)
	if err := gorm.G[Investment](r.db).Create(ctx, &inv); err != nil {
		return nil, err
	}
	return &inv.ID, nil
}

func (r *InvestmentRepository) Update(ctx context.Context, accountId domain.Id, id uuid.UUID, update domain.InvestmentPatch) error {
	assignments := map[string]interface{}{}
	if update.Name != nil {
		assignments["name"] = *update.Name
	}
	if update.Type != nil {
		assignments["type"] = *update.Type
	}
	if update.Currency != nil {
		assignments["currency"] = *update.Currency
	}
	if update.PurchaseDate != nil {
		assignments["purchase_date"] = *update.PurchaseDate
	}
	if update.MaturityDate != nil {
		assignments["maturity_date"] = *update.MaturityDate
	}
	if update.CurrentValue != nil {
		assignments["current_value"] = *update.CurrentValue
	}
	if len(assignments) == 0 {
		return nil
	}
	if rowsAffected, err := gorm.G[Investment](r.db).
		Where("id = ? AND account_id = ?", id, accountId).
		Set(clause.Assignments(assignments)).Update(ctx); err != nil {
		return err
	} else if rowsAffected == 0 {
		return app_errors.NewResourceNotFoundError("Investment", id.String())
	}
	return nil
}

func (r *InvestmentRepository) Get(ctx context.Context, accountId domain.Id, id uuid.UUID) (*domain.Investment, error) {
	if inv, err := gorm.G[Investment](r.db).Where("id = ? AND account_id = ?", id, accountId).First(ctx); err != nil {
		return nil, app_errors.NewResourceNotFoundError("Investment", id.String())
	} else {
		investment := dbInvestmentToDomainInvestment(inv)
		return &investment, nil
	}
}

func (r *InvestmentRepository) List(ctx context.Context, accountId domain.Id) ([]domain.Investment, error) {
	investments, err := gorm.G[Investment](r.db).Where("account_id = ?", accountId).Find(ctx)
	if err != nil {
		return []domain.Investment{}, err
	}
	modelInvestments := make([]domain.Investment, 0, len(investments))
	for _, inv := range investments {
		modelInvestments = append(modelInvestments, dbInvestmentToDomainInvestment(inv))
	}
	return modelInvestments, nil
}

func (r *InvestmentRepository) Delete(ctx context.Context, accountId domain.Id, id uuid.UUID) error {
	if rowsAffected, err := gorm.G[Investment](r.db).Where("id = ? AND account_id = ?", id, accountId).Delete(ctx); err != nil {
		return err
	} else if rowsAffected == 0 {
		return app_errors.NewResourceNotFoundError("Investment", id.String())
	}
	return nil
}

func dbInvestmentToDomainInvestment(i Investment) domain.Investment {
	d := domain.NewInvestment(
		i.ID,
		domain.Id(i.AccountId),
		i.Name,
		i.Type,
		i.Currency,
		i.PurchaseDate,
		i.MaturityDate,
	)
	d.CurrentValue = i.CurrentValue
	return d
}

func domainInvestmentToDBInvestment(i domain.Investment) Investment {
	return Investment{
		ID:           i.Id,
		AccountId:    uint(i.AccountId),
		Name:         i.Name,
		Type:         i.Type,
		CurrentValue: i.CurrentValue,
		Currency:     i.Currency,
		PurchaseDate: i.PurchaseDate,
		MaturityDate: i.MaturityDate,
	}
}
