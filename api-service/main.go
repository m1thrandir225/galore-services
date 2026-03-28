package main

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m1thrandir225/galore-services/api"
	cache2 "github.com/m1thrandir225/galore-services/internal/cache"
	categorizer2 "github.com/m1thrandir225/galore-services/internal/categorizer"
	"github.com/m1thrandir225/galore-services/internal/config"
	"github.com/m1thrandir225/galore-services/internal/db/sqlc"
	embedding2 "github.com/m1thrandir225/galore-services/internal/embedding"
	image2 "github.com/m1thrandir225/galore-services/internal/image"
	mail2 "github.com/m1thrandir225/galore-services/internal/mail"
	notifications2 "github.com/m1thrandir225/galore-services/internal/notifications"
	recipe2 "github.com/m1thrandir225/galore-services/internal/recipe"
	scheduler2 "github.com/m1thrandir225/galore-services/internal/scheduler"
	storage2 "github.com/m1thrandir225/galore-services/internal/storage"
	"github.com/m1thrandir225/galore-services/scheduler"
	pgxvector "github.com/pgvector/pgvector-go/pgx"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
)

type ginServerConfig struct {
	Config              config.Config
	Store               db.Store
	Storage             storage2.Service
	CacheStore          cache2.Store
	Embedding           embedding2.Service
	Categorizer         categorizer2.Service
	Scheduler           scheduler.SchedulerService
	MailService         mail2.Service
	NotificationService notifications2.Service
	CocktailGenerator   recipe2.Generator
	ImageGenerator      image2.Generator
}

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot load config")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot connect to database.")
	}
	// Connection pool to databse
	connPool.Config().AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxUUID.Register(conn.TypeMap())

		err = pgxvector.RegisterTypes(context.Background(), conn)
		if err != nil {
			log.Fatal().Err(err).Msg("There was an error registering vector types.")
		}

		return nil
	}

	/**
	FILE STORE
	*/
	localStorage := storage2.NewLocalStorage("./public")
	/**
	DB STORE
	*/
	store := db.NewStore(connPool)
	/**
	CACHE STORE
	*/
	cacheStore := cache2.NewRedisStore(
		config.CacheSource,
		config.CachePassword,
	)
	/**
	CATEGORIZER SERVICE
	*/
	categorizer := categorizer2.NewGaloreCategorizer(
		config.CategorizerServiceAddress,
		config.CategorizerServiceKey,
	)
	/**
	EMBEDDING SERVICE
	*/
	embeddingService := embedding2.NewGaloreEmbeddingService(
		config.EmbeddingServiceAddress,
		config.EmbeddingServiceKey,
	)
	/**
	BACKGROUND TASK PROCESSOR & SCHEDULER
	*/
	scheduler := scheduler2.NewWorkScheduler(
		"galore-work-pool",
		config.WorkerSource,
	)
	/**
	COCKTAIL RECIPE GENERATION SERVICE
	*/
	cocktailGenerator := recipe2.NewOpenAIPromptGenerator(
		config.OpenAIApiKey,
		config.OpenAIAssistantID,
		config.OpenAIThreadURL,
	)
	/**
	IMAGE GENERATION SERVICE
	*/
	imageGenerator := image2.NewStableDiffusionGenerator(
		config.StableDiffusionURL,
		config.StableDiffusionApiKey,
		"16:9",
		"png",
	)
	mailService := mail2.NewGenericMail(
		config.SMTPHost,
		config.SMTPPort,
		config.SMTPUser,
		config.SMTPPass,
	)
	/**
	NOTIFICATION SERVICE
	*/
	fcmNotifications, err := notifications2.NewFirebaseNotificator(config.FirebaseServiceKey)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot initialize firebase notifications.")
	}

	serverConfig := ginServerConfig{
		Config:              config,
		Store:               store,
		CacheStore:          cacheStore,
		Storage:             localStorage,
		Embedding:           embeddingService,
		Categorizer:         categorizer,
		Scheduler:           scheduler,
		MailService:         mailService,
		NotificationService: fcmNotifications,
		CocktailGenerator:   cocktailGenerator,
		ImageGenerator:      imageGenerator,
	}
	runGinServer(serverConfig)
}

func runGinServer(serverConfig ginServerConfig) {
	server, err := api.NewServer(
		serverConfig.Config,
		serverConfig.Store,
		serverConfig.Storage,
		serverConfig.CacheStore,
		serverConfig.Embedding,
		serverConfig.Categorizer,
		serverConfig.Scheduler,
		serverConfig.MailService,
		serverConfig.NotificationService,
		serverConfig.CocktailGenerator,
		serverConfig.ImageGenerator,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create server.")
	}

	err = server.Start(serverConfig.Config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot start server")
	}
}
