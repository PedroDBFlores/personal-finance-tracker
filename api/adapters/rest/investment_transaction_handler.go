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

type investmentTransactionCreator interface {
	Execute(ctx context.Context, accountId domain.Id, investmentId domain.Id, dto dtos.CreateInvestmentTransactionDto) (*domain.Id, error)
}

type investmentTransactionGetter interface {
	Execute(ctx context.Context, accountId domain.Id, investmentId, id domain.Id) (*dtos.InvestmentTransactionDto, error)
}

type investmentTransactionLister interface {
	Execute(ctx context.Context, accountId domain.Id, investmentId domain.Id) ([]dtos.InvestmentTransactionDto, error)
}

type investmentTransactionUpdater interface {
	Execute(ctx context.Context, accountId domain.Id, investmentId, id domain.Id, update domain.InvestmentTransactionPatch) error
}

type investmentTransactionDeleter interface {
	Execute(ctx context.Context, accountId domain.Id, investmentId, id domain.Id) error
}

// --- JSON wire types ---

type createInvestmentTransactionRequest struct {
	Type        string `json:"type"`
	Amount      string `json:"amount"`   // decimal string, e.g. "500.00" (Deposit/Withdraw)
	Quantity    string `json:"quantity"` // decimal string, e.g. "10" (Buy/Sell)
	Price       string `json:"price"`    // decimal string, e.g. "100.00" (Buy/Sell)
	Description string `json:"description"`
}

type updateInvestmentTransactionRequest struct {
	Type        *string `json:"type"`
	Amount      *string `json:"amount"`
	Quantity    *string `json:"quantity"`
	Price       *string `json:"price"`
	Description *string `json:"description"`
}

type investmentTransactionResponse struct {
	Id           string    `json:"id"`
	InvestmentId string    `json:"investment_id"`
	Type         string    `json:"type"`
	Amount       string    `json:"amount"` // formatted decimal, e.g. "500.00"
	Quantity     string    `json:"quantity"`
	Price        string    `json:"price"`
	Date         time.Time `json:"date"`
	Description  string    `json:"description"`
}

func toInvestmentTransactionResponse(t dtos.InvestmentTransactionDto) investmentTransactionResponse {
	return investmentTransactionResponse{
		Id:           t.Id.String(),
		InvestmentId: t.InvestmentId.String(),
		Type:         string(t.Type),
		Amount:       t.Amount.String(),
		Quantity:     t.Quantity.String(),
		Price:        t.Price.String(),
		Date:         t.Date,
		Description:  t.Description,
	}
}

// --- handler ---

type InvestmentTransactionHandler struct {
	create investmentTransactionCreator
	get    investmentTransactionGetter
	list   investmentTransactionLister
	update investmentTransactionUpdater
	del    investmentTransactionDeleter
	logger *slog.Logger
}

func NewInvestmentTransactionHandler(create investmentTransactionCreator, get investmentTransactionGetter, list investmentTransactionLister, update investmentTransactionUpdater, del investmentTransactionDeleter, logger *slog.Logger) *InvestmentTransactionHandler {
	return &InvestmentTransactionHandler{create, get, list, update, del, logger}
}

// Register mounts the transaction routes under prefix
// (e.g. "/accounts/{accountId}/investments/{investmentId}/transactions").
func (h *InvestmentTransactionHandler) Register(mux *http.ServeMux, prefix string) {
	mux.HandleFunc("POST "+prefix, h.Create)
	mux.HandleFunc("GET "+prefix, h.List)
	mux.HandleFunc("GET "+prefix+"/{transactionId}", h.Get)
	mux.HandleFunc("PATCH "+prefix+"/{transactionId}", h.Update)
	mux.HandleFunc("DELETE "+prefix+"/{transactionId}", h.Delete)
}

func (h *InvestmentTransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	accountId, ok := parseIdParam(w, r, "accountId")
	if !ok {
		return
	}
	investmentId, ok := parseIdParam(w, r, "investmentId")
	if !ok {
		return
	}
	var req createInvestmentTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	txType, err := domain.ParseInvestmentTransactionType(req.Type)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	amount, err := parseMoneyParam(req.Amount)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	quantity, err := parseMoneyParam(req.Quantity)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	price, err := parseMoneyParam(req.Price)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	dto := dtos.NewCreateInvestmentTransactionDto(txType, amount, quantity, price, time.Now(), req.Description)
	id, err := h.create.Execute(r.Context(), accountId, investmentId, dto)
	if err != nil {
		writeUsecaseError(w, err)
		return
	}
	tx, err := domain.NewInvestmentTransaction(*id, investmentId, txType, amount, quantity, price, dto.Date, req.Description)
	if err != nil {
		writeUsecaseError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toInvestmentTransactionResponse(dtos.NewInvestmentTransactionDto(tx)))
}

func (h *InvestmentTransactionHandler) Get(w http.ResponseWriter, r *http.Request) {
	accountId, ok := parseIdParam(w, r, "accountId")
	if !ok {
		return
	}
	investmentId, ok := parseIdParam(w, r, "investmentId")
	if !ok {
		return
	}
	id, ok := parseIdParam(w, r, "transactionId")
	if !ok {
		return
	}
	tx, err := h.get.Execute(r.Context(), accountId, investmentId, id)
	if err != nil {
		writeUsecaseError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toInvestmentTransactionResponse(*tx))
}

func (h *InvestmentTransactionHandler) List(w http.ResponseWriter, r *http.Request) {
	accountId, ok := parseIdParam(w, r, "accountId")
	if !ok {
		return
	}
	investmentId, ok := parseIdParam(w, r, "investmentId")
	if !ok {
		return
	}
	transactions, err := h.list.Execute(r.Context(), accountId, investmentId)
	if err != nil {
		writeUsecaseError(w, err)
		return
	}
	resp := make([]investmentTransactionResponse, 0, len(transactions))
	for _, t := range transactions {
		resp = append(resp, toInvestmentTransactionResponse(t))
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *InvestmentTransactionHandler) Update(w http.ResponseWriter, r *http.Request) {
	accountId, ok := parseIdParam(w, r, "accountId")
	if !ok {
		return
	}
	investmentId, ok := parseIdParam(w, r, "investmentId")
	if !ok {
		return
	}
	id, ok := parseIdParam(w, r, "transactionId")
	if !ok {
		return
	}
	var req updateInvestmentTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	patch := domain.InvestmentTransactionPatch{Description: req.Description}
	if req.Type != nil {
		txType, err := domain.ParseInvestmentTransactionType(*req.Type)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		patch.Type = &txType
	}
	if req.Amount != nil {
		amount, err := domain.ParseMoney(*req.Amount)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		patch.Amount = &amount
	}
	if req.Quantity != nil {
		quantity, err := domain.ParseMoney(*req.Quantity)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		patch.Quantity = &quantity
	}
	if req.Price != nil {
		price, err := domain.ParseMoney(*req.Price)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		patch.Price = &price
	}
	if err := h.update.Execute(r.Context(), accountId, investmentId, id, patch); err != nil {
		writeUsecaseError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *InvestmentTransactionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	accountId, ok := parseIdParam(w, r, "accountId")
	if !ok {
		return
	}
	investmentId, ok := parseIdParam(w, r, "investmentId")
	if !ok {
		return
	}
	id, ok := parseIdParam(w, r, "transactionId")
	if !ok {
		return
	}
	if err := h.del.Execute(r.Context(), accountId, investmentId, id); err != nil {
		writeUsecaseError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
