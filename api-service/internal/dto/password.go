package dto

import (
	"time"

	"github.com/google/uuid"
	db "github.com/m1thrandir225/galore-services/internal/db/sqlc"
)

type PasswordResetRequestDTO struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"user_id"`
	PasswordReset bool      `json:"password_reset"`
	ValidUntil    time.Time `json:"valid_until"`
}

func PasswordResetRequestDTOFromDB(passwordResetRequest db.ResetPasswordRequest) PasswordResetRequestDTO {
	return PasswordResetRequestDTO{
		ID:            passwordResetRequest.ID,
		UserID:        passwordResetRequest.UserID,
		PasswordReset: passwordResetRequest.PasswordReset,
		ValidUntil:    passwordResetRequest.ValidUntil,
	}
}
