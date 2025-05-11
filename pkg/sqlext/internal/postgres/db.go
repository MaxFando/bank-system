package postgres

import (
	"context"
	"database/sql"
	"github.com/MaxFando/bank-system/pkg/sqlext/transaction"

	"github.com/jmoiron/sqlx"
)

type DB struct {
	db                 *sqlx.DB
	transactionManager *TransactionManager
}

func NewDB(db *sqlx.DB) *DB {
	return &DB{
		db:                 db,
		transactionManager: NewTransactionManager(db),
	}
}

func (d *DB) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	executor := transaction.GetExecutor(ctx, d.db)
	if len(args) == 0 {
		return executor.GetContext(ctx, dest, query)
	}

	return executor.GetContext(ctx, dest, query, args...)
}

func (d *DB) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	executor := transaction.GetExecutor(ctx, d.db)
	if len(args) == 0 {
		return executor.SelectContext(ctx, dest, query)
	}

	return executor.SelectContext(ctx, dest, query, args...)
}

func (d *DB) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	executor := transaction.GetExecutor(ctx, d.db)
	if len(args) == 0 {
		return executor.ExecContext(ctx, query)
	}

	return executor.ExecContext(ctx, query, args...)
}

func (d *DB) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	executor := transaction.GetExecutor(ctx, d.db)
	if len(args) == 0 {
		return executor.QueryContext(ctx, query)
	}

	return executor.QueryContext(ctx, query, args...)
}

func (d *DB) NamedExec(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	executor := transaction.GetExecutor(ctx, d.db)
	return executor.NamedExecContext(ctx, query, arg)
}

func (d *DB) Rebind(query string) string {
	return d.db.Rebind(query)
}

func (d *DB) BindNamed(query string, arg interface{}) (string, []interface{}, error) {
	return d.db.BindNamed(query, arg)
}

func (d *DB) WithTx(ctx context.Context, fn transaction.AtomicFn, opts ...transaction.TxOption) error {
	return d.transactionManager.RunTransaction(ctx, fn, opts...)
}
