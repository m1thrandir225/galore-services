package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/galore-services/cache"
	db "github.com/m1thrandir225/galore-services/db/sqlc"
	embedding "github.com/m1thrandir225/galore-services/embedding_service"
	"github.com/m1thrandir225/galore-services/storage"
	"github.com/m1thrandir225/galore-services/util"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store, cacheStore cache.KeyValueStore) *Server {
	config := util.Config{
		TokenSymmetricKey:       util.RandomString(32),
		AccessTokenDuration:     time.Minute,
		EmbeddingServiceAddress: "http://localhost:8000",
		EmbeddingServiceKey:     "testing",
	}

	localStorage := storage.NewLocalStorage("./public")
	embeddingService := embedding.NewGaloreEmbeddingService(config.EmbeddingServiceAddress, config.EmbeddingServiceKey)

	server, err := NewServer(config, store, localStorage, cacheStore, embeddingService)
	require.NoError(t, err)
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
