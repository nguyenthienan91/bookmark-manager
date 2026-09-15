package repository

import (
	"context"
	"fmt"
	"time"

	pkgredis "github.com/nguyenthienan91/bookmark-manager/pkg/redis"
)

const redisKeyPrefix = "shorturl:"

// LinkRepository defines the contract for persisting shortened links.
//
//go:generate mockery --name LinkRepository --filename=link_repository.go --output ./mocks --outpkg mocks
type LinkRepository interface {
	// SaveLink stores the mapping code -> url with the given TTL.
	// Returns (true, nil) when saved, (false, nil) when code already exists,
	// and (false, err) on Redis error.
	SaveLink(ctx context.Context, code string, url string, expiration time.Duration) (bool, error)
}

type linkRepository struct {
	store pkgredis.Store
}

// NewLinkRepository creates a repository backed by the provided Redis store.
func NewLinkRepository(store pkgredis.Store) LinkRepository {
	return &linkRepository{store: store}
}

func (r *linkRepository) SaveLink(ctx context.Context, code string, url string, expiration time.Duration) (bool, error) {
	key := fmt.Sprintf("%s%s", redisKeyPrefix, code)
	created, err := r.store.SetNX(ctx, key, url, expiration)
	if err != nil {
		return false, err
	}
	return created, nil
}
