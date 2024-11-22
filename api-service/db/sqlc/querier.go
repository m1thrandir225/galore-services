// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateCategory(ctx context.Context, arg CreateCategoryParams) (Category, error)
	CreateCocktail(ctx context.Context, arg CreateCocktailParams) (Cocktail, error)
	CreateCocktailCategory(ctx context.Context, arg CreateCocktailCategoryParams) (CocktailCategory, error)
	CreateFCMToken(ctx context.Context, arg CreateFCMTokenParams) (FcmToken, error)
	CreateFlavour(ctx context.Context, name string) (Flavour, error)
	CreateNotification(ctx context.Context, arg CreateNotificationParams) (Notification, error)
	CreateNotificationType(ctx context.Context, arg CreateNotificationTypeParams) (NotificationType, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (CreateUserRow, error)
	CreateUserCocktail(ctx context.Context, arg CreateUserCocktailParams) (CreatedCocktail, error)
	DeleteCategory(ctx context.Context, id uuid.UUID) error
	DeleteCocktail(ctx context.Context, id uuid.UUID) error
	DeleteCocktailCategory(ctx context.Context, id uuid.UUID) error
	DeleteFCMToken(ctx context.Context, id uuid.UUID) error
	DeleteFlavour(ctx context.Context, id uuid.UUID) error
	DeleteNotificationType(ctx context.Context, id uuid.UUID) error
	DeleteSession(ctx context.Context, id uuid.UUID) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	DeleteUserCocktail(ctx context.Context, id uuid.UUID) error
	GetAllCategories(ctx context.Context) ([]Category, error)
	GetAllFlavours(ctx context.Context) ([]Flavour, error)
	GetAllTypes(ctx context.Context) ([]NotificationType, error)
	GetAllUserSessions(ctx context.Context, email string) ([]Session, error)
	GetCategoriesForCocktail(ctx context.Context, cocktailID uuid.UUID) ([]Category, error)
	GetCategoryById(ctx context.Context, id uuid.UUID) (Category, error)
	GetCategoryByTag(ctx context.Context, tag string) (Category, error)
	GetCocktail(ctx context.Context, id uuid.UUID) (Cocktail, error)
	GetCocktailAndSimilar(ctx context.Context, id uuid.UUID) ([]GetCocktailAndSimilarRow, error)
	GetCocktailCategory(ctx context.Context, id uuid.UUID) (CocktailCategory, error)
	GetCocktailsForCategory(ctx context.Context, categoryID uuid.UUID) ([]GetCocktailsForCategoryRow, error)
	GetFCMTokenById(ctx context.Context, id uuid.UUID) (FcmToken, error)
	GetFlavourId(ctx context.Context, id uuid.UUID) (Flavour, error)
	GetFlavourName(ctx context.Context, name string) (Flavour, error)
	GetLikedCocktail(ctx context.Context, arg GetLikedCocktailParams) (GetLikedCocktailRow, error)
	GetLikedCocktails(ctx context.Context, userID uuid.UUID) ([]GetLikedCocktailsRow, error)
	GetLikedFlavour(ctx context.Context, arg GetLikedFlavourParams) (Flavour, error)
	GetNotificationType(ctx context.Context, id uuid.UUID) (NotificationType, error)
	GetSession(ctx context.Context, id uuid.UUID) (Session, error)
	GetUser(ctx context.Context, id uuid.UUID) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserCocktail(ctx context.Context, id uuid.UUID) (CreatedCocktail, error)
	GetUserLikedFlavours(ctx context.Context, userID uuid.UUID) ([]Flavour, error)
	GetUserNotifications(ctx context.Context, userID uuid.UUID) ([]Notification, error)
	InvalidateSession(ctx context.Context, id uuid.UUID) (Session, error)
	LikeCocktail(ctx context.Context, arg LikeCocktailParams) (LikedCocktail, error)
	LikeFlavour(ctx context.Context, arg LikeFlavourParams) (LikedFlavour, error)
	SearchCocktails(ctx context.Context, dollar_1 string) ([]Cocktail, error)
	UnlikeCocktail(ctx context.Context, arg UnlikeCocktailParams) error
	UnlikeFlavour(ctx context.Context, arg UnlikeFlavourParams) error
	UpdateCategory(ctx context.Context, arg UpdateCategoryParams) (Category, error)
	UpdateCocktail(ctx context.Context, arg UpdateCocktailParams) (Cocktail, error)
	UpdateFlavour(ctx context.Context, arg UpdateFlavourParams) (Flavour, error)
	UpdateNotificationType(ctx context.Context, arg UpdateNotificationTypeParams) (NotificationType, error)
	UpdateUserEmailNotifications(ctx context.Context, arg UpdateUserEmailNotificationsParams) (bool, error)
	UpdateUserInformation(ctx context.Context, arg UpdateUserInformationParams) (UpdateUserInformationRow, error)
	UpdateUserNotification(ctx context.Context, arg UpdateUserNotificationParams) (Notification, error)
	UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) error
	UpdateUserPushNotifications(ctx context.Context, arg UpdateUserPushNotificationsParams) (bool, error)
}

var _ Querier = (*Queries)(nil)
