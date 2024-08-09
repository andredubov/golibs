package env

import (
	"net"
	"os"
	"strconv"
	"time"

	"github.com/andredubov/golibs/pkg/config"
	"github.com/pkg/errors"
)

const (
	redisHostEnvName              = "RD_HOST"
	redisPortEnvName              = "RD_PORT"
	redisConnectionTimeoutEnvName = "RD_CONNECTION_TIMEOUT_SEC"
	redisMaxIdleEnvName           = "RD_MAX_IDLE"
	redisMaxIdleTimeoutEnvName    = "RD_MAX_IDLE_TIMEOUT_SEC"
)

type redisConfig struct {
	host              string
	port              string
	connectionTimeout time.Duration
	maxIdle           int
	maxIdleTimeout    time.Duration
}

// NewRedisConfig returns a new instance of redisConfig struct
func NewRedisConfig() (config.RedisConfig, error) {
	host := os.Getenv(redisHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("redis host not found")
	}

	port := os.Getenv(redisPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("redis port not found")
	}

	connectionTimeoutStr := os.Getenv(redisConnectionTimeoutEnvName)
	if len(connectionTimeoutStr) == 0 {
		return nil, errors.New("redis connection timeout not found")
	}

	connectionTimeout, err := strconv.ParseInt(connectionTimeoutStr, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse connection timeout")
	}

	maxIdleStr := os.Getenv(redisMaxIdleEnvName)
	if len(maxIdleStr) == 0 {
		return nil, errors.New("redis max idle not found")
	}

	maxIdle, err := strconv.Atoi(maxIdleStr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse max idle")
	}

	maxIdleTimeoutStr := os.Getenv(redisMaxIdleTimeoutEnvName)
	if len(maxIdleTimeoutStr) == 0 {
		return nil, errors.New("redis idle timeout not found")
	}

	maxIdleTimeout, err := strconv.ParseInt(maxIdleTimeoutStr, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse idle timeout")
	}

	return &redisConfig{
		host:              host,
		port:              port,
		connectionTimeout: time.Duration(connectionTimeout) * time.Second,
		maxIdle:           maxIdle,
		maxIdleTimeout:    time.Duration(maxIdleTimeout) * time.Second,
	}, nil
}

func (cfg *redisConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}

func (cfg *redisConfig) ConnectionTimeout() time.Duration {
	return cfg.connectionTimeout
}

func (cfg *redisConfig) MaxIdle() int {
	return cfg.maxIdle
}

func (cfg *redisConfig) IdleTimeout() time.Duration {
	return cfg.maxIdleTimeout
}
