package main

import (
	"context"
	"errors"
	"time"

	tui "github.com/grindlemire/go-tui"
	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/dtos"
)

//go:generate go run github.com/grindlemire/go-tui/cmd/tui generate account_form.gsx

// accountCreator is the part of the create-account usecase the form needs.
type accountCreator interface {
	Execute(ctx context.Context, dto dtos.CreateAccountDto) (*domain.Id, error)
}

var accountTypes = []domain.AccountType{
	domain.AccountTypeBank,
	domain.AccountTypeBrokerage,
	domain.AccountTypeCreditCard,
}

var currencies = []domain.Currency{
	domain.CurrencyUSD,
	domain.CurrencyEUR,
	domain.CurrencyGBP,
}

// accountFormScreen is the create-account form. The name is a text input; type
// and currency are focusable selectors cycled with Left/Right. Tab moves focus.
type accountFormScreen struct {
	creator     accountCreator
	onDone      func()
	app         *tui.App
	name        *tui.State[string]
	accType     *tui.State[domain.AccountType]
	currency    *tui.State[domain.Currency]
	err         *tui.State[error]
	saving      *tui.State[bool]
	typeRef     *tui.Ref
	currencyRef *tui.Ref
}

func newAccountFormScreen(creator accountCreator, onDone func()) *accountFormScreen {
	return &accountFormScreen{
		creator:     creator,
		onDone:      onDone,
		name:        tui.NewState(""),
		accType:     tui.NewState(domain.AccountTypeBank),
		currency:    tui.NewState(domain.CurrencyUSD),
		err:         tui.NewState[error](nil),
		saving:      tui.NewState(false),
		typeRef:     tui.NewRef(),
		currencyRef: tui.NewRef(),
	}
}

// BindApp captures the app so the form can post async results to its event queue.
func (s *accountFormScreen) BindApp(app *tui.App) { s.app = app }

func cycle[T comparable](options []T, current T, dir int) T {
	for i, o := range options {
		if o == current {
			next := (i + dir) % len(options)
			if next < 0 {
				next += len(options)
			}
			return options[next]
		}
	}
	return options[0]
}

func (s *accountFormScreen) cycleType(dir int)     { s.accType.Set(cycle(accountTypes, s.accType.Get(), dir)) }
func (s *accountFormScreen) cycleCurrency(dir int) { s.currency.Set(cycle(currencies, s.currency.Get(), dir)) }

// selectorKey cycles the type/currency selector when its element is focused.
func (s *accountFormScreen) selectorKey(focused *tui.Element, key tui.Key) {
	if focused == nil {
		return
	}
	dir := 1
	if key == tui.KeyLeft {
		dir = -1
	}
	switch focused {
	case s.typeRef.El():
		s.cycleType(dir)
	case s.currencyRef.El():
		s.cycleCurrency(dir)
	}
}

// submit validates the form, runs the creator in the background and posts the
// result to the app's event queue.
func (s *accountFormScreen) submit(queue chan<- func()) {
	if s.saving.Get() {
		return
	}
	if s.name.Get() == "" {
		s.err.Set(errors.New("name is required"))
		return
	}
	s.err.Set(nil)
	s.saving.Set(true)
	dto := dtos.NewCreateAccountDto(s.name.Get(), s.accType.Get(), s.currency.Get(), time.Now())
	go func() {
		id, err := s.creator.Execute(context.Background(), dto)
		queue <- func() { s.applyResult(id, err) }
	}()
}

func (s *accountFormScreen) applyResult(_ *domain.Id, err error) {
	s.saving.Set(false)
	if err != nil {
		s.err.Set(err)
		return
	}
	s.err.Set(nil)
	if s.onDone != nil {
		s.onDone()
	}
}

func (s *accountFormScreen) KeyMap() tui.KeyMap {
	return tui.KeyMap{
		tui.OnStop(tui.KeyEscape, func(ke tui.KeyEvent) {
			if s.onDone != nil {
				s.onDone()
			}
		}),
		// Enter saves from any field; when the name input is focused its own
		// submit binding fires first and stops propagation.
		tui.OnStop(tui.KeyEnter, func(ke tui.KeyEvent) { s.submit(ke.App().EventQueue()) }),
		tui.On(tui.KeyTab, func(ke tui.KeyEvent) { ke.App().FocusNext() }),
		tui.On(tui.KeyLeft, func(ke tui.KeyEvent) {
			if f, ok := ke.App().Focused().(*tui.Element); ok {
				s.selectorKey(f, tui.KeyLeft)
			}
		}),
		tui.On(tui.KeyRight, func(ke tui.KeyEvent) {
			if f, ok := ke.App().Focused().(*tui.Element); ok {
				s.selectorKey(f, tui.KeyRight)
			}
		}),
	}
}

templ (s *accountFormScreen) Render() {
	<div class="flex-col gap-1 p-2">
		<span class="font-bold">New Account</span>
		<hr />
		<div class="flex gap-2 items-center">
			<span class="font-dim">Name:</span>
			<input
				value={s.name}
				placeholder="Account name..."
				width={30}
				border={tui.BorderRounded}
				autoFocus={true}
				onSubmit={func(string) { s.submit(s.app.EventQueue()) }}
			/>
		</div>
		<div class="flex gap-2 items-center">
			<span class="font-dim">Type:</span>
			<span ref={s.typeRef} class="focusable text-cyan border-rounded px-1">{string(s.accType.Get())}</span>
			<span class="font-dim">(left/right to change)</span>
		</div>
		<div class="flex gap-2 items-center">
			<span class="font-dim">Currency:</span>
			<span ref={s.currencyRef} class="focusable text-cyan border-rounded px-1">{string(s.currency.Get())}</span>
			<span class="font-dim">(left/right to change)</span>
		</div>
		if s.err.Get() != nil {
			<span class="text-red">{s.err.Get().Error()}</span>
		}
		if s.saving.Get() {
			<span class="font-dim">Saving...</span>
		}
		<br />
		<span class="font-dim">Tab to move, Enter to save, Esc to cancel</span>
	</div>
}
