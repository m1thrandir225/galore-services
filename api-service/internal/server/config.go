package httpserver

import (
	"github.com/m1thrandir225/galore-services/internal/cache"
	"github.com/m1thrandir225/galore-services/internal/config"
	db "github.com/m1thrandir225/galore-services/internal/db/sqlc"
)

type ServerConfig struct {
	AppConfig *config.Config
	Store     db.Store
	Cache     cache.Store
}

func NewServerConfig(config *config.Config, store db.Store, cache cache.Store) (*ServerConfig, error) {
	if config == nil {
		return nil, ErrInvalidGlobalConfig
	}
	if store == nil {
		return nil, ErrInvalidStore
	}

	if cache == nil {
		return nil, ErrInvalidCache
	}

	return &ServerConfig{
		AppConfig: config,
		Store:     store,
		Cache:     cache,
	}, nil
}
