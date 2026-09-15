package redis

import (
	"context"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

// Client wraps the go-redis client and implements the Store interface.
type Client struct {
	rdb *goredis.Client
}

// NewClient creates a new Redis client from the given config.
func NewClient(cfg *Config) (*Client, error) {
	rdb := goredis.NewClient(&goredis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	return &Client{rdb: rdb}, nil
}

// Ping checks that the Redis server is reachable.
func (c *Client) Ping(ctx context.Context) error {
	return c.rdb.Ping(ctx).Err()
}

// Close closes the underlying Redis connection.
func (c *Client) Close() error {
	return c.rdb.Close()
}

// SetNX implements Store. It sets key to value only when key does not exist.
func (c *Client) SetNX(ctx context.Context, key string, value string, expiration time.Duration) (bool, error) {
	return c.rdb.SetNX(ctx, key, value, expiration).Result()
}