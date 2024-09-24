package db

import (
	"context"
	"testing"

	"github.com/m1thrandir225/galore-services/util"
	"github.com/stretchr/testify/require"
)

func createRandomFlavour(t *testing.T) Flavour {
	flavourName := util.RandomString(8)
	flavour, err := testStore.CreateFlavour(context.Background(), flavourName)
	require.NoError(t, err)
	require.NotEmpty(t, flavour)

	require.Equal(t, flavour.Name, flavourName)
	return flavour
}

func TestCreateFlavour(t *testing.T) {
	createRandomFlavour(t)
}

func TestGetFlavourID(t *testing.T) {
	flavour := createRandomFlavour(t)

	selectedFlavour, err := testStore.GetFlavourId(context.Background(), flavour.ID)
	require.NoError(t, err)
	require.NotEmpty(t, selectedFlavour)

	require.Equal(t, selectedFlavour.ID, flavour.ID)
	require.Equal(t, selectedFlavour.Name, flavour.Name)
	require.Equal(t, selectedFlavour.CreatedAt, flavour.CreatedAt)
}

func TestGetFlavourName(t *testing.T) {
	flavour := createRandomFlavour(t)

	selectedFlavour, err := testStore.GetFlavourName(context.Background(), flavour.Name)
	require.NoError(t, err)
	require.NotEmpty(t, selectedFlavour)

	require.Equal(t, selectedFlavour.ID, flavour.ID)
	require.Equal(t, selectedFlavour.Name, flavour.Name)
	require.Equal(t, selectedFlavour.CreatedAt, flavour.CreatedAt)
}

func TestDeleteFlavour(t *testing.T) {
	flavour := createRandomFlavour(t)

	err := testStore.DeleteFlavour(context.Background(), flavour.ID)
	require.NoError(t, err)

	selectedFlavour, err := testStore.GetFlavourId(context.Background(), flavour.ID)
	require.Error(t, err)
	require.Empty(t, selectedFlavour)
	require.EqualError(t, err, ErrRecordNotFound.Error())
}

func TestGetAllFlavours(t *testing.T) {
	flavours, err := testStore.GetAllFlavours(context.Background())

	require.NoError(t, err)

	require.NotEmpty(t, flavours)

	require.True(t, len(flavours) > 0)
}
