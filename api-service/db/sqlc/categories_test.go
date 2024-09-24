package db

import (
	"context"
	"testing"

	"github.com/m1thrandir225/galore-services/util"
	"github.com/stretchr/testify/require"
)

func createRandomCategory(t *testing.T) Category {

	arg := CreateCategoryParams{
		Name: util.RandomString(48),
		Tag:  util.RandomString(12),
	}
	category, err := testStore.CreateCategory(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, category)
	require.Equal(t, arg.Name, category.Name)
	require.Equal(t, arg.Tag, category.Tag)

	return category
}

func TestCreateCategory(t *testing.T) {
	createRandomCategory(t)
}

func TestGetCategoryById(t *testing.T) {
	category := createRandomCategory(t)

	selectedCategory, err := testStore.GetCategoryById(context.Background(), category.ID)

	require.NoError(t, err)
	require.NotEmpty(t, selectedCategory)

	require.Equal(t, category.ID, selectedCategory.ID)
	require.Equal(t, category.Name, selectedCategory.Name)
	require.Equal(t, category.CreatedAt, selectedCategory.CreatedAt)
}

func TestGetCategoryByName(t *testing.T) {

	category := createRandomCategory(t)

	selectedCategory, err := testStore.GetCategoryByTag(context.Background(), category.Tag)

	require.NoError(t, err)
	require.NotEmpty(t, selectedCategory)

	require.Equal(t, category.ID, selectedCategory.ID)
	require.Equal(t, category.Name, selectedCategory.Name)
	require.Equal(t, category.Tag, selectedCategory.Tag)
	require.Equal(t, category.CreatedAt, selectedCategory.CreatedAt)
}

func TestUpdateCategory(t *testing.T) {
	category := createRandomCategory(t)

	newName := util.RandomString(48)
	newTag := util.RandomString(18)

	arg := UpdateCategoryParams{
		ID:   category.ID,
		Name: newName,
		Tag:  newTag,
	}

	updated, err := testStore.UpdateCategory(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, updated)

	require.Equal(t, category.ID, updated.ID)
	require.Equal(t, category.CreatedAt, updated.CreatedAt)
	require.NotEqual(t, category.Name, updated.Name)
	require.NotEqual(t, category.Tag, updated.Tag)
	require.Equal(t, arg.Tag, updated.Tag)
	require.Equal(t, arg.Name, updated.Name)
}

func TestDeleteCategory(t *testing.T) {
	category := createRandomCategory(t)

	err := testStore.DeleteCategory(context.Background(), category.ID)
	require.NoError(t, err)
	selectedCategory, err := testStore.GetCategoryById(context.Background(), category.ID)
	require.Error(t, err)
	require.Empty(t, selectedCategory)
	require.EqualError(t, err, ErrRecordNotFound.Error())
}

func TestGetAllCategories(t *testing.T) {
	category := createRandomCategory(t)

	categories, err := testStore.GetAllCategories(context.Background())

	require.NoError(t, err)
	require.NotEmpty(t, categories)
	require.Equal(t, category.ID, categories[0].ID)
	require.Equal(t, category.Name, categories[0].Name)
	require.Equal(t, category.CreatedAt, categories[0].CreatedAt)
	require.Equal(t, category.Tag, categories[0].Tag)
	require.Equal(t, category.CreatedAt, categories[0].CreatedAt)

	require.True(t, len(categories) > 0)
}
