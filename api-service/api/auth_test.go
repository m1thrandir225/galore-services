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
	"github.com/m1thrandir225/galore-services/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
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

			server := newTestServer(t, store, nil, storage)
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
	sessionId := uuid.New()

	testCases := []struct {
		name          string
		sessionId     string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			sessionId: sessionId.String(),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetSession(gomock.Any(), sessionId).Times(1).Return(db.Session{}, nil)
				store.EXPECT().DeleteSession(gomock.Any(), sessionId).Times(1).Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNoContent, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
		})
	}
}
