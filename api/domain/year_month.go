package domain

// YearMonth identifies a calendar month. Month is 1-based (1 = January).
type YearMonth struct {
	// Year is the full calendar year (e.g., 2025).
	Year int
	// Month is the calendar month, 1-12.
	Month int
}

// NewYearMonth creates a YearMonth, validating that month is in 1-12.
func NewYearMonth(year, month int) (YearMonth, error) {
	if month < 1 || month > 12 {
		return YearMonth{}, ErrInvalidMonth
	}
	return YearMonth{Year: year, Month: month}, nil
}
