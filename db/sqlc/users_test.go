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

//TODO: implement update user queries
