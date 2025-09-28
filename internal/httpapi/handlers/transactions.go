package handlers

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"

	"github.com/atalkhandelwal/transactions-api/internal/repository"
	"github.com/atalkhandelwal/transactions-api/internal/service"
	"github.com/jackc/pgx/v5"
)

type TransactionsHandler struct {
	Accounts     repository.AccountRepository
	Operations   repository.OperationTypeRepository
	Transactions repository.TransactionRepository
}

func NewTransactionsHandler(a repository.AccountRepository, o repository.OperationTypeRepository, t repository.TransactionRepository) *TransactionsHandler {
	return &TransactionsHandler{
		Accounts: a, Operations: o, Transactions: t,
	}
}

type createTransactionReq struct {
	AccountID       int64   `json:"account_id"`
	OperationTypeID int     `json:"operation_type_id"`
	Amount          float64 `json:"amount"`
}

func (h *TransactionsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req createTransactionReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if req.AccountID <= 0 || req.OperationTypeID <= 0 {
		writeError(w, http.StatusBadRequest, "account_id and operation_type_id are required")
		return
	}

	if !(req.Amount > 0) || math.IsNaN(req.Amount) || math.IsInf(req.Amount, 0) {
		writeError(w, http.StatusBadRequest, "amount is required and must be a positive number")
		return
	}

	if _, err := h.Accounts.GetByID(r.Context(), req.AccountID); err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "account not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to read account")
		return
	}

	exists, err := h.Operations.Exists(r.Context(), req.OperationTypeID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to read operation type")
		return
	}
	if !exists {
		writeError(w, http.StatusBadRequest, "invalid operation_type_id")
		return
	}

	normalized := service.NormalizeAmount(req.OperationTypeID, req.Amount)
	

	txn, err := h.Transactions.Create(r.Context(), repository.Transaction{
		AccountID:       req.AccountID,
		OperationTypeID: req.OperationTypeID,
		Amount:          normalized,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not create transaction")
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{
		"transaction_id":    txn.ID,
		"account_id":        txn.AccountID,
		"operation_type_id": txn.OperationTypeID,
		"amount":            txn.Amount,
		"event_date":        txn.EventDate,
	})
}
