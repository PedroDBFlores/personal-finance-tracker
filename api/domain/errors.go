package domain

import "errors"

// ErrInvalidMonth is returned when a month is outside the 1-12 range.
var ErrInvalidMonth = errors.New("invalid month: must be between 1 and 12")
