package domain

import "fmt"

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
