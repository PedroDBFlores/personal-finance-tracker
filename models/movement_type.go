package models

// MovementType defines the type of movement (either credit or debit)
type MovementType string

const (
	MovementTypeCredit MovementType = "Credit"
	MovementTypeDebit  MovementType = "Debit"
)
