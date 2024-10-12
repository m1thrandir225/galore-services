package api

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	mockdb "github.com/m1thrandir225/galore-services/db/mock"
	db "github.com/m1thrandir225/galore-services/db/sqlc"
	"github.com/m1thrandir225/galore-services/token"
	"github.com/m1thrandir225/galore-services/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetUserDetailsApi(t *testing.T) {
	user := randomUser(t)

	testCases := []struct {
		name          string
		userId        string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			userId: user.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}
}

func TestDeleteUserApi(t *testing.T) {}

func TestUpdateUserDetailsApi(t *testing.T) {}

func TestUpdateUserPasswordApi(t *testing.T) {}

func TestUpdateUserPushNotificationsApi(t *testing.T) {}

func TestUpdateUserEmailNotificationsApi(t *testing.T) {}

func randomUser(t *testing.T) db.User {
	var date pgtype.Date

	err := date.Scan(util.RandomDate())
	require.NoError(t, err)

	id := uuid.New()
	return db.User{
		ID:                        id,
		Name:                      util.RandomString(10),
		Email:                     util.RandomEmail(),
		Password:                  util.RandomString(10),
		AvatarUrl:                 util.RandomString(10),
		EnabledEmailNotifications: util.RandomBool(),
		EnabledPushNotifications:  util.RandomBool(),
		Birthday:                  date,
		CreatedAt:                 time.Now(),
	}
}
