package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/atalkhandelwal/transactions-api/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type AccountsHandler struct {
	Repo repository.AccountRepository
}

func NewAccountsHandler(repo repository.AccountRepository) *AccountsHandler {
	return &AccountsHandler{Repo: repo}
}

type createAccountReq struct {
	DocumentNumber string `json:"document_number"`
}

func (h *AccountsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.create(w, r)
	case http.MethodGet:
		h.getByID(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *AccountsHandler) create(w http.ResponseWriter, r *http.Request) {
	var req createAccountReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if req.DocumentNumber == "" {
		writeError(w, http.StatusBadRequest, "document_number is required")
		return
	}
	acc, err := h.Repo.Create(r.Context(), req.DocumentNumber)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			writeError(w, http.StatusConflict, "document_number already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, "could not create account")
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{
		"account_id":      acc.ID,
		"document_number": acc.DocumentNumber,
	})
}

func (h *AccountsHandler) getByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "accountId")
	if idStr == "" {
		writeError(w, http.StatusBadRequest, "accountId is required")
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid accountId")
		return
	}
	acc, err := h.Repo.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "account not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get account")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"account_id":      acc.ID,
		"document_number": acc.DocumentNumber,
	})

}
