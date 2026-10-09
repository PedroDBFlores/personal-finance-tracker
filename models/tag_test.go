package models

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewTag(t *testing.T) {
	id := uuid.New()
	name := "Test Tag"

	tag := NewTag(id, name)

	if tag.Id != id {
		t.Errorf("Expected Id %v, got %v", id, tag.Id)
	}
	if tag.Name != name {
		t.Errorf("Expected Name %s, got %s", name, tag.Name)
	}
}

func TestTagDefaults(t *testing.T) {
	id := uuid.Nil
	name := ""

	tag := Tag{
		Id:   id,
		Name: name,
	}

	if tag.Id != id {
		t.Errorf("Expected Id %v, got %v", id, tag.Id)
	}
	if tag.Name != name {
		t.Errorf("Expected Name %s, got %s", name, tag.Name)
	}
}
