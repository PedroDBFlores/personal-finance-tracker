package database

import (
	"fmt"
	"math/big"
	"time"
	"uuid"

	"github.com/dgraph-io/badger/v4"
	"github.com/pedrodbflores/personal-finance-tracker/models"
)

type movement struct {
	id           uuid.UUID   `json:"id"`
	description  string      `json:"description"`
	value        big.Int     `json:"value"`
	movementType string      `json:"type"`
	date         time.Time   `json:"date"`
	tags         []uuid.UUID `json:"tags"`
}

const movementKeyPrefix = "account:%s:movement"

func movementKey(accountId, movementId string) string {
	return fmt.Sprintf(movementKeyPrefix+":%s", accountId, movementId)
}

type MovementRepository struct {
	db *badger.DB
}

func NewMovementRepository(db *badger.DB) MovementRepository {
	return MovementRepository{db: db}
}

func (r *MovementRepository) Create(accountId uuid.UUID, movement models.Movement) (uuid.UUID, error) {
	return uuid.Nil(), nil
}

func (r *MovementRepository) Update(accountId uuid.UUID, movement models.Movement) error {
	return nil
}

func (r *MovementRepository) Get(accountId, id uuid.UUID) (*models.Movement, error) {
	return nil, nil
}

func (r *MovementRepository) List(accountId uuid.UUID) ([]models.Movement, error) {
	return nil, nil
}

func (r *MovementRepository) Delete(accountId, id uuid.UUID) error {
	return nil
}
