package db

import (
	"context"
	"testing"

	"github.com/m1thrandir225/galore-services/util"
	"github.com/stretchr/testify/require"
)

func createRandomNotificationType(t *testing.T) NotificationType {
	arg := CreateNotifcationTypeParams{
		Title:   util.RandomString(10),
		Content: util.RandomString(64),
		Tag:     util.RandomString(3),
	}

	notification_type, err := testStore.CreateNotifcationType(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, notification_type)

	require.Equal(t, arg.Title, notification_type.Title)
	require.Equal(t, arg.Content, notification_type.Content)
	require.Equal(t, arg.Tag, notification_type.Tag)

	return notification_type
}

func TestCreateNotificationType(t *testing.T) {
	createRandomNotificationType(t)
}

func TestGetAllTypes(t *testing.T) {
	createRandomNotificationType(t)

	all_types, err := testStore.GetAllTypes(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, all_types)
	require.Equal(t, len(all_types) > 0, true)
}

func TestGetNotificationType(t *testing.T) {
	notification_type := createRandomNotificationType(t)
	selected_type, err := testStore.GetNotificationType(context.Background(), notification_type.ID)
	require.NoError(t, err)
	require.NotEmpty(t, selected_type)

	require.Equal(t, notification_type.ID, selected_type.ID)
	require.Equal(t, notification_type.Tag, selected_type.Tag)
	require.Equal(t, notification_type.Content, selected_type.Content)
	require.Equal(t, notification_type.Title, selected_type.Title)
	require.Equal(t, notification_type.CreatedAt, selected_type.CreatedAt)
}

func TestDeleteNotificationType(t *testing.T) {
	notification_type := createRandomNotificationType(t)

	err := testStore.DeleteNotificationType(context.Background(), notification_type.ID)

	require.NoError(t, err)

	selected_type, err := testStore.GetNotificationType(context.Background(), notification_type.ID)

	require.Error(t, err)
	require.Empty(t, selected_type)
	require.EqualError(t, err, ErrRecordNotFound.Error())
}
