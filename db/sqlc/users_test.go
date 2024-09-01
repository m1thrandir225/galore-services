package db

import (
	"context"
	"testing"

	"github.com/m1thrandir225/galore-services/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) CreateUserRow {
	hashedPassword, err := util.HashPassowrd("lol1231231231")
	require.NoError(t, err)

	dbDate, err := util.TimeToDbDate("2000-01-01")

	require.NoError(t, err)

	arg := CreateUserParams{
		Email:     "james.brown@gmail.com",
		Name:      "James Brown",
		Birthday:  dbDate,
		AvatarUrl: "random.org",
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

func TestGetUser(t *testing.T) {}

func TestDeleteUser(t *testing.T) {}

//TODO: implement update user queries
