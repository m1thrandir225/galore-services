package api

import (
	"fmt"

	"github.com/m1thrandir225/galore-services/cache"
	embedding "github.com/m1thrandir225/galore-services/embedding_service"
	"github.com/m1thrandir225/galore-services/storage"

	"github.com/gin-gonic/gin"
	db "github.com/m1thrandir225/galore-services/db/sqlc"
	"github.com/m1thrandir225/galore-services/token"
	"github.com/m1thrandir225/galore-services/util"
)

type Server struct {
	config     util.Config
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
	storage    storage.FileService
	cache      cache.KeyValueStore
	embedding  embedding.EmbeddingService
}

func NewServer(config util.Config, store db.Store, storageService storage.FileService, cacheStore cache.KeyValueStore, embedding embedding.EmbeddingService) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
		storage:    storageService,
		cache:      cacheStore,
		embedding:  embedding,
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
