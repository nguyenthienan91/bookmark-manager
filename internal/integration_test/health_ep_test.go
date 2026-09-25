package integration_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nguyenthienan91/bookmark-manager/internal/api"
	repomocks "github.com/nguyenthienan91/bookmark-manager/internal/repository/mocks"
	healthservice "github.com/nguyenthienan91/bookmark-manager/internal/service"
	svcmocks "github.com/nguyenthienan91/bookmark-manager/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var errHealthEndpoint = errors.New("redis unavailable")

func TestHealthCheckEndpoint(t *testing.T) {
	t.Parallel()

	const (
		serviceName = "bookmark_service"
		instanceID  = "test-instance-id"
	)

	testCases := []struct {
		name                 string
		setupHealthService   func(t *testing.T) healthservice.HealthCheck
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "redis is reachable",
			setupHealthService: func(t *testing.T) healthservice.HealthCheck {
				repo := repomocks.NewHealthCheckRepository(t)
				repo.On("Ping", mock.Anything).Return(nil)
				return healthservice.NewHealthCheck(serviceName, instanceID, repo)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: `{
				"message": "OK",
				"service_name": "bookmark_service",
				"instance_id": "test-instance-id"
			}`,
		},
		{
			name: "redis is unavailable",
			setupHealthService: func(t *testing.T) healthservice.HealthCheck {
				repo := repomocks.NewHealthCheckRepository(t)
				repo.On("Ping", mock.Anything).Return(errHealthEndpoint)
				return healthservice.NewHealthCheck(serviceName, instanceID, repo)
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"redis unavailable"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			healthSvc := tc.setupHealthService(t)
			shortenSvc := svcmocks.NewShortenLink(t)
			apiEngine := api.NewEngine(&api.Config{
				ServiceName: serviceName,
				InstanceID:  instanceID,
			}, healthSvc, shortenSvc, testLogger)

			req := httptest.NewRequest(http.MethodGet, "/health-check", nil)
			rec := httptest.NewRecorder()

			apiEngine.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			assert.JSONEq(t, tc.expectedResponseBody, rec.Body.String())
		})
	}
}
