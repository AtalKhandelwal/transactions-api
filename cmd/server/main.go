package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/atalkhandelwal/transactions-api/internal/config"
	dbpkg "github.com/atalkhandelwal/transactions-api/internal/db"
	httpapi "github.com/atalkhandelwal/transactions-api/internal/httpapi"
	postgresrepo "github.com/atalkhandelwal/transactions-api/internal/repository/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.FromEnv()

	var pool *pgxpool.Pool
	var err error
	for i := 0; i < 20; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		pool, err = dbpkg.Connect(ctx, cfg.DSN())
		cancel()
		if err == nil {
			break
		}
		log.Printf("waiting for database... (%v)", err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	accountsRepo := postgresrepo.NewAccountRepo(pool)
	opsRepo := postgresrepo.NewOperationRepo(pool)
	txsRepo := postgresrepo.NewTransactionRepo(pool)

	//router constructs handlers internally
	r := httpapi.NewRouter(httpapi.Deps{
		Accounts: accountsRepo,
		Ops:      opsRepo,
		Tx:       txsRepo,
	})

	addr := ":" + cfg.Port
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
