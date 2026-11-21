package rest_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/pedrodbflores/personal-finance-tracker/adapters/rest"
	"github.com/pedrodbflores/personal-finance-tracker/app_errors"
	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/usecases/movement"
	"github.com/pedrodbflores/personal-finance-tracker/usecases/usecasesfakes"
	"github.com/stretchr/testify/require"
)

// newMovementRouter builds the movement routes with real usecases backed by fake ports.
func newMovementRouter(creator *usecasesfakes.FakeMovementCreator, getter *usecasesfakes.FakeMovementGetter, lister *usecasesfakes.FakeMovementLister, monthLister *usecasesfakes.FakeMovementMonthLister, byMonthLister *usecasesfakes.FakeMovementByMonthLister, updater *usecasesfakes.FakeMovementUpdater, deleter *usecasesfakes.FakeMovementDeleter) http.Handler {
	logger := testLogger()
	createUC := movement.NewCreateMovementUseCase(creator, logger)
	getUC := movement.NewGetMovementUseCase(getter, logger)
	listUC := movement.NewListMovementsUseCase(lister, logger)
	accountGetter := &usecasesfakes.FakeAccountGetter{}
	accountGetter.GetReturns(&domain.Account{Id: 1}, nil)
	monthsUC := movement.NewListMovementMonthsUseCase(accountGetter, monthLister, logger)
	byMonthUC := movement.NewListMovementsByMonthUseCase(accountGetter, byMonthLister, logger)
	updateUC := movement.NewUpdateMovementUseCase(updater, logger)
	deleteUC := movement.NewDeleteMovementUseCase(deleter, logger)
	h := rest.NewMovementHandler(&createUC, &getUC, &listUC, &monthsUC, &byMonthUC, &updateUC, &deleteUC, logger)
	mux := http.NewServeMux()
	h.Register(mux, "/accounts/{accountId}/movements")
	return mux
}

func sampleMovement(accountId domain.Id) domain.Movement {
	m := domain.NewMovement(accountId, "Salary", domain.Money(100000), domain.MovementTypeCredit, time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC))
	m.Id = 1
	return m
}

func TestMovementHandlerCreate(t *testing.T) {
	t.Run("creates a movement under an account", func(t *testing.T) {
		creator := &usecasesfakes.FakeMovementCreator{}
		id := domain.Id(9)
		creator.CreateReturns(&id, nil)
		r := newMovementRouter(creator, nil, nil, nil, nil, nil, nil)

		rec := doRequest(t, r, http.MethodPost, "/accounts/2/movements", map[string]string{
			"description": "Salary", "value": "1000.00", "type": "Credit",
		})

		require.Equal(t, http.StatusCreated, rec.Code)
		_, calledAccountId, modelMovement := creator.CreateArgsForCall(0)
		require.Equal(t, domain.Id(2), calledAccountId)
		require.Equal(t, "Salary", modelMovement.Description)
		require.Equal(t, domain.Money(100000), modelMovement.Value)
		require.Equal(t, domain.MovementTypeCredit, modelMovement.Type)

		var resp map[string]any
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		require.Equal(t, "9", resp["id"])
		require.Equal(t, "2", resp["account_id"])
		require.Equal(t, "1000.00", resp["value"])
	})

	t.Run("returns 400 on invalid amount", func(t *testing.T) {
		creator := &usecasesfakes.FakeMovementCreator{}
		r := newMovementRouter(creator, nil, nil, nil, nil, nil, nil)

		rec := doRequest(t, r, http.MethodPost, "/accounts/2/movements", map[string]string{
			"description": "Salary", "value": "abc", "type": "Credit",
		})
		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.Equal(t, 0, creator.CreateCallCount())
	})

	t.Run("returns 400 when description is missing", func(t *testing.T) {
		creator := &usecasesfakes.FakeMovementCreator{}
		r := newMovementRouter(creator, nil, nil, nil, nil, nil, nil)

		rec := doRequest(t, r, http.MethodPost, "/accounts/2/movements", map[string]string{
			"description": "", "value": "10", "type": "Credit",
		})
		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.Equal(t, 0, creator.CreateCallCount())
	})
}

func TestMovementHandlerGet(t *testing.T) {
	t.Run("returns the movement", func(t *testing.T) {
		getter := &usecasesfakes.FakeMovementGetter{}
		mov := sampleMovement(2)
		getter.GetReturns(&mov, nil)
		r := newMovementRouter(nil, getter, nil, nil, nil, nil, nil)

		rec := doRequest(t, r, http.MethodGet, "/accounts/2/movements/1", nil)

		require.Equal(t, http.StatusOK, rec.Code)
		_, calledAccountId, calledId := getter.GetArgsForCall(0)
		require.Equal(t, domain.Id(2), calledAccountId)
		require.Equal(t, domain.Id(1), calledId)
		var resp map[string]any
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		require.Equal(t, "Salary", resp["description"])
	})

	t.Run("returns 404 when not found", func(t *testing.T) {
		getter := &usecasesfakes.FakeMovementGetter{}
		getter.GetReturns(nil, app_errors.NewResourceNotFoundError("Movement", "7"))
		r := newMovementRouter(nil, getter, nil, nil, nil, nil, nil)

		rec := doRequest(t, r, http.MethodGet, "/accounts/2/movements/7", nil)
		require.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestMovementHandlerList(t *testing.T) {
	t.Run("lists movements of an account", func(t *testing.T) {
		lister := &usecasesfakes.FakeMovementLister{}
		lister.ListReturns([]domain.Movement{sampleMovement(2)}, nil)
		r := newMovementRouter(nil, nil, lister, nil, nil, nil, nil)

		rec := doRequest(t, r, http.MethodGet, "/accounts/2/movements", nil)

		require.Equal(t, http.StatusOK, rec.Code)
		_, calledAccountId := lister.ListArgsForCall(0)
		require.Equal(t, domain.Id(2), calledAccountId)
		var resp []map[string]any
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		require.Len(t, resp, 1)
	})
}

func TestMovementHandlerUpdate(t *testing.T) {
	t.Run("applies a partial update", func(t *testing.T) {
		updater := &usecasesfakes.FakeMovementUpdater{}
		updater.UpdateReturns(nil)
		r := newMovementRouter(nil, nil, nil, nil, nil, updater, nil)

		rec := doRequest(t, r, http.MethodPatch, "/accounts/2/movements/1", map[string]string{"value": "25.99"})

		require.Equal(t, http.StatusNoContent, rec.Code)
		_, calledAccountId, calledId, patch := updater.UpdateArgsForCall(0)
		require.Equal(t, domain.Id(2), calledAccountId)
		require.Equal(t, domain.Id(1), calledId)
		require.NotNil(t, patch.Value)
		require.Equal(t, domain.Money(2599), *patch.Value)
		require.Nil(t, patch.Description)
	})

	t.Run("returns 400 on invalid amount", func(t *testing.T) {
		updater := &usecasesfakes.FakeMovementUpdater{}
		r := newMovementRouter(nil, nil, nil, nil, nil, updater, nil)

		rec := doRequest(t, r, http.MethodPatch, "/accounts/2/movements/1", map[string]string{"value": "1.234"})
		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.Equal(t, 0, updater.UpdateCallCount())
	})
}

func TestMovementHandlerDelete(t *testing.T) {
	t.Run("deletes a movement", func(t *testing.T) {
		deleter := &usecasesfakes.FakeMovementDeleter{}
		deleter.DeleteReturns(nil)
		r := newMovementRouter(nil, nil, nil, nil, nil, nil, deleter)

		rec := doRequest(t, r, http.MethodDelete, "/accounts/2/movements/1", nil)

		require.Equal(t, http.StatusNoContent, rec.Code)
		_, calledAccountId, calledId := deleter.DeleteArgsForCall(0)
		require.Equal(t, domain.Id(2), calledAccountId)
		require.Equal(t, domain.Id(1), calledId)
	})

	t.Run("returns 404 when the movement doesn't exist", func(t *testing.T) {
		deleter := &usecasesfakes.FakeMovementDeleter{}
		deleter.DeleteReturns(app_errors.NewResourceNotFoundError("Movement", "1"))
		r := newMovementRouter(nil, nil, nil, nil, nil, nil, deleter)

		rec := doRequest(t, r, http.MethodDelete, "/accounts/2/movements/1", nil)
		require.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestMovementHandlerMonths(t *testing.T) {
	monthsRouter := func(accountGetter *usecasesfakes.FakeAccountGetter, monthLister *usecasesfakes.FakeMovementMonthLister) http.Handler {
		logger := testLogger()
		monthsUC := movement.NewListMovementMonthsUseCase(accountGetter, monthLister, logger)
		h := rest.NewMovementHandler(nil, nil, nil, &monthsUC, nil, nil, nil, logger)
		mux := http.NewServeMux()
		h.Register(mux, "/accounts/{accountId}/movements")
		return mux
	}
	byMonthRouter := func(accountGetter *usecasesfakes.FakeAccountGetter, byMonthLister *usecasesfakes.FakeMovementByMonthLister) http.Handler {
		logger := testLogger()
		byMonthUC := movement.NewListMovementsByMonthUseCase(accountGetter, byMonthLister, logger)
		h := rest.NewMovementHandler(nil, nil, nil, nil, &byMonthUC, nil, nil, logger)
		mux := http.NewServeMux()
		h.Register(mux, "/accounts/{accountId}/movements")
		return mux
	}
	existingAccount := func() *usecasesfakes.FakeAccountGetter {
		g := &usecasesfakes.FakeAccountGetter{}
		acc := domain.NewAccount(1, "Bank A", domain.AccountTypeBank, domain.CurrencyEUR, time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC))
		g.GetReturns(&acc, nil)
		return g
	}

	t.Run("GET months returns the year/month pairs", func(t *testing.T) {
		monthLister := &usecasesfakes.FakeMovementMonthLister{}
		monthLister.ListMonthsReturns([]domain.YearMonth{{Year: 2025, Month: 2}, {Year: 2025, Month: 1}}, nil)
		r := monthsRouter(existingAccount(), monthLister)

		rec := doRequest(t, r, http.MethodGet, "/accounts/1/movements/months", nil)

		require.Equal(t, http.StatusOK, rec.Code)
		_, calledAccId := monthLister.ListMonthsArgsForCall(0)
		require.Equal(t, domain.Id(1), calledAccId)
		var resp []map[string]any
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		require.Equal(t, []map[string]any{{"year": 2025.0, "month": 2.0}, {"year": 2025.0, "month": 1.0}}, resp)
	})

	t.Run("GET months returns 404 when the account does not exist", func(t *testing.T) {
		getter := &usecasesfakes.FakeAccountGetter{}
		getter.GetReturns(nil, app_errors.NewResourceNotFoundError("Account", "9"))
		r := monthsRouter(getter, &usecasesfakes.FakeMovementMonthLister{})

		rec := doRequest(t, r, http.MethodGet, "/accounts/9/movements/months", nil)

		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("GET by-month returns the movements of that month", func(t *testing.T) {
		byMonthLister := &usecasesfakes.FakeMovementByMonthLister{}
		byMonthLister.ListByMonthReturns([]domain.Movement{sampleMovement(1)}, nil)
		r := byMonthRouter(existingAccount(), byMonthLister)

		rec := doRequest(t, r, http.MethodGet, "/accounts/1/movements/by-month/2024/1", nil)

		require.Equal(t, http.StatusOK, rec.Code)
		_, calledAccId, year, month := byMonthLister.ListByMonthArgsForCall(0)
		require.Equal(t, domain.Id(1), calledAccId)
		require.Equal(t, 2024, year)
		require.Equal(t, 1, month)
		var resp []map[string]any
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		require.Len(t, resp, 1)
		require.Equal(t, "Salary", resp[0]["description"])
		require.Equal(t, "1000.00", resp[0]["value"])
	})

	t.Run("GET by-month returns 400 for a month outside 1-12", func(t *testing.T) {
		r := byMonthRouter(existingAccount(), &usecasesfakes.FakeMovementByMonthLister{})

		rec := doRequest(t, r, http.MethodGet, "/accounts/1/movements/by-month/2024/13", nil)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("GET by-month returns 400 for a non-numeric year", func(t *testing.T) {
		r := byMonthRouter(existingAccount(), &usecasesfakes.FakeMovementByMonthLister{})

		rec := doRequest(t, r, http.MethodGet, "/accounts/1/movements/by-month/abc/1", nil)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("GET by-month returns 404 when the account does not exist", func(t *testing.T) {
		getter := &usecasesfakes.FakeAccountGetter{}
		getter.GetReturns(nil, app_errors.NewResourceNotFoundError("Account", "9"))
		r := byMonthRouter(getter, &usecasesfakes.FakeMovementByMonthLister{})

		rec := doRequest(t, r, http.MethodGet, "/accounts/9/movements/by-month/2024/1", nil)

		require.Equal(t, http.StatusNotFound, rec.Code)
	})
}
