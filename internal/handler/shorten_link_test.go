package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nguyenthienan91/bookmark-manager/internal/service"
	svcmocks "github.com/nguyenthienan91/bookmark-manager/internal/service/mocks"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

var (
	errShortenSvc = errors.New("shorten service error")
	testLogger    = zerolog.New(os.Stdout)
)

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
			h := NewShortenLink(mockSvc, testLogger)
			h.ShortenLink(ctx)

			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			if tc.expectedResponseBody != "" {
				assert.JSONEq(t, tc.expectedResponseBody, rec.Body.String())
			}
		})
	}
}

func TestShortenLinkHandler_RedirectLink(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                 string
		code                 string
		setupMockService     func(ctx context.Context) *svcmocks.ShortenLink
		expectedStatusCode   int
		expectedResponseBody string
		expectedLocation     string
	}{
		{
			name: "success - 302 redirect",
			code: "abc1234",
			setupMockService: func(ctx context.Context) *svcmocks.ShortenLink {
				m := svcmocks.NewShortenLink(t)
				m.On("Resolve", ctx, "abc1234").Return("https://example.com", nil)
				return m
			},
			expectedStatusCode: http.StatusFound,
			expectedLocation:   "https://example.com",
		},
		{
			name: "empty code - returns 400",
			code: "",
			setupMockService: func(ctx context.Context) *svcmocks.ShortenLink {
				return svcmocks.NewShortenLink(t)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"invalid or missing code"}`,
		},
		{
			name: "invalid code format - returns 400",
			code: "abc-123!",
			setupMockService: func(ctx context.Context) *svcmocks.ShortenLink {
				return svcmocks.NewShortenLink(t)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"invalid or missing code"}`,
		},
		{
			name: "not found - returns 404",
			code: "xyz9999",
			setupMockService: func(ctx context.Context) *svcmocks.ShortenLink {
				m := svcmocks.NewShortenLink(t)
				m.On("Resolve", ctx, "xyz9999").Return("", service.ErrLinkNotFound)
				return m
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: `{"error":"short link not found"}`,
		},
		{
			name: "internal error - returns 500",
			code: "err5678",
			setupMockService: func(ctx context.Context) *svcmocks.ShortenLink {
				m := svcmocks.NewShortenLink(t)
				m.On("Resolve", ctx, "err5678").Return("", errors.New("redis connection lost"))
				return m
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"internal server error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/links/redirect/"+tc.code, nil)
			ctx.Params = gin.Params{{Key: "code", Value: tc.code}}

			mockSvc := tc.setupMockService(ctx.Request.Context())
			h := NewShortenLink(mockSvc, testLogger)
			h.RedirectLink(ctx)

			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			if tc.expectedResponseBody != "" {
				assert.JSONEq(t, tc.expectedResponseBody, rec.Body.String())
			}
			if tc.expectedLocation != "" {
				assert.Equal(t, tc.expectedLocation, rec.Header().Get("Location"))
			}
		})
	}
}


