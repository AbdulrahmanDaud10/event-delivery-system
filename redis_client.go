package main

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClientInterface interface represents the methods used by the redis client
type RedisClientInterface interface {
	Ping(ctx context.Context) *redis.StatusCmd
	RPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	LLen(ctx context.Context, key string) *redis.IntCmd
	ZAdd(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd
	ZRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.ZSliceCmd
	ZRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd
	BLPop(ctx context.Context, timeout time.Duration, keys ...string) *redis.StringSliceCmd
	ZCard(ctx context.Context, key string) *redis.IntCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Close() error
}

// Wrapping the actual client
type RedisClientWrapper struct {
	Client *redis.Client
}

func (wrapper *RedisClientWrapper) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	return wrapper.Client.Del(ctx, keys...)
}

func (wrapper *RedisClientWrapper) ZCard(ctx context.Context, key string) *redis.IntCmd {
	return wrapper.Client.ZCard(ctx, key)
}

func (wrapper *RedisClientWrapper) BLPop(ctx context.Context, timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	return wrapper.Client.BLPop(ctx, timeout, keys...)
}

func (wrapper *RedisClientWrapper) Close() error {
	return wrapper.Client.Close()
}

func (wrapper *RedisClientWrapper) Ping(ctx context.Context) *redis.StatusCmd {
	return wrapper.Client.Ping(ctx)
}

func (wrapper *RedisClientWrapper) RPush(ctx context.Context, keys string, values ...interface{}) *redis.IntCmd {
	return wrapper.Client.RPush(ctx, keys, values...)
}

func (wrapper *RedisClientWrapper) LLen(ctx context.Context, key string) *redis.IntCmd {
	return wrapper.Client.LLen(ctx, key)
}

func (wrapper *RedisClientWrapper) ZAdd(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd {
	// Unpack members and convert them to []redis.Z
	zMembers := make([]redis.Z, len(members))
	for i, member := range members {
		zMembers[i] = *member
	}

	// Call ZAdd with individual elements
	return wrapper.Client.ZAdd(ctx, key, zMembers...)
}

func (wrapper *RedisClientWrapper) ZRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.ZSliceCmd {
	return wrapper.Client.ZRangeByScoreWithScores(ctx, key, opt)
}

func (wrapper *RedisClientWrapper) ZRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	return wrapper.Client.ZRem(ctx, key, members...)
}
