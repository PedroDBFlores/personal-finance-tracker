package main

import (
	"context"
	"fmt"

	tui "github.com/grindlemire/go-tui"
	"github.com/pedrodbflores/personal-finance-tracker/dtos"
)

//go:generate go run github.com/grindlemire/go-tui/cmd/tui generate accounts.gsx

// accountLister is the part of the list-accounts usecase the screen needs.
type accountLister interface {
	Execute(ctx context.Context) ([]dtos.AccountDto, error)
}

// accountsScreen is the account list screen component.
type accountsScreen struct {
	lister   accountLister
	creator  accountCreator
	accounts *tui.State[[]dtos.AccountDto]
	err      *tui.State[error]
	loading  *tui.State[bool]
}

func newAccountsScreen(lister accountLister, creator accountCreator) *accountsScreen {
	return &accountsScreen{
		lister:   lister,
		creator:  creator,
		accounts: tui.NewState([]dtos.AccountDto{}),
		err:      tui.NewState[error](nil),
		loading:  tui.NewState(true),
	}
}

// load fetches the accounts in the background and posts the result to the
// app's event queue, so state is only mutated on the event loop.
func (s *accountsScreen) load(queue chan<- func()) {
	s.loading.Set(true)
	go func() {
		accounts, err := s.lister.Execute(context.Background())
		queue <- func() { s.applyResult(accounts, err) }
	}()
}

func (s *accountsScreen) applyResult(accounts []dtos.AccountDto, err error) {
	s.loading.Set(false)
	if err != nil {
		s.err.Set(err)
		return
	}
	s.err.Set(nil)
	s.accounts.Set(accounts)
}

func (s *accountsScreen) KeyMap() tui.KeyMap {
	return tui.KeyMap{
		tui.OnStop(tui.Rune('q'), func(ke tui.KeyEvent) { ke.App().Stop() }),
		tui.On(tui.Rune('r'), func(ke tui.KeyEvent) { s.load(ke.App().EventQueue()) }),
		tui.On(tui.Rune('c'), func(ke tui.KeyEvent) {
			app := ke.App()
			form := newAccountFormScreen(s.creator, func() {
				app.SetRootComponent(s)
				s.load(app.EventQueue())
			})
			app.SetRootComponent(form)
		}),
	}
}

templ (s *accountsScreen) Render() {
	<div class="flex-col gap-1 p-2">
		<span class="font-bold">Accounts</span>
		<hr />
		if s.loading.Get() {
			<span class="font-dim">Loading...</span>
		} else if s.err.Get() != nil {
			<span class="text-red">{s.err.Get().Error()}</span>
		} else if len(s.accounts.Get()) == 0 {
			<span class="font-dim">No accounts yet.</span>
		} else {
			for _, a := range s.accounts.Get() {
				<span>{fmt.Sprintf("- %s (%s, %s)", a.Name, a.Type, a.Currency)}</span>
			}
		}
		<br />
		<span class="font-dim">c new account, r reload, q quit</span>
	</div>
}
