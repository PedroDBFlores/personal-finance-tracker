package domain

// Tag represents a label for categorizing movements (e.g., "Groceries", "Utilities")
type Tag struct {
	// Id represents the unique tag identifier.
	Id Id
	// Name defines the tag name (e.g., "Food", "Bills")
	Name string
}

// NewTag creates a new Tag based on the provided parameters
func NewTag(id Id, name string) Tag {
	return Tag{
		Id:   id,
		Name: name,
	}
}
