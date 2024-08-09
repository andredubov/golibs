package cache

import (
	"context"
	"time"
)

// Cache interface for working with a cache
type Cache interface {
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string) (interface{}, error)
	HashSet(ctx context.Context, key string, values interface{}) error
	HashGetAll(ctx context.Context, key string) ([]interface{}, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
	Ping(ctx context.Context) error
	Close() error
}

// Client interface for working with a cache
type Client interface {
	Cache() Cache
	Close() error
}
