package postgres

import (
	"context"
	"errors"

	"github.com/atalkhandelwal/transactions-api/internal/repository"
	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OperationRepo struct {
	db *pgxpool.Pool
}

func NewOperationRepo(db *pgxpool.Pool) *OperationRepo {
	return &OperationRepo{db: db}
}

func (r *OperationRepo) Exists(ctx context.Context, id int) (bool, error) {
	const q = `SELECT 1 FROM operation_types WHERE id=$1`
	var one int
	if err := r.db.QueryRow(ctx, q, id).Scan(&one); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err 
	}
	return true, nil
}

var _ repository.OperationTypeRepository = (*OperationRepo)(nil)
