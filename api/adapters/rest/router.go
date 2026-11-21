package rest

import (
	"net/http"

	"github.com/rs/cors"
)

// NewRouter builds the full API router on a stdlib ServeMux: account routes plus
// the nested movement, investment and investment-transaction routes under
// /accounts/{accountId}/... Go 1.22+ method+path patterns carry the path params.
// CORS is open so the webapp dev server (e.g. :5173) can call the API.
func NewRouter(accounts *AccountHandler, movements *MovementHandler, investments *InvestmentHandler, transactions *InvestmentTransactionHandler) http.Handler {
	mux := http.NewServeMux()
	accounts.Register(mux, "/accounts")
	movements.Register(mux, "/accounts/{accountId}/movements")
	investments.Register(mux, "/accounts/{accountId}/investments")
	transactions.Register(mux, "/accounts/{accountId}/investments/{investmentId}/transactions")
	return cors.AllowAll().Handler(mux)
}
