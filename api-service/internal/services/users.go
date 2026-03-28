package services

import (
	"context"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/m1thrandir225/galore-services/internal/cache"
	db "github.com/m1thrandir225/galore-services/internal/db/sqlc"
	"github.com/m1thrandir225/galore-services/internal/security"
	"github.com/m1thrandir225/galore-services/pkg/shared"
)

type UserService interface {
	GetUserDetails(ctx context.Context, userId uuid.UUID) (db.User, error)
	GetUserDetailsByEmail(ctx context.Context, email string) (db.User, error)
	GetUserCounter(ctx context.Context, userId uuid.UUID) (int32, error)
	CreateHOTPCounter(ctx context.Context, userId uuid.UUID) error
	DeleteUser(ctx context.Context, userId uuid.UUID) error
	UpdateUserPassword(ctx context.Context, userId uuid.UUID, newPassword string) error
	UpdateUserInformation(ctx context.Context, userId uuid.UUID, email string, avatar *multipart.FileHeader, name string, birthday string) (db.UpdateUserInformationRow, error)
	UpdateUserPushNotificationsSettings(ctx context.Context, userId uuid.UUID, enabled bool) (bool, error)
	UpdateUserEmailNotificationsSettings(ctx context.Context, userId uuid.UUID, enabled bool) (bool, error)
}

type userService struct {
	store   db.Store
	cache   cache.Store
	storage StorageService
}

func NewUserService(
	store db.Store,
	cache cache.Store,
	storage StorageService,
) *userService {
	return &userService{
		store:   store,
		cache:   cache,
		storage: storage,
	}
}

func (s *userService) GetUserDetails(ctx context.Context, userId uuid.UUID) (db.User, error) {
	details, err := s.store.GetUser(ctx, userId)
	if err != nil {
		return db.User{}, err
	}

	return details, nil
}

func (s *userService) DeleteUser(ctx context.Context, userId uuid.UUID) error {
	err := s.store.DeleteUser(ctx, userId)
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) UpdateUserPassword(ctx context.Context, userId uuid.UUID, newPassword string) error {
	hashedPassword, err := security.HashPassword(newPassword)
	if err != nil {
		return err
	}

	args := db.UpdateUserPasswordParams{
		ID:       userId,
		Password: hashedPassword,
	}

	return s.store.UpdateUserPassword(ctx, args)
}

func (s *userService) UpdateUserInformation(
	ctx context.Context,
	userId uuid.UUID,
	email string,
	avatar *multipart.FileHeader,
	name string,
	birthday string,
) (db.UpdateUserInformationRow, error) {
	avatarURL := ""

	userInformation, err := s.GetUserDetails(ctx, userId)
	if err != nil {
		return db.UpdateUserInformationRow{}, err
	}

	newBirthday, err := shared.TimeToDbDate(birthday)
	if err != nil {
		return db.UpdateUserInformationRow{}, err
	}

	if avatar == nil {
		avatarURL = userInformation.AvatarUrl
	} else {
		newAvatarData, err := shared.BytesFromFile(avatar)
		if err != nil {
			return db.UpdateUserInformationRow{}, err
		}

		newFilePath, err := s.storage.ReplaceFile(userInformation.AvatarUrl, newAvatarData)
		if err != nil {
			return db.UpdateUserInformationRow{}, err
		}
		avatarURL = newFilePath
	}

	args := db.UpdateUserInformationParams{
		ID:        userId,
		Email:     email,
		AvatarUrl: avatarURL,
		Name:      name,
		Birthday:  newBirthday,
	}

	return s.store.UpdateUserInformation(ctx, args)
}

func (s *userService) UpdateUserPushNotificationsSettings(
	ctx context.Context,
	userId uuid.UUID,
	enabled bool,
) (bool, error) {
	args := db.UpdateUserPushNotificationsParams{
		ID:                       userId,
		EnabledPushNotifications: enabled,
	}
	return s.store.UpdateUserPushNotifications(ctx, args)
}

func (s *userService) UpdateUserEmailNotificationsSettings(
	ctx context.Context,
	userId uuid.UUID,
	enabled bool,
) (bool, error) {
	args := db.UpdateUserEmailNotificationsParams{
		ID:                        userId,
		EnabledEmailNotifications: enabled,
	}
	return s.store.UpdateUserEmailNotifications(ctx, args)
}

func (s *userService) GetUserDetailsByEmail(ctx context.Context, email string) (db.User, error) {
	return s.store.GetUserByEmail(ctx, email)
}

func (s *userService) GetUserCounter(ctx context.Context, userId uuid.UUID) (int32, error) {
	return s.store.GetCurrentCounter(ctx, userId)
}

func (s *userService) CreateHOTPCounter(ctx context.Context, userId uuid.UUID) error {
	params := db.CreateHotpCounterParams{
		UserID:  userId,
		Counter: 0,
	}
	return s.store.CreateHotpCounter(ctx, params)
}
