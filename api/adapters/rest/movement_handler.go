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

type movementCreator interface {
	Execute(ctx context.Context, accountId domain.Id, dto dtos.CreateMovementDto) (*domain.Id, error)
}

type movementGetter interface {
	Execute(ctx context.Context, accountId, id domain.Id) (*dtos.MovementDto, error)
}

type movementLister interface {
	Execute(ctx context.Context, accountId domain.Id) ([]dtos.MovementDto, error)
}

type movementUpdater interface {
	Execute(ctx context.Context, accountId, id domain.Id, update domain.MovementPatch) error
}

type movementDeleter interface {
	Execute(ctx context.Context, accountId, id domain.Id) error
}

// --- JSON wire types ---

type createMovementRequest struct {
	Description string `json:"description"`
	Value       string `json:"value"` // decimal string, e.g. "25.99"
	Type        string `json:"type"`
}

type updateMovementRequest struct {
	Description *string `json:"description"`
	Value       *string `json:"value"`
	Type        *string `json:"type"`
}

type movementResponse struct {
	Id          string    `json:"id"`
	AccountId   string    `json:"account_id"`
	Description string    `json:"description"`
	Value       string    `json:"value"` // formatted decimal, e.g. "25.99"
	Type        string    `json:"type"`
	Date        time.Time `json:"date"`
}

func toMovementResponse(m dtos.MovementDto) movementResponse {
	return movementResponse{
		Id:          m.Id.String(),
		AccountId:   m.AccountId.String(),
		Description: m.Description,
		Value:       m.Value.String(),
		Type:        string(m.Type),
		Date:        m.Date,
	}
}

// --- handler ---

type MovementHandler struct {
	create movementCreator
	get    movementGetter
	list   movementLister
	update movementUpdater
	del    movementDeleter
	logger *slog.Logger
}

func NewMovementHandler(create movementCreator, get movementGetter, list movementLister, update movementUpdater, del movementDeleter, logger *slog.Logger) *MovementHandler {
	return &MovementHandler{create, get, list, update, del, logger}
}

// Register mounts the movement routes under prefix (e.g. "/accounts/{accountId}/movements").
func (h *MovementHandler) Register(mux *http.ServeMux, prefix string) {
	mux.HandleFunc("POST "+prefix, h.Create)
	mux.HandleFunc("GET "+prefix, h.List)
	mux.HandleFunc("GET "+prefix+"/{movementId}", h.Get)
	mux.HandleFunc("PATCH "+prefix+"/{movementId}", h.Update)
	mux.HandleFunc("DELETE "+prefix+"/{movementId}", h.Delete)
}

func (h *MovementHandler) Create(w http.ResponseWriter, r *http.Request) {
	accountId, ok := parseIdParam(w, r, "accountId")
	if !ok {
		return
	}
	var req createMovementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if req.Description == "" {
		writeError(w, http.StatusBadRequest, "description is required")
		return
	}
	value, err := domain.ParseMoney(req.Value)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	dto := dtos.NewCreateMovementDto(req.Description, value, domain.MovementType(req.Type), time.Now())
	id, err := h.create.Execute(r.Context(), accountId, dto)
	if err != nil {
		writeUsecaseError(w, err)
		return
	}
	mov := domain.NewMovement(accountId, req.Description, value, domain.MovementType(req.Type), dto.Date)
	mov.Id = *id
	writeJSON(w, http.StatusCreated, toMovementResponse(dtos.NewMovementDto(mov)))
}

func (h *MovementHandler) Get(w http.ResponseWriter, r *http.Request) {
	accountId, ok := parseIdParam(w, r, "accountId")
	if !ok {
		return
	}
	id, ok := parseIdParam(w, r, "movementId")
	if !ok {
		return
	}
	mov, err := h.get.Execute(r.Context(), accountId, id)
	if err != nil {
		writeUsecaseError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toMovementResponse(*mov))
}

func (h *MovementHandler) List(w http.ResponseWriter, r *http.Request) {
	accountId, ok := parseIdParam(w, r, "accountId")
	if !ok {
		return
	}
	movements, err := h.list.Execute(r.Context(), accountId)
	if err != nil {
		writeUsecaseError(w, err)
		return
	}
	resp := make([]movementResponse, 0, len(movements))
	for _, m := range movements {
		resp = append(resp, toMovementResponse(m))
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *MovementHandler) Update(w http.ResponseWriter, r *http.Request) {
	accountId, ok := parseIdParam(w, r, "accountId")
	if !ok {
		return
	}
	id, ok := parseIdParam(w, r, "movementId")
	if !ok {
		return
	}
	var req updateMovementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	patch := domain.MovementPatch{Description: req.Description}
	if req.Value != nil {
		value, err := domain.ParseMoney(*req.Value)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		patch.Value = &value
	}
	if req.Type != nil {
		t := domain.MovementType(*req.Type)
		patch.Type = &t
	}
	if err := h.update.Execute(r.Context(), accountId, id, patch); err != nil {
		writeUsecaseError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *MovementHandler) Delete(w http.ResponseWriter, r *http.Request) {
	accountId, ok := parseIdParam(w, r, "accountId")
	if !ok {
		return
	}
	id, ok := parseIdParam(w, r, "movementId")
	if !ok {
		return
	}
	if err := h.del.Execute(r.Context(), accountId, id); err != nil {
		writeUsecaseError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
