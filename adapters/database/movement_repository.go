package database

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/pedrodbflores/personal-finance-tracker/app_errors"
	"github.com/pedrodbflores/personal-finance-tracker/domain"
)

type movement struct {
	Id           domain.Id           `json:"id"`
	Description  string              `json:"description"`
	Value        domain.Money        `json:"value"`
	MovementType domain.MovementType `json:"type"`
	Date         time.Time           `json:"date"`
}

const movementKeyPrefix = "account:%s:movement"

func movementKey(accountId, movementId domain.Id) []byte {
	return fmt.Appendf(nil, movementKeyPrefix+":%s", accountId, movementId)
}

type MovementRepository struct {
	db *badger.DB
}

func NewMovementRepository(db *badger.DB) MovementRepository {
	return MovementRepository{db: db}
}

func (r *MovementRepository) Create(_ context.Context, accountId domain.Id, movement domain.Movement) (*domain.Id, error) {
	id := domain.NewId()
	key := movementKey(accountId, id)
	movement.Id = id
	b, err := json.Marshal(domainToMovement(movement))
	if err != nil {
		return nil, err
	}
	if err := r.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, b)
	}); err != nil {
		return nil, err
	}

	return &id, nil
}

func (r *MovementRepository) Update(_ context.Context, accountId domain.Id, movement domain.Movement) error {
	key := movementKey(accountId, movement.Id)
	ok, err := IsKeyPresent(r.db, key)
	if err != nil {
		return err
	}
	if !ok {
		return app_errors.NewResourceNotFoundError("Movement", string(key))
	}
	return r.db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get(key); err != nil {
			return err
		}
		b, err := json.Marshal(domainToMovement(movement))
		if err != nil {
			return err
		}
		return txn.Set(key, b)
	})
}

func (r *MovementRepository) Get(_ context.Context, accountId, id domain.Id) (*domain.Movement, error) {
	key := movementKey(accountId, id)
	out, err := GetByKey[movement](r.db, key)
	if err != nil {
		return nil, err
	}
	modelMovement := movementToDomain(out)
	return &modelMovement, nil
}

func (r *MovementRepository) List(_ context.Context, accountId domain.Id) ([]domain.Movement, error) {
	movements, err := ListAllByKey[movement](r.db, fmt.Appendf(nil, movementKeyPrefix, accountId))
	if err != nil {
		return []domain.Movement{}, err
	}
	var modelMovements []domain.Movement
	for _, movements := range movements {
		modelMovements = append(modelMovements, movementToDomain(movements))
	}

	return modelMovements, nil
}

func (r *MovementRepository) Delete(_ context.Context, accountId, id domain.Id) error {
	accKey := accountKey(accountId)
	ok, err := IsKeyPresent(r.db, accKey)
	if err != nil {
		return err
	}
	if !ok {
		return app_errors.NewResourceNotFoundError("Account", string(accKey))
	}
	key := movementKey(accountId, id)
	ok, err = IsKeyPresent(r.db, key)
	if err != nil {
		return err
	}
	if !ok {
		return app_errors.NewResourceNotFoundError("Movement", string(key))
	}
	return r.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}

func movementToDomain(m movement) domain.Movement {
	d := domain.NewMovement(
		m.Description,
		m.Value,
		m.MovementType,
		m.Date,
	)
	d.Id = m.Id
	return d
}

func domainToMovement(m domain.Movement) movement {
	return movement{
		Id:           m.Id,
		Description:  m.Description,
		Value:        m.Value,
		MovementType: m.Type,
		Date:         m.Date,
	}
}
