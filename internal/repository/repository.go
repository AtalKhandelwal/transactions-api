package repository

import (
	"context"
	"errors"
	"time"
)

type Account struct {
	ID             int64  `json:"account_id"`
	DocumentNumber string `json:"document_number"`
	CreatedAt      time.Time `json:"-"`
}

type Transaction struct {
	ID              int64     `json:"transaction_id"`
	AccountID       int64     `json:"account_id"`
	OperationTypeID int       `json:"operation_type_id"`
	Amount          float64   `json:"amount"`
	EventDate       time.Time `json:"event_date"`
}

type AccountRepository interface {
	Create(ctx context.Context,documentNumber string) (Account, error)
	GetByID(ctx context.Context,id int64) (Account, error)
}

type OperationTypeRepository interface {
	Exists(ctx context.Context,id int) (bool, error)
}

type TransactionRepository interface {
	Create(ctx context.Context,t Transaction) (Transaction, error)
}

var ErrDuplicate = errors.New("duplicate")