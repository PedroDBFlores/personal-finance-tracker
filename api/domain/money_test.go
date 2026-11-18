package domain_test

import (
	"testing"

	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/stretchr/testify/require"
)

func TestMoney(t *testing.T) {
	var m domain.Money = 100
	require.Equal(t, "1.00", m.String())
}
