package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m1thrandir225/galore-services/internal/config"
)

var testStore Store

func TestMain(m *testing.M) {
	config, err := config.LoadConfig("../..")
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}
	connPool, err := pgxpool.New(context.Background(), config.TestingDBSource)
	if err != nil {
		log.Fatal("Cannot connect to db: ", err)
	}

	testStore = NewStore(connPool)

	os.Exit(m.Run())
}
