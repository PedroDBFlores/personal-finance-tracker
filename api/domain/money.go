package domain

import (
	"fmt"
	"strconv"
	"strings"
)

// Money represents a monetary amount in a specific currency, stored as an integer number of minor units (e.g., cents).
type Money int64

func (m Money) Add(other Money) Money {
	return m + other
}

func (m Money) Subtract(other Money) Money {
	return m - other
}

func (m Money) Multiply(quantity int) Money {
	return m * Money(quantity)
}

func (m Money) IsNegative() bool {
	return m < 0
}

func (m Money) String() string {
	curr := m / 100
	cents := m % 100
	if cents < 0 {
		cents = -cents
	}
	return fmt.Sprintf("%d.%02d", curr, cents)
}

// ParseMoney parses a decimal string like "12.34" into Money minor units (1234).
func ParseMoney(s string) (Money, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("amount is required")
	}
	parts := strings.Split(s, ".")
	switch len(parts) {
	case 1:
		whole, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid amount %q", s)
		}
		return Money(whole * 100), nil
	case 2:
		if len(parts[1]) > 2 {
			return 0, fmt.Errorf("invalid amount %q: at most 2 decimal places", s)
		}
		whole, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid amount %q", s)
		}
		frac, err := strconv.ParseInt(parts[1]+strings.Repeat("0", 2-len(parts[1])), 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid amount %q", s)
		}
		return Money(whole*100 + frac), nil
	default:
		return 0, fmt.Errorf("invalid amount %q", s)
	}
}
