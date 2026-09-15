package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	svcmocks "github.com/nguyenthienan91/bookmark-manager/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

var errShortenSvc = errors.New("shorten service error")

func TestShortenLinkHandler_ShortenLink(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                 string
		requestBody          string
		setupMockService     func(ctx context.Context) *svcmocks.ShortenLink
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "success",
			requestBody: `{"url":"https://example.com","exp":604800}`,
			setupMockService: func(ctx context.Context) *svcmocks.ShortenLink {
				m := svcmocks.NewShortenLink(t)
				m.On("Shorten", ctx, "https://example.com", 604800).Return("abc1234", nil)
				return m
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"code":"abc1234","message":"Shorten URL generated successfully!"}`,
		},
		{
			name:        "success - no expiration",
			requestBody: `{"url":"https://example.com"}`,
			setupMockService: func(ctx context.Context) *svcmocks.ShortenLink {
				m := svcmocks.NewShortenLink(t)
				m.On("Shorten", ctx, "https://example.com", 0).Return("xyz7890", nil)
				return m
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"code":"xyz7890","message":"Shorten URL generated successfully!"}`,
		},
		{
			name:        "missing url - returns 400",
			requestBody: `{"exp":100}`,
			setupMockService: func(ctx context.Context) *svcmocks.ShortenLink {
				return svcmocks.NewShortenLink(t) // never called
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"Key: 'ShortenLinkRequest.URL' Error:Field validation for 'URL' failed on the 'required' tag"}`,
		},
		{
			name:        "invalid json - returns 400",
			requestBody: `not-json`,
			setupMockService: func(ctx context.Context) *svcmocks.ShortenLink {
				return svcmocks.NewShortenLink(t)
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:        "invalid url format - returns 400",
			requestBody: `{"url":"not-a-url"}`,
			setupMockService: func(ctx context.Context) *svcmocks.ShortenLink {
				return svcmocks.NewShortenLink(t)
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:        "negative exp - returns 400",
			requestBody: `{"url":"https://example.com","exp":-1}`,
			setupMockService: func(ctx context.Context) *svcmocks.ShortenLink {
				return svcmocks.NewShortenLink(t)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"exp must be non-negative"}`,
		},
		{
			name:        "service error - returns 500",
			requestBody: `{"url":"https://example.com","exp":0}`,
			setupMockService: func(ctx context.Context) *svcmocks.ShortenLink {
				m := svcmocks.NewShortenLink(t)
				m.On("Shorten", ctx, "https://example.com", 0).Return("", errShortenSvc)
				return m
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"shorten service error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/links/shorten",
				bytes.NewBufferString(tc.requestBody))
			ctx.Request.Header.Set("Content-Type", "application/json")

			mockSvc := tc.setupMockService(ctx.Request.Context())
			h := NewShortenLink(mockSvc)
			h.ShortenLink(ctx)

			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			if tc.expectedResponseBody != "" {
				assert.JSONEq(t, tc.expectedResponseBody, rec.Body.String())
			}
		})
	}
}

