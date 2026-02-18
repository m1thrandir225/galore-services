package services

import (
	"context"
	"errors"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"github.com/m1thrandir225/galore-services/internal/cache"
	db "github.com/m1thrandir225/galore-services/internal/db/sqlc"
	"github.com/m1thrandir225/galore-services/internal/security"
	"github.com/m1thrandir225/galore-services/pkg/shared"
)

type TokenType string

const (
	AccessToken  TokenType = "access_token"
	RefreshToken TokenType = "refresh_token"
)

type AuthService interface {
	RegisterUser(
		ctx context.Context,
		name,
		email,
		password,
		birthday string,
		avatarURL *multipart.FileHeader,
	) (*db.User, error)

	LoginUser(ctx context.Context, email, password string) (*db.User, error)
	RefreshToken(ctx context.Context, refreshToken, sessionID string) (string, error)
	ForgotPassword(ctx context.Context, email string) (string, time.Time, error)
	VerifyOTP(ctx context.Context, email, otp string) (string, db.ResetPasswordRequest, error)
	ResetPassword(ctx context.Context, passwordRequestID, newPassword, confirmPassword string) error
	VerifyToken(ctx context.Context, refreshToken string) (*security.Payload, error)

	CreateToken(ctx context.Context, userID uuid.UUID, duration time.Duration, tokenType TokenType) (string, security.Payload, error)
}

type authService struct {
	store                db.Store
	cache                cache.Store
	tokenMaker           security.Maker
	userService          UserService
	sessionService       SessionService
	storage              StorageService
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func NewAuthService(
	store db.Store,
	cache cache.Store,
	tokenMaker security.Maker,
	userService UserService,
	sessionService SessionService,
	storage StorageService,
	accessTokenDuration time.Duration,
	refreshTokenDuration time.Duration,
) *authService {
	return &authService{
		store:                store,
		cache:                cache,
		tokenMaker:           tokenMaker,
		userService:          userService,
		sessionService:       sessionService,
		storage:              storage,
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}
}
func (s *authService) RegisterUser(
	ctx context.Context,
	name,
	email,
	password,
	birthday string,
	avatarURL *multipart.FileHeader,
	clientIP,
	userAgent string,
) (*db.User, error) {
	userId, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	hashedPassword, err := security.HashPassword(password)
	if err != nil {
		return nil, err
	}
	imageData, err := shared.BytesFromFile(avatarURL)
	if err != nil {
		return nil, err
	}

	avatarFilePath, err := s.storage.UploadFile(imageData, userId.String(), avatarURL.Filename)
	if err != nil {
		return nil, err
	}
	otpSecret, err := security.GenerateOTPSecret()
	if err != nil {
		return nil, err
	}

	userBirthday, err := shared.ParseDate(birthday)
	if err != nil {
		return nil, err
	}

	arg := db.CreateUserParams{
		ID:         userId,
		Email:      email,
		Password:   hashedPassword,
		AvatarUrl:  avatarFilePath,
		Birthday:   userBirthday,
		HotpSecret: otpSecret,
	}

	userData, err := s.store.CreateUser(ctx, arg)
	if err != nil {
		return nil, err
	}

	return &userData, err
}

func (s *authService) LoginUser(ctx context.Context, email, password string) (*db.User, error) {
	user, err := s.store.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	err = security.ComparePassword(user.Password, password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string, sessionID uuid.UUID) (string, error) {
	panic("not implemented") // TODO: Implement
}

func (s *authService) ForgotPassword(ctx context.Context, email string) error {
	userData, err := s.userService.GetUserDetailsByEmail(ctx, email)
	if err != nil {
		return err
	}

	currentUserCounter, err := s.userService.GetUserCounter(ctx, userData.ID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {

			err = s.userService.CreateHOTPCounter(ctx, userData.ID)
			if err != nil {
				return ErrorCreatingHOTPCounter
			}
			currentUserCounter = 0
		}
		return ErrorGettingHOTPCounter
	}

	_, err = security.GenerateHOTP(userData.HotpSecret, uint64(currentUserCounter))
	if err != nil {
		return err
	}
	//TODO: send OTP code to mail
	return nil
}

func (s *authService) VerifyOTP(ctx context.Context) {
	panic("not implemented") // TODO: Implement
}

func (s *authService) VerifyOTPResponse(ctx context.Context) {
	panic("not implemented") // TODO: Implement
}

func (s *authService) ResetPassword(ctx context.Context) {
	panic("not implemented") // TODO: Implement
}

func (s *authService) VerifyToken(ctx context.Context, refreshToken string) (*security.Payload, error) {
	return s.tokenMaker.VerifyToken(refreshToken)
}
