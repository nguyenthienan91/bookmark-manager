package integration_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nguyenthienan91/bookmark-manager/internal/api"
	"github.com/stretchr/testify/assert"
)

func TestGenPassEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                 string
		setupTestHTTP        func(api api.Engine) *httptest.ResponseRecorder
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "success",
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/generate-password", nil)
				responseRecorder := httptest.NewRecorder()
				api.ServeHTTP(responseRecorder, req)
				return responseRecorder
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `"password"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			apiEngine := api.NewEngine(&api.Config{})
			responseRecorder := tc.setupTestHTTP(apiEngine)

			assert.Equal(t, tc.expectedStatusCode, responseRecorder.Code)
			assert.Contains(t, responseRecorder.Body.String(), tc.expectedResponseBody)
		})
	}
}