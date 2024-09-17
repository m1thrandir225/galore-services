package cache

import "context"

/*
 * Interface for a key-value store
 */
type KeyValueStore interface {
	SaveItem(ctx context.Context, name string) error
	GetItem(ctx context.Context, name string) (string, error)
	DeleteItem(ctx context.Context, name string) error
}
