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

// MovementMonthDto is a year/month pair that has movements on an account.
type MovementMonthDto struct {
	// Year is the full calendar year (e.g., 2025).
	Year int
	// Month is the calendar month, 1-12.
	Month int
}

func NewMovementMonthDto(ym domain.YearMonth) MovementMonthDto {
	return MovementMonthDto{Year: ym.Year, Month: ym.Month}
}
