package service

import (
	"context"
	"errors"
	"testing"

	repomocks "github.com/nguyenthienan91/bookmark-manager/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var errHealthRepository = errors.New("redis unavailable")

func TestHealthCheckService_Check(t *testing.T) {
	t.Parallel()

	const (
		serviceName = "bookmark_service"
		instanceID  = "test-instance-id"
	)

	testCases := []struct {
		name          string
		setupMockRepo func(t *testing.T) *repomocks.HealthCheckRepository
		expectedErr   error
	}{
		{
			name: "redis is reachable",
			setupMockRepo: func(t *testing.T) *repomocks.HealthCheckRepository {
				repo := repomocks.NewHealthCheckRepository(t)
				repo.On("Ping", mock.Anything).Return(nil)
				return repo
			},
		},
		{
			name: "redis is unavailable",
			setupMockRepo: func(t *testing.T) *repomocks.HealthCheckRepository {
				repo := repomocks.NewHealthCheckRepository(t)
				repo.On("Ping", mock.Anything).Return(errHealthRepository)
				return repo
			},
			expectedErr: errHealthRepository,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := tc.setupMockRepo(t)
			healthSvc := NewHealthCheck(serviceName, instanceID, repo)

			response, err := healthSvc.Check(context.Background())

			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
				assert.Empty(t, response)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, "OK", response.Message)
			assert.Equal(t, serviceName, response.ServiceName)
			assert.Equal(t, instanceID, response.InstanceID)
		})
	}
}
