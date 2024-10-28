package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	mockcache "github.com/m1thrandir225/galore-services/cache/mock"
	mockdb "github.com/m1thrandir225/galore-services/db/mock"
	"github.com/m1thrandir225/galore-services/dto"
	"github.com/m1thrandir225/galore-services/token"
	"github.com/pgvector/pgvector-go"
	"github.com/stretchr/testify/require"

	"github.com/google/uuid"
	db "github.com/m1thrandir225/galore-services/db/sqlc"
	"github.com/m1thrandir225/galore-services/util"
)

func TestCreateCocktailApi(t *testing.T) {}

func TestGetCocktailApi(t *testing.T) {
	userId := uuid.New()
	cocktail := randomCocktail(t)

	testCases := []struct {
		name          string
		cocktailId    string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore, cache *mockcache.MockKeyValueStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:       "OK",
			cocktailId: cocktail.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, userId, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockKeyValueStore) {
				//
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
		})
	}
}

func TestUpdateCocktailApi(t *testing.T) {}

func TestDeleteCocktailApi(t *testing.T) {}

func requireBodyMatchCocktail(t *testing.T, body *bytes.Buffer, cocktail db.Cocktail) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var reqCocktail db.Cocktail
	err = json.Unmarshal(data, &reqCocktail)
	require.NoError(t, err)

	require.Equal(t, reqCocktail.ID, cocktail.ID)
	require.Equal(t, reqCocktail.Glass, cocktail.Glass)
	require.Equal(t, reqCocktail.Name, cocktail.Name)
	require.Equal(t, reqCocktail.IsAlcoholic, cocktail.IsAlcoholic)
	require.Equal(t, reqCocktail.Embedding, cocktail.Embedding)
	require.Equal(t, reqCocktail.Ingredients, cocktail.Ingredients)
	require.Equal(t, reqCocktail.Instructions, cocktail.Instructions)
}

func randomCocktail(t *testing.T) db.Cocktail {
	var isAlcoholic pgtype.Bool
	err := isAlcoholic.Scan(util.RandomBool())
	require.NoError(t, err)

	var embedding pgvector.Vector

	err = embedding.Scan(util.RandomFloatArray(10, 150, 768))
	require.NoError(t, err)

	return db.Cocktail{
		ID:        uuid.New(),
		Name:      util.RandomString(10),
		CreatedAt: time.Now(),
		Image:     util.RandomString(10),
		Ingredients: dto.IngredientDto{
			Ingredients: util.RandomIngredients(),
		},
		Instructions: strings.Join(util.RandomInstructions(), ","),
		Glass:        util.RandomString(10),
		IsAlcoholic:  isAlcoholic,
		Embedding:    embedding,
	}
}
