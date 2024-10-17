package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
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
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateNotificationTypeApi(t *testing.T) {
	notificationType := randomNotificationType(t)

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
					Tag:     notificationType.Tag,
					Content: notificationType.Content,
				}
				store.EXPECT().CreateNotificationType(gomock.Any(), arg).Times(1).Return(notificationType, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchNotificationType(t, notificationType, recorder.Body)
			},
		},
		{
			name: "Unauthorized",
			body: gin.H{},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateNotificationType(gomock.Any(), gomock.Any()).Times(0).Return(db.NotificationType{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)

			},
		},
		{
			name: "Bad Request",
			body: gin.H{},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateNotificationTypeParams{
					Title:   notificationType.Title,
					Content: notificationType.Content,
					Tag:     notificationType.Tag,
				}
				store.EXPECT().CreateNotificationType(gomock.Any(), arg).Times(0).Return(db.NotificationType{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Internal Server Error",
			body: gin.H{
				"title":   notificationType.Title,
				"content": notificationType.Content,
				"tag":     notificationType.Tag,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateNotificationType(gomock.Any(), gomock.Any()).Times(1).Return(db.NotificationType{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		testCase := testCases[i]

		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			testCase.buildStubs(store)

			server := newTestServer(t, store, nil)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(testCase.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/api/v1/notification_types")
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			testCase.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(recorder)
		})
	}
}
func TestGetNotificationTypeApi(t *testing.T) {
	userId := uuid.New()
	notificationType := randomNotificationType(t)

	testCases := []struct {
		name               string
		notificationTypeId string
		setupAuth          func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs         func(store *mockdb.MockStore)
		checkResponse      func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:               "OK",
			notificationTypeId: notificationType.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetNotificationType(gomock.Any(), notificationType.ID).Times(1).Return(notificationType, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchNotificationType(t, notificationType, recorder.Body)
			},
		},
		{
			name:               "Not found",
			notificationTypeId: uuid.New().String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetNotificationType(gomock.Any(), gomock.Any()).Times(1).Return(db.NotificationType{}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:               "Bad Request",
			notificationTypeId: util.RandomString(2),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetNotificationType(gomock.Any(), gomock.Any()).Times(0).Return(db.NotificationType{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:               "Unauthorized",
			notificationTypeId: notificationType.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetNotificationType(gomock.Any(), gomock.Any()).Times(0).Return(notificationType, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:               "Internal Server Error",
			notificationTypeId: notificationType.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetNotificationType(gomock.Any(), gomock.Any()).Times(1).Return(notificationType, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		testCase := testCases[i]

		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			testCase.buildStubs(store)

			server := newTestServer(t, store, nil)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/api/v1/notification_types/%s", testCase.notificationTypeId)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			testCase.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(recorder)
		})
	}
}
func TestGetNotificationTypesApi(t *testing.T) {
	notificationTypes := make([]db.NotificationType, 5)
	for _ = range notificationTypes {
		notificationType := randomNotificationType(t)
		notificationTypes = append(notificationTypes, notificationType)
	}

	userId := uuid.New()

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAllTypes(gomock.Any()).Times(1).Return(notificationTypes, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Unauthorized",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAllTypes(gomock.Any()).Times(0).Return(notificationTypes, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Not Found",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAllTypes(gomock.Any()).Times(1).Return(notificationTypes, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "Internal Server Error",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAllTypes(gomock.Any()).Times(1).Return(notificationTypes, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
	for i := range testCases {
		testCase := testCases[i]

		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			testCase.buildStubs(store)

			server := newTestServer(t, store, nil)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/api/v1/notification_types")
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			testCase.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(recorder)
		})
	}
}
func TestDeleteNotificationTypeApi(t *testing.T) {
	notificationType := randomNotificationType(t)
	userId := uuid.New()
	testCases := []struct {
		name               string
		notificationTypeId string
		setupAuth          func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs         func(store *mockdb.MockStore)
		checkResponse      func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:               "OK",
			notificationTypeId: notificationType.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteNotificationType(gomock.Any(), notificationType.ID).Times(1).Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:               "Not Found",
			notificationTypeId: notificationType.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteNotificationType(gomock.Any(), gomock.Any()).Times(1).Return(sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:               "Bad Request",
			notificationTypeId: util.RandomString(10),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteNotificationType(gomock.Any(), gomock.Any()).Times(0).Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:               "Unauthorized",
			notificationTypeId: notificationType.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteNotificationType(gomock.Any(), gomock.Any()).Times(0).Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:               "Internal Server Error",
			notificationTypeId: notificationType.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteNotificationType(gomock.Any(), notificationType.ID).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		testCase := testCases[i]

		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			testCase.buildStubs(store)

			server := newTestServer(t, store, nil)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/api/v1/notification_types/%s", testCase.notificationTypeId)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			testCase.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(recorder)
		})
	}
}
func TestUpdateNotificationTypeApi(t *testing.T) {
	notificationType := randomNotificationType(t)

	updatedNotification := db.NotificationType{
		ID:        notificationType.ID,
		Title:     util.RandomString(10),
		Tag:       util.RandomString(5),
		Content:   util.RandomString(12),
		CreatedAt: notificationType.CreatedAt,
	}
	userId := uuid.New()
	testCases := []struct {
		name               string
		notificationTypeId string
		body               gin.H
		setupAuth          func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs         func(store *mockdb.MockStore)
		checkResponse      func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:               "OK",
			notificationTypeId: notificationType.ID.String(),
			body: gin.H{
				"title":   updatedNotification.Title,
				"content": updatedNotification.Content,
				"tag":     updatedNotification.Tag,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateNotificationTypeParams{
					ID:      notificationType.ID,
					Title:   updatedNotification.Title,
					Content: updatedNotification.Content,
					Tag:     updatedNotification.Tag,
				}
				store.EXPECT().UpdateNotificationType(gomock.Any(), arg).Times(1).Return(updatedNotification, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:               "Bad Request",
			notificationTypeId: notificationType.ID.String(),
			body:               gin.H{},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateNotificationType(gomock.Any(), gomock.Any()).Times(0).Return(updatedNotification, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:               "Not Found",
			notificationTypeId: notificationType.ID.String(),
			body: gin.H{
				"title":   updatedNotification.Title,
				"content": updatedNotification.Content,
				"tag":     updatedNotification.Tag,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateNotificationTypeParams{
					ID:      notificationType.ID,
					Title:   updatedNotification.Title,
					Content: updatedNotification.Content,
					Tag:     updatedNotification.Tag,
				}
				store.EXPECT().UpdateNotificationType(gomock.Any(), arg).Times(1).Return(updatedNotification, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:               "Bad Request - Wrong ID",
			notificationTypeId: util.RandomString(5),
			body: gin.H{
				"title":   updatedNotification.Title,
				"content": updatedNotification.Content,
				"tag":     updatedNotification.Tag,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateNotificationType(gomock.Any(), gomock.Any()).Times(0).Return(db.NotificationType{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:               "Internal Server Error",
			notificationTypeId: notificationType.ID.String(),
			body: gin.H{
				"title":   updatedNotification.Title,
				"content": updatedNotification.Content,
				"tag":     updatedNotification.Tag,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateNotificationTypeParams{
					ID:      notificationType.ID,
					Title:   updatedNotification.Title,
					Content: updatedNotification.Content,
					Tag:     updatedNotification.Tag,
				}
				store.EXPECT().UpdateNotificationType(gomock.Any(), arg).Times(1).Return(updatedNotification, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
	for i := range testCases {
		testCase := testCases[i]

		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			testCase.buildStubs(store)

			server := newTestServer(t, store, nil)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(testCase.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/api/v1/notification_types/%s", testCase.notificationTypeId)
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			testCase.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(recorder)
		})
	}
}

func requireBodyMatchNotificationType(t *testing.T, notificationType db.NotificationType, body *bytes.Buffer) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var reqNotificationType db.NotificationType
	err = json.Unmarshal(data, &reqNotificationType)

	require.NoError(t, err)
	require.Equal(t, notificationType.ID, reqNotificationType.ID)
	require.Equal(t, notificationType.Title, reqNotificationType.Title)
	require.Equal(t, notificationType.Content, reqNotificationType.Content)
	require.Equal(t, notificationType.Tag, reqNotificationType.Tag)
	require.WithinDuration(t, notificationType.CreatedAt, reqNotificationType.CreatedAt, time.Second)
}

func randomNotificationType(t *testing.T) db.NotificationType {
	return db.NotificationType{
		ID:        uuid.New(),
		Title:     util.RandomString(10),
		Tag:       util.RandomString(5),
		Content:   util.RandomString(100),
		CreatedAt: time.Now(),
	}
}
