package env

import (
	"fmt"
	"os"

	"github.com/andredubov/golibs/pkg/config"
)

const (
	hostEnvName     = "PG_HOST"
	portEnvName     = "PG_PORT"
	dbnameEnvName   = "PG_DB"
	userEnvName     = "PG_USER"
	passwordEnvName = "PG_PASSWORD"
	sslmodeEnvName  = "PG_SSL_MODE"
)

type postgresConfig struct {
	host     string
	port     string
	dbname   string
	user     string
	password string
	sslmode  string
}

// NewPostgresConfig returns an intance of postgresConfig struct
func NewPostgresConfig() (config.PostgresConfig, error) {
	host := os.Getenv(hostEnvName)
	if len(host) == 0 {
		return nil, fmt.Errorf("%s", "postgres host not found")
	}

	port := os.Getenv(portEnvName)
	if len(port) == 0 {
		return nil, fmt.Errorf("%s", "postgres port not found")
	}

	dbname := os.Getenv(dbnameEnvName)
	if len(port) == 0 {
		return nil, fmt.Errorf("%s", "postgres database name not found")
	}

	user := os.Getenv(userEnvName)
	if len(port) == 0 {
		return nil, fmt.Errorf("%s", "postgres user not found")
	}

	password := os.Getenv(passwordEnvName)
	if len(port) == 0 {
		return nil, fmt.Errorf("%s", "postgres password not found")
	}

	sslmode := os.Getenv(sslmodeEnvName)
	if len(port) == 0 {
		return nil, fmt.Errorf("%s", "postgres ssl mode not found")
	}

	return &postgresConfig{
		host:     host,
		port:     port,
		dbname:   dbname,
		user:     user,
		password: password,
		sslmode:  sslmode,
	}, nil
}

// DSN returns postgres database connecton string
func (cfg *postgresConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		cfg.host,
		cfg.port,
		cfg.dbname,
		cfg.user,
		cfg.password,
		cfg.sslmode,
	)
}
