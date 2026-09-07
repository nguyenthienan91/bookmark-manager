package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nguyenthienan91/bookmark-manager/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

var testError = errors.New("test error")

func TestGeneratePasswordHandler_GeneratePassword(t *testing.T) {

	t.Parallel()
	testCases := []struct {
		name string
		setupRequest func(ctx *gin.Context)
	  setupMockService func(ctx context.Context) *mocks.GenPass
		expectedStatusCode int
		expectedResponseBody string
	}{
			{
				name: "success",
				setupRequest: func(ctx *gin.Context) {
					ctx.Request = httptest.NewRequest(http.MethodGet, "/generate-password", nil)
			},
				setupMockService: func(ctx context.Context) *mocks.GenPass {
					service := mocks.NewGenPass(t)
					service.On("GeneratePassword", 12).Return("123456789012", nil)
					return service
				},
				expectedStatusCode: http.StatusOK,
				expectedResponseBody: `{"password":"123456789012"}`,
			},
			{
				name: "service failed",
				setupRequest: func(ctx *gin.Context) {
					ctx.Request = httptest.NewRequest(http.MethodGet, "/generate-password", nil)
			},
				setupMockService: func(ctx context.Context) *mocks.GenPass {
					service := mocks.NewGenPass(t)
					service.On("GeneratePassword", 12).Return("", testError)
					return service
				},
				expectedStatusCode: http.StatusInternalServerError,
				expectedResponseBody: `{"error":"Failed to generate password"}`,
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				rec := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(rec)
				tc.setupRequest(ctx)
				mockService := tc.setupMockService(ctx)
				testHandler := NewGenPass(mockService)
				testHandler.GeneratePassword(ctx)

				assert.Equal(t, tc.expectedStatusCode, rec.Code)
				assert.JSONEq(t, tc.expectedResponseBody, rec.Body.String())

				
			})
		}
	}