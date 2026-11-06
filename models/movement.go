package models

import (
	"time"

	"uuid"
)

// Movement represents
type Movement struct {
	// Id represents the unique movement identifier.
	Id uuid.UUID
	// Description defines what was the movement about (salary, bills, etc...)
	Description string
	// Value's last two digits are always the cents
	Value int64
	// Type defines the type of movement
	Type MovementType
	// Date defines the movement's date
	Date time.Time
}

// NewMovement creates a new Movement based on the provided parameters
func NewMovement(id uuid.UUID, description string, value int64, movementType MovementType, date time.Time) Movement {
	return Movement{
		Id:          id,
		Description: description,
		Value:       value,
		Type:        movementType,
		Date:        date,
	}
}
