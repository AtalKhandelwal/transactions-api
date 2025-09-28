package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/atalkhandelwal/transactions-api/internal/httpapi/handlers"
	"github.com/atalkhandelwal/transactions-api/internal/repository"
	"github.com/jackc/pgx/v5"
)

// --- stubs ---
type stubAccRepo struct {
	getResp repository.Account
	getErr  error
}

func (s *stubAccRepo) Create(ctx context.Context, _ string) (repository.Account, error) { return repository.Account{}, nil }

func (s *stubAccRepo) GetByID(ctx context.Context, id int64) (repository.Account, error) {
	if s.getResp.ID == 0 && s.getErr == nil {
		return repository.Account{ID: id, DocumentNumber: "X"}, nil
	}
	return s.getResp, s.getErr
}

type stubOpsRepo struct {
	exists bool
}

func (s *stubOpsRepo) Exists(ctx context.Context, id int) (bool, error) {
	if id >= 1 && id <= 4 {
		return s.exists || true, nil
	}
	return false, nil
}

type stubTxRepo struct {
	got repository.Transaction
}

func (s *stubTxRepo) Create(ctx context.Context, t repository.Transaction) (repository.Transaction, error) {
	s.got = t
	t.ID = 99
	t.EventDate = time.Now()
	return t, nil
}

// --- tests ---

func TestCreateTransaction_Normalizes_Negative_ForPurchase(t *testing.T) {
	a := &stubAccRepo{}
	o := &stubOpsRepo{exists: true}
	tx := &stubTxRepo{}
	h := handlers.NewTransactionsHandler(a, o, tx)

	body := `{"account_id":1,"operation_type_id":1,"amount":50}`
	req := httptest.NewRequest(http.MethodPost, "/transactions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("got %d, want 201", rr.Code)
	}
	if tx.got.Amount != -50.0 {
		t.Fatalf("amount not normalized. got %v, want -50", tx.got.Amount)
	}
	var resp map[string]any
	_ = json.NewDecoder(rr.Body).Decode(&resp)
	if int(resp["transaction_id"].(float64)) != 99 {
		t.Fatalf("unexpected transaction_id: %v", resp["transaction_id"])
	}
}

func TestCreateTransaction_InvalidOperationType(t *testing.T) {
	a := &stubAccRepo{}
	o := &stubOpsRepo{exists: false}
	tx := &stubTxRepo{}
	h := handlers.NewTransactionsHandler(a, o, tx)

	body := `{"account_id":1,"operation_type_id":999,"amount":10}`
	req := httptest.NewRequest(http.MethodPost, "/transactions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("got %d, want 400 for invalid operation type", rr.Code)
	}
}

func TestCreateTransaction_AccountNotFound(t *testing.T) {
	a := &stubAccRepo{getErr: pgx.ErrNoRows}
	o := &stubOpsRepo{exists: true}
	tx := &stubTxRepo{}
	h := handlers.NewTransactionsHandler(a, o, tx)

	body := `{"account_id":999,"operation_type_id":1,"amount":10}`
	req := httptest.NewRequest(http.MethodPost, "/transactions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("got %d, want 404 when account not found", rr.Code)
	}
}

func TestCreateTransaction_RejectsNegativeAmountFromClient(t *testing.T) {
	a := &stubAccRepo{}
	o := &stubOpsRepo{exists: true}
	tx := &stubTxRepo{}
	h := handlers.NewTransactionsHandler(a, o, tx)

	body := `{"account_id":1,"operation_type_id":1,"amount":-50}`
	req := httptest.NewRequest(http.MethodPost, "/transactions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("got %d; want 400 because amount must be positive", rr.Code)
	}
}
