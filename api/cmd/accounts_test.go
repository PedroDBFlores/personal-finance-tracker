package main

import (
	"context"
	"errors"
	"testing"
	"time"

	tui "github.com/grindlemire/go-tui"
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

func TestAccountsScreenLoad(t *testing.T) {
	t.Run("applies loaded accounts to state", func(t *testing.T) {
		accounts := sampleAccounts()
		s := newAccountsScreen(&fakeAccountLister{accounts: accounts}, &fakeAccountCreator{})
		queue := make(chan func(), 1)

		s.load(queue)
		fn := <-queue // background fetch finished
		fn()          // apply on the event loop

		require.False(t, s.loading.Get())
		require.NoError(t, s.err.Get())
		require.Equal(t, accounts, s.accounts.Get())
	})

	t.Run("applies load errors to state", func(t *testing.T) {
		s := newAccountsScreen(&fakeAccountLister{err: errors.New("db down")}, &fakeAccountCreator{})
		queue := make(chan func(), 1)

		s.load(queue)
		(<-queue)()

		require.False(t, s.loading.Get())
		require.EqualError(t, s.err.Get(), "db down")
	})
}

func TestAccountsScreenKeyMap(t *testing.T) {
	s := newAccountsScreen(&fakeAccountLister{}, &fakeAccountCreator{})
	km := s.KeyMap()

	assertBound := func(key tui.KeyEvent) {
		t.Helper()
		for _, b := range km {
			if key.IsRune() && b.Pattern.Rune == key.Rune {
				return
			}
			if !key.IsRune() && b.Pattern.Key == key.Key {
				return
			}
		}
		t.Errorf("no key binding for %v", key)
	}

	assertBound(tui.KeyEvent{Key: tui.KeyRune, Rune: 'r'})
	assertBound(tui.KeyEvent{Key: tui.KeyRune, Rune: 'q'})
}
