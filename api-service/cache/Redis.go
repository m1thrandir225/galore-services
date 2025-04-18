package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Description:
// A key-value redis compoliant store
type RedisStore struct {
	store *redis.Client
}

// Description:
// Returns a new object of a RedisStore
//
// Parameters:
// address: string,
// password: string
//
// Return:
// *RedisStore
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

// Description:
// Get an item from the store based on the given key
//
// Parameters:
// ctx: context.Context,
// key: string
//
// Return:
// string,
// error
func (redis *RedisStore) GetItem(ctx context.Context, key string) (string, error) {
	value, err := redis.store.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return value, nil
}

// Description:
// Save the item into the current store
//
// Parameters:
// ctx: context.Context,
// key: string,
// value: string //FIXME: this should be interface? and then json serialized
//
// Return:
// error
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

// Description;
// Delete an item based on its key. Error if it doesn't exist.
//
// Parameters:
// ctx: context.Context,
// key: string
//
// Return:
// error
func (redis *RedisStore) DeleteItem(ctx context.Context, key string) error {
	err := redis.store.Del(ctx, key)
	if err != nil {
		return err.Err()
	}
	return nil
}
