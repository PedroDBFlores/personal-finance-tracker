package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
)

// NewRouter builds the full API router: account routes plus the nested
// movement routes under /accounts/{accountId}/movements. CORS is open so the
// webapp dev server (e.g. :5173) can call the API.
func NewRouter(accounts *AccountHandler, movements *MovementHandler) http.Handler {
	r := chi.NewRouter()
	r.Route("/accounts", func(r chi.Router) {
		r.Mount("/", accounts.Router())
		r.Route("/{accountId}/movements", func(r chi.Router) {
			r.Mount("/", movements.Router())
		})
	})
	return cors.AllowAll().Handler(r)
}
