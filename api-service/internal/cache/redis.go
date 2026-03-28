package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisStore implements the Store interface
type RedisStore struct {
	store *redis.Client
}

// NewRedisStore returns a RedisStore implementation
func NewRedisStore(address, password string) (*RedisStore, error) {
	if len(address) == 0 {
		return nil, errors.New("redis address is empty")
	}
	if len(password) == 0 {
		return nil, errors.New("redis password is empty")
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	})

	return &RedisStore{
		store: rdb,
	}, nil
}

func (redis *RedisStore) Get(ctx context.Context, key string) (string, error) {
	value, err := redis.store.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return value, nil
}

func (redis *RedisStore) Save(ctx context.Context, key, value string) error {
	startTime := time.Now()

	endTime := startTime.AddDate(0, 0, 2)

	duration := endTime.Sub(startTime)

	err := redis.store.Set(ctx, key, value, duration)
	if err != nil {
		return err.Err()
	}
	return nil
}

func (redis *RedisStore) Delete(ctx context.Context, key string) error {
	err := redis.store.Del(ctx, key)
	if err != nil {
		return err.Err()
	}
	return nil
}
