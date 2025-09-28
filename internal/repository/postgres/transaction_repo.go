package postgres

import (
	"context"

	"github.com/atalkhandelwal/transactions-api/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepo struct {
	db *pgxpool.Pool
}

func NewTransactionRepo(db *pgxpool.Pool) *TransactionRepo {
	return &TransactionRepo{db: db}
}

func (r *TransactionRepo) Create(ctx context.Context,t repository.Transaction) (repository.Transaction, error) {
	const q = `
		INSERT INTO transactions (account_id, operation_type_id, amount)
		VALUES ($1, $2, $3)
		RETURNING id, account_id, operation_type_id, amount, event_date
	`
	var out repository.Transaction
	err := r.db.QueryRow(ctx, q, t.AccountID, t.OperationTypeID, t.Amount).
		Scan(&out.ID, &out.AccountID, &out.OperationTypeID, &out.Amount, &out.EventDate)
	return out, err
}

var _ repository.TransactionRepository = (*TransactionRepo)(nil)
