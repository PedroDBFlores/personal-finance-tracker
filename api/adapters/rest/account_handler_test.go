package rest_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pedrodbflores/personal-finance-tracker/adapters/rest"
	"github.com/pedrodbflores/personal-finance-tracker/app_errors"
	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/usecases/account"
	"github.com/pedrodbflores/personal-finance-tracker/usecases/usecasesfakes"
	"github.com/stretchr/testify/require"
)

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

// newAccountRouter builds the account routes with real usecases backed by fake ports.
func newAccountRouter(creator *usecasesfakes.FakeAccountCreator, getter *usecasesfakes.FakeAccountGetter, lister *usecasesfakes.FakeAccountLister, updater *usecasesfakes.FakeAccountUpdater, deleter *usecasesfakes.FakeAccountDeleter, balance *account.GetAccountBalanceUseCase) http.Handler {
	logger := testLogger()
	createUC := account.NewCreateAccountUseCase(creator, logger)
	getUC := account.NewGetAccountUseCase(getter, logger)
	listUC := account.NewListAccountsUseCase(lister, logger)
	updateUC := account.NewUpdateAccountUseCase(updater, logger)
	deleteUC := account.NewDeleteAccountUseCase(deleter, logger)
	h := rest.NewAccountHandler(&createUC, &getUC, &listUC, &updateUC, &deleteUC, balance, logger)
	mux := http.NewServeMux()
	h.Register(mux, "/accounts")
	return mux
}

func sampleAccount() domain.Account {
	return domain.NewAccount(1, "Bank A", domain.AccountTypeBank, domain.CurrencyEUR, time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC))
}

func doRequest(t *testing.T, handler http.Handler, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		require.NoError(t, json.NewEncoder(&buf).Encode(body))
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

func TestAccountHandlerCreate(t *testing.T) {
	t.Run("creates an account and returns 201", func(t *testing.T) {
		creator := &usecasesfakes.FakeAccountCreator{}
		id := domain.Id(5)
		creator.CreateReturns(&id, nil)
		r := newAccountRouter(creator, nil, nil, nil, nil, nil)

		rec := doRequest(t, r, http.MethodPost, "/accounts", map[string]string{
			"name": "Bank A", "type": "Bank", "currency": "EUR",
		})

		require.Equal(t, http.StatusCreated, rec.Code)
		_, modelAccount := creator.CreateArgsForCall(0)
		require.Equal(t, "Bank A", modelAccount.Name)
		require.Equal(t, domain.AccountTypeBank, modelAccount.Type)
		require.Equal(t, domain.CurrencyEUR, modelAccount.Currency)

		var resp map[string]any
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		require.Equal(t, "5", resp["id"])
		require.Equal(t, "Bank A", resp["name"])
	})

	t.Run("returns 400 when name is missing", func(t *testing.T) {
		creator := &usecasesfakes.FakeAccountCreator{}
		r := newAccountRouter(creator, nil, nil, nil, nil, nil)

		rec := doRequest(t, r, http.MethodPost, "/accounts", map[string]string{"name": "", "type": "Bank", "currency": "EUR"})
		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.Equal(t, 0, creator.CreateCallCount())
	})

	t.Run("returns 500 when the usecase fails", func(t *testing.T) {
		creator := &usecasesfakes.FakeAccountCreator{}
		creator.CreateReturns(nil, errors.New("db down"))
		r := newAccountRouter(creator, nil, nil, nil, nil, nil)

		rec := doRequest(t, r, http.MethodPost, "/accounts", map[string]string{"name": "Bank A", "type": "Bank", "currency": "EUR"})
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestAccountHandlerGet(t *testing.T) {
	t.Run("returns the account", func(t *testing.T) {
		getter := &usecasesfakes.FakeAccountGetter{}
		acc := sampleAccount()
		getter.GetReturns(&acc, nil)
		r := newAccountRouter(nil, getter, nil, nil, nil, nil)

		rec := doRequest(t, r, http.MethodGet, "/accounts/1", nil)

		require.Equal(t, http.StatusOK, rec.Code)
		_, calledId := getter.GetArgsForCall(0)
		require.Equal(t, domain.Id(1), calledId)
		var resp map[string]any
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		require.Equal(t, "Bank A", resp["name"])
		require.Equal(t, "Bank", resp["type"])
		require.Equal(t, "EUR", resp["currency"])
	})

	t.Run("returns 404 when not found", func(t *testing.T) {
		getter := &usecasesfakes.FakeAccountGetter{}
		getter.GetReturns(nil, app_errors.NewResourceNotFoundError("Account", "9"))
		r := newAccountRouter(nil, getter, nil, nil, nil, nil)

		rec := doRequest(t, r, http.MethodGet, "/accounts/9", nil)
		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("returns 400 on non-numeric id", func(t *testing.T) {
		getter := &usecasesfakes.FakeAccountGetter{}
		r := newAccountRouter(nil, getter, nil, nil, nil, nil)

		rec := doRequest(t, r, http.MethodGet, "/accounts/abc", nil)
		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.Equal(t, 0, getter.GetCallCount())
	})
}

func TestAccountHandlerList(t *testing.T) {
	t.Run("lists accounts", func(t *testing.T) {
		lister := &usecasesfakes.FakeAccountLister{}
		lister.ListReturns([]domain.Account{sampleAccount()}, nil)
		r := newAccountRouter(nil, nil, lister, nil, nil, nil)

		rec := doRequest(t, r, http.MethodGet, "/accounts", nil)

		require.Equal(t, http.StatusOK, rec.Code)
		var resp []map[string]any
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		require.Len(t, resp, 1)
		require.Equal(t, "Bank A", resp[0]["name"])
	})
}

func TestAccountHandlerUpdate(t *testing.T) {
	t.Run("applies a partial update", func(t *testing.T) {
		updater := &usecasesfakes.FakeAccountUpdater{}
		updater.UpdateReturns(nil)
		r := newAccountRouter(nil, nil, nil, updater, nil, nil)

		rec := doRequest(t, r, http.MethodPatch, "/accounts/3", map[string]string{"name": "New Name"})

		require.Equal(t, http.StatusNoContent, rec.Code)
		_, calledId, patch := updater.UpdateArgsForCall(0)
		require.Equal(t, domain.Id(3), calledId)
		require.NotNil(t, patch.Name)
		require.Equal(t, "New Name", *patch.Name)
		require.Nil(t, patch.Currency)
	})

	t.Run("returns 404 when the account doesn't exist", func(t *testing.T) {
		updater := &usecasesfakes.FakeAccountUpdater{}
		updater.UpdateReturns(app_errors.NewResourceNotFoundError("Account", "3"))
		r := newAccountRouter(nil, nil, nil, updater, nil, nil)

		rec := doRequest(t, r, http.MethodPatch, "/accounts/3", map[string]string{"name": "X"})
		require.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestAccountHandlerDelete(t *testing.T) {
	t.Run("deletes an account", func(t *testing.T) {
		deleter := &usecasesfakes.FakeAccountDeleter{}
		deleter.DeleteReturns(nil)
		r := newAccountRouter(nil, nil, nil, nil, deleter, nil)

		rec := doRequest(t, r, http.MethodDelete, "/accounts/4", nil)

		require.Equal(t, http.StatusNoContent, rec.Code)
		_, calledId := deleter.DeleteArgsForCall(0)
		require.Equal(t, domain.Id(4), calledId)
	})

	t.Run("returns 404 when the account doesn't exist", func(t *testing.T) {
		deleter := &usecasesfakes.FakeAccountDeleter{}
		deleter.DeleteReturns(app_errors.NewResourceNotFoundError("Account", "4"))
		r := newAccountRouter(nil, nil, nil, nil, deleter, nil)

		rec := doRequest(t, r, http.MethodDelete, "/accounts/4", nil)
		require.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestAccountHandlerBalance(t *testing.T) {
	now := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	money := func(s string) domain.Money {
		m, err := domain.ParseMoney(s)
		require.NoError(t, err)
		return m
	}

	balanceRouter := func(getter *usecasesfakes.FakeAccountGetter, movements *usecasesfakes.FakeMovementLister, investments *usecasesfakes.FakeInvestmentLister, txs *usecasesfakes.FakeInvestmentTransactionLister) http.Handler {
		logger := testLogger()
		balanceUC := account.NewGetAccountBalanceUseCase(getter, movements, investments, txs, logger)
		return newAccountRouter(nil, nil, nil, nil, nil, &balanceUC)
	}

	t.Run("returns the combined balance", func(t *testing.T) {
		getter := &usecasesfakes.FakeAccountGetter{}
		movements := &usecasesfakes.FakeMovementLister{}
		investments := &usecasesfakes.FakeInvestmentLister{}
		txs := &usecasesfakes.FakeInvestmentTransactionLister{}

		acc := sampleAccount()
		getter.GetReturns(&acc, nil)
		movements.ListReturns([]domain.Movement{
			domain.NewMovement(1, "Salary", money("1000.00"), domain.MovementTypeCredit, now),
			domain.NewMovement(1, "Groceries", money("25.99"), domain.MovementTypeDebit, now),
		}, nil)
		inv := domain.NewInvestment(2, 1, "Vanguard ETF", domain.ETF, domain.CurrencyEUR, now, nil)
		investments.ListReturns([]domain.Investment{inv}, nil)
		txs.ListReturns([]domain.InvestmentTransaction{
			domain.NewInvestmentTransaction(10, 2, domain.Buy, 0, money("10.00"), money("100.00"), now, "Initial Buy"),
		}, nil)

		rec := doRequest(t, balanceRouter(getter, movements, investments, txs), http.MethodGet, "/accounts/1/balance", nil)

		require.Equal(t, http.StatusOK, rec.Code)
		_, calledId := getter.GetArgsForCall(0)
		require.Equal(t, domain.Id(1), calledId)
		var resp map[string]any
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		require.Equal(t, "1", resp["account_id"])
		require.Equal(t, "EUR", resp["currency"])
		require.Equal(t, "974.01", resp["movements_balance"])
		require.Equal(t, "1000.00", resp["investments_balance"])
		require.Equal(t, "1974.01", resp["total_balance"])
	})

	t.Run("returns 404 when the account does not exist", func(t *testing.T) {
		getter := &usecasesfakes.FakeAccountGetter{}
		getter.GetReturns(nil, app_errors.NewResourceNotFoundError("Account", "9"))
		r := balanceRouter(getter, &usecasesfakes.FakeMovementLister{}, &usecasesfakes.FakeInvestmentLister{}, &usecasesfakes.FakeInvestmentTransactionLister{})

		rec := doRequest(t, r, http.MethodGet, "/accounts/9/balance", nil)

		require.Equal(t, http.StatusNotFound, rec.Code)
		var resp map[string]any
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		require.Contains(t, resp["error"], "Account")
	})

	t.Run("returns 400 for an invalid id", func(t *testing.T) {
		r := balanceRouter(&usecasesfakes.FakeAccountGetter{}, &usecasesfakes.FakeMovementLister{}, &usecasesfakes.FakeInvestmentLister{}, &usecasesfakes.FakeInvestmentTransactionLister{})

		rec := doRequest(t, r, http.MethodGet, "/accounts/abc/balance", nil)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
