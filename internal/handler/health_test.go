package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nguyenthienan91/bookmark-manager/internal/model"
	"github.com/stretchr/testify/assert"
)

type fakeHealthCheckService struct{}

func (f *fakeHealthCheckService) Check() model.HealthCheckResponse {
	return model.HealthCheckResponse{
		Message:     "OK",
		ServiceName: "bookmark_service",
		InstanceID:  "test-instance-id",
	}
}

func TestHealthCheckHandler_HealthCheck(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/health-check", nil)

	handler := NewHealthCheck(&fakeHealthCheckService{})

	handler.HealthCheck(ctx)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{
		"message": "OK",
		"service_name": "bookmark_service",
		"instance_id": "test-instance-id"
	}`, rec.Body.String())
}