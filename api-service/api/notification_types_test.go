package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	mockdb "github.com/m1thrandir225/galore-services/db/mock"
	db "github.com/m1thrandir225/galore-services/db/sqlc"
	"github.com/m1thrandir225/galore-services/token"
	"github.com/m1thrandir225/galore-services/util"
	"go.uber.org/mock/gomock"
)

//TODO: write the tests for the notification types api

func TestCreateNotificationTypeApi(t *testing.T) {
	notificationType := randomNotificationType()

	userId := uuid.New()
	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"title":   notificationType.Title,
				"content": notificationType.Content,
				"tag":     notificationType.Tag,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateNotificationTypeParams{
					Title:   notificationType.Title,
					Content: notificationType.Content,
					Tag:     notificationType.Tag,
				}
				store.EXPECT().CreateNotificationType(gomock.Any(), arg).Return(notificationType, nil)
			},
		},
	}

	for i := range testCases {
		testCase := testCases[i]

		t.Run(testCase.name, func(t *testing.T) {
			//TODO: implement run func()
		})
	}
}
func TestGetNotificationTypesApi(t *testing.T)   {}
func TestGetNotificationTypeApi(t *testing.T)    {}
func TestDeleteNotificationTypeApi(t *testing.T) {}
func TestUpdateNotificationTypeApi(t *testing.T) {}

func randomNotificationType() db.NotificationType {
	return db.NotificationType{
		ID:        uuid.New(),
		Title:     util.RandomString(10),
		Tag:       util.RandomString(5),
		Content:   util.RandomString(100),
		CreatedAt: time.Now(),
	}
}
