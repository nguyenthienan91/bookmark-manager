package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	redismocks "github.com/nguyenthienan91/bookmark-manager/pkg/redis/mocks"
	"github.com/stretchr/testify/assert"
)

var errRedis = errors.New("redis error")

func TestLinkRepository_SaveLink(t *testing.T) {
	t.Parallel()

	const (
		testCode = "abc1234"
		testURL  = "https://example.com"
		testKey  = "shorturl:" + testCode
	)
	testExp := 7 * 24 * time.Hour

	testCases := []struct {
		name            string
		setupMock       func(t *testing.T) *redismocks.Store
		code            string
		url             string
		expiration      time.Duration
		expectedCreated bool
		expectedErr     error
	}{
		{
			name: "save successfully",
			setupMock: func(t *testing.T) *redismocks.Store {
				m := redismocks.NewStore(t)
				m.On("SetNX", context.Background(), testKey, testURL, testExp).Return(true, nil)
				return m
			},
			code:            testCode,
			url:             testURL,
			expiration:      testExp,
			expectedCreated: true,
			expectedErr:     nil,
		},
		{
			name: "key already exists",
			setupMock: func(t *testing.T) *redismocks.Store {
				m := redismocks.NewStore(t)
				m.On("SetNX", context.Background(), testKey, testURL, testExp).Return(false, nil)
				return m
			},
			code:            testCode,
			url:             testURL,
			expiration:      testExp,
			expectedCreated: false,
			expectedErr:     nil,
		},
		{
			name: "redis error",
			setupMock: func(t *testing.T) *redismocks.Store {
				m := redismocks.NewStore(t)
				m.On("SetNX", context.Background(), testKey, testURL, testExp).Return(false, errRedis)
				return m
			},
			code:            testCode,
			url:             testURL,
			expiration:      testExp,
			expectedCreated: false,
			expectedErr:     errRedis,
		},
		{
			name: "key has correct prefix",
			setupMock: func(t *testing.T) *redismocks.Store {
				m := redismocks.NewStore(t)
				// assert the key passed to Redis includes the shorturl: prefix
				m.On("SetNX", context.Background(), "shorturl:"+testCode, testURL, time.Duration(0)).Return(true, nil)
				return m
			},
			code:            testCode,
			url:             testURL,
			expiration:      0,
			expectedCreated: true,
			expectedErr:     nil,
		},
		{
			name: "expiration passed correctly",
			setupMock: func(t *testing.T) *redismocks.Store {
				m := redismocks.NewStore(t)
				m.On("SetNX", context.Background(), testKey, testURL, 30*time.Second).Return(true, nil)
				return m
			},
			code:            testCode,
			url:             testURL,
			expiration:      30 * time.Second,
			expectedCreated: true,
			expectedErr:     nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockStore := tc.setupMock(t)
			repo := NewLinkRepository(mockStore)

			created, err := repo.SaveLink(context.Background(), tc.code, tc.url, tc.expiration)

			assert.ErrorIs(t, err, tc.expectedErr)
			assert.Equal(t, tc.expectedCreated, created)
		})
	}
}

