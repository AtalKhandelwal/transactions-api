package postgres

import (
	"context"
	"errors"

	"github.com/atalkhandelwal/transactions-api/internal/repository"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountRepo struct {
	db *pgxpool.Pool
}

func NewAccountRepo(db *pgxpool.Pool) *AccountRepo {
	return &AccountRepo{db: db}
}

func (r *AccountRepo) Create(ctx context.Context,documentNumber string) (repository.Account, error) {
	const q = `INSERT INTO accounts (document_number) VALUES ($1) RETURNING id, document_number, created_at`
	var a repository.Account
	err := r.db.QueryRow(ctx, q, documentNumber).Scan(&a.ID, &a.DocumentNumber, &a.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return repository.Account{}, repository.ErrDuplicate
		}
		return repository.Account{}, err
	}
	return a, err
}

func (r *AccountRepo) GetByID(ctx context.Context,id int64) (repository.Account, error) {
	const q = `SELECT id, document_number, created_at FROM accounts WHERE id=$1`
	var a repository.Account
	err := r.db.QueryRow(ctx, q, id).Scan(&a.ID, &a.DocumentNumber, &a.CreatedAt)
	return a, err
}

var _ repository.AccountRepository = (*AccountRepo)(nil)
