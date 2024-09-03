package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func createLikedFlavour(t *testing.T) LikedFlavour {
	user := createRandomUser(t)
	flavour := createRandomFlavour(t)

	arg := LikeFlavourParams{
		UserID:    user.ID,
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
	createLikedFlavour(t)
}

func TestGetLikedFlavour(t *testing.T) {
	likedFlavour := createLikedFlavour(t)

	get_args := GetLikedFlavourParams{
		UserID:    likedFlavour.UserID,
		FlavourID: likedFlavour.FlavourID,
	}

	selected_flavour, err := testStore.GetLikedFlavour(context.Background(), get_args)
	require.NoError(t, err)
	require.NotEmpty(t, selected_flavour)

	require.Equal(t, likedFlavour.UserID, selected_flavour.UserID)
	require.Equal(t, likedFlavour.FlavourID, selected_flavour.FlavourID)

}

func TestDeleteLikeFlavour(t *testing.T) {
	likedFlavour := createLikedFlavour(t)

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
