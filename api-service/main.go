package main

import (
	"context"
	"os"

	"github.com/m1thrandir225/galore-services/cache"
	categorizer "github.com/m1thrandir225/galore-services/categorizer_service"
	embedding "github.com/m1thrandir225/galore-services/embedding_service"
	"github.com/m1thrandir225/galore-services/storage"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m1thrandir225/galore-services/api"
	db "github.com/m1thrandir225/galore-services/db/sqlc"
	"github.com/m1thrandir225/galore-services/util"
	pgxvector "github.com/pgvector/pgvector-go/pgx"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
)

type ginServerConfig struct {
	Config      util.Config
	Store       db.Store
	Storage     storage.FileService
	CacheStore  cache.KeyValueStore
	Embedding   embedding.EmbeddingService
	Categorizer categorizer.CategorizerService
}

func main() {
	config, err := util.LoadConfig(".")
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

	localStorage := storage.NewLocalStorage("./public")

	store := db.NewStore(connPool)
	cacheStore := cache.NewRedisStore(config.CacheSource, config.CachePassword)
	categorizer := categorizer.NewGaloreCategorizer(config.CategorizerServiceAddress, config.CategorizerServiceKey)
	embeddingService := embedding.NewGaloreEmbeddingService(config.EmbeddingServiceAddress, config.EmbeddingServiceKey)

	serverConfig := ginServerConfig{
		Config:      config,
		Store:       store,
		CacheStore:  cacheStore,
		Storage:     localStorage,
		Embedding:   embeddingService,
		Categorizer: categorizer,
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
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create server.")
	}

	err = server.Start(serverConfig.Config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot start server")
	}
}
