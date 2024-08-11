package redis

import (
	"context"
	"log"
	"time"

	"github.com/andredubov/golibs/pkg/client/cache"
	"github.com/andredubov/golibs/pkg/config"
	redigo "github.com/gomodule/redigo/redis"
)

type handler func(ctx context.Context, conn redigo.Conn) error

type rd struct {
	connectionPool *redigo.Pool
	config         config.RedisConfig
}

// NewCache returns a new instance of redis struct
func NewCache(connectionPool *redigo.Pool, config config.RedisConfig) cache.Cache {
	return &rd{
		connectionPool,
		config,
	}
}

// Set binds key and its value
func (r *rd) Set(ctx context.Context, key string, value interface{}) error {
	err := r.execute(ctx, func(ctx context.Context, conn redigo.Conn) error {
		if _, err := conn.Do("SET", redigo.Args{key}.Add(value)...); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// Get returns value by its key from cache
func (r *rd) Get(ctx context.Context, key string) (interface{}, error) {
	var value interface{}
	err := r.execute(ctx, func(ctx context.Context, conn redigo.Conn) error {
		var errEx error
		value, errEx = conn.Do("GET", key)
		if errEx != nil {
			return errEx
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return value, nil
}

// HashSet binds key and value pair to the hash into cache
func (r *rd) HashSet(ctx context.Context, hash string, values interface{}) error {
	err := r.execute(ctx, func(ctx context.Context, conn redigo.Conn) error {
		if _, err := conn.Do("HSET", redigo.Args{hash}.AddFlat(values)...); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// HashGetAll returns a sequence of key-value pairs of the corresponding hash
func (r *rd) HashGetAll(ctx context.Context, hash string) ([]interface{}, error) {
	var values []interface{}
	err := r.execute(ctx, func(ctx context.Context, conn redigo.Conn) error {
		var errEx error
		values, errEx = redigo.Values(conn.Do("HGETALL", hash))
		if errEx != nil {
			return errEx
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return values, nil
}

// Expire sets time to expire key into cache
func (r *rd) Expire(ctx context.Context, key string, expiration time.Duration) error {
	err := r.execute(ctx, func(ctx context.Context, conn redigo.Conn) error {
		if _, err := conn.Do("EXPIRE", key, int(expiration.Seconds())); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// Delete a value by its key into cache
func (r *rd) Delete(ctx context.Context, key string) error {
	err := r.execute(ctx, func(ctx context.Context, conn redigo.Conn) error {
		_, err := conn.Do("DEL", key)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// Ping tests cache connection
func (r *rd) Ping(ctx context.Context) error {
	err := r.execute(ctx, func(ctx context.Context, conn redigo.Conn) error {
		if _, err := conn.Do("PING"); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// Close closes cache connection
func (r *rd) Close() error {
	return r.connectionPool.Close()
}

func (r *rd) execute(ctx context.Context, handler handler) error {
	connection, err := r.getConnect(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err = connection.Close(); err != nil {
			log.Printf("failed to close redis connection: %v\n", err)
		}
	}()

	if err = handler(ctx, connection); err != nil {
		return err
	}

	return nil
}

func (r *rd) getConnect(ctx context.Context) (redigo.Conn, error) {
	getConnTimeoutCtx, cancel := context.WithTimeout(ctx, r.config.ConnectionTimeout())
	defer cancel()

	conn, err := r.connectionPool.GetContext(getConnTimeoutCtx)
	if err != nil {
		log.Printf("failed to get redis connection: %v\n", err)
		if err := conn.Close(); err != nil {
			return nil, err
		}

		return nil, err
	}

	return conn, nil
}
