package domain

import (
	"strings"
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

// NewMovement creates a new Movement based on the provided parameters.
// It returns a ValidationError if the description is empty, the movement type
// is not supported, or the value is negative.
func NewMovement(accountId Id, description string, value Money, movementType MovementType, date time.Time) (Movement, error) {
	if strings.TrimSpace(description) == "" {
		return Movement{}, newValidationError("movement description is required")
	}
	switch movementType {
	case MovementTypeCredit, MovementTypeDebit:
	default:
		return Movement{}, newValidationError("unsupported movement type %q", movementType)
	}
	if value < 0 {
		return Movement{}, newValidationError("movement value must not be negative")
	}
	return Movement{
		AccountId:   accountId,
		Description: description,
		Value:       value,
		Type:        movementType,
		Date:        date,
	}, nil
}

// MovementPatch represents the fields to modify on a Movement. Nil fields are left unchanged.
type MovementPatch struct {
	// Description replaces the movement description when set.
	Description *string
	// Value replaces the movement value when set.
	Value *Money
	// Type replaces the movement type when set.
	Type *MovementType
	// Date replaces the movement date when set.
	Date *time.Time
}
