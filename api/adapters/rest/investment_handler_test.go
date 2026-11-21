package rest_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pedrodbflores/personal-finance-tracker/adapters/rest"
	"github.com/pedrodbflores/personal-finance-tracker/app_errors"
	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/usecases/investment"
	"github.com/pedrodbflores/personal-finance-tracker/usecases/usecasesfakes"
	"github.com/stretchr/testify/require"
)

// newInvestmentRouter builds the investment routes with real usecases backed by fake ports.
func newInvestmentRouter(creator *usecasesfakes.FakeInvestmentCreator, getter *usecasesfakes.FakeInvestmentGetter, lister *usecasesfakes.FakeInvestmentLister, updater *usecasesfakes.FakeInvestmentUpdater, deleter *usecasesfakes.FakeInvestmentDeleter, txLister *usecasesfakes.FakeInvestmentTransactionLister) http.Handler {
	logger := testLogger()
	createUC := investment.NewCreateInvestmentUseCase(creator, logger)
	getUC := investment.NewGetInvestmentUseCase(getter, txLister, logger)
	listUC := investment.NewListInvestmentsUseCase(lister, txLister, logger)
	updateUC := investment.NewUpdateInvestmentUseCase(updater, logger)
	deleteUC := investment.NewDeleteInvestmentUseCase(deleter, logger)
	h := rest.NewInvestmentHandler(&createUC, &getUC, &listUC, &updateUC, &deleteUC, logger)
	mux := http.NewServeMux()
	h.Register(mux, "/accounts/{accountId}/investments")
	return mux
}

func sampleInvestment(accountId domain.Id, id uuid.UUID) domain.Investment {
	return domain.NewInvestment(id, accountId, "Vanguard ETF", domain.ETF, domain.CurrencyUSD, time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC), nil)
}

func TestInvestmentHandlerCreate(t *testing.T) {
	t.Run("creates an investment under an account", func(t *testing.T) {
		creator := &usecasesfakes.FakeInvestmentCreator{}
		id := uuid.New()
		creator.CreateReturns(&id, nil)
		r := newInvestmentRouter(creator, nil, nil, nil, nil, nil)

		rec := doRequest(t, r, http.MethodPost, "/accounts/2/investments", map[string]string{
			"name": "Vanguard ETF", "type": "ETF", "currency": "USD",
		})

		require.Equal(t, http.StatusCreated, rec.Code)
		_, calledAccountId, modelInvestment := creator.CreateArgsForCall(0)
		require.Equal(t, domain.Id(2), calledAccountId)
		require.Equal(t, "Vanguard ETF", modelInvestment.Name)
		require.Equal(t, domain.ETF, modelInvestment.Type)
		require.Equal(t, domain.CurrencyUSD, modelInvestment.Currency)

		var resp map[string]any
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		require.Equal(t, id.String(), resp["id"])
		require.Equal(t, "2", resp["account_id"])
		require.Equal(t, "0.00", resp["current_value"])
	})

	t.Run("returns 400 when name is missing", func(t *testing.T) {
		creator := &usecasesfakes.FakeInvestmentCreator{}
		r := newInvestmentRouter(creator, nil, nil, nil, nil, nil)

		rec := doRequest(t, r, http.MethodPost, "/accounts/2/investments", map[string]string{
			"name": "", "type": "ETF", "currency": "USD",
		})
		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.Equal(t, 0, creator.CreateCallCount())
	})

	t.Run("returns 400 on invalid type", func(t *testing.T) {
		creator := &usecasesfakes.FakeInvestmentCreator{}
		r := newInvestmentRouter(creator, nil, nil, nil, nil, nil)

		rec := doRequest(t, r, http.MethodPost, "/accounts/2/investments", map[string]string{
			"name": "Crypto", "type": "Crypto", "currency": "USD",
		})
		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.Equal(t, 0, creator.CreateCallCount())
	})
}

func TestInvestmentHandlerGet(t *testing.T) {
	t.Run("returns the investment with the current value calculated from its transactions", func(t *testing.T) {
		getter := &usecasesfakes.FakeInvestmentGetter{}
		txLister := &usecasesfakes.FakeInvestmentTransactionLister{}
		investmentId := uuid.New()
		inv := sampleInvestment(2, investmentId)
		getter.GetReturns(&inv, nil)
		buy := domain.NewInvestmentTransaction(uuid.New(), investmentId, domain.Buy, 0, domain.Money(1000), domain.Money(10000), time.Now(), "Buy")
		txLister.ListReturns([]domain.InvestmentTransaction{buy}, nil)
		r := newInvestmentRouter(nil, getter, nil, nil, nil, txLister)

		rec := doRequest(t, r, http.MethodGet, "/accounts/2/investments/"+investmentId.String(), nil)

		require.Equal(t, http.StatusOK, rec.Code)
		_, calledAccountId, calledId := getter.GetArgsForCall(0)
		require.Equal(t, domain.Id(2), calledAccountId)
		require.Equal(t, investmentId, calledId)
		var resp map[string]any
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		require.Equal(t, "Vanguard ETF", resp["name"])
		require.Equal(t, "1000.00", resp["current_value"])
	})

	t.Run("returns 404 when not found", func(t *testing.T) {
		getter := &usecasesfakes.FakeInvestmentGetter{}
		txLister := &usecasesfakes.FakeInvestmentTransactionLister{}
		investmentId := uuid.New()
		getter.GetReturns(nil, app_errors.NewResourceNotFoundError("Investment", investmentId.String()))
		r := newInvestmentRouter(nil, getter, nil, nil, nil, txLister)

		rec := doRequest(t, r, http.MethodGet, "/accounts/2/investments/"+investmentId.String(), nil)
		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("returns 400 on invalid investment id", func(t *testing.T) {
		r := newInvestmentRouter(nil, nil, nil, nil, nil, nil)
		rec := doRequest(t, r, http.MethodGet, "/accounts/2/investments/not-a-uuid", nil)
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestInvestmentHandlerList(t *testing.T) {
	t.Run("lists investments of an account", func(t *testing.T) {
		lister := &usecasesfakes.FakeInvestmentLister{}
		txLister := &usecasesfakes.FakeInvestmentTransactionLister{}
		lister.ListReturns([]domain.Investment{sampleInvestment(2, uuid.New())}, nil)
		txLister.ListReturns([]domain.InvestmentTransaction{}, nil)
		r := newInvestmentRouter(nil, nil, lister, nil, nil, txLister)

		rec := doRequest(t, r, http.MethodGet, "/accounts/2/investments", nil)

		require.Equal(t, http.StatusOK, rec.Code)
		_, calledAccountId := lister.ListArgsForCall(0)
		require.Equal(t, domain.Id(2), calledAccountId)
		var resp []map[string]any
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		require.Len(t, resp, 1)
	})
}

func TestInvestmentHandlerUpdate(t *testing.T) {
	t.Run("applies a partial update", func(t *testing.T) {
		updater := &usecasesfakes.FakeInvestmentUpdater{}
		updater.UpdateReturns(nil)
		r := newInvestmentRouter(nil, nil, nil, updater, nil, nil)
		investmentId := uuid.New()

		rec := doRequest(t, r, http.MethodPatch, "/accounts/2/investments/"+investmentId.String(), map[string]string{"name": "Renamed ETF"})

		require.Equal(t, http.StatusNoContent, rec.Code)
		_, calledAccountId, calledId, patch := updater.UpdateArgsForCall(0)
		require.Equal(t, domain.Id(2), calledAccountId)
		require.Equal(t, investmentId, calledId)
		require.NotNil(t, patch.Name)
		require.Equal(t, "Renamed ETF", *patch.Name)
		require.Nil(t, patch.Type)
	})

	t.Run("returns 400 on invalid type", func(t *testing.T) {
		updater := &usecasesfakes.FakeInvestmentUpdater{}
		r := newInvestmentRouter(nil, nil, nil, updater, nil, nil)

		rec := doRequest(t, r, http.MethodPatch, "/accounts/2/investments/"+uuid.New().String(), map[string]string{"type": "Crypto"})
		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.Equal(t, 0, updater.UpdateCallCount())
	})
}

func TestInvestmentHandlerDelete(t *testing.T) {
	t.Run("deletes an investment", func(t *testing.T) {
		deleter := &usecasesfakes.FakeInvestmentDeleter{}
		deleter.DeleteReturns(nil)
		r := newInvestmentRouter(nil, nil, nil, nil, deleter, nil)
		investmentId := uuid.New()

		rec := doRequest(t, r, http.MethodDelete, "/accounts/2/investments/"+investmentId.String(), nil)

		require.Equal(t, http.StatusNoContent, rec.Code)
		_, calledAccountId, calledId := deleter.DeleteArgsForCall(0)
		require.Equal(t, domain.Id(2), calledAccountId)
		require.Equal(t, investmentId, calledId)
	})

	t.Run("returns 404 when the investment doesn't exist", func(t *testing.T) {
		deleter := &usecasesfakes.FakeInvestmentDeleter{}
		investmentId := uuid.New()
		deleter.DeleteReturns(app_errors.NewResourceNotFoundError("Investment", investmentId.String()))
		r := newInvestmentRouter(nil, nil, nil, nil, deleter, nil)

		rec := doRequest(t, r, http.MethodDelete, "/accounts/2/investments/"+investmentId.String(), nil)
		require.Equal(t, http.StatusNotFound, rec.Code)
	})
}
