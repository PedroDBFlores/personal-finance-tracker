package models

import "uuid"

// Tag represents a label for categorizing movements (e.g., "Groceries", "Utilities")
type Tag struct {
	// Id represents the unique tag identifier.
	Id uuid.UUID
	// Name defines the tag name (e.g., "Food", "Bills")
	Name string
}

// NewTag creates a new Tag based on the provided parameters
func NewTag(id uuid.UUID, name string) Tag {
	return Tag{
		Id:   id,
		Name: name,
	}
}
