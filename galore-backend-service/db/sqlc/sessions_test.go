package db

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/m1thrandir225/galore-services/util"
	"github.com/stretchr/testify/require"
)

func createRandomSession(userEmail string, t *testing.T) Session {
	arg := CreateSessionParams{
		ID:           uuid.New(),
		Email:        userEmail,
		RefreshToken: util.RandomString(5),
		UserAgent:    util.RandomString(8),
		ClientIp:     util.RandomString(11),
		IsBlocked:    false,
	}

	session, err := testStore.CreateSession(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, session)

	require.Equal(t, arg.Email, session.Email)
	require.Equal(t, arg.ClientIp, session.ClientIp)
	require.Equal(t, arg.UserAgent, session.UserAgent)
	require.Equal(t, arg.IsBlocked, session.IsBlocked)
	require.Equal(t, arg.RefreshToken, session.RefreshToken)

	return session
}

func TestCreateSession(t *testing.T) {
	user := createRandomUser(t)
	createRandomSession(user.Email, t)
}

func TestGetSession(t *testing.T) {
	user := createRandomUser(t)
	session := createRandomSession(user.Email, t)

	selectedSession, err := testStore.GetSession(context.Background(), session.ID)

	require.NoError(t, err)
	require.NotEmpty(t, selectedSession)

	require.Equal(t, session.ID, selectedSession.ID)
	require.Equal(t, session.RefreshToken, selectedSession.RefreshToken)
	require.Equal(t, session.Email, selectedSession.Email)
	require.Equal(t, session.ClientIp, selectedSession.ClientIp)
	require.Equal(t, session.UserAgent, selectedSession.UserAgent)
	require.Equal(t, session.IsBlocked, selectedSession.IsBlocked)
	require.Equal(t, session.CreatedAt, selectedSession.CreatedAt)
}

func TestGetAllUserSessions(t *testing.T) {

}

func TestDeleteSession(t *testing.T) {}
