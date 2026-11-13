package domain

import (
	"testing"
	"time"
)

func TestNewMovement(t *testing.T) {
	accountId := NewId()
	description := "Test Movement"
	value := Money(10000)
	movementType := MovementTypeCredit
	date := time.Now()

	movement := NewMovement(accountId, description, value, movementType, date)

	if movement.AccountId != accountId {
		t.Errorf("Expected AccountId %s, got %s", accountId.String(), movement.AccountId.String())
	}

	if movement.Description != description {
		t.Errorf("Expected Description %s, got %s", description, movement.Description)
	}
	if movement.Value != value {
		t.Errorf("Expected Value %v, got %v", value, movement.Value)
	}
	if movement.Type != movementType {
		t.Errorf("Expected Type %s, got %s", movementType, movement.Type)
	}
	if !movement.Date.Equal(date) {
		t.Errorf("Expected Date %v, got %v", date, movement.Date)
	}
}

func TestMovementDefaults(t *testing.T) {
	id := NewId()
	accountId := NewId()
	description := "Test Movement"
	value := Money(0)
	movementType := MovementTypeDebit
	date := time.Time{}

	movement := Movement{
		AccountId:   accountId,
		Id:          id,
		Description: description,
		Value:       value,
		Type:        movementType,
		Date:        date,
	}

	if movement.Id != id {
		t.Errorf("Expected Id %v, got %v", id, movement.Id)
	}
	if movement.Description != description {
		t.Errorf("Expected Description %s, got %s", description, movement.Description)
	}
	if movement.Value != value {
		t.Errorf("Expected Value %v, got %v", value, movement.Value)
	}
	if movement.Type != movementType {
		t.Errorf("Expected Type %s, got %s", movementType, movement.Type)
	}
	if !movement.Date.Equal(date) {
		t.Errorf("Expected Date %v, got %v", date, movement.Date)
	}
}
