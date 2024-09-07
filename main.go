package main

import (
	"context"
	"github.com/m1thrandir225/galore-services/storage"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m1thrandir225/galore-services/api"
	db "github.com/m1thrandir225/galore-services/db/sqlc"
	"github.com/m1thrandir225/galore-services/util"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
)

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

	connPool.Config().AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxUUID.Register(conn.TypeMap())
		return nil
	}

	localStorage := storage.NewLocalStorage("/public")

	store := db.NewStore(connPool)

	runGinServer(config, store, localStorage)

}

func runGinServer(config util.Config, store db.Store, storage storage.FileService) {
	server, err := api.NewServer(config, store)

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create server.")
	}

	err = server.Start(config.HTTPServerAddress)

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot start server")
	}
}
