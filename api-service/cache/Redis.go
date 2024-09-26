package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

/**
* Store that should work with any db that supports the redis protocol (Dragonfly
  * etc..)
*/
type RedisStore struct {
	store *redis.Client
}

func NewRedisStore(address, password string) *RedisStore {
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	})

	return &RedisStore{
		store: rdb,
	}
}

func (redis *RedisStore) GetItem(ctx context.Context, key string) (string, error) {
	value, err := redis.store.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return value, nil
}

func (redis *RedisStore) SaveItem(ctx context.Context, key, value string) error {
	startTime := time.Now()

	endTime := startTime.AddDate(0, 0, 2)

	duration := endTime.Sub(startTime)

	err := redis.store.Set(ctx, key, value, duration)
	if err != nil {
		return err.Err()
	}
	return nil
}

func (redis *RedisStore) DeleteItem(ctx context.Context, key string) error {
	err := redis.store.Del(ctx, key)
	if err != nil {
		return err.Err()
	}
	return nil
}
