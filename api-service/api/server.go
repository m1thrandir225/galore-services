package api

import (
	"fmt"
	"github.com/m1thrandir225/galore-services/cocktail_gen"
	"github.com/m1thrandir225/galore-services/image_gen"
	"github.com/m1thrandir225/galore-services/notifications"
	"log"

	"github.com/m1thrandir225/galore-services/cache"
	categorizer "github.com/m1thrandir225/galore-services/categorizer_service"
	embedding "github.com/m1thrandir225/galore-services/embedding_service"
	"github.com/m1thrandir225/galore-services/mail"
	"github.com/m1thrandir225/galore-services/scheduler"
	"github.com/m1thrandir225/galore-services/storage"

	"github.com/gin-gonic/gin"
	db "github.com/m1thrandir225/galore-services/db/sqlc"
	"github.com/m1thrandir225/galore-services/token"
	"github.com/m1thrandir225/galore-services/util"
)

type Server struct {
	config              util.Config
	store               db.Store
	router              *gin.Engine
	tokenMaker          token.Maker
	storage             storage.FileService
	cache               cache.KeyValueStore
	embedding           embedding.EmbeddingService
	categorizer         categorizer.CategorizerService
	scheduler           scheduler.SchedulerService
	mailService         mail.MailService
	notificationService notifications.NotificationService
	cocktailGenerator   cocktail_gen.CocktailGenerator
	imageGenerator      image_gen.ImageGenerator
}

func NewServer(
	config util.Config,
	store db.Store,
	storageService storage.FileService,
	cacheStore cache.KeyValueStore,
	embedding embedding.EmbeddingService,
	categorizer categorizer.CategorizerService,
	scheduler scheduler.SchedulerService,
	mailService mail.MailService,
	notificationService notifications.NotificationService,
	cocktailGenerator cocktail_gen.CocktailGenerator,
	imageGenerator image_gen.ImageGenerator,

) (*Server, error) {
	log.Println(config)
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
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
