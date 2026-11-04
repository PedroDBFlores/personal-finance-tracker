package app_errors

import (
	"fmt"
)

type ResourceNotFoundError struct {
	resource, id string
}

func NewResourceNotFoundError(resource, id string) error {
	return ResourceNotFoundError{resource, id}
}

func (e ResourceNotFoundError) Error() string {
	return fmt.Sprintf("resource %s with id %s not found", e.resource, e.id)
}
