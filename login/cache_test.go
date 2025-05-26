package login

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRedisClient struct {
	mock.Mock
	*redis.Client
}

func NewMockRedisClient() *MockRedisClient {
	return &MockRedisClient{
		Client: redis.NewClient(&redis.Options{}),
	}
}

func (m *MockRedisClient) SMembers(ctx context.Context, key string) *redis.StringSliceCmd {
	args := m.Called(ctx, key)
	cmd := redis.NewStringSliceCmd(ctx)
	if args.Get(0) != nil {
		cmd.SetVal(args.Get(0).([]string))
	}
	cmd.SetErr(args.Error(1))
	return cmd
}

func (m *MockRedisClient) SAdd(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	args := m.Called(ctx, key, members)
	cmd := redis.NewIntCmd(ctx)
	cmd.SetVal(args.Get(0).(int64))
	cmd.SetErr(args.Error(1))
	return cmd
}

func (m *MockRedisClient) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	args := m.Called(ctx, key, expiration)
	cmd := redis.NewBoolCmd(ctx)
	cmd.SetVal(args.Get(0).(bool))
	cmd.SetErr(args.Error(1))
	return cmd
}

func TestInitCache(t *testing.T) {
	InitCache("localhost:6379")
	assert.NotNil(t, rdb)
}

func TestGetUserPermissionsFromCache_Success(t *testing.T) {
	// Setup mock
	mockClient := NewMockRedisClient()
	rdb = mockClient.Client

	// Setup test context
	ctx := context.WithValue(context.Background(), "username", "testuser")

	// Mock expected calls
	mockClient.On("SMembers", ctx, "user_perms:testuser").
		Return(redis.NewStringSliceCmd(ctx, "perm1", "perm2"))

	// Test
	perms := GetUserPermissionsFromCache(ctx)
	assert.Equal(t, []string{"perm1", "perm2"}, perms)
	mockClient.AssertExpectations(t)
}

func TestCacheUserPermission_Success(t *testing.T) {
	// Setup mock
	mockClient := NewMockRedisClient()
	rdb = mockClient.Client

	// Setup test context
	ctx := context.WithValue(context.Background(), "username", "testuser")

	// Mock expected calls
	mockClient.On("SAdd", ctx, "user_perms:testuser", []interface{}{"perm1"}).
		Return(redis.NewIntCmd(ctx))
	mockClient.On("Expire", ctx, "user_perms:testuser", cacheTTL).
		Return(redis.NewBoolCmd(ctx))

	// Test
	CacheUserPermission(ctx, "perm1")
	mockClient.AssertExpectations(t)
}
