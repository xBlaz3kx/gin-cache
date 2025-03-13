package persist

import (
	"context"
	"errors"
	"time"

	redisv9 "github.com/redis/go-redis/v9"
)

type RedisV9Store struct {
	RedisClient *redisv9.Client
}

// NewRedisV9Store create a redis memory store with redis client
func NewRedisV9Store(redisClient *redisv9.Client) *RedisV9Store {
	return &RedisV9Store{
		RedisClient: redisClient,
	}
}

// Set put key value pair to redis, and expire after expireDuration
func (store *RedisV9Store) Set(ctx context.Context, key string, value interface{}, expire time.Duration) error {
	payload, err := Serialize(value)
	if err != nil {
		return err
	}

	return store.RedisClient.Set(ctx, key, payload, expire).Err()
}

// Delete remove key in redis, do nothing if key doesn't exist
func (store *RedisV9Store) Delete(ctx context.Context, key string) error {
	return store.RedisClient.Del(ctx, key).Err()
}

// Get retrieves an item from redis, if key doesn't exist, return ErrCacheMiss
func (store *RedisV9Store) Get(ctx context.Context, key string, value interface{}) error {
	payload, err := store.RedisClient.Get(ctx, key).Bytes()
	if errors.Is(err, redisv9.Nil) {
		return ErrCacheMiss
	}

	if err != nil {
		return err
	}
	return Deserialize(payload, value)
}
