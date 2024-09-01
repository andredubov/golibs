package env

import (
	"fmt"
	"net"
	"os"

	"github.com/andredubov/golibs/pkg/config"
)

const (
	httpHostEnvName = "HTTP_HOST"
	httpPortEnvName = "HTTP_PORT"
)

type httpConfig struct {
	host string
	port string
}

// NewHTTPConfig returns an instance of httpConfig struct
func NewHTTPConfig() (config.HTTPConfig, error) {
	host := os.Getenv(httpHostEnvName)
	if len(host) == 0 {
		return nil, fmt.Errorf("%s", "http host not found")
	}

	port := os.Getenv(httpPortEnvName)
	if len(port) == 0 {
		return nil, fmt.Errorf("%s", "http port not found")
	}

	return &httpConfig{
		host: host,
		port: port,
	}, nil
}

// Address returns grpc server address
func (cfg *httpConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
