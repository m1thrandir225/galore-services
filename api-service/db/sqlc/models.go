// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	dto "github.com/m1thrandir225/galore-services/dto"
	"github.com/pgvector/pgvector-go"
)

type Category struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Tag       string    `json:"tag"`
	CreatedAt time.Time `json:"created_at"`
}

type Cocktail struct {
	ID           uuid.UUID          `json:"id"`
	Name         string             `json:"name"`
	IsAlcoholic  pgtype.Bool        `json:"is_alcoholic"`
	Glass        string             `json:"glass"`
	Image        string             `json:"image"`
	Instructions dto.InstructionDto `json:"instructions"`
	Ingredients  dto.IngredientDto  `json:"ingredients"`
	Embedding    pgvector.Vector    `json:"embedding"`
	CreatedAt    time.Time          `json:"created_at"`
}

type CocktailCategory struct {
	ID         uuid.UUID `json:"id"`
	CocktailID uuid.UUID `json:"cocktail_id"`
	CategoryID uuid.UUID `json:"category_id"`
}

type CreatedCocktail struct {
	ID           uuid.UUID            `json:"id"`
	Name         string               `json:"name"`
	Image        string               `json:"image"`
	Ingredients  dto.IngredientDto    `json:"ingredients"`
	Instructions dto.AiInstructionDto `json:"instructions"`
	Description  string               `json:"description"`
	UserID       uuid.UUID            `json:"user_id"`
	Embedding    pgvector.Vector      `json:"embedding"`
	CreatedAt    time.Time            `json:"created_at"`
}

type FcmToken struct {
	ID        uuid.UUID `json:"id"`
	Token     string    `json:"token"`
	DeviceID  string    `json:"device_id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Flavour struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type LikedCocktail struct {
	ID         uuid.UUID `json:"id"`
	CocktailID uuid.UUID `json:"cocktail_id"`
	UserID     uuid.UUID `json:"user_id"`
}

type LikedFlavour struct {
	FlavourID uuid.UUID `json:"flavour_id"`
	UserID    uuid.UUID `json:"user_id"`
}

type Notification struct {
	ID                 uuid.UUID `json:"id"`
	UserID             uuid.UUID `json:"user_id"`
	NotificationTypeID uuid.UUID `json:"notification_type_id"`
	Opened             bool      `json:"opened"`
	CreatedAt          time.Time `json:"created_at"`
}

type NotificationType struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tag       string    `json:"tag"`
	CreatedAt time.Time `json:"created_at"`
}

type Session struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type User struct {
	ID                        uuid.UUID   `json:"id"`
	Email                     string      `json:"email"`
	Name                      string      `json:"name"`
	Password                  string      `json:"password"`
	AvatarUrl                 string      `json:"avatar_url"`
	EnabledPushNotifications  bool        `json:"enabled_push_notifications"`
	EnabledEmailNotifications bool        `json:"enabled_email_notifications"`
	CreatedAt                 time.Time   `json:"created_at"`
	Birthday                  pgtype.Date `json:"birthday"`
}
