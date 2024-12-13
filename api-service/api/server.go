package api

import (
	"fmt"
	"net/http"

	"github.com/m1thrandir225/galore-services/background_jobs"
	"github.com/m1thrandir225/galore-services/cache"
	categorizer "github.com/m1thrandir225/galore-services/categorizer_service"
	embedding "github.com/m1thrandir225/galore-services/embedding_service"
	"github.com/m1thrandir225/galore-services/storage"

	"github.com/gin-gonic/gin"
	db "github.com/m1thrandir225/galore-services/db/sqlc"
	"github.com/m1thrandir225/galore-services/token"
	"github.com/m1thrandir225/galore-services/util"
)

type Server struct {
	config      util.Config
	store       db.Store
	router      *gin.Engine
	tokenMaker  token.Maker
	storage     storage.FileService
	cache       cache.KeyValueStore
	embedding   embedding.EmbeddingService
	categorizer categorizer.CategorizerService
	scheduler   background_jobs.SchedulerService
}

type HealthResponse struct {
	Status string `json:"status"`
}

func NewServer(
	config util.Config,
	store db.Store,
	storageService storage.FileService,
	cacheStore cache.KeyValueStore,
	embedding embedding.EmbeddingService,
	categorizer categorizer.CategorizerService,
	scheduler background_jobs.SchedulerService,
) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:      config,
		store:       store,
		tokenMaker:  tokenMaker,
		storage:     storageService,
		cache:       cacheStore,
		embedding:   embedding,
		categorizer: categorizer,
		scheduler:   scheduler,
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) Start(address string) error {
	go server.scheduler.Start()
	defer server.scheduler.Stop()

	return server.router.Run(address)
}

func (server *Server) checkService(ctx *gin.Context) {
	message := HealthResponse{
		Status: "health",
	}
	ctx.JSON(http.StatusOK, message)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
