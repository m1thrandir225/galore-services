// Package server all http related logic and routing
package server

import (
	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/galore-services/internal/cache"
	"github.com/m1thrandir225/galore-services/internal/config"
	db "github.com/m1thrandir225/galore-services/internal/db/sqlc"
	"github.com/m1thrandir225/galore-services/internal/security"
)

type Server struct {
	config     *config.Config
	store      db.Store
	cache      cache.Store
	tokenMaker security.Maker
	router     *gin.Engine
}

func New(
	config *config.Config,
	store db.Store,
	cache cache.Store,
	tokenMaker security.Maker,
) (*Server, error) {
	if config == nil {
		return nil, ErrInvalidServerConfig
	}
	if store == nil {
		return nil, ErrInvalidStore
	}

	if cache == nil {
		return nil, ErrInvalidCache
	}

	if tokenMaker == nil {
		return nil, ErrInvalidTokenMaker
	}

	router := gin.Default()

	return &Server{
		config:     config,
		store:      store,
		cache:      cache,
		tokenMaker: tokenMaker,
		router:     router,
	}, nil
}

func (s *Server) Run() error {
	return s.router.Run(s.config.HTTPServerAddress)
}
