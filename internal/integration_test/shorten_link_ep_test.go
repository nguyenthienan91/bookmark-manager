package integration_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nguyenthienan91/bookmark-manager/internal/api"
	svcmocks "github.com/nguyenthienan91/bookmark-manager/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestShortenLinkEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                 string
		requestBody          string
		setupMockService     func(t *testing.T) *svcmocks.ShortenLink
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "success",
			requestBody: `{"url":"https://example.com/article","exp":604800}`,
			setupMockService: func(t *testing.T) *svcmocks.ShortenLink {
				m := svcmocks.NewShortenLink(t)
				m.On("Shorten", context.Background(), "https://example.com/article", 604800).
					Return("abc1234", nil)
				return m
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"code":"abc1234","message":"Shorten URL generated successfully!"}`,
		},
		{
			name:        "missing url - 400",
			requestBody: `{"exp":100}`,
			setupMockService: func(t *testing.T) *svcmocks.ShortenLink {
				return svcmocks.NewShortenLink(t)
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:        "invalid url format - 400",
			requestBody: `{"url":"not-a-url"}`,
			setupMockService: func(t *testing.T) *svcmocks.ShortenLink {
				return svcmocks.NewShortenLink(t)
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:        "negative exp - 400",
			requestBody: `{"url":"https://example.com","exp":-5}`,
			setupMockService: func(t *testing.T) *svcmocks.ShortenLink {
				return svcmocks.NewShortenLink(t)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"exp must be non-negative"}`,
		},
		{
			name:        "service error - 500",
			requestBody: `{"url":"https://example.com","exp":0}`,
			setupMockService: func(t *testing.T) *svcmocks.ShortenLink {
				m := svcmocks.NewShortenLink(t)
				m.On("Shorten", context.Background(), "https://example.com", 0).
					Return("", errors.New("internal error"))
				return m
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"internal error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockSvc := tc.setupMockService(t)
			apiEngine := api.NewEngine(&api.Config{
				ServiceName: "bookmark_service",
				InstanceID:  "test-instance-id",
			}, mockSvc)

			req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten",
				bytes.NewBufferString(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			apiEngine.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			if tc.expectedResponseBody != "" {
				assert.JSONEq(t, tc.expectedResponseBody, rec.Body.String())
			}
		})
	}
}

