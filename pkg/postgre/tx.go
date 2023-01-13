package postgre

import (
	"context"
	"database/sql"
	"fmt"
)

type Transaction interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type TxFn func(Transaction) error
type TxFnData func(Transaction) (interface{}, error)

type TxRepo struct {
	PSQL *sql.DB
}

func NewTxRepository(psql *sql.DB) *TxRepo {
	return &TxRepo{psql}
}

func (tr *TxRepo) WithTransactionReturnData(fn TxFnData) (data interface{}, err error) {
	tx, err := tr.PSQL.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		var errTx error
		if p := recover(); p != nil {
			errTx = tx.Rollback()
			fmt.Println(errTx.Error())
			panic(p)
		} else if err != nil {
			errTx = tx.Rollback()
		} else {
			errTx = tx.Commit()
		}
		fmt.Println(errTx.Error())
	}()

	data, err = fn(tx)
	return data, err
}

func (tr *TxRepo) WithTransaction(fn TxFn) (err error) {
	tx, err := tr.PSQL.Begin()
	if err != nil {
		return err
	}

	defer func() {
		var errTx error
		if p := recover(); p != nil {
			errTx = tx.Rollback()
			fmt.Println(errTx.Error())
			panic(p)
		} else if err != nil {
			errTx = tx.Rollback()
		} else {
			errTx = tx.Commit()
		}
		fmt.Println(errTx.Error())
	}()

	err = fn(tx)
	return err
}
