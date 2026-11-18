package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
)

// NewRouter builds the full API router: account routes plus the nested
// movement, investment and investment-transaction routes under
// /accounts/{accountId}/... CORS is open so the webapp dev server
// (e.g. :5173) can call the API.
func NewRouter(accounts *AccountHandler, movements *MovementHandler, investments *InvestmentHandler, transactions *InvestmentTransactionHandler) http.Handler {
	r := chi.NewRouter()
	r.Route("/accounts", func(r chi.Router) {
		r.Mount("/", accounts.Router())
		r.Route("/{accountId}/movements", func(r chi.Router) {
			r.Mount("/", movements.Router())
		})
		r.Route("/{accountId}/investments", func(r chi.Router) {
			r.Mount("/", investments.Router())
			r.Route("/{investmentId}/transactions", func(r chi.Router) {
				r.Mount("/", transactions.Router())
			})
		})
	})
	return cors.AllowAll().Handler(r)
}
