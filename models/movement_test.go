package models

import (
	"math/big"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewMovement(t *testing.T) {
	id := uuid.New()
	description := "Test Movement"
	value := big.NewInt(10000) // $100.00
	movementType := Credit
	date := time.Now()
	accountId := uuid.New()
	tags := []uuid.UUID{uuid.New(), uuid.New()}

	movement := NewMovement(id, description, *value, movementType, date, accountId, tags)

	if movement.Id != id {
		t.Errorf("Expected Id %v, got %v", id, movement.Id)
	}
	if movement.Description != description {
		t.Errorf("Expected Description %s, got %s", description, movement.Description)
	}
	if movement.Value.Cmp(value) != 0 {
		t.Errorf("Expected Value %v, got %v", value, movement.Value)
	}
	if movement.Type != movementType {
		t.Errorf("Expected Type %s, got %s", movementType, movement.Type)
	}
	if !movement.Date.Equal(date) {
		t.Errorf("Expected Date %v, got %v", date, movement.Date)
	}
	if movement.AccountId != accountId {
		t.Errorf("Expected AccountId %v, got %v", accountId, movement.AccountId)
	}
	if len(movement.Tags) != len(tags) {
		t.Errorf("Expected Tags length %d, got %d", len(tags), len(movement.Tags))
	}
	for i, tag := range tags {
		if movement.Tags[i] != tag {
			t.Errorf("Expected Tag at index %d to be %v, got %v", i, tag, movement.Tags[i])
		}
	}
}

func TestMovementDefaults(t *testing.T) {
	id := uuid.New()
	description := "Test Movement"
	value := big.NewInt(0)
	movementType := Debit
	date := time.Time{}
	accountId := uuid.Nil
	tags := []uuid.UUID{}

	movement := Movement{
		Id:          id,
		Description: description,
		Value:       *value,
		Type:        movementType,
		Date:        date,
		AccountId:   accountId,
		Tags:        tags,
	}

	if movement.Id != id {
		t.Errorf("Expected Id %v, got %v", id, movement.Id)
	}
	if movement.Description != description {
		t.Errorf("Expected Description %s, got %s", description, movement.Description)
	}
	if movement.Value.Cmp(value) != 0 {
		t.Errorf("Expected Value %v, got %v", value, movement.Value)
	}
	if movement.Type != movementType {
		t.Errorf("Expected Type %s, got %s", movementType, movement.Type)
	}
	if movement.AccountId != accountId {
		t.Errorf("Expected AccountId %v, got %v", accountId, movement.AccountId)
	}
	if len(movement.Tags) != len(tags) {
		t.Errorf("Expected Tags length %d, got %d", len(tags), len(movement.Tags))
	}
}
