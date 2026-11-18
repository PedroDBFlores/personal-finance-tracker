package rest_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/pedrodbflores/personal-finance-tracker/adapters/rest"
	"github.com/pedrodbflores/personal-finance-tracker/app_errors"
	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/usecases/investmenttransaction"
	"github.com/pedrodbflores/personal-finance-tracker/usecases/usecasesfakes"
	"github.com/stretchr/testify/require"
)

// newInvestmentTransactionRouter builds the transaction routes with real usecases backed by fake ports.
func newInvestmentTransactionRouter(creator *usecasesfakes.FakeInvestmentTransactionCreator, getter *usecasesfakes.FakeInvestmentTransactionGetter, lister *usecasesfakes.FakeInvestmentTransactionLister, updater *usecasesfakes.FakeInvestmentTransactionUpdater, deleter *usecasesfakes.FakeInvestmentTransactionDeleter) chi.Router {
	logger := testLogger()
	createUC := investmenttransaction.NewCreateInvestmentTransactionUseCase(creator, logger)
	getUC := investmenttransaction.NewGetInvestmentTransactionUseCase(getter, logger)
	listUC := investmenttransaction.NewListInvestmentTransactionsUseCase(lister, logger)
	updateUC := investmenttransaction.NewUpdateInvestmentTransactionUseCase(updater, logger)
	deleteUC := investmenttransaction.NewDeleteInvestmentTransactionUseCase(deleter, logger)
	h := rest.NewInvestmentTransactionHandler(&createUC, &getUC, &listUC, &updateUC, &deleteUC, logger)
	r := chi.NewRouter()
	r.Route("/accounts/{accountId}/investments/{investmentId}/transactions", func(r chi.Router) {
		r.Mount("/", h.Router())
	})
	return r
}

func sampleTransaction(investmentId uuid.UUID) domain.InvestmentTransaction {
	return domain.NewInvestmentTransaction(uuid.New(), investmentId, domain.Buy, 0, domain.Money(1000), domain.Money(10000), time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC), "Initial Buy")
}

func TestInvestmentTransactionHandlerCreate(t *testing.T) {
	investmentId := uuid.New()

	t.Run("creates a buy transaction under an investment", func(t *testing.T) {
		creator := &usecasesfakes.FakeInvestmentTransactionCreator{}
		id := uuid.New()
		creator.CreateReturns(&id, nil)
		r := newInvestmentTransactionRouter(creator, nil, nil, nil, nil)

		rec := doRequest(t, r, http.MethodPost, "/accounts/2/investments/"+investmentId.String()+"/transactions", map[string]string{
			"type": "Buy", "quantity": "10", "price": "100.00", "description": "Initial Buy",
		})

		require.Equal(t, http.StatusCreated, rec.Code)
		_, calledInvestmentId, modelTx := creator.CreateArgsForCall(0)
		require.Equal(t, investmentId, calledInvestmentId)
		require.Equal(t, domain.Buy, modelTx.Type)
		require.Equal(t, domain.Money(1000), modelTx.Quantity)
		require.Equal(t, domain.Money(10000), modelTx.Price)
		require.Equal(t, "Initial Buy", modelTx.Description)

		var resp map[string]any
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		require.Equal(t, id.String(), resp["id"])
		require.Equal(t, investmentId.String(), resp["investment_id"])
		require.Equal(t, "100.00", resp["price"])
	})

	t.Run("creates a deposit transaction with amount", func(t *testing.T) {
		creator := &usecasesfakes.FakeInvestmentTransactionCreator{}
		id := uuid.New()
		creator.CreateReturns(&id, nil)
		r := newInvestmentTransactionRouter(creator, nil, nil, nil, nil)

		rec := doRequest(t, r, http.MethodPost, "/accounts/2/investments/"+investmentId.String()+"/transactions", map[string]string{
			"type": "Deposit", "amount": "500.00",
		})

		require.Equal(t, http.StatusCreated, rec.Code)
		_, _, modelTx := creator.CreateArgsForCall(0)
		require.Equal(t, domain.Deposit, modelTx.Type)
		require.Equal(t, domain.Money(50000), modelTx.Amount)
	})

	t.Run("returns 400 on invalid type", func(t *testing.T) {
		creator := &usecasesfakes.FakeInvestmentTransactionCreator{}
		r := newInvestmentTransactionRouter(creator, nil, nil, nil, nil)

		rec := doRequest(t, r, http.MethodPost, "/accounts/2/investments/"+investmentId.String()+"/transactions", map[string]string{
			"type": "Transfer", "amount": "10",
		})
		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.Equal(t, 0, creator.CreateCallCount())
	})

	t.Run("returns 400 on invalid amount", func(t *testing.T) {
		creator := &usecasesfakes.FakeInvestmentTransactionCreator{}
		r := newInvestmentTransactionRouter(creator, nil, nil, nil, nil)

		rec := doRequest(t, r, http.MethodPost, "/accounts/2/investments/"+investmentId.String()+"/transactions", map[string]string{
			"type": "Deposit", "amount": "abc",
		})
		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.Equal(t, 0, creator.CreateCallCount())
	})

	t.Run("returns 400 on invalid investment id", func(t *testing.T) {
		creator := &usecasesfakes.FakeInvestmentTransactionCreator{}
		r := newInvestmentTransactionRouter(creator, nil, nil, nil, nil)

		rec := doRequest(t, r, http.MethodPost, "/accounts/2/investments/nope/transactions", map[string]string{
			"type": "Deposit", "amount": "10",
		})
		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.Equal(t, 0, creator.CreateCallCount())
	})
}

func TestInvestmentTransactionHandlerGet(t *testing.T) {
	investmentId := uuid.New()

	t.Run("returns the transaction", func(t *testing.T) {
		getter := &usecasesfakes.FakeInvestmentTransactionGetter{}
		tx := sampleTransaction(investmentId)
		getter.GetReturns(&tx, nil)
		r := newInvestmentTransactionRouter(nil, getter, nil, nil, nil)

		rec := doRequest(t, r, http.MethodGet, "/accounts/2/investments/"+investmentId.String()+"/transactions/"+tx.Id.String(), nil)

		require.Equal(t, http.StatusOK, rec.Code)
		_, calledInvestmentId, calledId := getter.GetArgsForCall(0)
		require.Equal(t, investmentId, calledInvestmentId)
		require.Equal(t, tx.Id, calledId)
		var resp map[string]any
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		require.Equal(t, "Buy", resp["type"])
	})

	t.Run("returns 404 when not found", func(t *testing.T) {
		getter := &usecasesfakes.FakeInvestmentTransactionGetter{}
		txId := uuid.New()
		getter.GetReturns(nil, app_errors.NewResourceNotFoundError("InvestmentTransaction", txId.String()))
		r := newInvestmentTransactionRouter(nil, getter, nil, nil, nil)

		rec := doRequest(t, r, http.MethodGet, "/accounts/2/investments/"+investmentId.String()+"/transactions/"+txId.String(), nil)
		require.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestInvestmentTransactionHandlerList(t *testing.T) {
	t.Run("lists transactions of an investment", func(t *testing.T) {
		lister := &usecasesfakes.FakeInvestmentTransactionLister{}
		investmentId := uuid.New()
		lister.ListReturns([]domain.InvestmentTransaction{sampleTransaction(investmentId)}, nil)
		r := newInvestmentTransactionRouter(nil, nil, lister, nil, nil)

		rec := doRequest(t, r, http.MethodGet, "/accounts/2/investments/"+investmentId.String()+"/transactions", nil)

		require.Equal(t, http.StatusOK, rec.Code)
		_, calledInvestmentId := lister.ListArgsForCall(0)
		require.Equal(t, investmentId, calledInvestmentId)
		var resp []map[string]any
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		require.Len(t, resp, 1)
	})
}

func TestInvestmentTransactionHandlerUpdate(t *testing.T) {
	investmentId := uuid.New()

	t.Run("applies a partial update", func(t *testing.T) {
		updater := &usecasesfakes.FakeInvestmentTransactionUpdater{}
		updater.UpdateReturns(nil)
		r := newInvestmentTransactionRouter(nil, nil, nil, updater, nil)
		txId := uuid.New()

		rec := doRequest(t, r, http.MethodPatch, "/accounts/2/investments/"+investmentId.String()+"/transactions/"+txId.String(), map[string]string{"price": "120.00"})

		require.Equal(t, http.StatusNoContent, rec.Code)
		_, calledInvestmentId, calledId, patch := updater.UpdateArgsForCall(0)
		require.Equal(t, investmentId, calledInvestmentId)
		require.Equal(t, txId, calledId)
		require.NotNil(t, patch.Price)
		require.Equal(t, domain.Money(12000), *patch.Price)
		require.Nil(t, patch.Description)
	})

	t.Run("returns 400 on invalid price", func(t *testing.T) {
		updater := &usecasesfakes.FakeInvestmentTransactionUpdater{}
		r := newInvestmentTransactionRouter(nil, nil, nil, updater, nil)

		rec := doRequest(t, r, http.MethodPatch, "/accounts/2/investments/"+investmentId.String()+"/transactions/"+uuid.New().String(), map[string]string{"price": "1.234"})
		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.Equal(t, 0, updater.UpdateCallCount())
	})
}

func TestInvestmentTransactionHandlerDelete(t *testing.T) {
	investmentId := uuid.New()

	t.Run("deletes a transaction", func(t *testing.T) {
		deleter := &usecasesfakes.FakeInvestmentTransactionDeleter{}
		deleter.DeleteReturns(nil)
		r := newInvestmentTransactionRouter(nil, nil, nil, nil, deleter)
		txId := uuid.New()

		rec := doRequest(t, r, http.MethodDelete, "/accounts/2/investments/"+investmentId.String()+"/transactions/"+txId.String(), nil)

		require.Equal(t, http.StatusNoContent, rec.Code)
		_, calledInvestmentId, calledId := deleter.DeleteArgsForCall(0)
		require.Equal(t, investmentId, calledInvestmentId)
		require.Equal(t, txId, calledId)
	})

	t.Run("returns 404 when the transaction doesn't exist", func(t *testing.T) {
		deleter := &usecasesfakes.FakeInvestmentTransactionDeleter{}
		txId := uuid.New()
		deleter.DeleteReturns(app_errors.NewResourceNotFoundError("InvestmentTransaction", txId.String()))
		r := newInvestmentTransactionRouter(nil, nil, nil, nil, deleter)

		rec := doRequest(t, r, http.MethodDelete, "/accounts/2/investments/"+investmentId.String()+"/transactions/"+txId.String(), nil)
		require.Equal(t, http.StatusNotFound, rec.Code)
	})
}
