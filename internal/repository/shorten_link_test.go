package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	redismocks "github.com/nguyenthienan91/bookmark-manager/pkg/redis/mocks"
	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

const (
	testKey = "shorturl:abc1234"
	testURL = "https://example.com"
)

var errRedis = errors.New("redis error")

func newBoolCommand(ctx context.Context, value bool, err error) *goredis.BoolCmd {
	cmd := goredis.NewBoolCmd(ctx)
	cmd.SetVal(value)
	cmd.SetErr(err)
	return cmd
}

func TestLinkRepository_SaveLink(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		setupMock       func(t *testing.T) *redismocks.Store
		key             string
		url             string
		expiration      time.Duration
		expectedCreated bool
		expectedErr     error
	}{
		{
			name: "save successfully",
			setupMock: func(t *testing.T) *redismocks.Store {
				m := redismocks.NewStore(t)
				m.On("SetNX", context.Background(), testKey, testURL, 7*24*time.Hour).
					Return(newBoolCommand(context.Background(), true, nil))
				return m
			},
			key:             testKey,
			url:             testURL,
			expiration:      7 * 24 * time.Hour,
			expectedCreated: true,
		},
		{
			name: "key already exists",
			setupMock: func(t *testing.T) *redismocks.Store {
				m := redismocks.NewStore(t)
				m.On("SetNX", context.Background(), testKey, testURL, 7*24*time.Hour).
					Return(newBoolCommand(context.Background(), false, nil))
				return m
			},
			key:             testKey,
			url:             testURL,
			expiration:      7 * 24 * time.Hour,
			expectedCreated: false,
		},
		{
			name: "redis error",
			setupMock: func(t *testing.T) *redismocks.Store {
				m := redismocks.NewStore(t)
				m.On("SetNX", context.Background(), testKey, testURL, 7*24*time.Hour).
					Return(newBoolCommand(context.Background(), false, errRedis))
				return m
			},
			key:         testKey,
			url:         testURL,
			expiration:  7 * 24 * time.Hour,
			expectedErr: errRedis,
		},
		{
			name: "uses the provided key",
			setupMock: func(t *testing.T) *redismocks.Store {
				m := redismocks.NewStore(t)
				m.On("SetNX", context.Background(), "custom:key", testURL, time.Duration(0)).
					Return(newBoolCommand(context.Background(), true, nil))
				return m
			},
			key:             "custom:key",
			url:             testURL,
			expectedCreated: true,
		},
		{
			name: "passes expiration correctly",
			setupMock: func(t *testing.T) *redismocks.Store {
				m := redismocks.NewStore(t)
				m.On("SetNX", context.Background(), testKey, testURL, 30*time.Second).
					Return(newBoolCommand(context.Background(), true, nil))
				return m
			},
			key:             testKey,
			url:             testURL,
			expiration:      30 * time.Second,
			expectedCreated: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			store := tc.setupMock(t)
			repo := NewLinkRepository(store)

			created, err := repo.SaveLink(
				context.Background(),
				tc.key,
				tc.url,
				tc.expiration,
			)

			assert.Equal(t, tc.expectedCreated, created)
			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func newStringCommand(ctx context.Context, value string, err error) *goredis.StringCmd {
	cmd := goredis.NewStringCmd(ctx)
	cmd.SetVal(value)
	cmd.SetErr(err)
	return cmd
}

func TestLinkRepository_GetLink(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		setupMock   func(t *testing.T) *redismocks.Store
		key         string
		expectedURL string
		expectedErr error
	}{
		{
			name: "success",
			setupMock: func(t *testing.T) *redismocks.Store {
				m := redismocks.NewStore(t)
				m.On("Get", context.Background(), testKey).
					Return(newStringCommand(context.Background(), testURL, nil))
				return m
			},
			key:         testKey,
			expectedURL: testURL,
		},
		{
			name: "key not found - redis.Nil",
			setupMock: func(t *testing.T) *redismocks.Store {
				m := redismocks.NewStore(t)
				m.On("Get", context.Background(), testKey).
					Return(newStringCommand(context.Background(), "", goredis.Nil))
				return m
			},
			key:         testKey,
			expectedErr: goredis.Nil,
		},
		{
			name: "redis error",
			setupMock: func(t *testing.T) *redismocks.Store {
				m := redismocks.NewStore(t)
				m.On("Get", context.Background(), testKey).
					Return(newStringCommand(context.Background(), "", errRedis))
				return m
			},
			key:         testKey,
			expectedErr: errRedis,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			store := tc.setupMock(t)
			repo := NewLinkRepository(store)

			url, err := repo.GetLink(context.Background(), tc.key)

			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedURL, url)
		})
	}
}

