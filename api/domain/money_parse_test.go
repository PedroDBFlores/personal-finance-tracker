package domain_test

import (
	"testing"

	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/stretchr/testify/require"
)

func TestParseMoney(t *testing.T) {
	tests := map[string]struct {
		in      string
		want    domain.Money
		wantErr bool
	}{
		"whole units":       {"1234", domain.Money(123400), false},
		"with dot":          {"12.34", domain.Money(1234), false},
		"one decimal":       {"12.3", domain.Money(1230), false},
		"zero":              {"0", domain.Money(0), false},
		"zero with dot":     {"0.00", domain.Money(0), false},
		"large":             {"1000000.99", domain.Money(100000099), false},
		"whitespace":        {" 12.34 ", domain.Money(1234), false},
		"empty":             {"", 0, true},
		"not a number":      {"abc", 0, true},
		"too many decimals": {"1.234", 0, true},
		"two dots":          {"1.2.3", 0, true},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := domain.ParseMoney(tt.in)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
