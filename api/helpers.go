package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/m1thrandir225/galore-services/token"
)

type UriId struct {
	ID string `uri:"id" binding:"required,uuid"`
}

func getPayloadFromContext(c *gin.Context) (*token.Payload, error) {
	data, exists := c.Get(authorizationPayloadKey)
	if !exists {
		return nil, errors.New("authorization header not found")
	}
	payload := data.(*token.Payload)
	return payload, nil
}

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
