package cache

import "context"

// Description:
// Generic Key-Value Store interface
type KeyValueStore interface {
	SaveItem(ctx context.Context, key, value string) error
	GetItem(ctx context.Context, key string) (string, error)
	DeleteItem(ctx context.Context, key string) error
}
