package service

import (
	"context"

	"github.com/nguyenthienan91/bookmark-manager/internal/model"
	"github.com/nguyenthienan91/bookmark-manager/internal/repository"
)

// HealthCheck is the interface for health-check logic.
//
//go:generate mockery --name HealthCheck --filename=health_check.go
type HealthCheck interface {
	Check(ctx context.Context) (model.HealthCheckResponse, error)
}

type healthCheckService struct {
	serviceName string
	instanceID  string
	repo        repository.HealthCheckRepository
}

// NewHealthCheck creates a health-check service with a Redis health dependency.
func NewHealthCheck(serviceName string, instanceID string, repo repository.HealthCheckRepository) HealthCheck {
	return &healthCheckService{
		serviceName: serviceName,
		instanceID:  instanceID,
		repo:        repo,
	}
}

// Check verifies Redis connectivity before returning service information.
func (h *healthCheckService) Check(ctx context.Context) (model.HealthCheckResponse, error) {
	if err := h.repo.Ping(ctx); err != nil {
		return model.HealthCheckResponse{}, err
	}

	return model.HealthCheckResponse{
		Message:     "OK",
		ServiceName: h.serviceName,
		InstanceID:  h.instanceID,
	}, nil
}
