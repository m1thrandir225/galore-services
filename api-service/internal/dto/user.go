// Package DTO defines Data Transfer Objects returned by the API.
package dto

import (
	"time"

	"github.com/google/uuid"
	db "github.com/m1thrandir225/galore-services/internal/db/sqlc"
)

type UserDTO struct {
	ID                        uuid.UUID `json:"id"`
	Name                      string    `json:"name"`
	Email                     string    `json:"email"`
	AvatarURL                 string    `json:"avatar_url"`
	Role                      string    `json:"role"`
	Birthday                  string    `json:"birthday"`
	EnabledPushNotifications  bool      `json:"enabled_push_notifications"`
	EnabledEmailNotifications bool      `json:"enabled_email_notifications"`
}

func UserDTOFromDB(user db.User) UserDTO {
	birthday := user.Birthday.Format(time.RFC3339)
	return UserDTO{
		ID:                        user.ID,
		Name:                      user.Name,
		Email:                     user.Email,
		AvatarURL:                 user.AvatarUrl,
		Role:                      string(user.Role),
		Birthday:                  birthday,
		EnabledPushNotifications:  user.EnabledPushNotifications,
		EnabledEmailNotifications: user.EnabledEmailNotifications,
	}
}
