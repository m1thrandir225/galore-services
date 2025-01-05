package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	mockcategorize "github.com/m1thrandir225/galore-services/categorizer_service/mock"
	mockembedding "github.com/m1thrandir225/galore-services/embedding_service/mock"
	mockstorage "github.com/m1thrandir225/galore-services/storage/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	mockdb "github.com/m1thrandir225/galore-services/db/mock"
	db "github.com/m1thrandir225/galore-services/db/sqlc"
	"github.com/m1thrandir225/galore-services/token"
	"github.com/m1thrandir225/galore-services/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetUserDetailsApi(t *testing.T) {
	password := util.RandomString(25)
	user := randomUser(t, password)

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
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "Not Found - Mismatch Context",
			userId: "3a5b94c0-517f-4140-87e3-b2755e9372c5",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:   "Not Found",
			userId: user.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "Internal Error",
			userId: user.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
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
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:   "Bad Request",
			userId: "aa",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		testCase := testCases[i]

		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			storage := mockstorage.NewMockFileService(ctrl)
			categorizer := mockcategorize.NewMockCategorizerService(ctrl)
			embeddingService := mockembedding.NewMockEmbeddingService(ctrl)

			testCase.buildStubs(store)

			server := newTestServer(t, ctrl, store, nil, storage, categorizer, embeddingService)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/api/v1/users/%s", testCase.userId)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			testCase.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(recorder)
		})
	}
}

func TestDeleteUserApi(t *testing.T) {
	password := util.RandomString(25)
	user := randomUser(t, password)

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
				store.EXPECT().DeleteUser(gomock.Any(), user.ID).Times(1).Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "Internal Error",
			userId: user.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteUser(gomock.Any(), user.ID).Times(1).Return(sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "Bad Request - uuid parse",
			userId: "aaaaaa",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteUser(gomock.Any(), gomock.Any()).Times(0)
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
				store.EXPECT().
					DeleteUser(gomock.Any(), user.ID).Times(1).Return(sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "Unauthorized",
			userId: user.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			storage := mockstorage.NewMockFileService(ctrl)
			categorizer := mockcategorize.NewMockCategorizerService(ctrl)
			embeddingService := mockembedding.NewMockEmbeddingService(ctrl)

			testCase.buildStubs(store)

			server := newTestServer(t, ctrl, store, nil, storage, categorizer, embeddingService)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/api/v1/users/%s", testCase.userId)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			testCase.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(recorder)
		})
	}
}

func TestUpdateUserDetailsApi(t *testing.T) {
	password := util.RandomString(25)
	user := randomUser(t, password)

	var newDate pgtype.Date
	newDateString := util.RandomDate()
	avatarUrl := util.RandomString(10)
	err := newDate.Scan(newDateString)
	if err != nil {
		return
	}

	arg := db.UpdateUserInformationParams{
		ID:        user.ID,
		Name:      "James Brown",
		AvatarUrl: avatarUrl,
		Email:     "jamesbrown@gmail.com",
		Birthday:  newDate,
	}

	testCases := []struct {
		name          string
		userId        string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore, storage *mockstorage.MockFileService)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK-NoFile",
			userId: user.ID.String(),
			body: gin.H{
				"id":       user.ID,
				"name":     arg.Name,
				"email":    arg.Email,
				"birthday": newDateString,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore, storage *mockstorage.MockFileService) {
				gomock.InOrder(
					store.EXPECT().GetUser(gomock.Any(), arg.ID).Times(1).Return(user, nil),
					//storage.EXPECT().ReplaceFile(gomock.Any(), gomock.Any()).Times(1).Return("", nil),
					store.EXPECT().UpdateUserInformation(gomock.Any(), arg).Times(1).Return(db.UpdateUserInformationRow{}, nil),
				)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "OK-NoFile",
			userId: user.ID.String(),
			body: gin.H{
				"id":       user.ID,
				"name":     arg.Name,
				"email":    arg.Email,
				"birthday": newDateString,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore, storage *mockstorage.MockFileService) {
				gomock.InOrder(
					store.EXPECT().GetUser(gomock.Any(), arg.ID).Times(1).Return(user, nil),
					//storage.EXPECT().ReplaceFile(gomock.Any(), gomock.Any()).Times(1).Return("", nil),
					store.EXPECT().UpdateUserInformation(gomock.Any(), arg).Times(1).Return(db.UpdateUserInformationRow{}, nil),
				)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}

	for i := range testCases {
		testCase := testCases[i]

		t.Run(testCase.name, func(t *testing.T) {
			//TODO: finish implementing this
			//ctrl := gomock.NewController(t)
			//defer ctrl.Finish()
			//
			//body := new(bytes.Buffer)
			//writer := multipart.NewWriter(body)
			//
			//_ = writer.WriteField("email", testCase.body["email"].(string))
			//_ = writer.WriteField("name", testCase.body["name"].(string))
			//_ = writer.WriteField("birthday", testCase.body["birthday"].(string))
			//if testCase.name == "OK-WithFile" {
			//	fileWriter, err := writer.CreateFormFile("avatar_url", avatarUrl)
			//	require.NoError(t, err)
			//	_, err = fileWriter.Write([]byte(util.RandomString(512)))
			//	require.NoError(t, err)
			//}
			//
			//err := writer.Close()
			//require.NoError(t, err)
			//
			//url := fmt.Sprintf("/api/v1/users/%s", testCase.userId)
			//request, err := http.NewRequest(http.MethodPost, url, body)
			//require.NoError(t, err)
			//
			//request.Header.Set("Content-Type", writer.FormDataContentType())
			//
			//store := mockdb.NewMockStore(ctrl)
			//storage := mockstorage.NewMockFileService(ctrl)
			//testCase.buildStubs(store, storage)
			//
			//server := newTestServer(t, store, nil, storage)
			//recorder := httptest.NewRecorder()
			//
			//testCase.setupAuth(t, request, server.tokenMaker)
			//server.router.ServeHTTP(recorder, request)
			//testCase.checkResponse(recorder)
		})
	}
}

func TestUpdateUserPasswordApi(t *testing.T) {
	password := util.RandomString(25)
	user := randomUser(t, password)

	newPassword := util.RandomString(10)

	testCases := []struct {
		name          string
		userId        string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			userId: user.ID.String(),
			body: gin.H{
				"new_password": newPassword,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				//TODO: As the hashing is implemented into the function itself I can't replicate the hash and give it the proper arguments, Any will do
				store.EXPECT().UpdateUserPassword(gomock.Any(), gomock.Any()).Times(1).Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "Unauthorized",
			userId: user.ID.String(),
			body: gin.H{
				"new_password": newPassword,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateUserPassword(gomock.Any(), gomock.Any()).Times(0).Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:   "Not Found",
			userId: user.ID.String(),
			body: gin.H{
				"new_password": newPassword,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateUserPassword(gomock.Any(), gomock.Any()).Times(1).Return(sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "Bad Request",
			userId: "aaaaa",
			body:   gin.H{},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateUserPasswordParams{
					ID:       user.ID,
					Password: util.RandomString(10),
				}
				store.EXPECT().UpdateUserPassword(gomock.Any(), arg).Times(0).Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "Internal Server Error",
			userId: user.ID.String(),
			body: gin.H{
				"new_password": newPassword,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateUserPassword(gomock.Any(), gomock.Any()).Times(1).Return(sql.ErrConnDone)
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
			storage := mockstorage.NewMockFileService(ctrl)
			categorizer := mockcategorize.NewMockCategorizerService(ctrl)
			embeddingService := mockembedding.NewMockEmbeddingService(ctrl)
			testCase.buildStubs(store)

			server := newTestServer(t, ctrl, store, nil, storage, categorizer, embeddingService)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(testCase.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/api/v1/users/%s/password", testCase.userId)
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			testCase.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(recorder)
		})
	}
}

func TestUpdateUserPushNotificationsApi(t *testing.T) {
	password := util.RandomString(25)
	user := randomUser(t, password)

	testCases := []struct {
		name          string
		userId        string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			userId: user.ID.String(),
			body: gin.H{
				"enabled": true,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateUserPushNotificationsParams{
					ID:                       user.ID,
					EnabledPushNotifications: true,
				}
				store.EXPECT().
					UpdateUserPushNotifications(gomock.Any(), arg).
					Times(1).
					Return(true, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				require.Equal(t, strconv.FormatBool(true), recorder.Body.String())
			},
		},
		{
			name:   "Unauthorized",
			userId: user.ID.String(),
			body: gin.H{
				"enabled": true,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateUserPushNotificationsParams{
					ID:                       user.ID,
					EnabledPushNotifications: true,
				}
				store.EXPECT().
					UpdateUserPushNotifications(gomock.Any(), arg).
					Times(0).
					Return(true, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:   "Bad Request",
			userId: user.ID.String(),
			body:   gin.H{},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUserPushNotifications(gomock.Any(), gomock.Any()).
					Times(0)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "Bad Request",
			userId: util.RandomString(10),
			body: gin.H{
				"enabled": true,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUserPushNotifications(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "Not Found",
			userId: user.ID.String(),
			body: gin.H{
				"enabled": true,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateUserPushNotificationsParams{
					ID:                       user.ID,
					EnabledPushNotifications: true,
				}
				store.EXPECT().
					UpdateUserPushNotifications(gomock.Any(), arg).
					Times(1).Return(true, sql.ErrNoRows)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "Internal Server Error",
			userId: user.ID.String(),
			body: gin.H{
				"enabled": true,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateUserPushNotificationsParams{
					ID:                       user.ID,
					EnabledPushNotifications: true,
				}
				store.EXPECT().
					UpdateUserPushNotifications(gomock.Any(), arg).
					Times(1).Return(true, sql.ErrConnDone)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "Invalid ID",
			userId: util.RandomString(10),
			body: gin.H{
				"enabled": true,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateUserPushNotificationsParams{
					ID:                       user.ID,
					EnabledPushNotifications: true,
				}
				store.EXPECT().
					UpdateUserPushNotifications(gomock.Any(), arg).
					Times(0).Return(true, sql.ErrConnDone)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		testCase := testCases[i]

		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			categorizer := mockcategorize.NewMockCategorizerService(ctrl)
			embeddingService := mockembedding.NewMockEmbeddingService(ctrl)

			testCase.buildStubs(store)

			server := newTestServer(t, ctrl, store, nil, nil, categorizer, embeddingService)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(testCase.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/api/v1/users/%s/push-notifications", testCase.userId)
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			testCase.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(recorder)
		})
	}
}

func TestUpdateUserEmailNotificationsApi(t *testing.T) {
	password := util.RandomString(25)
	user := randomUser(t, password)
	testCases := []struct {
		name          string
		userId        string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			userId: user.ID.String(),
			body: gin.H{
				"enabled": true,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateUserEmailNotificationsParams{
					ID:                        user.ID,
					EnabledEmailNotifications: true,
				}
				store.EXPECT().
					UpdateUserEmailNotifications(gomock.Any(), arg).
					Times(1).
					Return(true, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "Unauthorized",
			userId: user.ID.String(),
			body: gin.H{
				"enabled": true,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateUserEmailNotificationsParams{
					ID:                        user.ID,
					EnabledEmailNotifications: true,
				}
				store.EXPECT().
					UpdateUserEmailNotifications(gomock.Any(), arg).
					Times(0).
					Return(true, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:   "Bad Request",
			userId: user.ID.String(),
			body:   gin.H{},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateUserEmailNotificationsParams{
					ID:                        user.ID,
					EnabledEmailNotifications: true,
				}
				store.EXPECT().
					UpdateUserEmailNotifications(gomock.Any(), arg).
					Times(0).
					Return(true, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "Not Found",
			userId: user.ID.String(),
			body: gin.H{
				"enabled": true,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateUserEmailNotificationsParams{
					ID:                        user.ID,
					EnabledEmailNotifications: true,
				}
				store.EXPECT().
					UpdateUserEmailNotifications(gomock.Any(), arg).
					Times(1).
					Return(true, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "Internal Server Error",
			userId: user.ID.String(),
			body: gin.H{
				"enabled": true,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateUserEmailNotificationsParams{
					ID:                        user.ID,
					EnabledEmailNotifications: true,
				}
				store.EXPECT().
					UpdateUserEmailNotifications(gomock.Any(), arg).
					Times(1).
					Return(true, sql.ErrConnDone)
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
			storage := mockstorage.NewMockFileService(ctrl)
			categorizer := mockcategorize.NewMockCategorizerService(ctrl)
			embeddingService := mockembedding.NewMockEmbeddingService(ctrl)

			testCase.buildStubs(store)

			server := newTestServer(t, ctrl, store, nil, storage, categorizer, embeddingService)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(testCase.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/api/v1/users/%s/email-notifications", testCase.userId)
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			testCase.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(recorder)
		})
	}
}

func randomUser(t *testing.T, password string) db.User {
	var date pgtype.Date

	err := date.Scan(util.RandomDate())
	require.NoError(t, err)
	hashedPassword, err := util.HashPassowrd(password)
	require.NoError(t, err)

	id := uuid.New()
	return db.User{
		ID:                        id,
		Name:                      util.RandomString(10),
		Email:                     util.RandomEmail(),
		Password:                  hashedPassword,
		AvatarUrl:                 util.RandomString(10),
		EnabledEmailNotifications: util.RandomBool(),
		EnabledPushNotifications:  util.RandomBool(),
		Birthday:                  date,
		CreatedAt:                 time.Now(),
	}
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, expectedUser db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var reqUser db.User
	err = json.Unmarshal(data, &reqUser)

	require.NoError(t, err)
	require.Equal(t, expectedUser.ID, reqUser.ID)
	require.Equal(t, expectedUser.Name, reqUser.Name)
	require.Equal(t, expectedUser.Email, reqUser.Email)
	require.WithinDuration(t, reqUser.CreatedAt, reqUser.CreatedAt, time.Second)
	require.Equal(t, reqUser.Password, reqUser.Password)
}
