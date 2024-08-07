package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/andredubov/golibs/pkg/client/database"
	"github.com/andredubov/golibs/pkg/client/database/prettier"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type key string

const (
	// TxKey context key
	TxKey key = "tx"
)

type pg struct {
	dbConnection *pgxpool.Pool
}

// NewDB ...
func NewDB(dbConnection *pgxpool.Pool) database.Database {
	return &pg{
		dbConnection: dbConnection,
	}
}

// ScanOneContext ...
func (p *pg) ScanOneContext(ctx context.Context, dest interface{}, q database.Query, args ...interface{}) error {
	logQuery(ctx, q, args...)

	row, err := p.QueryContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanOne(dest, row)
}

// ScanAllContext ...
func (p *pg) ScanAllContext(ctx context.Context, dest interface{}, q database.Query, args ...interface{}) error {
	logQuery(ctx, q, args...)

	rows, err := p.QueryContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanAll(dest, rows)
}

// ExecContext ...
func (p *pg) ExecContext(ctx context.Context, query database.Query, args ...interface{}) (pgconn.CommandTag, error) {
	logQuery(ctx, query, args...)

	if tx, ok := ctx.Value(TxKey).(pgx.Tx); ok {
		return tx.Exec(ctx, query.QueryRaw, args...)
	}

	return p.dbConnection.Exec(ctx, query.QueryRaw, args...)
}

// QueryContext ...
func (p *pg) QueryContext(ctx context.Context, q database.Query, args ...interface{}) (pgx.Rows, error) {
	logQuery(ctx, q, args...)

	if tx, ok := ctx.Value(TxKey).(pgx.Tx); ok {
		return tx.Query(ctx, q.QueryRaw, args...)
	}

	return p.dbConnection.Query(ctx, q.QueryRaw, args...)
}

// QueryRowContext ...
func (p *pg) QueryRowContext(ctx context.Context, q database.Query, args ...interface{}) pgx.Row {
	logQuery(ctx, q, args...)

	if tx, ok := ctx.Value(TxKey).(pgx.Tx); ok {
		return tx.QueryRow(ctx, q.QueryRaw, args...)
	}

	return p.dbConnection.QueryRow(ctx, q.QueryRaw, args...)
}

// BeginTx ...
func (p *pg) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return p.dbConnection.BeginTx(ctx, txOptions)
}

// Ping ...
func (p *pg) Ping(ctx context.Context) error {
	return p.dbConnection.Ping(ctx)
}

// Close ...
func (p *pg) Close() {
	p.dbConnection.Close()
}

// MakeContextTx ...
func MakeContextTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, TxKey, tx)
}

func logQuery(ctx context.Context, q database.Query, args ...interface{}) {
	prettyQuery := prettier.Pretty(q.QueryRaw, prettier.PlaceholderDollar, args...)
	log.Println(
		ctx,
		fmt.Sprintf("sql: %s", q.Name),
		fmt.Sprintf("query: %s", prettyQuery),
	)
}
