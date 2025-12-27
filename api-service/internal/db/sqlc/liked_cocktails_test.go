package db

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func createRandomLikedCocktail(userId, cocktailId uuid.UUID, t *testing.T) LikedCocktail {
	arg := LikeCocktailParams{
		UserID:     userId,
		CocktailID: cocktailId,
	}

	likedCocktail, err := testStore.LikeCocktail(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, likedCocktail)

	require.Equal(t, arg.CocktailID, likedCocktail.CocktailID)
	require.Equal(t, arg.UserID, likedCocktail.UserID)

	return likedCocktail
}

func TestLikeCocktail(t *testing.T) {
	user := createRandomUser(t)
	cocktail := createRandomCocktail(t)
	likedCocktail := createRandomLikedCocktail(user.ID, cocktail.ID, t)

	require.NotEmpty(t, likedCocktail)
	require.Equal(t, cocktail.ID, likedCocktail.CocktailID)
	require.Equal(t, user.ID, likedCocktail.UserID)
	require.NotEmpty(t, likedCocktail.ID)
}

func TestGetLikedCocktail(t *testing.T) {
	user := createRandomUser(t)
	cocktail := createRandomCocktail(t)
	likedCocktail := createRandomLikedCocktail(user.ID, cocktail.ID, t)

	arg := GetLikedCocktailParams{
		UserID:     user.ID,
		CocktailID: cocktail.ID,
	}

	selected, err := testStore.GetLikedCocktail(context.Background(), arg)
	require.NoError(t, err)

	require.Equal(t, likedCocktail.ID, selected.ID_2)
	require.Equal(t, user.ID, selected.UserID)
	require.Equal(t, cocktail.ID, selected.CocktailID)
	require.Equal(t, cocktail.ID, selected.ID)
	require.Equal(t, cocktail.Name, selected.Name)
	require.WithinDuration(t, cocktail.CreatedAt, selected.CreatedAt, time.Millisecond)
	require.Equal(t, cocktail.Glass, selected.Glass)
	require.Equal(t, cocktail.Ingredients, selected.Ingredients)
	require.Equal(t, cocktail.Instructions, selected.Instructions)
	require.Equal(t, cocktail.IsAlcoholic, selected.IsAlcoholic)
}

func TestGetLikedCocktails(t *testing.T) {
	user := createRandomUser(t)
	cocktail := createRandomCocktail(t)

	var likedCocktails []LikedCocktail
	for i := 0; i < 10; i++ {
		liked := createRandomLikedCocktail(user.ID, cocktail.ID, t)
		likedCocktails = append(likedCocktails, liked)
	}

	selected, err := testStore.GetLikedCocktails(context.Background(), user.ID)
	require.NoError(t, err)

	require.Equal(t, len(likedCocktails), len(selected))
	require.NotEmpty(t, selected)

	for i := 0; i < len(likedCocktails); i++ {
		require.NotEmpty(t, selected[i])
		require.NotEmpty(t, selected[i].ID)
		require.NotEmpty(t, selected[i].Name)
		require.NotEmpty(t, selected[i].Glass)
		require.NotEmpty(t, selected[i].Instructions)
		require.NotEmpty(t, selected[i].IsAlcoholic)
		require.NotEmpty(t, selected[i].Ingredients)
		require.NotEmpty(t, selected[i].Image)
		require.NotEmpty(t, selected[i].CreatedAt)
		require.Equal(t, selected[i].CocktailID, cocktail.ID)
		require.Equal(t, selected[i].UserID, user.ID)
	}
}

func TestUnlikeCocktail(t *testing.T) {
	user := createRandomUser(t)
	cocktail := createRandomCocktail(t)
	createRandomLikedCocktail(user.ID, cocktail.ID, t)
	arg := UnlikeCocktailParams{
		UserID:     user.ID,
		CocktailID: cocktail.ID,
	}
	err := testStore.UnlikeCocktail(context.Background(), arg)
	require.NoError(t, err)

	selected, err := testStore.GetLikedCocktail(context.Background(), GetLikedCocktailParams{
		UserID:     user.ID,
		CocktailID: cocktail.ID,
	})
	require.Error(t, err)
	require.Empty(t, selected)
	require.EqualError(t, err, ErrRecordNotFound.Error())
}
