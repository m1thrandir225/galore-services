package db

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func createRandomLikedCocktail(user_id, cocktail_id uuid.UUID, t *testing.T) LikedCocktail {
	arg := LikeCocktailParams{
		UserID:     user_id,
		CocktailID: cocktail_id,
	}

	liked_cocktail, err := testStore.LikeCocktail(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, liked_cocktail)

	require.Equal(t, arg.CocktailID, liked_cocktail.CocktailID)
	require.Equal(t, arg.UserID, liked_cocktail.UserID)

	return liked_cocktail
}

func TestLikeCocktail(t *testing.T) {
	user := createRandomUser(t)
	cocktail := createRandomCocktail(t)
	createRandomLikedCocktail(user.ID, cocktail.ID, t)
}

func TestGetLikedCocktails(t *testing.T) {
	user := createRandomUser(t)
	cocktail := createRandomCocktail(t)
	createRandomLikedCocktail(user.ID, cocktail.ID, t)
}

func TestGetLikedCocktail(t *testing.T) {
	user := createRandomUser(t)
	cocktail := createRandomCocktail(t)
	createRandomLikedCocktail(user.ID, cocktail.ID, t)
}

func TestUnlikeCocktail(t *testing.T) {
	user := createRandomUser(t)
	cocktail := createRandomCocktail(t)
	createRandomLikedCocktail(user.ID, cocktail.ID, t)

}
