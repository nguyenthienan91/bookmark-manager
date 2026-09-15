package redis

import (
	"context"
	"time"
)

// Store defines the Redis operations used by this application.
//
//go:generate mockery --name Store --filename=store.go --output ./mocks --outpkg mocks
type Store interface {
	// SetNX sets key to value only if the key does not exist.
	// Returns true if the key was set (created), false if it already existed.
	SetNX(ctx context.Context, key string, value string, expiration time.Duration) (bool, error)
}
