package login

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/exp/slog"
)

var (
	rdb      *redis.Client
	cacheTTL = 5 * time.Minute
)

func InitCache(redisAddr string) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	slog.Info("Redis cache initialized", "addr", redisAddr)
}

func GetUserPermissionsFromCache(ctx context.Context) []string {
	if rdb == nil {
		slog.Info("Redis client not initialized")
		return nil
	}

	username, ok := ctx.Value("username").(string)
	if !ok {
		slog.Info("Username not found in context", "context", ctx)
		return nil
	}

	key := "user_perms:" + username
	perms, err := rdb.SMembers(ctx, key).Result()
	if err != nil {
		slog.Error("Failed to get permissions from cache", "error", err, "username", username, "key", key)
		return nil
	}

	slog.Info("Retrieved permissions from cache", "username", username, "permissions", perms)
	return perms
}

func CacheUserPermission(ctx context.Context, permission string) {
	if rdb == nil {
		slog.Info("Redis client not initialized")
		return
	}

	username, ok := ctx.Value("username").(string)
	if !ok {
		slog.Info("Username not found in context", "context", ctx)
		return
	}

	key := "user_perms:" + username
	err := rdb.SAdd(ctx, key, permission).Err()
	if err != nil {
		slog.Error("Failed to cache permission", "error", err, "username", username, "permission", permission)
		return
	}

	err = rdb.Expire(ctx, key, cacheTTL).Err()
	if err != nil {
		slog.Error("Failed to set cache expiry", "error", err, "username", username, "key", key)
		return
	}

	slog.Info("Cached user permission", "username", username, "permission", permission, "ttl", cacheTTL)
}
