package db

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func createRandomCocktailCategory(cocktail_id, category_id uuid.UUID, t *testing.T) CocktailCategory {

	arg := CreateCocktailCategoryParams{
		CocktailID: cocktail_id,
		CategoryID: category_id,
	}

	c_c, err := testStore.CreateCocktailCategory(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, c_c)

	require.Equal(t, arg.CategoryID, c_c.CategoryID)
	require.Equal(t, arg.CocktailID, c_c.CocktailID)

	return c_c
}

func TestCreateCocktailCategory(t *testing.T) {
	cocktail := createRandomCocktail(t)
	category := createRandomCategory(t)
	createRandomCocktailCategory(cocktail.ID, category.ID, t)
}
func TestGetCocktailCategory(t *testing.T) {
	cocktail := createRandomCocktail(t)
	category := createRandomCategory(t)
	cocktail_category := createRandomCocktailCategory(cocktail.ID, category.ID, t)

	selected, err := testStore.GetCocktailCategory(context.Background(), cocktail_category.ID)

	require.NoError(t, err)

	require.NotEmpty(t, selected)

	require.Equal(t, cocktail_category.ID, selected.ID)
	require.Equal(t, cocktail_category.CategoryID, selected.CategoryID)
	require.Equal(t, cocktail_category.CocktailID, selected.CocktailID)

}

func TestGetCategoriesForCocktail(t *testing.T) {
	cocktail := createRandomCocktail(t)
	var categories []Category
	for i := 0; i < 10; i++ {
		category := createRandomCategory(t)
		createRandomCocktailCategory(cocktail.ID, category.ID, t)
		categories = append(categories, category)
	}

	categories_for_cocktail, err := testStore.GetCategoriesForCocktail(context.Background(), cocktail.ID)
	require.NoError(t, err)
	require.NotEmpty(t, categories_for_cocktail)

	require.Equal(t, len(categories), len(categories_for_cocktail))
}

func TestGetCocktailsForCategory(t *testing.T) {
	category := createRandomCategory(t)

	var cocktails []Cocktail

	for i := 0; i < 10; i++ {
		cocktail := createRandomCocktail(t)

		createRandomCocktailCategory(cocktail.ID, category.ID, t)
		cocktails = append(cocktails, cocktail)
	}

	cocktails_for_category, err := testStore.GetCocktailsForCategory(context.Background(), category.ID)

	require.NoError(t, err)
	require.NotEmpty(t, cocktails_for_category)

	require.Equal(t, len(cocktails), len(cocktails_for_category))

}

func TestDeleteCocktailCategory(t *testing.T) {
	cocktail := createRandomCocktail(t)
	category := createRandomCategory(t)
	cocktail_category := createRandomCocktailCategory(cocktail.ID, category.ID, t)

	err := testStore.DeleteCocktailCategory(context.Background(), cocktail_category.ID)
	require.NoError(t, err)

	selected, err := testStore.GetCocktailCategory(context.Background(), cocktail_category.ID)
	require.Error(t, err)
	require.Empty(t, selected)
	require.EqualError(t, err, ErrRecordNotFound.Error())
}
