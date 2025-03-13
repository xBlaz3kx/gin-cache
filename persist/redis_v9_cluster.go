package persist

import (
	"context"
	"errors"
	"time"

	redisv9 "github.com/redis/go-redis/v9"
)

type RedisClusterStore struct {
	RedisClient *redisv9.ClusterClient `inject:""`
}

// NewRedisClusterStore create a redis memory store with redis client
func NewRedisClusterStore(redisClient *redisv9.ClusterClient) *RedisClusterStore {
	return &RedisClusterStore{
		RedisClient: redisClient,
	}
}

// Set put key value pair to redis, and expire after expireDuration
func (store *RedisClusterStore) Set(ctx context.Context, key string, value interface{}, expire time.Duration) error {
	payload, err := Serialize(value)
	if err != nil {
		return err
	}

	return store.RedisClient.Set(ctx, key, payload, expire).Err()
}

// Delete remove key in redis, do nothing if key doesn't exist
func (store *RedisClusterStore) Delete(ctx context.Context, key string) error {
	return store.RedisClient.Del(ctx, key).Err()
}

// Get retrieves an item from redis, if key doesn't exist, return ErrCacheMiss
func (store *RedisClusterStore) Get(ctx context.Context, key string, value interface{}) error {
	payload, err := store.RedisClient.Get(ctx, key).Bytes()
	if errors.Is(err, redisv9.Nil) {
		return ErrCacheMiss
	}

	if err != nil {
		return err
	}
	return Deserialize(payload, value)
}
