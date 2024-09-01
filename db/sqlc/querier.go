// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateFlavour(ctx context.Context, name string) (Flavour, error)
	CreateNotifcationType(ctx context.Context, arg CreateNotifcationTypeParams) (NotificationType, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (CreateUserRow, error)
	CreateUserCocktail(ctx context.Context, arg CreateUserCocktailParams) (CreatedCocktail, error)
	DeleteFlavour(ctx context.Context, id uuid.UUID) error
	DeleteNotificationType(ctx context.Context, id uuid.UUID) error
	DeleteSession(ctx context.Context, id uuid.UUID) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	DeleteUserCocktail(ctx context.Context, id uuid.UUID) error
	GetAllTypes(ctx context.Context) ([]NotificationType, error)
	GetCreatedCocktails(ctx context.Context, userID uuid.UUID) ([]CreatedCocktail, error)
	GetFlavourId(ctx context.Context, id uuid.UUID) (Flavour, error)
	GetFlavourName(ctx context.Context, name string) (Flavour, error)
	GetNotificationType(ctx context.Context, id uuid.UUID) (NotificationType, error)
	GetSession(ctx context.Context, id uuid.UUID) (Session, error)
	GetUser(ctx context.Context, id uuid.UUID) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserCocktail(ctx context.Context, id uuid.UUID) (CreatedCocktail, error)
	GetUserFCMTokens(ctx context.Context, userID uuid.UUID) ([]FcmToken, error)
	UpdateFlavour(ctx context.Context, arg UpdateFlavourParams) (Flavour, error)
}

var _ Querier = (*Queries)(nil)
