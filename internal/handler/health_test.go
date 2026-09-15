package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nguyenthienan91/bookmark-manager/internal/model"
	"github.com/nguyenthienan91/bookmark-manager/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var errHealthCheck = errors.New("redis unavailable")

func TestHealthCheckHandler_HealthCheck(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                 string
		setupMockService     func(t *testing.T) *mocks.HealthCheck
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "redis is reachable",
			setupMockService: func(t *testing.T) *mocks.HealthCheck {
				mockSvc := mocks.NewHealthCheck(t)
				mockSvc.On("Check", mock.Anything).Return(model.HealthCheckResponse{
					Message:     "OK",
					ServiceName: "bookmark_service",
					InstanceID:  "test-instance-id",
				}, nil)
				return mockSvc
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: `{
				"message": "OK",
				"service_name": "bookmark_service",
				"instance_id":  "test-instance-id"
			}`,
		},
		{
			name: "redis is unavailable",
			setupMockService: func(t *testing.T) *mocks.HealthCheck {
				mockSvc := mocks.NewHealthCheck(t)
				mockSvc.On("Check", mock.Anything).Return(model.HealthCheckResponse{}, errHealthCheck)
				return mockSvc
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"redis unavailable"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			ctx.Request = httptest.NewRequest(http.MethodGet, "/health-check", nil).
				WithContext(context.Background())

			mockSvc := tc.setupMockService(t)
			h := NewHealthCheck(mockSvc)
			h.HealthCheck(ctx)

			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			assert.JSONEq(t, tc.expectedResponseBody, rec.Body.String())
		})
	}
}
