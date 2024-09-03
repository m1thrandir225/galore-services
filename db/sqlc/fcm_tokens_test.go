package db

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/m1thrandir225/galore-services/util"
	"github.com/stretchr/testify/require"
)

func createRandomFCMToken(user_id uuid.UUID, t *testing.T) FcmToken {
	arg := CreateFCMTokenParams{
		Token:    util.RandomString(12),
		DeviceID: util.RandomString(12),
		UserID:   user_id,
	}

	fcm_token, err := testStore.CreateFCMToken(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, fcm_token)

	require.Equal(t, arg.Token, fcm_token.Token)
	require.Equal(t, arg.UserID, fcm_token.UserID)
	require.Equal(t, arg.DeviceID, fcm_token.DeviceID)

	return fcm_token
}

func TestCreateFcmToken(t *testing.T) {
	user := createRandomUser(t)
	createRandomFCMToken(user.ID, t)
}

func TestGetFcmTokenById(t *testing.T) {
	user := createRandomUser(t)

	token := createRandomFCMToken(user.ID, t)

	selected_token, err := testStore.GetFCMTokenById(context.Background(), token.ID)
	require.NoError(t, err)
	require.NotEmpty(t, selected_token)
}

func TestDeleteFcmToken(t *testing.T) {
	user := createRandomUser(t)

	token := createRandomFCMToken(user.ID, t)

	err := testStore.DeleteFCMToken(context.Background(), token.ID)

	require.NoError(t, err)

	selected_token, err := testStore.GetFCMTokenById(context.Background(), token.ID)

	require.Error(t, err)
	require.Empty(t, selected_token)
	require.EqualError(t, err, ErrRecordNotFound.Error())
}
