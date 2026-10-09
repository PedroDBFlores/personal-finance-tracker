package models

import (
	"math/big"
	"time"

	"github.com/google/uuid"
)

// Movement represents
type Movement struct {
	// Id represents the unique movement identifier.
	Id uuid.UUID
	// Description defines what was the movement about (salary, bills, etc...)
	Description string
	// Value's last two digits are always the cents
	Value big.Int
	// Type defines the type of movement
	Type MovementType
	// Date defines the movement's date
	Date time.Time
	// AccountId links the movement to an account
	AccountId uuid.UUID
	// Tags links the movement to tags for categorization
	Tags []uuid.UUID
}

// NewMovement creates a new Movement based on the provided parameters
func NewMovement(id uuid.UUID, description string, value big.Int, movementType MovementType, date time.Time, accountId uuid.UUID, tags []uuid.UUID) Movement {
	return Movement{
		Id:          id,
		Description: description,
		Value:       value,
		Type:        movementType,
		Date:        date,
		AccountId:   accountId,
		Tags:        tags,
	}
}
