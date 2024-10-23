package db

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"testing"

	"github.com/m1thrandir225/galore-services/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) CreateUserRow {
	hashedPassword, err := util.HashPassowrd("lol1231231231")
	require.NoError(t, err)

	dbDate, err := util.TimeToDbDate(util.RandomDate())

	require.NoError(t, err)

	name := util.RandomString(12)

	arg := CreateUserParams{
		ID:        uuid.New(),
		Email:     util.RandomEmail(),
		Name:      name,
		Birthday:  dbDate,
		AvatarUrl: fmt.Sprintf("https://api.dicebear.com/9.x/pixel-art/svg?seed=%s", name),
		Password:  hashedPassword,
	}

	user, err := testStore.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, user.AvatarUrl, arg.AvatarUrl)
	require.Equal(t, user.Name, arg.Name)
	require.Equal(t, user.Email, arg.Email)
	require.Equal(t, hashedPassword, arg.Password)
	require.Equal(t, user.Birthday, arg.Birthday)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user := createRandomUser(t)

	selectedUser, err := testStore.GetUser(context.Background(), user.ID)

	require.NoError(t, err)

	require.NotEmpty(t, selectedUser)
	require.Equal(t, user.ID, selectedUser.ID)
	require.Equal(t, user.Name, selectedUser.Name)
	require.Equal(t, user.Email, selectedUser.Email)
	require.Equal(t, user.AvatarUrl, selectedUser.AvatarUrl)
	require.Equal(t, user.Birthday, selectedUser.Birthday)

}

func TestGetUserByEmail(t *testing.T) {
	user := createRandomUser(t)

	selectedUser, err := testStore.GetUserByEmail(context.Background(), user.Email)

	require.NoError(t, err)
	require.NotEmpty(t, selectedUser)
	require.Equal(t, user.ID, selectedUser.ID)
	require.Equal(t, user.Name, selectedUser.Name)
	require.Equal(t, user.Email, selectedUser.Email)
	require.Equal(t, user.AvatarUrl, selectedUser.AvatarUrl)
	require.Equal(t, user.Birthday, selectedUser.Birthday)
}

func TestDeleteUser(t *testing.T) {
	user := createRandomUser(t)

	err := testStore.DeleteUser(context.Background(), user.ID)
	require.NoError(t, err)

	selectedUser, err := testStore.GetUser(context.Background(), user.ID)
	require.Error(t, err)
	require.EqualError(t, err, ErrRecordNotFound.Error())
	require.Empty(t, selectedUser)
}

// TODO: implement update user queries
func TestUpdateEmailNotifications(t *testing.T) {
	user := createRandomUser(t)

	arg := UpdateUserEmailNotificationsParams{
		ID:                        user.ID,
		EnabledEmailNotifications: !user.EnabledEmailNotifications,
	}
	updatedValue, err := testStore.UpdateUserEmailNotifications(context.Background(), arg)

	require.NoError(t, err)

	require.NotEqual(t, user.EnabledEmailNotifications, updatedValue)
	require.NotEmpty(t, updatedValue)

	updatedUser, err := testStore.GetUser(context.Background(), user.ID)

	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.Equal(t, arg.EnabledEmailNotifications, updatedUser.EnabledEmailNotifications)
}

func TestUpdatePushNotifications(t *testing.T) {

	user := createRandomUser(t)

	arg := UpdateUserPushNotificationsParams{
		ID:                       user.ID,
		EnabledPushNotifications: !user.EnabledPushNotifications,
	}
	updatedValue, err := testStore.UpdateUserPushNotifications(context.Background(), arg)

	require.NoError(t, err)

	require.NotEqual(t, user.EnabledPushNotifications, updatedValue)
	require.NotEmpty(t, updatedValue)

	updatedUser, err := testStore.GetUser(context.Background(), user.ID)

	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.Equal(t, arg.EnabledPushNotifications, updatedUser.EnabledPushNotifications)

}

func TestUpdateUserInformation(t *testing.T) {
	user := createRandomUser(t)
	newBirthday, err := util.TimeToDbDate(util.RandomDate())

	arg := UpdateUserInformationParams{
		ID:        user.ID,
		Email:     util.RandomEmail(),
		AvatarUrl: util.RandomString(46),
		Name:      util.RandomString(36),
		Birthday:  newBirthday,
	}

	updatedUser, err := testStore.UpdateUserInformation(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)

	require.Equal(t, user.ID, updatedUser.ID)
	require.Equal(t, user.CreatedAt, updatedUser.CreatedAt)
	require.Equal(t, user.EnabledEmailNotifications, updatedUser.EnabledEmailNotifications)
	require.Equal(t, user.EnabledPushNotifications, updatedUser.EnabledPushNotifications)
	require.NotEqual(t, user.Email, updatedUser.Email)
	require.NotEqual(t, user.AvatarUrl, updatedUser.AvatarUrl)
	require.NotEqual(t, user.Name, updatedUser.Name)
	require.NotEqual(t, user.Birthday, updatedUser.Birthday)

	require.Equal(t, arg.Name, updatedUser.Name)
	require.Equal(t, arg.Email, updatedUser.Email)
	require.Equal(t, arg.AvatarUrl, updatedUser.AvatarUrl)
	require.Equal(t, arg.Birthday, updatedUser.Birthday)
}

func TestUpdateUserPassword(t *testing.T) {
	user := createRandomUser(t)

	newPassword, err := util.HashPassowrd(util.RandomString(32))
	require.NoError(t, err)

	arg := UpdateUserPasswordParams{
		ID:       user.ID,
		Password: newPassword,
	}

	err = testStore.UpdateUserPassword(context.Background(), arg)

	// The new password should not be returned as the user doesn't need to know, if there is a problem
	// Then we just check for that
	require.NoError(t, err)

}
