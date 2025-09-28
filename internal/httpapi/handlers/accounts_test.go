package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/atalkhandelwal/transactions-api/internal/httpapi/handlers"
	"github.com/atalkhandelwal/transactions-api/internal/repository"
	"github.com/jackc/pgx/v5"
)

// --- stubs ---
type stubAccountsRepo struct {
	createResp repository.Account
	createErr  error
	getResp    repository.Account
	getErr     error
}

func (s *stubAccountsRepo) Create(ctx context.Context, doc string) (repository.Account, error) {
	if s.createErr != nil {
		return repository.Account{}, s.createErr
	}
	if s.createResp.DocumentNumber == "" {
		s.createResp.DocumentNumber = doc
	}
	return s.createResp, s.createErr
}
func (s *stubAccountsRepo) GetByID(ctx context.Context, id int64) (repository.Account, error) {
	return s.getResp, s.getErr
}

// --- tests ---

func TestCreateAccount_Success(t *testing.T) {
	repo := &stubAccountsRepo{createResp: repository.Account{ID: 1}}
	h := handlers.NewAccountsHandler(repo)

	req := httptest.NewRequest(http.MethodPost, "/accounts", strings.NewReader(`{"document_number":"12345678900"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("got %d, want 201", rr.Code)
	}
	var body map[string]any
	_ = json.NewDecoder(rr.Body).Decode(&body)
	if body["account_id"].(float64) != 1 {
		t.Fatalf("unexpected account_id: %v", body["account_id"])
	}
	if body["document_number"].(string) != "12345678900" {
		t.Fatalf("unexpected document_number: %v", body["document_number"])
	}
}

func TestGetAccount_BadID_BadRequest(t *testing.T) {
	repo := &stubAccountsRepo{getErr: pgx.ErrNoRows}
	h := handlers.NewAccountsHandler(repo)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/accounts/abc", nil) // invalid id → 400

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("got %d, want 400 due to invalid accountId format", rr.Code)
	}
}

func TestCreateAccount_Duplicate_Returns409(t *testing.T) {
	repo := &stubAccountsRepo{createErr: repository.ErrDuplicate}
	h := handlers.NewAccountsHandler(repo)

	req := httptest.NewRequest(http.MethodPost, "/accounts",
		strings.NewReader(`{"document_number":"12345678900"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusConflict {
		t.Fatalf("got %d, want 409 Conflict", rr.Code)
	}
}
