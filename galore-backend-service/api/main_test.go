package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/galore-services/cache"
	db "github.com/m1thrandir225/galore-services/db/sqlc"
	"github.com/m1thrandir225/galore-services/storage"
	"github.com/m1thrandir225/galore-services/util"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store, cacheStore cache.KeyValueStore) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	localStorage := storage.NewLocalStorage("./public")

	server, err := NewServer(config, store, localStorage, cacheStore)
	require.NoError(t, err)
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
