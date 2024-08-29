package api

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/galore-services/token"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(authorizationHeaderKey)

		if len(authHeader) == 0 {
			err := errors.New("Authorization Header not provided.")
			ctx.AbortWithStatusJSON(401, errorResponse(err))
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			err := errors.New("Invalid header format")
			ctx.AbortWithStatusJSON(401, errorResponse(err))
			return
		}

		token := fields[1]

		payload, err := tokenMaker.VerifyToken(token)

		if err != nil {
			ctx.AbortWithStatusJSON(401, errorResponse(err))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
