package api

import (
	"os"
	"testing"
	"time"

	"github.com/m1thrandir225/galore-services/internal/cache"
	"github.com/m1thrandir225/galore-services/internal/categorizer"
	"github.com/m1thrandir225/galore-services/internal/config"
	"github.com/m1thrandir225/galore-services/internal/db/sqlc"
	"github.com/m1thrandir225/galore-services/internal/embedding"
	"github.com/m1thrandir225/galore-services/internal/image/mock"
	"github.com/m1thrandir225/galore-services/internal/mail/mock"
	"github.com/m1thrandir225/galore-services/internal/notifications/mock"
	"github.com/m1thrandir225/galore-services/internal/recipe/mock"
	"github.com/m1thrandir225/galore-services/internal/scheduler/mock"
	"github.com/m1thrandir225/galore-services/internal/storage"
	"github.com/m1thrandir225/galore-services/internal/util"
	"go.uber.org/mock/gomock"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newTestServer(
	t *testing.T,
	ctrl *gomock.Controller,
	store db.Store,
	cacheStore cache.Store,
	fileStorage storage.Service,
	categorizer categorizer.Service,
	embeddingService embedding.Service,
) *Server {
	config := config.Config{
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
