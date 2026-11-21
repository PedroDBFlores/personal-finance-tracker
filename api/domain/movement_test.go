package domain

import (
	"errors"
	"testing"
	"time"
)

func TestNewMovement(t *testing.T) {
	accountId := NewId()
	description := "Test Movement"
	value := Money(10000)
	movementType := MovementTypeCredit
	date := time.Now()

	movement, err := NewMovement(accountId, description, value, movementType, date)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

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

func TestNewMovementValidation(t *testing.T) {
	now := time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)
	expectValidationError := func(t *testing.T, err error) {
		t.Helper()
		var ve ValidationError
		if !errors.As(err, &ve) {
			t.Fatalf("expected ValidationError, got %v", err)
		}
	}

	t.Run("accepts a valid credit", func(t *testing.T) {
		m, err := NewMovement(1, "Salary", Money(100000), MovementTypeCredit, now)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if m.Type != MovementTypeCredit {
			t.Errorf("Expected type Credit, got %s", m.Type)
		}
	})

	t.Run("rejects an empty description", func(t *testing.T) {
		_, err := NewMovement(1, "", Money(100), MovementTypeCredit, now)
		expectValidationError(t, err)
	})

	t.Run("rejects an unsupported movement type", func(t *testing.T) {
		_, err := NewMovement(1, "Salary", Money(100), MovementType("Transfer"), now)
		expectValidationError(t, err)
	})

	t.Run("rejects a negative value", func(t *testing.T) {
		_, err := NewMovement(1, "Refund", Money(-100), MovementTypeDebit, now)
		expectValidationError(t, err)
	})
}
