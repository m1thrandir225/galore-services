package db

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func createLikedFlavour(user_id uuid.UUID, t *testing.T) LikedFlavour {
	flavour := createRandomFlavour(t)

	arg := LikeFlavourParams{
		UserID:    user_id,
		FlavourID: flavour.ID,
	}

	likedFlavour, err := testStore.LikeFlavour(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, likedFlavour)

	require.Equal(t, arg.UserID, likedFlavour.UserID)
	require.Equal(t, arg.FlavourID, likedFlavour.FlavourID)

	return likedFlavour
}

func TestLikeFlavour(t *testing.T) {
	user := createRandomUser(t)
	createLikedFlavour(user.ID, t)
}

func TestGetLikedFlavour(t *testing.T) {
	user := createRandomUser(t)
	likedFlavour := createLikedFlavour(user.ID, t)

	get_args := GetLikedFlavourParams{
		UserID:    likedFlavour.UserID,
		FlavourID: likedFlavour.FlavourID,
	}

	selected_flavour, err := testStore.GetFlavourId(context.Background(), likedFlavour.FlavourID)

	require.NoError(t, err)
	require.NotEmpty(t, selected_flavour)

	selected_liked_flavour, err := testStore.GetLikedFlavour(context.Background(), get_args)
	require.NoError(t, err)
	require.NotEmpty(t, selected_liked_flavour)

	require.Equal(t, selected_flavour.ID, selected_liked_flavour.ID)
	require.Equal(t, selected_flavour.Name, selected_liked_flavour.Name)
	require.Equal(t, selected_flavour.CreatedAt, selected_liked_flavour.CreatedAt)
}

func TestGetUserLikedFlavours(t *testing.T) {
	user := createRandomUser(t)

	var flavours []Flavour
	for i := 0; i < 10; i++ {
		liked_flavour := createLikedFlavour(user.ID, t)

		flavour, err := testStore.GetFlavourId(context.Background(), liked_flavour.FlavourID)
		require.NoError(t, err)
		require.NotEmpty(t, flavour)

		flavours = append(flavours, flavour)
	}

	liked_flavours, err := testStore.GetUserLikedFlavours(context.Background(), user.ID)

	require.NoError(t, err)

	require.NotEmpty(t, liked_flavours)

	require.Equal(t, len(flavours), len(liked_flavours))
}

func TestUnikeFlavour(t *testing.T) {
	user := createRandomUser(t)
	likedFlavour := createLikedFlavour(user.ID, t)

	arg := UnlikeFlavourParams{
		UserID:    likedFlavour.UserID,
		FlavourID: likedFlavour.FlavourID,
	}

	err := testStore.UnlikeFlavour(context.Background(), arg)
	require.NoError(t, err)

	get_args := GetLikedFlavourParams{
		UserID:    likedFlavour.UserID,
		FlavourID: likedFlavour.FlavourID,
	}

	selected_flavour, err := testStore.GetLikedFlavour(context.Background(), get_args)

	require.Error(t, err)
	require.Empty(t, selected_flavour)

	require.EqualError(t, err, ErrRecordNotFound.Error())

}
