package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/m1thrandir225/galore-services/token"
)

func verifyUserIdWithToken(userId uuid.UUID, ctx *gin.Context) error {
	data, exists := ctx.Get(authorizationPayloadKey)
	if !exists {
		return errors.New(authorizationPayloadKey + " was not found in the request context")
	}
	payload := data.(*token.Payload)

	if payload.UserId != userId {
		return errors.New(authorizationPayloadKey + " was not found in the request context")
	}
	return nil
}
