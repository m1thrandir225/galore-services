package cache

import "context"

/*
 * Interface for a key-value store
 */
type KeyValueStore interface {
	SaveItem(ctx context.Context, key, value string) error
	GetItem(ctx context.Context, key string) (string, error)
	DeleteItem(ctx context.Context, key string) error
}
