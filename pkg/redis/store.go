package redis

import (
	"context"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

// Store defines the Redis operation used by the application.
//
//go:generate mockery --name Store --filename=store.go --output ./mocks --outpkg mocks
type Store interface {
	// Ping checks whether the Redis server is reachable.
	Ping(ctx context.Context) *goredis.StatusCmd

	// SetNX sets key to value only if the key does not exist.
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *goredis.BoolCmd
}
