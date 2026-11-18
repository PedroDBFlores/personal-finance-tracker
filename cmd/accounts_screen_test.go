package main

import (
	"context"
	"errors"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/dtos"
	"github.com/stretchr/testify/require"
)

type fakeAccountLister struct {
	accounts []dtos.AccountDto
	err      error
}

func (f *fakeAccountLister) Execute(_ context.Context) ([]dtos.AccountDto, error) {
	return f.accounts, f.err
}

func sampleAccounts() []dtos.AccountDto {
	return []dtos.AccountDto{
		dtos.NewAccountDto(domain.NewAccount(1, "Bank A", domain.AccountTypeBank, domain.CurrencyEUR, time.Now())),
		dtos.NewAccountDto(domain.NewAccount(2, "Broker B", domain.AccountTypeBrokerage, domain.CurrencyUSD, time.Now())),
	}
}

func TestAccountsScreen(t *testing.T) {
	t.Run("loads accounts on init and renders them", func(t *testing.T) {
		lister := &fakeAccountLister{accounts: sampleAccounts()}
		m := newAccountsScreen(lister)

		cmd := m.Init()
		require.NotNil(t, cmd)
		msg := cmd()
		loaded, ok := msg.(accountsLoadedMsg)
		require.True(t, ok, "expected accountsLoadedMsg, got %T", msg)
		require.Equal(t, lister.accounts, loaded.accounts)

		updated, _ := m.Update(loaded)
		screen, ok := updated.(accountsScreen)
		require.True(t, ok)
		require.False(t, screen.loading)
		require.NoError(t, screen.err)
		require.Len(t, screen.accounts, 2)
		require.Contains(t, screen.View().Content, "Bank A")
		require.Contains(t, screen.View().Content, "Broker B")
	})

	t.Run("surfaces load errors", func(t *testing.T) {
		lister := &fakeAccountLister{err: errors.New("db down")}
		m := newAccountsScreen(lister)

		msg := m.Init()()
		failed, ok := msg.(loadFailedMsg)
		require.True(t, ok, "expected loadFailedMsg, got %T", msg)

		updated, _ := m.Update(failed)
		screen := updated.(accountsScreen)
		require.False(t, screen.loading)
		require.EqualError(t, screen.err, "db down")
		require.Contains(t, screen.View().Content, "db down")
	})

	t.Run("reloads on r key", func(t *testing.T) {
		lister := &fakeAccountLister{accounts: sampleAccounts()}
		m := newAccountsScreen(lister)

		updated, cmd := m.Update(tea.KeyPressMsg{Code: 'r'})
		screen := updated.(accountsScreen)
		require.True(t, screen.loading)
		require.NotNil(t, cmd)
	})

	t.Run("quits on q key", func(t *testing.T) {
		m := newAccountsScreen(&fakeAccountLister{})
		_, cmd := m.Update(tea.KeyPressMsg{Code: 'q'})
		require.NotNil(t, cmd)
	})
}
