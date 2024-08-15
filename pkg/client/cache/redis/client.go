package redis

import (
	"context"

	"github.com/andredubov/golibs/pkg/client/cache"
	"github.com/andredubov/golibs/pkg/config"
	redigo "github.com/gomodule/redigo/redis"
)

type redisClient struct {
	masterCache cache.Cache
}

// New returns a new instance of redisClient struct
func New(ctx context.Context, cfg config.RedisConfig) (cache.Client, error) {
	connectionPool := &redigo.Pool{
		MaxIdle:     cfg.MaxIdle(),
		IdleTimeout: cfg.IdleTimeout(),
		DialContext: func(ctx context.Context) (redigo.Conn, error) {
			return redigo.DialContext(ctx, "tcp", cfg.Address())
		},
	}

	return &redisClient{
		masterCache: &rd{connectionPool, cfg},
	}, nil
}

// Cache returns cache
func (r *redisClient) Cache() cache.Cache {
	return r.masterCache
}

// Close cache connection
func (r *redisClient) Close() error {
	if r.masterCache != nil {
		return r.masterCache.Close()
	}

	return nil
}
