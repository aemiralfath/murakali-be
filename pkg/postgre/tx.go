package postgre

import (
	"context"
	"database/sql"
)

type Transaction interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type TxFn func(Transaction) error

type TxRepo struct {
	PSQL *sql.DB
}

func NewTxRepository(psql *sql.DB) *TxRepo {
	return &TxRepo{psql}
}

func (tr *TxRepo) WithTransaction(fn TxFn) (err error) {
	tx, err := tr.PSQL.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}