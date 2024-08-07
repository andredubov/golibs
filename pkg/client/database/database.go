package database

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

// Handler is a function that is executed in a transaction
type Handler func(ctx context.Context) error

// Client a client for working with a database
type Client interface {
	Database() Database
	Close() error
}

// Transactor is a transaction manager that executes a user-specified handler in a transaction
type Transactor interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

// TxManager is a transaction manager that executes a user-specified handler in a transaction
type TxManager interface {
	ReadCommitted(ctx context.Context, f Handler) error
}

// Query wrapper above the query, storing the query name and the query itself
// The request name is used for logging and can potentially be used somewhere else, for example, for tracking
type Query struct {
	Name     string
	QueryRaw string
}

// SQLExecer combines NamedExecer and QueryExecer
type SQLExecer interface {
	NamedExecer
	QueryExecer
}

// NamedExecer interface for working with named queries using tags in structures
type NamedExecer interface {
	ScanOneContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
	ScanAllContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
}

// QueryExecer interface for working with regular queries
type QueryExecer interface {
	ExecContext(ctx context.Context, q Query, args ...interface{}) (pgconn.CommandTag, error)
	QueryContext(ctx context.Context, q Query, args ...interface{}) (pgx.Rows, error)
	QueryRowContext(ctx context.Context, q Query, args ...interface{}) pgx.Row
}

// Pinger interface for checking database connection
type Pinger interface {
	Ping(ctx context.Context) error
}

// Database interface for working with a database
type Database interface {
	SQLExecer
	Transactor
	Pinger
	Close()
}
