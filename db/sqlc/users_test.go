package db

import (
	"context"
	"fmt"
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
	updated_value, err := testStore.UpdateUserEmailNotifications(context.Background(), arg)

	require.NoError(t, err)

	require.NotEqual(t, user.EnabledEmailNotifications, updated_value)
	require.NotEmpty(t, updated_value)

	updated_user, err := testStore.GetUser(context.Background(), user.ID)

	require.NoError(t, err)
	require.NotEmpty(t, updated_user)
	require.Equal(t, arg.EnabledEmailNotifications, updated_user.EnabledEmailNotifications)
}

func TestUpdatePushNotifications(t *testing.T) {

	user := createRandomUser(t)

	arg := UpdateUserPushNotificationsParams{
		ID:                       user.ID,
		EnabledPushNotifications: !user.EnabledPushNotifications,
	}
	updated_value, err := testStore.UpdateUserPushNotifications(context.Background(), arg)

	require.NoError(t, err)

	require.NotEqual(t, user.EnabledPushNotifications, updated_value)
	require.NotEmpty(t, updated_value)

	updated_user, err := testStore.GetUser(context.Background(), user.ID)

	require.NoError(t, err)
	require.NotEmpty(t, updated_user)
	require.Equal(t, arg.EnabledPushNotifications, updated_user.EnabledPushNotifications)

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

	updated_user, err := testStore.UpdateUserInformation(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updated_user)

	require.Equal(t, user.ID, updated_user.ID)
	require.Equal(t, user.CreatedAt, updated_user.CreatedAt)
	require.Equal(t, user.EnabledEmailNotifications, updated_user.EnabledEmailNotifications)
	require.Equal(t, user.EnabledPushNotifications, updated_user.EnabledPushNotifications)
	require.NotEqual(t, user.Email, updated_user.Email)
	require.NotEqual(t, user.AvatarUrl, updated_user.AvatarUrl)
	require.NotEqual(t, user.Name, updated_user.Name)
	require.NotEqual(t, user.Birthday, updated_user.Birthday)

	require.Equal(t, arg.Name, updated_user.Name)
	require.Equal(t, arg.Email, updated_user.Email)
	require.Equal(t, arg.AvatarUrl, updated_user.AvatarUrl)
	require.Equal(t, arg.Birthday, updated_user.Birthday)
}

func TestUpdateUserPassword(t *testing.T) {
	user := createRandomUser(t)

	new_password, err := util.HashPassowrd(util.RandomString(32))
	require.NoError(t, err)

	arg := UpdateUserPasswordParams{
		ID:       user.ID,
		Password: new_password,
	}

	err = testStore.UpdateUserPassword(context.Background(), arg)

	// The new password should not be returned as the user doesn't need to know, if there is a problem
	// Then we just check for that
	require.NoError(t, err)

}
