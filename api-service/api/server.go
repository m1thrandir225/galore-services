package api

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/galore-services/internal/cache"
	"github.com/m1thrandir225/galore-services/internal/categorizer"
	"github.com/m1thrandir225/galore-services/internal/config"
	"github.com/m1thrandir225/galore-services/internal/db/sqlc"
	"github.com/m1thrandir225/galore-services/internal/embedding"
	"github.com/m1thrandir225/galore-services/internal/image"
	"github.com/m1thrandir225/galore-services/internal/mail"
	"github.com/m1thrandir225/galore-services/internal/notifications"
	"github.com/m1thrandir225/galore-services/internal/recipe"
	"github.com/m1thrandir225/galore-services/internal/storage"
	token2 "github.com/m1thrandir225/galore-services/internal/token"
	"github.com/m1thrandir225/galore-services/scheduler"
)

type Server struct {
	config              config.Config
	store               db.Store
	router              *gin.Engine
	tokenMaker          token2.Maker
	storage             storage.Service
	cache               cache.Store
	embedding           embedding.Service
	categorizer         categorizer.Service
	scheduler           scheduler.SchedulerService
	mailService         mail.Service
	notificationService notifications.Service
	cocktailGenerator   recipe.Generator
	imageGenerator      image.Generator
}

func NewServer(
	config config.Config,
	store db.Store,
	storageService storage.Service,
	cacheStore cache.Store,
	embedding embedding.Service,
	categorizer categorizer.Service,
	scheduler scheduler.SchedulerService,
	mailService mail.Service,
	notificationService notifications.Service,
	cocktailGenerator recipe.Generator,
	imageGenerator image.Generator,

) (*Server, error) {
	log.Println(config)
	tokenMaker, err := token2.NewPASETOMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:              config,
		store:               store,
		tokenMaker:          tokenMaker,
		storage:             storageService,
		cache:               cacheStore,
		embedding:           embedding,
		categorizer:         categorizer,
		scheduler:           scheduler,
		mailService:         mailService,
		notificationService: notificationService,
		cocktailGenerator:   cocktailGenerator,
		imageGenerator:      imageGenerator,
	}
	if server.config.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	server.setupRouter()
	server.registerBackgroundHandlers()

	return server, nil
}

func (server *Server) Start(address string) error {
	go server.scheduler.Start()
	defer server.scheduler.Stop()

	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
