package main

import (
	"context"
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/pedrodbflores/personal-finance-tracker/dtos"
)

// accountLister is the part of the list-accounts usecase the screen needs.
type accountLister interface {
	Execute(ctx context.Context) ([]dtos.AccountDto, error)
}

type accountsLoadedMsg struct{ accounts []dtos.AccountDto }
type loadFailedMsg struct{ err error }

func loadAccounts(lister accountLister) tea.Cmd {
	return func() tea.Msg {
		accounts, err := lister.Execute(context.Background())
		if err != nil {
			return loadFailedMsg{err}
		}
		return accountsLoadedMsg{accounts}
	}
}

// accountsScreen is the account list screen.
type accountsScreen struct {
	lister   accountLister
	accounts []dtos.AccountDto
	err      error
	loading  bool
}

func newAccountsScreen(lister accountLister) accountsScreen {
	return accountsScreen{lister: lister, loading: true}
}

func (m accountsScreen) Init() tea.Cmd {
	return loadAccounts(m.lister)
}

func (m accountsScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case accountsLoadedMsg:
		m.accounts = msg.accounts
		m.loading = false
		return m, nil
	case loadFailedMsg:
		m.err = msg.err
		m.loading = false
		return m, nil
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "r":
			m.loading = true
			return m, loadAccounts(m.lister)
		}
	}
	return m, nil
}

func (m accountsScreen) View() tea.View {
	s := "Accounts\n\n"
	switch {
	case m.loading:
		s += "Loading...\n"
	case m.err != nil:
		s += fmt.Sprintf("Error: %v\n", m.err)
	case len(m.accounts) == 0:
		s += "No accounts yet.\n"
	default:
		for _, a := range m.accounts {
			s += fmt.Sprintf("- %s (%s, %s)\n", a.Name, a.Type, a.Currency)
		}
	}
	s += "\nPress r to reload, q to quit.\n"
	return tea.NewView(s)
}
