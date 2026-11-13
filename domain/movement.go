package domain

import (
	"time"
)

// Movement represents
type Movement struct {
	// Id represents the unique movement identifier.
	Id Id
	// Id represents the unique movement identifier.
	AccountId Id
	// Description defines what was the movement about (salary, bills, etc...)
	Description string
	// Value's last two digits are always the cents
	Value Money
	// Type defines the type of movement
	Type MovementType
	// Date defines the movement's date
	Date time.Time
}

// NewMovement creates a new Movement based on the provided parameters
func NewMovement(accountId Id, description string, value Money, movementType MovementType, date time.Time) Movement {
	return Movement{
		AccountId:   accountId,
		Description: description,
		Value:       value,
		Type:        movementType,
		Date:        date,
	}
}
