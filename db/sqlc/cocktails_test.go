package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/m1thrandir225/galore-services/dto"
	"github.com/m1thrandir225/galore-services/util"
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

	arg := CreateCocktailParams{
		Name:         util.RandomString(40),
		Image:        util.RandomString(80),
		Glass:        util.RandomString(12),
		Instructions: instructions,
		Ingredients:  ingredients,
		IsAlcoholic:  is_alcoholic,
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

	return cocktail
}

func TestCreateCocktail(t *testing.T) {
	createRandomCocktail(t)
}
