package auth

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"github.com/m1thrandir225/galore-services/internal/dto"
)

type registerUserRequest struct {
	Name      string                `form:"name" json:"name" binding:"required"`
	Email     string                `form:"email" json:"email" binding:"required,email"`
	AvatarUrl *multipart.FileHeader `form:"avatar_url" json:"avatar_url" binding:"required"`
	Password  string                `form:"password" json:"password" binding:"required"`
	Birthday  string                `form:"birthday" json:"birthday" binding:"required"`
}

type loginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// authenticatedUserResponse is the same response for register and login
type authenticatedUserResponse struct {
	User                  dto.UserDTO `json:"user"`
	AccessToken           string      `json:"access_token"`
	AccessTokenExpiresAt  time.Time   `json:"access_token_expires_at"`
	RefreshToken          string      `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time   `json:"refresh_token_expires_at"`
	SessionID             uuid.UUID   `json:"session_id"`
}

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
	SessionID    string `json:"session_id" binding:"required,uuid"`
}

type refreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type logoutRequest struct {
	SessionID uuid.UUID `json:"session_id" binding:"required"`
}

type forgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type forgotPasswordResponse struct {
	OTP        string    `json:"otp"`
	ValidUntil time.Time `json:"valid_until"`
}

type verifyOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required"`
}

type verifyOTPResponse struct {
	ResetPasswordRequest dto.PasswordResetRequestDTO `json:"reset_password_request"`
	Email                string                      `json:"email"`
}

type resetPasswordRequest struct {
	PasswordRequestId string `json:"password_request_id" binding:"required"`
	NewPassword       string `json:"new_password" binding:"required"`
	ConfirmPassword   string `json:"confirm_password" binding:"required"`
}
