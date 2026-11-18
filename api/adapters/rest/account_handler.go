package rest

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/pedrodbflores/personal-finance-tracker/app_errors"
	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/dtos"
)

// Usecase interfaces the handlers depend on (satisfied by the usecase structs).

type accountCreator interface {
	Execute(ctx context.Context, dto dtos.CreateAccountDto) (*domain.Id, error)
}

type accountGetter interface {
	Execute(ctx context.Context, id domain.Id) (*dtos.AccountDto, error)
}

type accountLister interface {
	Execute(ctx context.Context) ([]dtos.AccountDto, error)
}

type accountUpdater interface {
	Execute(ctx context.Context, id domain.Id, update domain.AccountPatch) error
}

type accountDeleter interface {
	Execute(ctx context.Context, id domain.Id) error
}

// --- JSON wire types ---

type createAccountRequest struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Currency string `json:"currency"`
}

type updateAccountRequest struct {
	Name     *string `json:"name"`
	Type     *string `json:"type"`
	Currency *string `json:"currency"`
}

type accountResponse struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

func toAccountResponse(a dtos.AccountDto) accountResponse {
	return accountResponse{
		Id:        a.Id.String(),
		Name:      a.Name,
		Type:      string(a.Type),
		Currency:  string(a.Currency),
		CreatedAt: a.CreatedAt,
	}
}

// --- handler ---

type AccountHandler struct {
	create accountCreator
	get    accountGetter
	list   accountLister
	update accountUpdater
	del    accountDeleter
	logger *slog.Logger
}

func NewAccountHandler(create accountCreator, get accountGetter, list accountLister, update accountUpdater, del accountDeleter, logger *slog.Logger) *AccountHandler {
	return &AccountHandler{create, get, list, update, del, logger}
}

// Router mounts the account routes.
func (h *AccountHandler) Router() chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.Create)
	r.Get("/", h.List)
	r.Get("/{id}", h.Get)
	r.Patch("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
	return r
}

func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	dto := dtos.NewCreateAccountDto(req.Name, domain.AccountType(req.Type), domain.Currency(req.Currency), time.Now())
	id, err := h.create.Execute(r.Context(), dto)
	if err != nil {
		writeUsecaseError(w, err)
		return
	}
	acc := dtos.NewAccountDto(domain.NewAccount(*id, req.Name, domain.AccountType(req.Type), domain.Currency(req.Currency), dto.CreatedAt))
	writeJSON(w, http.StatusCreated, toAccountResponse(acc))
}

func (h *AccountHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIdParam(w, r, "id")
	if !ok {
		return
	}
	acc, err := h.get.Execute(r.Context(), id)
	if err != nil {
		writeUsecaseError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toAccountResponse(*acc))
}

func (h *AccountHandler) List(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.list.Execute(r.Context())
	if err != nil {
		writeUsecaseError(w, err)
		return
	}
	resp := make([]accountResponse, 0, len(accounts))
	for _, a := range accounts {
		resp = append(resp, toAccountResponse(a))
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *AccountHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIdParam(w, r, "id")
	if !ok {
		return
	}
	var req updateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	patch := domain.AccountPatch{
		Name:     req.Name,
		Currency: (*domain.Currency)(req.Currency),
	}
	if req.Type != nil {
		t := domain.AccountType(*req.Type)
		patch.Type = &t
	}
	if err := h.update.Execute(r.Context(), id, patch); err != nil {
		writeUsecaseError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *AccountHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIdParam(w, r, "id")
	if !ok {
		return
	}
	if err := h.del.Execute(r.Context(), id); err != nil {
		writeUsecaseError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- shared helpers ---

func parseIdParam(w http.ResponseWriter, r *http.Request, name string) (domain.Id, bool) {
	raw := chi.URLParam(r, name)
	n, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid "+name+" parameter")
		return 0, false
	}
	return domain.Id(n), true
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func writeUsecaseError(w http.ResponseWriter, err error) {
	var notFound app_errors.ResourceNotFoundError
	if errors.As(err, &notFound) {
		writeError(w, http.StatusNotFound, notFound.Error())
		return
	}
	writeError(w, http.StatusInternalServerError, err.Error())
}
