package db

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/m1thrandir225/galore-services/dto"
	"github.com/m1thrandir225/galore-services/util"
	"github.com/stretchr/testify/require"
)

func createRandomUserCocktail(user_id uuid.UUID, t *testing.T) CreatedCocktail {

	ingredients := dto.IngredientDto{
		Ingredients: util.RandomIngredients(),
	}

	instructions := dto.AiInstructionDto{
		Instructions: util.RandomAiInstructions(),
	}

	arg := CreateUserCocktailParams{
		Name:         util.RandomString(10),
		Image:        util.RandomString(24),
		Ingredients:  ingredients,
		Instructions: instructions,
		UserID:       user_id,
		Description:  util.RandomString(256),
	}

	cocktail, err := testStore.CreateUserCocktail(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, cocktail)

	require.Equal(t, arg.Name, cocktail.Name)
	require.Equal(t, arg.Image, cocktail.Image)
	require.Equal(t, arg.Ingredients, cocktail.Ingredients)
	require.Equal(t, arg.Instructions, cocktail.Instructions)
	require.Equal(t, arg.UserID, cocktail.UserID)
	require.Equal(t, arg.Description, cocktail.Description)

	return cocktail
}

func TestCreateUserCocktail(t *testing.T) {
	user := createRandomUser(t)
	createRandomUserCocktail(user.ID, t)
}

func TestGetUserCocktail(t *testing.T) {
	user := createRandomUser(t)

	cocktail := createRandomUserCocktail(user.ID, t)

	selected_cocktail, err := testStore.GetUserCocktail(context.Background(), cocktail.ID)

	require.NoError(t, err)
	require.NotEmpty(t, selected_cocktail)

	require.Equal(t, selected_cocktail.Name, cocktail.Name)
	require.Equal(t, selected_cocktail.Image, cocktail.Image)
	require.Equal(t, selected_cocktail.Ingredients, cocktail.Ingredients)
	require.Equal(t, selected_cocktail.Instructions, cocktail.Instructions)
	require.Equal(t, selected_cocktail.UserID, cocktail.UserID)
	require.Equal(t, selected_cocktail.Description, cocktail.Description)
}

func TestDeleteUserCocktail(t *testing.T) {
	user := createRandomUser(t)

	cocktail := createRandomUserCocktail(user.ID, t)

	err := testStore.DeleteUserCocktail(context.Background(), cocktail.ID)

	require.NoError(t, err)

	selected_ccktail, err := testStore.GetUserCocktail(context.Background(), cocktail.ID)

	require.Error(t, err)
	require.Empty(t, selected_ccktail)
	require.EqualError(t, err, ErrRecordNotFound.Error())
}
