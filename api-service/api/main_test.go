package api

import (
	categorizer "github.com/m1thrandir225/galore-services/categorizer_service"
	mockcocktailgen "github.com/m1thrandir225/galore-services/cocktail_gen/mock"
	mockimagegen "github.com/m1thrandir225/galore-services/image_gen/mock"
	mockmail "github.com/m1thrandir225/galore-services/mail/mock"
	mocknotifications "github.com/m1thrandir225/galore-services/notifications/mock"
	mockscheduler "github.com/m1thrandir225/galore-services/scheduler/mock"
	"go.uber.org/mock/gomock"
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

func newTestServer(
	t *testing.T,
	ctrl *gomock.Controller,
	store db.Store,
	cacheStore cache.KeyValueStore,
	fileStorage storage.FileService,
	categorizer categorizer.CategorizerService,
	embeddingService embedding.EmbeddingService,
) *Server {
	config := util.Config{
		Environment:             "testing",
		EmbeddingServiceAddress: "http://localhost:8000",
		EmbeddingServiceKey:     "testing",
		TokenSymmetricKey:       util.RandomString(32),
		AccessTokenDuration:     time.Minute,
		RefreshTokenDuration:    time.Minute,
		TOTPSecret:              util.RandomString(32),
	}
	mockScheduler := mockscheduler.NewMockSchedulerService(ctrl)
	mockMail := mockmail.NewMockMailService(ctrl)
	mockNotifications := mocknotifications.NewMockNotificationService(ctrl)
	mockCocktailGenerator := mockcocktailgen.NewMockCocktailGenerator(ctrl)
	mockImageGenerator := mockimagegen.NewMockImageGenerator(ctrl)

	server, err := NewServer(
		config,
		store,
		fileStorage,
		cacheStore,
		embeddingService,
		categorizer,
		mockScheduler,
		mockMail,
		mockNotifications,
		mockCocktailGenerator,
		mockImageGenerator,
	)
	require.NoError(t, err)
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
