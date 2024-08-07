package postgres

import (
	"context"

	"github.com/andredubov/golibs/pkg/client/database"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type pgClient struct {
	masterDBC database.Database
}

// New returns an instance of pgClient struct
func New(ctx context.Context, connectionString string) (database.Client, error) {
	dbConnection, err := pgxpool.Connect(ctx, connectionString)
	if err != nil {
		return nil, errors.Errorf("failed to connect to database: %v", err)
	}

	return &pgClient{
		masterDBC: &pg{dbConnection: dbConnection},
	}, nil
}

// Database returns database
func (c *pgClient) Database() database.Database {
	return c.masterDBC
}

// Close database connection
func (c *pgClient) Close() error {
	if c.masterDBC != nil {
		c.masterDBC.Close()
	}

	return nil
}
