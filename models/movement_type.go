package models

// MovementType defines the type of movement (either credit or debit)
type MovementType string

const (
	Credit MovementType = "Credit"
	Debit  MovementType = "Debit"
)
