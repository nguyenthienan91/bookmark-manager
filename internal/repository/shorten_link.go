package repository

import (
	"context"
	"time"

	pkgredis "github.com/nguyenthienan91/bookmark-manager/pkg/redis"
)

// LinkRepository defines the contract for persisting shortened links.
//
//go:generate mockery --name LinkRepository --filename=link_repository.go --output ./mocks --outpkg mocks
type LinkRepository interface {
	// SaveLink stores the mapping key -> URL with the given TTL.
	// Returns (true, nil) when saved, (false, nil) when key already exists,
	// and (false, err) on Redis error.
	SaveLink(ctx context.Context, key string, url string, expiration time.Duration) (bool, error)

	// GetLink retrieves the original URL stored under key.
	// Returns ("", redis.Nil) when the key does not exist.
	GetLink(ctx context.Context, key string) (string, error)
}

type linkRepository struct {
	store pkgredis.Store
}

// NewLinkRepository creates a repository backed by the provided Redis store.
func NewLinkRepository(store pkgredis.Store) LinkRepository {
	return &linkRepository{store: store}
}

// SaveLink stores a URL under key using Redis SETNX and returns whether the key was created.
func (r *linkRepository) SaveLink(ctx context.Context, key string, url string, expiration time.Duration) (bool, error) {
	return r.store.SetNX(ctx, key, url, expiration).Result()
}

// GetLink retrieves the original URL stored under key.
func (r *linkRepository) GetLink(ctx context.Context, key string) (string, error) {
	return r.store.Get(ctx, key).Result()
}
