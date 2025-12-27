// Package cache provides utilities for managing key-value cache saving, validation & invalidation
package cache

import "context"

type Store interface {
	Save(ctx context.Context, key, value string) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}
