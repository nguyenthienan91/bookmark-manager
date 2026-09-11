package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheckService_Check(t *testing.T) {
	t.Parallel()

	service := NewHealthCheck("bookmark_service", "test-instance-id")

	response := service.Check()

	assert.Equal(t, "OK", response.Message)
	assert.Equal(t, "bookmark_service", response.ServiceName)
	assert.Equal(t, "test-instance-id", response.InstanceID)
}