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
	"github.com/m1thrandir225/galore-services/token"
	"github.com/m1thrandir225/galore-services/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestCreateCategoryApi(t *testing.T) {
	userId, err := uuid.NewRandom()
	require.NoError(t, err)
	category := randomCategory(t)
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
				"name": category.Name,
				"tag":  category.Tag,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateCategoryParams{
					Name: category.Name,
					Tag:  category.Tag,
				}
				store.EXPECT().
					CreateCategory(gomock.Any(), arg).
					Times(1).
					Return(category, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchCategory(t, recorder.Body, category)
			},
		},
		{
			name: "Internal Error",
			body: gin.H{
				"name": category.Name,
				"tag":  category.Tag,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateCategoryParams{
					Tag:  category.Tag,
					Name: category.Name,
				}
				store.EXPECT().
					CreateCategory(gomock.Any(), arg).
					Times(1).
					Return(db.Category{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Bad Request",
			body: gin.H{},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateCategory(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Unauthorized",
			body: gin.H{
				"name": category.Name,
				"tag":  category.Tag,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateCategory(gomock.Any(), gomock.Any()).
					Times(0)
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
			testCase.buildStubs(store)

			server := newTestServer(t, store, nil)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(testCase.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/api/v1/categories")
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			testCase.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(recorder)
		})
	}
}
func TestGetCategoriesApi(t *testing.T) {
	userId, err := uuid.NewRandom()
	require.NoError(t, err)

	categories := make([]db.Category, 5)
	for i := range categories {
		categories[i] = randomCategory(t)
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
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAllCategories(gomock.Any()).
					Times(1).
					Return(categories, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Internal Error",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAllCategories(gomock.Any()).
					Times(1).Return(nil, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Not found",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAllCategories(gomock.Any()).
					Times(1).Return(nil, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "Unauthorized",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAllCategories(gomock.Any()).
					Times(0)
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
			testCase.buildStubs(store)

			server := newTestServer(t, store, nil)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/api/v1/categories")
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			testCase.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(recorder)
		})
	}
}
func TestGetCategoryByIdApi(t *testing.T) {
	userId, err := uuid.NewRandom()
	require.NoError(t, err)

	category := randomCategory(t)

	testCases := []struct {
		name          string
		categoryId    string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:       "OK",
			categoryId: category.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetCategoryById(gomock.Any(), category.ID).
					Times(1).
					Return(category, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchCategory(t, recorder.Body, category)
			},
		},
		{
			name:       "Not Valid Id",
			categoryId: "test",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetCategoryById(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:       "Not Found",
			categoryId: "30bb0f4e-a2ad-4405-8d8a-17f32deb9f6a",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetCategoryById(gomock.Any(), gomock.Any()).Times(1).Return(db.Category{}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:       "Internal Error",
			categoryId: category.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetCategoryById(gomock.Any(), gomock.Any()).Times(1).Return(db.Category{}, sql.ErrConnDone)
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

			url := fmt.Sprintf("/api/v1/categories/%s", testCase.categoryId)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			testCase.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(recorder)
		})
	}
}
func TestDeleteCategoryApi(t *testing.T) {
	userId, err := uuid.NewRandom()
	require.NoError(t, err)

	category := randomCategory(t)

	testCases := []struct {
		name          string
		categoryId    string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:       "OK",
			categoryId: category.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteCategory(gomock.Any(), category.ID).Times(1).Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

			},
		},
		{
			name:       "Not Found",
			categoryId: "30bb0f4e-a2ad-4405-8d8a-17f32deb9f6a",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteCategory(gomock.Any(), gomock.Any()).Times(1).Return(sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)

			},
		},
		{
			name:       "Bad Request",
			categoryId: util.RandomString(10),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteCategory(gomock.Any(), gomock.Any()).Times(0).Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:       "Internal Error",
			categoryId: category.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteCategory(gomock.Any(), gomock.Any()).Times(1).Return(sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)

			},
		},
		{
			name:       "Unauthorized",
			categoryId: category.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteCategory(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)

			},
		},
	}
	for i := range testCases {
		testCase := testCases[i]

		t.Run(testCase.name, func(t *testing.T) {
			t.Run(testCase.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				store := mockdb.NewMockStore(ctrl)
				testCase.buildStubs(store)

				server := newTestServer(t, store, nil)
				recorder := httptest.NewRecorder()

				url := fmt.Sprintf("/api/v1/categories/%s", testCase.categoryId)
				request, err := http.NewRequest(http.MethodDelete, url, nil)
				require.NoError(t, err)

				testCase.setupAuth(t, request, server.tokenMaker)
				server.router.ServeHTTP(recorder, request)
				testCase.checkResponse(recorder)
			})
		})
	}
}
func TestUpdateCategoryApi(t *testing.T) {
	userId, err := uuid.NewRandom()
	require.NoError(t, err)

	category := randomCategory(t)
	updatedName := util.RandomString(16)

	testCases := []struct {
		name          string
		categoryId    string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:       "OK",
			categoryId: category.ID.String(),
			body: gin.H{
				"name": updatedName,
				"tag":  strings.ToLower(updatedName),
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateCategoryParams{
					ID:   category.ID,
					Name: updatedName,
					Tag:  strings.ToLower(updatedName),
				}
				store.EXPECT().UpdateCategory(gomock.Any(), arg).Times(1).Return(db.Category{
					ID:        category.ID,
					Name:      updatedName,
					Tag:       strings.ToLower(updatedName),
					CreatedAt: category.CreatedAt,
				}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:       "Not Found",
			categoryId: "30bb0f4e-a2ad-4405-8d8a-17f32deb9f6a",
			body: gin.H{
				"name": updatedName,
				"tag":  strings.ToLower(updatedName),
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateCategory(gomock.Any(), gomock.Any()).Times(1).Return(db.Category{}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:       "Bad Request - No Id",
			categoryId: "aaa",
			body: gin.H{
				"name": updatedName,
				"tag":  strings.ToLower(updatedName),
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateCategory(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:       "Bad Request - No Data",
			categoryId: category.ID.String(),
			body:       gin.H{},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateCategory(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:       "Internal Error",
			categoryId: category.ID.String(),
			body: gin.H{
				"name": updatedName,
				"tag":  strings.ToLower(updatedName),
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateCategory(gomock.Any(), gomock.Any()).Times(1).Return(db.Category{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		testCase := testCases[i]

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		store := mockdb.NewMockStore(ctrl)
		testCase.buildStubs(store)

		server := newTestServer(t, store, nil)
		recorder := httptest.NewRecorder()

		data, err := json.Marshal(testCase.body)
		require.NoError(t, err)

		url := fmt.Sprintf("/api/v1/categories/%s", testCase.categoryId)
		request, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(data))
		require.NoError(t, err)

		testCase.setupAuth(t, request, server.tokenMaker)
		server.router.ServeHTTP(recorder, request)
		testCase.checkResponse(recorder)

	}
}

func randomCategory(t *testing.T) db.Category {
	name := util.RandomString(10)
	id, err := uuid.NewRandom()
	require.NoError(t, err)
	return db.Category{
		ID:        id,
		Name:      name,
		Tag:       strings.ToLower(name),
		CreatedAt: time.Now(),
	}
}

func requireBodyMatchCategory(t *testing.T, body *bytes.Buffer, category db.Category) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var reqCategory db.Category
	err = json.Unmarshal(data, &reqCategory)

	require.NoError(t, err)
	require.Equal(t, category.Name, reqCategory.Name)
	require.Equal(t, category.Tag, reqCategory.Tag)
	require.WithinDuration(t, category.CreatedAt, reqCategory.CreatedAt, time.Second)
	require.Equal(t, category.ID, reqCategory.ID)
}
