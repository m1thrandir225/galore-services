package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	mockstorage "github.com/m1thrandir225/galore-services/storage/mock"
	"io"
	"net/http"
	"net/http/httptest"
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
	user := randomUser(t)

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
			testCase.buildStubs(store)

			server := newTestServer(t, store, nil, storage)
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
	user := randomUser(t)

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
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.name, func(t *testing.T) {})
	}
}

func TestUpdateUserDetailsApi(t *testing.T) {
	user := randomUser(t)

	var newDate pgtype.Date
	newDateString := util.RandomDate()
	avatarUrl := util.RandomString(10)
	newDate.Scan(newDateString)

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
	user := randomUser(t)

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
				arg := db.UpdateUserPasswordParams{
					ID:       user.ID,
					Password: util.RandomString(10),
				}
				store.EXPECT().UpdateUserPassword(gomock.Any(), arg).Times(1).Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "Unauthorized",
			userId: user.ID.String(),
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
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateUserPasswordParams{
					ID:       user.ID,
					Password: util.RandomString(10),
				}
				store.EXPECT().UpdateUserPassword(gomock.Any(), arg).Times(1).Return(sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "Bad Request",
			userId: user.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateUserPasswordParams{
					ID:       user.ID,
					Password: util.RandomString(10),
				}
				store.EXPECT().UpdateUserPassword(gomock.Any(), arg).Times(1).Return(sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Internal Server Error",
		},
	}

	for i := range testCases {
		testCase := testCases[i]

		t.Run(testCase.name, func(t *testing.T) {
		})
	}
}

func TestUpdateUserPushNotificationsApi(t *testing.T) {}

func TestUpdateUserEmailNotificationsApi(t *testing.T) {}

func randomUser(t *testing.T) db.User {
	var date pgtype.Date

	err := date.Scan(util.RandomDate())
	require.NoError(t, err)

	id := uuid.New()
	return db.User{
		ID:                        id,
		Name:                      util.RandomString(10),
		Email:                     util.RandomEmail(),
		Password:                  util.RandomString(10),
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
