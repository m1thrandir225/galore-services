package db

import (
	"context"
	"testing"

	"github.com/m1thrandir225/galore-services/util"
	"github.com/stretchr/testify/require"
)

func createRandomCategory(t *testing.T) Category {
	name := util.RandomString(12)

	category, err := testStore.CreateCategory(context.Background(), name)

	require.NoError(t, err)
	require.NotEmpty(t, category)
	require.Equal(t, category.Name, name)

	return category
}

func TestCreateCategory(t *testing.T) {
	createRandomCategory(t)
}

func TestGetCategoryById(t *testing.T) {
	category := createRandomCategory(t)

	selected_category, err := testStore.GetCategoryById(context.Background(), category.ID)

	require.NoError(t, err)
	require.NotEmpty(t, selected_category)

	require.Equal(t, category.ID, selected_category.ID)
	require.Equal(t, category.Name, selected_category.Name)
	require.Equal(t, category.CreatedAt, selected_category.CreatedAt)
}

func TestGetCategoryByName(t *testing.T) {

	category := createRandomCategory(t)

	selected_category, err := testStore.GetCategoryByName(context.Background(), category.Name)

	require.NoError(t, err)
	require.NotEmpty(t, selected_category)

	require.Equal(t, category.ID, selected_category.ID)
	require.Equal(t, category.Name, selected_category.Name)
	require.Equal(t, category.CreatedAt, selected_category.CreatedAt)
}

func TestUpdateCategory(t *testing.T) {
	category := createRandomCategory(t)

	new_name := util.RandomString(12)
	arg := UpdateCategoryParams{
		ID:   category.ID,
		Name: new_name,
	}

	updated, err := testStore.UpdateCategory(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, updated)

	require.Equal(t, category.ID, updated.ID)
	require.Equal(t, category.CreatedAt, updated.CreatedAt)
	require.NotEqual(t, category.Name, updated.Name)
	require.Equal(t, new_name, updated.Name)
}

func TestDeleteCategory(t *testing.T) {
	category := createRandomCategory(t)

	err := testStore.DeleteCategory(context.Background(), category.ID)
	require.NoError(t, err)
	selected_category, err := testStore.GetCategoryById(context.Background(), category.ID)
	require.Error(t, err)
	require.Empty(t, selected_category)
	require.EqualError(t, err, ErrRecordNotFound.Error())
}
