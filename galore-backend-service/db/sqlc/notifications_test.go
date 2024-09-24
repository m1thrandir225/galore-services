package db

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func createRandomNotification(user_id, notification_type_id uuid.UUID, t *testing.T) Notification {
	arg := CreateNotificationParams{
		UserID:             user_id,
		NotificationTypeID: notification_type_id,
	}

	notification, err := testStore.CreateNotification(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, notification)

	require.Equal(t, arg.UserID, notification.UserID)
	require.Equal(t, arg.NotificationTypeID, notification.NotificationTypeID)

	return notification
}

func TestCreateNotification(t *testing.T) {
	user := createRandomUser(t)
	notification_type := createRandomNotificationType(t)
	createRandomNotification(user.ID, notification_type.ID, t)
}

func TestGetUserNotifications(t *testing.T) {
	user := createRandomUser(t)

	var notifications []Notification

	for i := 0; i < 10; i++ {
		notification_type := createRandomNotificationType(t)
		notification := createRandomNotification(user.ID, notification_type.ID, t)

		notifications = append(notifications, notification)
	}

	user_notifications, err := testStore.GetUserNotifications(context.Background(), user.ID)
	require.NoError(t, err)

	require.NotEmpty(t, user_notifications)

	require.Equal(t, len(notifications), len(user_notifications))
}

func TestUpdateUserNotification(t *testing.T) {
	user := createRandomUser(t)
	notification_type := createRandomNotificationType(t)

	notification := createRandomNotification(user.ID, notification_type.ID, t)

	arg := UpdateUserNotificationParams{
		Opened: true,
		ID:     notification.ID,
	}

	updated, err := testStore.UpdateUserNotification(context.Background(), arg)
	require.NoError(t, err)

	require.NotEmpty(t, updated)

	require.Equal(t, notification.ID, updated.ID)
	require.Equal(t, notification.CreatedAt, updated.CreatedAt)
	require.Equal(t, notification.UserID, updated.UserID)
	require.Equal(t, notification.NotificationTypeID, updated.NotificationTypeID)

	require.NotEqual(t, notification.Opened, updated.Opened)

	require.Equal(t, arg.Opened, updated.Opened)
}
