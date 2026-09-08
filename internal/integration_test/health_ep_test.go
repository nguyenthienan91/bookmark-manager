package integration_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nguyenthienan91/bookmark-manager/internal/api"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheckEndpoint(t *testing.T) {
	t.Parallel()

	apiEngine := api.NewEngine(&api.Config{
		ServiceName: "bookmark_service",
		InstanceID:  "test-instance-id",
	})

	req := httptest.NewRequest(http.MethodGet, "/health-check", nil)
	rec := httptest.NewRecorder()

	apiEngine.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{
		"message": "OK",
		"service_name": "bookmark_service",
		"instance_id": "test-instance-id"
	}`, rec.Body.String())
}