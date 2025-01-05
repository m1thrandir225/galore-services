package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	mockdb "github.com/m1thrandir225/galore-services/db/mock"
	db "github.com/m1thrandir225/galore-services/db/sqlc"
	mockstorage "github.com/m1thrandir225/galore-services/storage/mock"
	"github.com/m1thrandir225/galore-services/token"
	"github.com/m1thrandir225/galore-services/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRegisterUserApi(t *testing.T) {
	//TODO: implement
	//user := randomUser(t)

}

func TestLoginUserApi(t *testing.T) {
	password := util.RandomString(25)
	user := randomUser(t, password)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"email":    user.Email,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserByEmail(gomock.Any(), user.Email).Times(1).Return(user, nil)
				store.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Times(1).Return(db.Session{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Bad Request",
			body: gin.H{},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserByEmail(gomock.Any(), user.Email).Times(0).Return(user, nil)
				store.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Bad Password",
			body: gin.H{
				"email":    user.Email,
				"password": "random",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserByEmail(gomock.Any(), user.Email).Times(1).Return(user, nil)
				store.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Times(0).Return(db.Session{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Not Found",
			body: gin.H{
				"email":    "random@gmail.com",
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserByEmail(gomock.Any(), "random@gmail.com").Times(1).Return(user, sql.ErrNoRows)
				store.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Times(0).Return(db.Session{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "Internal Server Error - Getting User",
			body: gin.H{
				"email":    user.Email,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserByEmail(gomock.Any(), user.Email).Times(1).Return(user, nil)
				store.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Times(1).Return(db.Session{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Internal Server Error - Creating Session",
			body: gin.H{
				"email":    user.Email,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserByEmail(gomock.Any(), user.Email).Times(1).Return(user, nil)
				store.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Times(1).Return(db.Session{}, sql.ErrConnDone)
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
			storage := mockstorage.NewMockFileService(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, ctrl, store, nil, storage, nil, nil)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/api/v1/login")
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestLogoutApi(t *testing.T) {
	userId := uuid.New()
	sessionId := uuid.New()

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
				"session_id": sessionId,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().InvalidateSession(gomock.Any(), sessionId).Times(1).Return(db.Session{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Bad Request",
			body: gin.H{},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().InvalidateSession(gomock.Any(), sessionId).Times(0).Return(db.Session{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Not Found",
			body: gin.H{
				"session_id": sessionId,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().InvalidateSession(gomock.Any(), sessionId).Times(1).Return(db.Session{}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "Unauthorized",
			body: gin.H{
				"session_id": sessionId,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().InvalidateSession(gomock.Any(), sessionId).Times(0).Return(db.Session{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Internal Server Error - Invalidating Session",
			body: gin.H{
				"session_id": sessionId,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().InvalidateSession(gomock.Any(), sessionId).Times(1).Return(db.Session{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				log.Print(recorder.Body.String())
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

			server := newTestServer(t, ctrl, store, nil, nil, nil, nil)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/api/v1/logout")
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestRefreshTokenApi(t *testing.T) {
	userId := uuid.New()
	sessionId := uuid.New()

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "",
			body: gin.H{
				"session_id": sessionId,
				"user_id":    userId,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Invalidate Session",
			body: gin.H{
				"session_id": "",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		{
			name: "Invalidate Session",
			body: gin.H{
				"session_id": sessionId,
				"user_id":    userId,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().InvalidateSession(gomock.Any(), sessionId).Times(1).Return(db.Session{}, nil)

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
			server := newTestServer(t, ctrl, store, nil, nil, nil, nil)
			recorder := httptest.NewRecorder()

			var refreshToken string
			if tc.name == "Invalidate Session" {
				refreshToken = util.RandomString(10)
			} else {
				rT, _, err := server.tokenMaker.CreateToken(userId, time.Hour)
				require.NoError(t, err)
				refreshToken = rT

			}
			data, err := json.Marshal(gin.H{
				"session_id":    tc.body["session_id"],
				"refresh_token": refreshToken,
			})
			require.NoError(t, err)

			url := fmt.Sprintf("/api/v1/refresh")
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

//func TestRegisterApi(t *testing.T) {
//	name := util.RandomString(10)
//	email := util.RandomEmail()
//	avatarName := util.RandomString(10)
//	password := util.RandomString(10)
//	birthday := util.RandomDate()
//
//	testCases := []struct {
//		name          string
//		body          gin.H
//		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
//		buildStubs    func(store *mockdb.MockStore)
//		checkResponse func(recorder *httptest.ResponseRecorder)
//	}{
//		{
//			name: "OK",
//			body: gin.H{
//				"name":        name,
//				"email":       email,
//				"avatar_name": avatarName,
//				"birthday":    birthday,
//				"password":    password,
//			},
//		},
//	}
//
//	for i := range testCases {
//		tc := testCases[i]
//
//		t.Run(tc.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//			defer ctrl.Finish()
//
//			store := mockdb.NewMockStore(ctrl)
//			storage := mockstorage.NewMockFileService(ctrl)
//
//			tc.buildStubs(store)
//
//			server := newTestServer(t, store, nil, nil, nil, nil)
//
//			recorder := httptest.NewRecorder()
//
//			url := fmt.Sprintf("/api/v1/register")
//
//			request, err := http.NewRequest(http.MethodPost, url, nil)
//		})
//	}
//}
