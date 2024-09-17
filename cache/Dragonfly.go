package cache

import (
	"context"

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

func (dragonfly *DragonflyStore) GetItem(ctx context.Context, name string) (string, error) {
  err := dragonfly.store.Set(ctx, , value interface{}, expiration time.Duration)
	return "", nil
}

func (dragonfly *DragonflyStore) SaveItem(ctx context.Context, name string) error {
	return nil
}

func (dragonfly *DragonflyStore) DeleteItem(ctx context.Context name string) error {
	return nil
}
