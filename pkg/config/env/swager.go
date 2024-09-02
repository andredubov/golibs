package env

import (
	"net"
	"os"

	"github.com/andredubov/golibs/pkg/config"
	"github.com/pkg/errors"
)

const (
	swaggerHostEnvName = "SWGR_HOST"
	swaggerPortEnvName = "SWGR_PORT"
)

type swaggerConfig struct {
	host string
	port string
}

// NewSwaggerConfig returns new swagger config
func NewSwaggerConfig() (config.SwaggerConfig, error) {
	host := os.Getenv(swaggerHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("swagger host not found")
	}

	port := os.Getenv(swaggerPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("swagger port not found")
	}

	return &swaggerConfig{
		host: host,
		port: port,
	}, nil
}

// Address returns address from host and port
func (cfg *swaggerConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
