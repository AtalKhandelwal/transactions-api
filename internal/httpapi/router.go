package http

import (
	"net/http"

	chI "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/atalkhandelwal/transactions-api/internal/httpapi/handlers"
	"github.com/atalkhandelwal/transactions-api/internal/repository"
)

type Deps struct {
	Accounts repository.AccountRepository
	Ops      repository.OperationTypeRepository
	Tx       repository.TransactionRepository
}

func NewRouter(d Deps) *chI.Mux {
	r := chI.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Health
	r.Get("/healthc", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// Construct handlers here from dependencies
	accountsHandler := handlers.NewAccountsHandler(d.Accounts)
	transactionsHandler := handlers.NewTransactionsHandler(d.Accounts, d.Ops, d.Tx)

	// Accounts
	r.Method(http.MethodPost, "/accounts", accountsHandler)
	r.Method(http.MethodGet, "/accounts/{accountId}", accountsHandler)

	// Transactions
	r.Method(http.MethodPost, "/transactions", transactionsHandler)

	return r
}
