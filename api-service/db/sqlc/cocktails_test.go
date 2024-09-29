package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/m1thrandir225/galore-services/dto"
	"github.com/m1thrandir225/galore-services/util"
	"github.com/pgvector/pgvector-go"
	"github.com/stretchr/testify/require"
)

func createRandomCocktail(t *testing.T) Cocktail {
	var is_alcoholic pgtype.Bool

	is_alcoholic.Scan(util.RandomBool())

	ingredients := dto.IngredientDto{
		Ingredients: util.RandomIngredients(),
	}

	instructions := dto.InstructionDto{
		Instructions: util.RandomInstructions(),
	}

	floatArr := util.RandomFloatArray(0.1, 1.0, 768)

	embedding := pgvector.NewVector(floatArr)

	arg := CreateCocktailParams{
		Name:         util.RandomString(40),
		Image:        util.RandomString(80),
		Glass:        util.RandomString(12),
		Instructions: instructions,
		Ingredients:  ingredients,
		IsAlcoholic:  is_alcoholic,
		Embedding:    embedding,
	}

	cocktail, err := testStore.CreateCocktail(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, cocktail)

	require.Equal(t, arg.Name, cocktail.Name)
	require.Equal(t, arg.Image, cocktail.Image)
	require.Equal(t, arg.Glass, cocktail.Glass)
	require.Equal(t, arg.Instructions, cocktail.Instructions)
	require.Equal(t, arg.Ingredients, cocktail.Ingredients)
	require.Equal(t, arg.IsAlcoholic, cocktail.IsAlcoholic)
	require.Equal(t, arg.Embedding, cocktail.Embedding)

	return cocktail
}

func TestCreateCocktail(t *testing.T) {
	createRandomCocktail(t)
}

func TestGetCocktail(t *testing.T) {
	cocktail := createRandomCocktail(t)

	selected, err := testStore.GetCocktail(context.Background(), cocktail.ID)

	require.NoError(t, err)
	require.NotEmpty(t, selected)

	require.Equal(t, cocktail.ID, selected.ID)
	require.Equal(t, cocktail.Name, selected.Name)
	require.Equal(t, cocktail.Glass, selected.Glass)
	require.Equal(t, cocktail.IsAlcoholic, selected.IsAlcoholic)
	require.Equal(t, cocktail.Instructions, selected.Instructions)
	require.Equal(t, cocktail.Ingredients, selected.Ingredients)
	require.WithinDuration(t, cocktail.CreatedAt, selected.CreatedAt, time.Second)
}

func TestDeleteCocktail(t *testing.T) {
	cocktail := createRandomCocktail(t)

	err := testStore.DeleteCocktail(context.Background(), cocktail.ID)

	require.NoError(t, err)

	selected, err := testStore.GetCocktail(context.Background(), cocktail.ID)

	require.Error(t, err)
	require.Empty(t, selected)
	require.EqualError(t, err, ErrRecordNotFound.Error())
}

func TestUpdateCocktail(t *testing.T) {
	cocktail := createRandomCocktail(t)

	var isAlcoholic pgtype.Bool
	ingredients := dto.IngredientDto{
		Ingredients: util.RandomIngredients(),
	}

	instructions := dto.InstructionDto{
		Instructions: util.RandomInstructions(),
	}
	err := isAlcoholic.Scan(!cocktail.IsAlcoholic.Bool)
	require.NoError(t, err)

	arg := UpdateCocktailParams{
		ID:           cocktail.ID,
		Name:         util.RandomString(48),
		Glass:        util.RandomString(12),
		IsAlcoholic:  isAlcoholic,
		Image:        util.RandomString(80),
		Instructions: instructions,
		Ingredients:  ingredients,
	}

	updated, err := testStore.UpdateCocktail(context.Background(), arg)

	require.NoError(t, err)

	require.NotEmpty(t, updated)

	require.Equal(t, arg.ID, updated.ID)
	require.Equal(t, arg.Name, updated.Name)
	require.Equal(t, arg.Glass, updated.Glass)
	require.Equal(t, arg.IsAlcoholic, updated.IsAlcoholic)
	require.Equal(t, arg.Instructions, updated.Instructions)
	require.Equal(t, arg.Ingredients, updated.Ingredients)
	require.Equal(t, arg.IsAlcoholic, updated.IsAlcoholic)
	require.WithinDuration(t, cocktail.CreatedAt, updated.CreatedAt, time.Second)

	require.Equal(t, cocktail.ID, updated.ID)

	require.NotEqual(t, cocktail.Glass, updated.Glass)
	require.NotEqual(t, cocktail.Ingredients, updated.Ingredients)
	require.NotEqual(t, cocktail.Instructions, updated.Instructions)
	require.NotEqual(t, cocktail.Glass, updated.Glass)
	require.NotEqual(t, cocktail.IsAlcoholic, updated.IsAlcoholic)
	require.NotEqual(t, cocktail.Image, updated.Image)
}
