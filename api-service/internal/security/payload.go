package security

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Payload represents the data that is stored in a given token
type Payload struct {
	ID        uuid.UUID `json:"id"`
	UserId    uuid.UUID `json:"user_id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// NewPayload returns a payload
func NewPayload(userId uuid.UUID, duration time.Duration) (*Payload, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenId,
		UserId:    userId,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

// Valid checks if the current payload is expired
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

func GetPayloadFromContext(ctx *gin.Context) (*Payload, error) {
	data, exists := ctx.Get(AuthorizationPayloadKey)
	if !exists {
		return nil, ErrPayloadNotFound
	}
	payload := data.(*Payload)
	return payload, nil
}

func VerifyUserIDWithToken(userId uuid.UUID, ctx *gin.Context) error {
	payload, err := GetPayloadFromContext(ctx)
	if err != nil {
		return err
	}
	if payload.UserId != userId {
		return ErrUserIDMismatch
	}
	return nil
}
