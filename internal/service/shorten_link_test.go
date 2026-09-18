package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	repomocks "github.com/nguyenthienan91/bookmark-manager/internal/repository/mocks"
	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var errRepo = errors.New("repository error")

func shortCodeKeyMatcher() interface{} {
	return mock.MatchedBy(func(key string) bool {
		if !strings.HasPrefix(key, shortCodeKeyPrefix) {
			return false
		}

		code := strings.TrimPrefix(key, shortCodeKeyPrefix)
		if len(code) != shortCodeLength {
			return false
		}

		for _, ch := range code {
			if !strings.ContainsRune(shortCodeCharset, ch) {
				return false
			}
		}
		return true
	})
}

func TestShortenLinkService_Shorten(t *testing.T) {
	t.Parallel()

	const testURL = "https://example.com/article"

	testCases := []struct {
		name        string
		url         string
		expSeconds  int
		setupMock   func(t *testing.T) *repomocks.LinkRepository
		assertCode  func(t *testing.T, code string)
		expectedErr error
	}{
		{
			name:       "success - code has 7 chars",
			url:        testURL,
			expSeconds: 0,
			setupMock: func(t *testing.T) *repomocks.LinkRepository {
				m := repomocks.NewLinkRepository(t)
				m.On("SaveLink", context.Background(), shortCodeKeyMatcher(), testURL, time.Duration(0)).
					Return(true, nil)
				return m
			},
			assertCode: func(t *testing.T, code string) {
				assert.Len(t, code, 7)
			},
			expectedErr: nil,
		},
		{
			name:       "success - code is alphanumeric only",
			url:        testURL,
			expSeconds: 0,
			setupMock: func(t *testing.T) *repomocks.LinkRepository {
				m := repomocks.NewLinkRepository(t)
				m.On("SaveLink", context.Background(), shortCodeKeyMatcher(), testURL, time.Duration(0)).
					Return(true, nil)
				return m
			},
			assertCode: func(t *testing.T, code string) {
				for _, ch := range code {
					assert.True(t,
						(ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9'),
						"unexpected char %q in code %q", ch, code,
					)
				}
			},
			expectedErr: nil,
		},
		{
			name:       "success - expiration converted from seconds",
			url:        testURL,
			expSeconds: 604800, // 7 days
			setupMock: func(t *testing.T) *repomocks.LinkRepository {
				m := repomocks.NewLinkRepository(t)
				m.On("SaveLink", context.Background(), shortCodeKeyMatcher(), testURL, 604800*time.Second).
					Return(true, nil)
				return m
			},
			assertCode:  func(t *testing.T, code string) { assert.Len(t, code, 7) },
			expectedErr: nil,
		},
		{
			name:       "success - zero exp means no TTL",
			url:        testURL,
			expSeconds: 0,
			setupMock: func(t *testing.T) *repomocks.LinkRepository {
				m := repomocks.NewLinkRepository(t)
				m.On("SaveLink", context.Background(), shortCodeKeyMatcher(), testURL, time.Duration(0)).
					Return(true, nil)
				return m
			},
			assertCode:  func(t *testing.T, code string) { assert.Len(t, code, 7) },
			expectedErr: nil,
		},
		{
			name:       "repository returns error",
			url:        testURL,
			expSeconds: 0,
			setupMock: func(t *testing.T) *repomocks.LinkRepository {
				m := repomocks.NewLinkRepository(t)
				m.On("SaveLink", context.Background(), shortCodeKeyMatcher(), testURL, time.Duration(0)).
					Return(false, errRepo)
				return m
			},
			assertCode:  func(t *testing.T, code string) { assert.Empty(t, code) },
			expectedErr: errRepo,
		},
		{
			name:       "retry on collision - succeeds on second attempt",
			url:        testURL,
			expSeconds: 0,
			setupMock: func(t *testing.T) *repomocks.LinkRepository {
				m := repomocks.NewLinkRepository(t)
				// first call: collision; second call: success
				m.On("SaveLink", context.Background(), shortCodeKeyMatcher(), testURL, time.Duration(0)).
					Return(false, nil).Once()
				m.On("SaveLink", context.Background(), shortCodeKeyMatcher(), testURL, time.Duration(0)).
					Return(true, nil).Once()
				return m
			},
			assertCode:  func(t *testing.T, code string) { assert.Len(t, code, 7) },
			expectedErr: nil,
		},
		{
			name:       "retry exhausted - all attempts collide",
			url:        testURL,
			expSeconds: 0,
			setupMock: func(t *testing.T) *repomocks.LinkRepository {
				m := repomocks.NewLinkRepository(t)
				// all maxRetryCount attempts return collision
				m.On("SaveLink", context.Background(), shortCodeKeyMatcher(), testURL, time.Duration(0)).
					Return(false, nil).Times(maxRetryCount)
				return m
			},
			assertCode:  func(t *testing.T, code string) { assert.Empty(t, code) },
			expectedErr: errors.New("failed to generate a unique short code after 10 retries"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockRepo := tc.setupMock(t)
			svc := NewShortenLink(mockRepo)

			code, err := svc.Shorten(context.Background(), tc.url, tc.expSeconds)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
			tc.assertCode(t, code)
		})
	}
}

func TestShortenLinkService_Resolve(t *testing.T) {
	t.Parallel()

	const testURL = "https://example.com/article"

	testCases := []struct {
		name        string
		code        string
		setupMock   func(t *testing.T) *repomocks.LinkRepository
		expectedURL string
		expectedErr error
	}{
		{
			name: "success - resolves URL",
			code: "abc1234",
			setupMock: func(t *testing.T) *repomocks.LinkRepository {
				m := repomocks.NewLinkRepository(t)
				m.On("GetLink", context.Background(), "shorturl:abc1234").Return(testURL, nil)
				return m
			},
			expectedURL: testURL,
		},
		{
			name: "not found - redis.Nil becomes ErrLinkNotFound",
			code: "missing1",
			setupMock: func(t *testing.T) *repomocks.LinkRepository {
				m := repomocks.NewLinkRepository(t)
				m.On("GetLink", context.Background(), "shorturl:missing1").Return("", goredis.Nil)
				return m
			},
			expectedErr: ErrLinkNotFound,
		},
		{
			name: "repository error - propagated",
			code: "err0001",
			setupMock: func(t *testing.T) *repomocks.LinkRepository {
				m := repomocks.NewLinkRepository(t)
				m.On("GetLink", context.Background(), "shorturl:err0001").Return("", errRepo)
				return m
			},
			expectedErr: errRepo,
		},
		{
			name: "builds correct key with prefix",
			code: "XyZ9876",
			setupMock: func(t *testing.T) *repomocks.LinkRepository {
				m := repomocks.NewLinkRepository(t)
				m.On("GetLink", context.Background(), "shorturl:XyZ9876").Return(testURL, nil)
				return m
			},
			expectedURL: testURL,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockRepo := tc.setupMock(t)
			svc := NewShortenLink(mockRepo)

			url, err := svc.Resolve(context.Background(), tc.code)

			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
				assert.Empty(t, url)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedURL, url)
		})
	}
}

