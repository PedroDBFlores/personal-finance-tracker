package domain

import (
	"errors"
	"fmt"
)

// ErrInvalidMonth is returned when a month is outside the 1-12 range.
var ErrInvalidMonth = errors.New("invalid month: must be between 1 and 12")

// ValidationError is returned by constructors when an input value is invalid.
type ValidationError struct {
	Msg string
}

func (e ValidationError) Error() string { return e.Msg }

func newValidationError(format string, args ...any) ValidationError {
	return ValidationError{Msg: fmt.Sprintf(format, args...)}
}
