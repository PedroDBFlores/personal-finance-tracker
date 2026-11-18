package main

import (
	"context"
	"errors"
	"testing"

	tui "github.com/grindlemire/go-tui"
	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/dtos"
	"github.com/stretchr/testify/require"
)

type fakeAccountCreator struct {
	id  *domain.Id
	err error

	called   bool
	lastDto  dtos.CreateAccountDto
}

func (f *fakeAccountCreator) Execute(_ context.Context, dto dtos.CreateAccountDto) (*domain.Id, error) {
	f.called = true
	f.lastDto = dto
	return f.id, f.err
}

func TestAccountFormSubmit(t *testing.T) {
	t.Run("creates the account and calls onDone", func(t *testing.T) {
		id := domain.Id(42)
		creator := &fakeAccountCreator{id: &id}
		done := false
		s := newAccountFormScreen(creator, func() { done = true })
		s.name.Set("New Bank")
		s.accType.Set(domain.AccountTypeBank)
		s.currency.Set(domain.CurrencyEUR)

		queue := make(chan func(), 1)
		s.submit(queue)
		fn := <-queue
		fn()

		require.True(t, creator.called)
		require.Equal(t, "New Bank", creator.lastDto.Name)
		require.Equal(t, domain.AccountTypeBank, creator.lastDto.Type)
		require.Equal(t, domain.CurrencyEUR, creator.lastDto.Currency)
		require.False(t, s.saving.Get())
		require.NoError(t, s.err.Get())
		require.True(t, done, "onDone should be called after a successful create")
	})

	t.Run("rejects an empty name without calling the creator", func(t *testing.T) {
		creator := &fakeAccountCreator{}
		s := newAccountFormScreen(creator, func() {})

		queue := make(chan func(), 1)
		s.submit(queue)

		require.False(t, creator.called)
		require.EqualError(t, s.err.Get(), "name is required")
		require.Len(t, queue, 0)
	})

	t.Run("surfaces creator errors and does not call onDone", func(t *testing.T) {
		creator := &fakeAccountCreator{err: errors.New("db write failed")}
		done := false
		s := newAccountFormScreen(creator, func() { done = true })
		s.name.Set("New Bank")

		queue := make(chan func(), 1)
		s.submit(queue)
		(<-queue)()

		require.False(t, s.saving.Get())
		require.EqualError(t, s.err.Get(), "db write failed")
		require.False(t, done)
	})
}

func TestAccountFormKeyMap(t *testing.T) {
	s := newAccountFormScreen(&fakeAccountCreator{}, func() {})
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

	// Tab moves focus between fields; arrows cycle the focused selector.
	assertBound(tui.KeyEvent{Key: tui.KeyTab})
	assertBound(tui.KeyEvent{Key: tui.KeyLeft})
	assertBound(tui.KeyEvent{Key: tui.KeyRight})
	assertBound(tui.KeyEvent{Key: tui.KeyEscape})
}

func TestAccountFormSelectorKeys(t *testing.T) {
	s := newAccountFormScreen(&fakeAccountCreator{}, func() {})
	typeEl := tui.New()
	currencyEl := tui.New()
	s.typeRef.Set(typeEl)
	s.currencyRef.Set(currencyEl)

	// Left/Right on the type selector cycle the account type.
	s.selectorKey(typeEl, tui.KeyLeft)
	require.Equal(t, domain.AccountTypeCreditCard, s.accType.Get())
	s.selectorKey(typeEl, tui.KeyRight)
	require.Equal(t, domain.AccountTypeBank, s.accType.Get())

	// Left/Right on the currency selector cycle the currency.
	s.selectorKey(currencyEl, tui.KeyRight)
	require.Equal(t, domain.CurrencyEUR, s.currency.Get())

	// Left/Right on any other element (e.g. the name input) do nothing.
	s.selectorKey(nil, tui.KeyRight)
	require.Equal(t, domain.AccountTypeBank, s.accType.Get())
	require.Equal(t, domain.CurrencyEUR, s.currency.Get())
}

func TestAccountFormCycles(t *testing.T) {
	s := newAccountFormScreen(&fakeAccountCreator{}, func() {})

	s.cycleType(1)
	require.Equal(t, domain.AccountTypeBrokerage, s.accType.Get())
	s.cycleType(1)
	require.Equal(t, domain.AccountTypeCreditCard, s.accType.Get())
	s.cycleType(1)
	require.Equal(t, domain.AccountTypeBank, s.accType.Get())

	s.cycleCurrency(1)
	require.Equal(t, domain.CurrencyEUR, s.currency.Get())
	s.cycleCurrency(1)
	require.Equal(t, domain.CurrencyGBP, s.currency.Get())
	s.cycleCurrency(1)
	require.Equal(t, domain.CurrencyUSD, s.currency.Get())
}
