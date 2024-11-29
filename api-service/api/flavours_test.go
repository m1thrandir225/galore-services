package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	mockcategorize "github.com/m1thrandir225/galore-services/categorizer_service/mock"
	mockembedding "github.com/m1thrandir225/galore-services/embedding_service/mock"
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

func TestCreateFlavourApi(t *testing.T) {
	flavour := randomFlavour(t)
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
				"name": flavour.Name,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateFlavour(gomock.Any(), flavour.Name).Times(1).Return(flavour, nil)
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
				store.EXPECT().CreateFlavour(gomock.Any(), flavour.Name).Times(0).Return(flavour, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Internal Server Error",
			body: gin.H{
				"name": flavour.Name,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateFlavour(gomock.Any(), flavour.Name).Times(1).Return(flavour, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Unauthorized",
			body: gin.H{
				"name": flavour.Name,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateFlavour(gomock.Any(), flavour.Name).Times(0).Return(flavour, nil)
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
			categorizer := mockcategorize.NewMockCategorizerService(ctrl)
			embeddingService := mockembedding.NewMockEmbeddingService(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store, nil, nil, categorizer, embeddingService)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/api/v1/flavours")
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestGetFlavourApi(t *testing.T) {
	flavour := randomFlavour(t)
	userId := uuid.New()

	testCases := []struct {
		name          string
		flavourId     string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			flavourId: flavour.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetFlavourId(gomock.Any(), flavour.ID).Times(1).Return(flavour, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				requireBodyMatchFlavour(t, recorder.Body, flavour)
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:      "Bad Request",
			flavourId: util.RandomString(10),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetFlavourId(gomock.Any(), flavour.ID).Times(0).Return(flavour, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "Not Found",
			flavourId: flavour.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetFlavourId(gomock.Any(), flavour.ID).Times(1).Return(flavour, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "Internal Server Error",
			flavourId: flavour.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetFlavourId(gomock.Any(), flavour.ID).Times(1).Return(flavour, sql.ErrConnDone)
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
			categorizer := mockcategorize.NewMockCategorizerService(ctrl)
			embeddingService := mockembedding.NewMockEmbeddingService(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store, nil, nil, categorizer, embeddingService)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/api/v1/flavours/%s", tc.flavourId)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestDeleteFlavourApi(t *testing.T) {
	flavour := randomFlavour(t)
	userId := uuid.New()

	testCases := []struct {
		name          string
		flavourId     string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			flavourId: flavour.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteFlavour(gomock.Any(), flavour.ID).Times(1).Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:      "Bad Request",
			flavourId: util.RandomString(10),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteFlavour(gomock.Any(), flavour.ID).Times(0).Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "Internal Server Error",
			flavourId: flavour.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteFlavour(gomock.Any(), flavour.ID).Times(1).Return(sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "Not Found",
			flavourId: flavour.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteFlavour(gomock.Any(), flavour.ID).Times(1).Return(sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "Unauthorized",
			flavourId: flavour.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteFlavour(gomock.Any(), flavour.ID).Times(0).Return(nil)
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
			categorizer := mockcategorize.NewMockCategorizerService(ctrl)
			embeddingService := mockembedding.NewMockEmbeddingService(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store, nil, nil, categorizer, embeddingService)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/api/v1/flavours/%s", tc.flavourId)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestUpdateFlavourApi(t *testing.T) {
	flavour := randomFlavour(t)
	flavourWithNewName := randomFlavour(t)
	userId := uuid.New()

	testCases := []struct {
		name          string
		flavourId     string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			flavourId: flavour.ID.String(),
			body: gin.H{
				"name": flavourWithNewName.Name,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateFlavourParams{
					ID:   flavour.ID,
					Name: flavourWithNewName.Name,
				}
				store.EXPECT().UpdateFlavour(gomock.Any(), arg).Times(1).Return(db.Flavour{
					Name:      flavourWithNewName.Name,
					ID:        flavour.ID,
					CreatedAt: flavour.CreatedAt,
				}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:      "Bad Request",
			flavourId: flavour.ID.String(),
			body:      gin.H{},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateFlavourParams{
					ID:   flavour.ID,
					Name: flavourWithNewName.Name,
				}
				store.EXPECT().UpdateFlavour(gomock.Any(), arg).Times(0).Return(db.Flavour{
					Name:      flavourWithNewName.Name,
					ID:        flavour.ID,
					CreatedAt: flavour.CreatedAt,
				}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "Bad Request",
			flavourId: util.RandomString(10),
			body:      gin.H{},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateFlavourParams{
					ID:   flavour.ID,
					Name: flavourWithNewName.Name,
				}
				store.EXPECT().UpdateFlavour(gomock.Any(), arg).Times(0).Return(db.Flavour{
					Name:      flavourWithNewName.Name,
					ID:        flavour.ID,
					CreatedAt: flavour.CreatedAt,
				}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "Not Found",
			flavourId: flavour.ID.String(),
			body: gin.H{
				"name": flavourWithNewName.Name,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateFlavourParams{
					ID:   flavour.ID,
					Name: flavourWithNewName.Name,
				}
				store.EXPECT().UpdateFlavour(gomock.Any(), arg).Times(1).Return(db.Flavour{
					Name:      flavourWithNewName.Name,
					ID:        flavour.ID,
					CreatedAt: flavour.CreatedAt,
				}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "Internal Server Error",
			flavourId: flavour.ID.String(),
			body: gin.H{
				"name": flavourWithNewName.Name,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateFlavourParams{
					ID:   flavour.ID,
					Name: flavourWithNewName.Name,
				}
				store.EXPECT().UpdateFlavour(gomock.Any(), arg).Times(1).Return(db.Flavour{
					Name:      flavourWithNewName.Name,
					ID:        flavour.ID,
					CreatedAt: flavour.CreatedAt,
				}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "Unauthorized",
			flavourId: flavour.ID.String(),
			body: gin.H{
				"name": flavourWithNewName.Name,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateFlavourParams{
					ID:   flavour.ID,
					Name: flavourWithNewName.Name,
				}
				store.EXPECT().UpdateFlavour(gomock.Any(), arg).Times(0).Return(db.Flavour{
					Name:      flavourWithNewName.Name,
					ID:        flavour.ID,
					CreatedAt: flavour.CreatedAt,
				}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		store := mockdb.NewMockStore(ctrl)
		categorizer := mockcategorize.NewMockCategorizerService(ctrl)
		embeddingService := mockembedding.NewMockEmbeddingService(ctrl)
		tc.buildStubs(store)

		server := newTestServer(t, store, nil, nil, categorizer, embeddingService)
		recorder := httptest.NewRecorder()

		data, err := json.Marshal(tc.body)
		require.NoError(t, err)

		url := fmt.Sprintf("/api/v1/flavours/%s", tc.flavourId)
		request, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(data))
		require.NoError(t, err)

		tc.setupAuth(t, request, server.tokenMaker)
		server.router.ServeHTTP(recorder, request)
		tc.checkResponse(recorder)
	}
}

func TestGetAllFlavoursApi(t *testing.T) {
	user := uuid.New()
	var flavours []db.Flavour
	for i := 1; i <= 10; i++ {
		flavour := randomFlavour(t)
		flavours = append(flavours, flavour)
	}

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAllFlavours(gomock.Any()).Times(1).Return(flavours, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Not Found",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAllFlavours(gomock.Any()).Times(1).Return(flavours, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "Internal Server Error",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAllFlavours(gomock.Any()).Times(1).Return(flavours, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Unauthorized",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAllFlavours(gomock.Any()).Times(0).Return(flavours, nil)
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
			categorizer := mockcategorize.NewMockCategorizerService(ctrl)
			embeddingService := mockembedding.NewMockEmbeddingService(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store, nil, nil, categorizer, embeddingService)

			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/api/v1/flavours")
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func requireBodyMatchFlavour(t *testing.T, body *bytes.Buffer, flavour db.Flavour) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var reqFlavour db.Flavour
	err = json.Unmarshal(data, &reqFlavour)
	require.NoError(t, err)

	require.Equal(t, flavour.ID, reqFlavour.ID)
	require.Equal(t, flavour.Name, reqFlavour.Name)
	require.WithinDuration(t, flavour.CreatedAt, reqFlavour.CreatedAt, time.Millisecond)
}

func randomFlavour(t *testing.T) db.Flavour {
	return db.Flavour{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		Name:      util.RandomString(10),
	}
}
