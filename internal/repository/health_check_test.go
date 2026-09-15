package repository

import (
	"context"
	"errors"
	"testing"

	redismocks "github.com/nguyenthienan91/bookmark-manager/pkg/redis/mocks"
	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

var errRedisPing = errors.New("redis ping error")

func newStatusCommand(ctx context.Context, err error) *goredis.StatusCmd {
	cmd := goredis.NewStatusCmd(ctx)
	cmd.SetVal("PONG")
	cmd.SetErr(err)
	return cmd
}

func TestHealthCheckRepository_Ping(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		setupMock   func(t *testing.T) *redismocks.Store
		expectedErr error
	}{
		{
			name: "redis is reachable",
			setupMock: func(t *testing.T) *redismocks.Store {
				store := redismocks.NewStore(t)
				store.On("Ping", context.Background()).
					Return(newStatusCommand(context.Background(), nil))
				return store
			},
		},
		{
			name: "redis is unavailable",
			setupMock: func(t *testing.T) *redismocks.Store {
				store := redismocks.NewStore(t)
				store.On("Ping", context.Background()).
					Return(newStatusCommand(context.Background(), errRedisPing))
				return store
			},
			expectedErr: errRedisPing,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			store := tc.setupMock(t)
			repo := NewHealthCheckRepository(store)

			err := repo.Ping(context.Background())

			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
				return
			}
			assert.NoError(t, err)
		})
	}
}
