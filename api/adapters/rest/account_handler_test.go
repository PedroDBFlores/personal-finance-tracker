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
func newAccountRouter(creator *usecasesfakes.FakeAccountCreator, getter *usecasesfakes.FakeAccountGetter, lister *usecasesfakes.FakeAccountLister, updater *usecasesfakes.FakeAccountUpdater, deleter *usecasesfakes.FakeAccountDeleter) http.Handler {
	logger := testLogger()
	createUC := account.NewCreateAccountUseCase(creator, logger)
	getUC := account.NewGetAccountUseCase(getter, logger)
	listUC := account.NewListAccountsUseCase(lister, logger)
	updateUC := account.NewUpdateAccountUseCase(updater, logger)
	deleteUC := account.NewDeleteAccountUseCase(deleter, logger)
	h := rest.NewAccountHandler(&createUC, &getUC, &listUC, &updateUC, &deleteUC, logger)
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
		r := newAccountRouter(creator, nil, nil, nil, nil)

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
		r := newAccountRouter(creator, nil, nil, nil, nil)

		rec := doRequest(t, r, http.MethodPost, "/accounts", map[string]string{"name": "", "type": "Bank", "currency": "EUR"})
		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.Equal(t, 0, creator.CreateCallCount())
	})

	t.Run("returns 500 when the usecase fails", func(t *testing.T) {
		creator := &usecasesfakes.FakeAccountCreator{}
		creator.CreateReturns(nil, errors.New("db down"))
		r := newAccountRouter(creator, nil, nil, nil, nil)

		rec := doRequest(t, r, http.MethodPost, "/accounts", map[string]string{"name": "Bank A", "type": "Bank", "currency": "EUR"})
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestAccountHandlerGet(t *testing.T) {
	t.Run("returns the account", func(t *testing.T) {
		getter := &usecasesfakes.FakeAccountGetter{}
		acc := sampleAccount()
		getter.GetReturns(&acc, nil)
		r := newAccountRouter(nil, getter, nil, nil, nil)

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
		r := newAccountRouter(nil, getter, nil, nil, nil)

		rec := doRequest(t, r, http.MethodGet, "/accounts/9", nil)
		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("returns 400 on non-numeric id", func(t *testing.T) {
		getter := &usecasesfakes.FakeAccountGetter{}
		r := newAccountRouter(nil, getter, nil, nil, nil)

		rec := doRequest(t, r, http.MethodGet, "/accounts/abc", nil)
		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.Equal(t, 0, getter.GetCallCount())
	})
}

func TestAccountHandlerList(t *testing.T) {
	t.Run("lists accounts", func(t *testing.T) {
		lister := &usecasesfakes.FakeAccountLister{}
		lister.ListReturns([]domain.Account{sampleAccount()}, nil)
		r := newAccountRouter(nil, nil, lister, nil, nil)

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
		r := newAccountRouter(nil, nil, nil, updater, nil)

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
		r := newAccountRouter(nil, nil, nil, updater, nil)

		rec := doRequest(t, r, http.MethodPatch, "/accounts/3", map[string]string{"name": "X"})
		require.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestAccountHandlerDelete(t *testing.T) {
	t.Run("deletes an account", func(t *testing.T) {
		deleter := &usecasesfakes.FakeAccountDeleter{}
		deleter.DeleteReturns(nil)
		r := newAccountRouter(nil, nil, nil, nil, deleter)

		rec := doRequest(t, r, http.MethodDelete, "/accounts/4", nil)

		require.Equal(t, http.StatusNoContent, rec.Code)
		_, calledId := deleter.DeleteArgsForCall(0)
		require.Equal(t, domain.Id(4), calledId)
	})

	t.Run("returns 404 when the account doesn't exist", func(t *testing.T) {
		deleter := &usecasesfakes.FakeAccountDeleter{}
		deleter.DeleteReturns(app_errors.NewResourceNotFoundError("Account", "4"))
		r := newAccountRouter(nil, nil, nil, nil, deleter)

		rec := doRequest(t, r, http.MethodDelete, "/accounts/4", nil)
		require.Equal(t, http.StatusNotFound, rec.Code)
	})
}
