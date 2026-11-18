package dtos

import (
	"time"

	"github.com/pedrodbflores/personal-finance-tracker/domain"
)

type CreateMovementDto struct {
	// Description defines what the movement was about (salary, bills, etc...)
	Description string
	// Value's last two digits are always the cents
	Value domain.Money
	// Type defines the type of movement
	Type domain.MovementType
	// Date defines the movement's date
	Date time.Time
}

func NewCreateMovementDto(description string, value domain.Money, movementType domain.MovementType, date time.Time) CreateMovementDto {
	return CreateMovementDto{description, value, movementType, date}
}

type MovementDto struct {
	domain.Movement
}

func NewMovementDto(m domain.Movement) MovementDto {
	return MovementDto{m}
}
