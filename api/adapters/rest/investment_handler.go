package rest

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/pedrodbflores/personal-finance-tracker/domain"
	"github.com/pedrodbflores/personal-finance-tracker/dtos"
)

// Usecase interfaces the handlers depend on (satisfied by the usecase structs).

type investmentCreator interface {
	Execute(ctx context.Context, accountId domain.Id, dto dtos.CreateInvestmentDto) (*domain.Id, error)
}

type investmentGetter interface {
	Execute(ctx context.Context, accountId domain.Id, id domain.Id) (*dtos.InvestmentDto, error)
}

type investmentLister interface {
	Execute(ctx context.Context, accountId domain.Id) ([]dtos.InvestmentDto, error)
}

type investmentUpdater interface {
	Execute(ctx context.Context, accountId domain.Id, id domain.Id, update domain.InvestmentPatch) error
}

type investmentDeleter interface {
	Execute(ctx context.Context, accountId domain.Id, id domain.Id) error
}

// --- JSON wire types ---

type createInvestmentRequest struct {
	Name         string     `json:"name"`
	Type         string     `json:"type"`
	Currency     string     `json:"currency"`
	PurchaseDate time.Time  `json:"purchase_date"`
	MaturityDate *time.Time `json:"maturity_date,omitempty"`
}

type updateInvestmentRequest struct {
	Name         *string     `json:"name"`
	Type         *string     `json:"type"`
	Currency     *string     `json:"currency"`
	PurchaseDate *time.Time  `json:"purchase_date"`
	MaturityDate **time.Time `json:"maturity_date"`
}

type investmentResponse struct {
	Id           string     `json:"id"`
	AccountId    string     `json:"account_id"`
	Name         string     `json:"name"`
	Type         string     `json:"type"`
	CurrentValue string     `json:"current_value"` // formatted decimal, e.g. "1000.00"
	Currency     string     `json:"currency"`
	PurchaseDate time.Time  `json:"purchase_date"`
	MaturityDate *time.Time `json:"maturity_date,omitempty"`
}

func toInvestmentResponse(i dtos.InvestmentDto) investmentResponse {
	return investmentResponse{
		Id:           i.Id.String(),
		AccountId:    i.AccountId.String(),
		Name:         i.Name,
		Type:         string(i.Type),
		CurrentValue: i.CurrentValue.String(),
		Currency:     string(i.Currency),
		PurchaseDate: i.PurchaseDate,
		MaturityDate: i.MaturityDate,
	}
}

// --- handler ---

type InvestmentHandler struct {
	create investmentCreator
	get    investmentGetter
	list   investmentLister
	update investmentUpdater
	del    investmentDeleter
	logger *slog.Logger
}

func NewInvestmentHandler(create investmentCreator, get investmentGetter, list investmentLister, update investmentUpdater, del investmentDeleter, logger *slog.Logger) *InvestmentHandler {
	return &InvestmentHandler{create, get, list, update, del, logger}
}

// Register mounts the investment routes under prefix (e.g. "/accounts/{accountId}/investments").
func (h *InvestmentHandler) Register(mux *http.ServeMux, prefix string) {
	mux.HandleFunc("POST "+prefix, h.Create)
	mux.HandleFunc("GET "+prefix, h.List)
	mux.HandleFunc("GET "+prefix+"/{investmentId}", h.Get)
	mux.HandleFunc("PATCH "+prefix+"/{investmentId}", h.Update)
	mux.HandleFunc("DELETE "+prefix+"/{investmentId}", h.Delete)
}

func (h *InvestmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	accountId, ok := parseIdParam(w, r, "accountId")
	if !ok {
		return
	}
	var req createInvestmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	investmentType, err := domain.ParseInvestmentType(req.Type)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	purchaseDate := req.PurchaseDate
	if purchaseDate.IsZero() {
		purchaseDate = time.Now()
	}
	dto := dtos.NewCreateInvestmentDto(req.Name, investmentType, domain.Currency(req.Currency), purchaseDate, req.MaturityDate)
	id, err := h.create.Execute(r.Context(), accountId, dto)
	if err != nil {
		writeUsecaseError(w, err)
		return
	}
	inv, err := domain.NewInvestment(*id, accountId, req.Name, investmentType, domain.Currency(req.Currency), purchaseDate, req.MaturityDate)
	if err != nil {
		writeUsecaseError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toInvestmentResponse(dtos.NewInvestmentDto(inv)))
}

func (h *InvestmentHandler) Get(w http.ResponseWriter, r *http.Request) {
	accountId, ok := parseIdParam(w, r, "accountId")
	if !ok {
		return
	}
	id, ok := parseIdParam(w, r, "investmentId")
	if !ok {
		return
	}
	inv, err := h.get.Execute(r.Context(), accountId, id)
	if err != nil {
		writeUsecaseError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toInvestmentResponse(*inv))
}

func (h *InvestmentHandler) List(w http.ResponseWriter, r *http.Request) {
	accountId, ok := parseIdParam(w, r, "accountId")
	if !ok {
		return
	}
	investments, err := h.list.Execute(r.Context(), accountId)
	if err != nil {
		writeUsecaseError(w, err)
		return
	}
	resp := make([]investmentResponse, 0, len(investments))
	for _, i := range investments {
		resp = append(resp, toInvestmentResponse(i))
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *InvestmentHandler) Update(w http.ResponseWriter, r *http.Request) {
	accountId, ok := parseIdParam(w, r, "accountId")
	if !ok {
		return
	}
	id, ok := parseIdParam(w, r, "investmentId")
	if !ok {
		return
	}
	var req updateInvestmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	patch := domain.InvestmentPatch{
		Name:         req.Name,
		Currency:     (*domain.Currency)(req.Currency),
		PurchaseDate: req.PurchaseDate,
		MaturityDate: req.MaturityDate,
	}
	if req.Type != nil {
		investmentType, err := domain.ParseInvestmentType(*req.Type)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		patch.Type = &investmentType
	}
	if err := h.update.Execute(r.Context(), accountId, id, patch); err != nil {
		writeUsecaseError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *InvestmentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	accountId, ok := parseIdParam(w, r, "accountId")
	if !ok {
		return
	}
	id, ok := parseIdParam(w, r, "investmentId")
	if !ok {
		return
	}
	if err := h.del.Execute(r.Context(), accountId, id); err != nil {
		writeUsecaseError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- shared helpers ---

// parseMoneyParam parses a decimal string into Money; an empty string means zero.
func parseMoneyParam(s string) (domain.Money, error) {
	if s == "" {
		return 0, nil
	}
	return domain.ParseMoney(s)
}
