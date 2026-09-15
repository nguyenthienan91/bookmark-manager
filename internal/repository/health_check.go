package repository

import (
	"context"

	pkgredis "github.com/nguyenthienan91/bookmark-manager/pkg/redis"
)

// HealthCheckRepository defines the persistence dependency used by health checks.
//
//go:generate mockery --name HealthCheckRepository --filename=health_check.go --output ./mocks --outpkg mocks
type HealthCheckRepository interface {
	// Ping checks whether the backing Redis service is reachable.
	Ping(ctx context.Context) error
}

type healthCheckRepository struct {
	store pkgredis.Store
}

// NewHealthCheckRepository creates a health-check repository backed by Redis.
func NewHealthCheckRepository(store pkgredis.Store) HealthCheckRepository {
	return &healthCheckRepository{store: store}
}

// Ping checks Redis connectivity and returns any error from the Redis client.
func (r *healthCheckRepository) Ping(ctx context.Context) error {
	return r.store.Ping(ctx).Err()
}
