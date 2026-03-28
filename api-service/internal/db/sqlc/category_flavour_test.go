package db

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func createRandomCategoryFlavour(t *testing.T, categoryID, flavourID uuid.UUID) CategoryFlavour {
	arg := CreateCategoryFlavourParams{
		CategoryID: categoryID,
		FlavourID:  flavourID,
	}

	categoryFlavour, err := testStore.CreateCategoryFlavour(context.Background(), arg)

	require.NoError(t, err)

	return categoryFlavour
}

func TestCreateCategoryFlavour(t *testing.T) {
	category := createRandomCategory(t)
	flavour := createRandomFlavour(t)
	createRandomCategoryFlavour(t, category.ID, flavour.ID)
}

func TestGetCategoryFlavour(t *testing.T) {
	category := createRandomCategory(t)
	flavour := createRandomFlavour(t)
	categoryFlavour := createRandomCategoryFlavour(t, category.ID, flavour.ID)

	selected, err := testStore.GetCategoryFlavour(context.Background(), categoryFlavour.ID)
	require.NoError(t, err)

	require.Equal(t, categoryFlavour.ID, selected.ID)
	require.Equal(t, categoryFlavour.CategoryID, selected.CategoryID)
	require.Equal(t, categoryFlavour.FlavourID, selected.FlavourID)
}

func TestDeleteCategoryFlavour(t *testing.T) {
	category := createRandomCategory(t)
	flavour := createRandomFlavour(t)
	categoryFlavour := createRandomCategoryFlavour(t, category.ID, flavour.ID)

	err := testStore.DeleteCategoryFlavour(context.Background(), categoryFlavour.ID)
	require.NoError(t, err)

	deletedItem, err := testStore.GetCategoryFlavour(context.Background(), categoryFlavour.ID)
	require.Error(t, err)
	require.Empty(t, deletedItem)

}

func TestGetCategoriesFromLikedFlavours(t *testing.T) {
	//TODO: implement
}
