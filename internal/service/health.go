package service

import "github.com/nguyenthienan91/bookmark-manager/internal/model"

type HealthCheck interface {
	Check() model.HealthCheckResponse
}

type healthCheckService struct {
	serviceName string
	instanceID  string
}

func NewHealthCheck(serviceName string, instanceID string) HealthCheck {
	return &healthCheckService{
		serviceName: serviceName,
		instanceID:  instanceID,
	}
}

func (h *healthCheckService) Check() model.HealthCheckResponse {
	return model.HealthCheckResponse{
		Message:     "OK",
		ServiceName: h.serviceName,
		InstanceID:  h.instanceID,
	}
}