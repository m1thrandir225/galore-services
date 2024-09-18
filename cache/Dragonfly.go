package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type DragonflyStore struct {
	store *redis.Client
}

func NewDragonflyStore(address, password string) *DragonflyStore {
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	})

	return &DragonflyStore{
		store: rdb,
	}
}

func (dragonfly *DragonflyStore) GetItem(ctx context.Context, key string) (string, error) {
	value, err := dragonfly.store.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return value, nil
}

func (dragonfly *DragonflyStore) SaveItem(ctx context.Context, key, value string) error {
	startTime := time.Now()

	endTime := startTime.AddDate(0, 0, 2)

	duration := endTime.Sub(startTime)

	err := dragonfly.store.Set(ctx, key, value, duration)
	if err != nil {
		return err.Err()
	}
	return nil
}

func (dragonfly *DragonflyStore) DeleteItem(ctx context.Context, key string) error {
	err := dragonfly.store.Del(ctx, key)
	if err != nil {
		return err.Err()
	}
	return nil
}
