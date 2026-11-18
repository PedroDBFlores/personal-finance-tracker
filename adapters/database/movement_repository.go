package database

import (
	"context"

	"github.com/pedrodbflores/personal-finance-tracker/app_errors"
	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MovementRepository struct {
	db *gorm.DB
}

func NewMovementRepository(db *gorm.DB) MovementRepository {
	return MovementRepository{db: db}
}

func (r *MovementRepository) Create(ctx context.Context, accountId domain.Id, movement domain.Movement) (*domain.Id, error) {
	movement.AccountId = accountId
	mov := domainMovementToDBMovement(movement)
	if err := gorm.G[Movement](r.db).Create(ctx, &mov); err != nil {
		return nil, err
	}
	id := domain.Id(mov.ID)
	return &id, nil
}

func (r *MovementRepository) Update(ctx context.Context, accountId, id domain.Id, update domain.MovementPatch) error {
	assignments := map[string]interface{}{}
	if update.Description != nil {
		assignments["description"] = *update.Description
	}
	if update.Value != nil {
		assignments["value"] = *update.Value
	}
	if update.Type != nil {
		assignments["movement_type"] = *update.Type
	}
	if update.Date != nil {
		assignments["date"] = *update.Date
	}
	if len(assignments) == 0 {
		return nil
	}
	if rowsAffected, err := gorm.G[Movement](r.db).
		Where("id = ? AND account_id = ?", id, accountId).
		Set(clause.Assignments(assignments)).Update(ctx); err != nil {
		return err
	} else if rowsAffected == 0 {
		return app_errors.NewResourceNotFoundError("Movement", id.String())
	}
	return nil
}

func (r *MovementRepository) Get(ctx context.Context, accountId, id domain.Id) (*domain.Movement, error) {
	if mov, err := gorm.G[Movement](r.db).Where("id = ? AND account_id = ?", id, accountId).First(ctx); err != nil {
		return nil, app_errors.NewResourceNotFoundError("Movement", id.String())
	} else {
		movement := dbMovementToDomainMovement(mov)
		return &movement, nil
	}
}

func (r *MovementRepository) List(ctx context.Context, accountId domain.Id) ([]domain.Movement, error) {
	movements, err := gorm.G[Movement](r.db).Where("account_id = ?", accountId).Find(ctx)
	if err != nil {
		return []domain.Movement{}, err
	}
	var modelMovements []domain.Movement
	for _, movements := range movements {
		modelMovements = append(modelMovements, dbMovementToDomainMovement(movements))
	}

	return modelMovements, nil
}

func (r *MovementRepository) Delete(ctx context.Context, accountId, id domain.Id) error {
	if rowsAffected, err := gorm.G[Movement](r.db).Where("id = ? AND account_id = ?", id, accountId).Delete(ctx); err != nil {
		return err
	} else if rowsAffected == 0 {
		return app_errors.NewResourceNotFoundError("Movement", id.String())
	}
	return nil
}

func dbMovementToDomainMovement(m Movement) domain.Movement {
	d := domain.NewMovement(
		domain.Id(m.AccountId),
		m.Description,
		m.Value,
		m.MovementType,
		m.Date,
	)
	d.Id = domain.Id(m.ID)
	return d
}

func domainMovementToDBMovement(m domain.Movement) Movement {
	return Movement{
		ID:           uint(m.Id),
		AccountId:    uint(m.AccountId),
		Description:  m.Description,
		Value:        m.Value,
		MovementType: m.Type,
		Date:         m.Date,
	}
}
