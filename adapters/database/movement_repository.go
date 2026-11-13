package database

import (
	"context"

	"github.com/pedrodbflores/personal-finance-tracker/app_errors"
	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"gorm.io/gorm"
)

type MovementRepository struct {
	db *gorm.DB
}

func NewMovementRepository(db *gorm.DB) MovementRepository {
	return MovementRepository{db: db}
}

func (r *MovementRepository) Create(ctx context.Context, movement domain.Movement) (*domain.Id, error) {
	mov := domainMovementToDBMovement(movement)
	if err := gorm.G[Movement](r.db).Create(ctx, &mov); err != nil {
		return nil, err
	}
	id := domain.Id(mov.ID)
	return &id, nil
}

func (r *MovementRepository) Update(ctx context.Context, movement domain.Movement) error {
	if rowsAffected, err := gorm.G[Movement](r.db).Updates(ctx, domainMovementToDBMovement(movement)); rowsAffected == 0 {
		return app_errors.NewResourceNotFoundError("Movement", movement.Id.String())
	} else if err != nil {
		return err
	}
	return nil
}

func (r *MovementRepository) Get(ctx context.Context, id domain.Id) (*domain.Movement, error) {
	if acc, err := gorm.G[Movement](r.db).Where("id = ?", id).First(ctx); err != nil {
		return nil, app_errors.NewResourceNotFoundError("Movement", id.String())
	} else {
		account := dbMovementToDomainMovement(acc)
		return &account, nil
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

func (r *MovementRepository) Delete(ctx context.Context, id domain.Id) error {
	if rowsAffected, err := gorm.G[Movement](r.db).Where("id = ?", id).Delete(ctx); err != nil {
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
