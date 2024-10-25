package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	mockdb "github.com/m1thrandir225/galore-services/db/mock"
	db "github.com/m1thrandir225/galore-services/db/sqlc"
	"github.com/m1thrandir225/galore-services/token"
	"github.com/m1thrandir225/galore-services/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateNotification(t *testing.T) {
	notification := randomNotification(t)

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
				"user_id":              notification.UserID,
				"notification_type_id": notification.NotificationTypeID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, notification.UserID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateNotificationParams{
					UserID:             notification.UserID,
					NotificationTypeID: notification.NotificationTypeID,
				}
				store.EXPECT().CreateNotification(gomock.Any(), arg).Times(1).Return(notification, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchNotification(t, recorder.Body, notification)
			},
		},
		{
			name: "Unauthorized",
			body: gin.H{
				"user_id":              notification.UserID,
				"notification_type_id": notification.NotificationTypeID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateNotificationParams{
					UserID:             notification.UserID,
					NotificationTypeID: notification.NotificationTypeID,
				}
				store.EXPECT().CreateNotification(gomock.Any(), arg).Times(0).Return(notification, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Bad Request",
			body: gin.H{},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, notification.UserID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateNotificationParams{
					UserID:             notification.UserID,
					NotificationTypeID: notification.NotificationTypeID,
				}
				store.EXPECT().CreateNotification(gomock.Any(), arg).Times(0).Return(notification, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Internal Server Error",
			body: gin.H{
				"user_id":              notification.UserID,
				"notification_type_id": notification.NotificationTypeID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, notification.UserID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateNotificationParams{
					UserID:             notification.UserID,
					NotificationTypeID: notification.NotificationTypeID,
				}
				store.EXPECT().CreateNotification(gomock.Any(), arg).Times(1).Return(notification, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store, nil, nil)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/api/v1/notifications")
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestUpdateNotificationStatus(t *testing.T) {
	notification := randomNotification(t)

	testCases := []struct {
		name           string
		notificationId string
		body           gin.H
		setupAuth      func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs     func(store *mockdb.MockStore)
		checkResponse  func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:           "OK",
			notificationId: notification.ID.String(),
			body: gin.H{
				"opened": true,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, notification.UserID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateUserNotificationParams{
					ID:     notification.ID,
					Opened: true,
				}
				store.EXPECT().UpdateUserNotification(gomock.Any(), arg).Times(1).Return(db.Notification{
					ID:                 notification.ID,
					Opened:             true,
					NotificationTypeID: notification.NotificationTypeID,
					CreatedAt:          notification.CreatedAt,
					UserID:             notification.UserID,
				}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				requireUpdatedNotMatchInitial(t, recorder.Body, notification)
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:           "Bad Request",
			notificationId: notification.ID.String(),
			body:           gin.H{},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, notification.UserID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateUserNotificationParams{
					ID:     notification.ID,
					Opened: true,
				}
				store.EXPECT().UpdateUserNotification(gomock.Any(), arg).Times(0).Return(db.Notification{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:           "Bad Request - URI",
			notificationId: "aaaaa",
			body: gin.H{
				"opened": true,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, notification.UserID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateUserNotificationParams{
					ID:     notification.ID,
					Opened: true,
				}
				store.EXPECT().UpdateUserNotification(gomock.Any(), arg).Times(0).Return(db.Notification{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:           "Not Found",
			notificationId: notification.ID.String(),
			body: gin.H{
				"opened": true,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, notification.UserID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateUserNotificationParams{
					ID:     notification.ID,
					Opened: true,
				}
				store.EXPECT().UpdateUserNotification(gomock.Any(), arg).Times(1).Return(db.Notification{}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:           "Unauthorized",
			notificationId: notification.ID.String(),
			body: gin.H{
				"opened": true,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateUserNotificationParams{
					ID:     notification.ID,
					Opened: true,
				}
				store.EXPECT().UpdateUserNotification(gomock.Any(), arg).Times(0).Return(db.Notification{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:           "Internal Server Error",
			notificationId: notification.ID.String(),
			body: gin.H{
				"opened": true,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, notification.UserID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateUserNotificationParams{
					ID:     notification.ID,
					Opened: true,
				}
				store.EXPECT().UpdateUserNotification(gomock.Any(), arg).Times(1).Return(db.Notification{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store, nil, nil)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/api/v1/notifications/%s", tc.notificationId)
			request, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestGetUserNotifications(t *testing.T) {
	var notifications []db.Notification
	user := randomUser(t, util.RandomString(10))
	for i := 0; i < 10; i++ {
		notification := db.Notification{
			UserID:             user.ID,
			Opened:             false,
			ID:                 uuid.New(),
			CreatedAt:          time.Now(),
			NotificationTypeID: randomNotificationType(t).ID,
		}
		notifications = append(notifications, notification)
	}

	testCases := []struct {
		name          string
		userId        string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			userId: user.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserNotifications(gomock.Any(), user.ID).Times(1).Return(notifications, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "Bad Request",
			userId: "aaaa",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserNotifications(gomock.Any(), user.ID).Times(0).Return(notifications, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "Not Found",
			userId: user.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserNotifications(gomock.Any(), user.ID).Times(1).Return(notifications, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "Internal Server Error",
			userId: user.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserNotifications(gomock.Any(), user.ID).Times(1).Return(notifications, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "Unauthorized",
			userId: user.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserNotifications(gomock.Any(), user.ID).Times(0).Return(notifications, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store, nil, nil)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/api/v1/users/%s/notifications", tc.userId)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func requireBodyMatchNotification(t *testing.T, body *bytes.Buffer, expectedNotification db.Notification) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var reqNotification db.Notification
	err = json.Unmarshal(data, &reqNotification)

	require.NoError(t, err)
	require.Equal(t, expectedNotification.ID.String(), reqNotification.ID.String())
	require.Equal(t, expectedNotification.NotificationTypeID, reqNotification.NotificationTypeID)
	require.Equal(t, expectedNotification.Opened, reqNotification.Opened)
	require.WithinDuration(t, expectedNotification.CreatedAt, reqNotification.CreatedAt, time.Millisecond)
}

func requireUpdatedNotMatchInitial(t *testing.T, body *bytes.Buffer, expectedNotification db.Notification) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var reqNotification db.Notification
	err = json.Unmarshal(data, &reqNotification)

	require.NoError(t, err)
	require.Equal(t, expectedNotification.ID, reqNotification.ID)
	require.Equal(t, expectedNotification.NotificationTypeID, reqNotification.NotificationTypeID)
	require.Equal(t, expectedNotification.UserID, reqNotification.UserID)
	require.WithinDuration(t, expectedNotification.CreatedAt, reqNotification.CreatedAt, time.Millisecond)
	require.NotEqual(t, expectedNotification.Opened, reqNotification.Opened)
}

func randomNotification(t *testing.T) db.Notification {
	user := randomUser(t, util.RandomString(10))
	notificationType := randomNotificationType(t)

	notification := db.Notification{
		UserID:             user.ID,
		NotificationTypeID: notificationType.ID,
		ID:                 uuid.New(),
		CreatedAt:          time.Now(),
		Opened:             false,
	}
	return notification
}
