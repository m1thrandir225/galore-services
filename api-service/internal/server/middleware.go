package server

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	db "github.com/m1thrandir225/galore-services/internal/db/sqlc"
	"github.com/m1thrandir225/galore-services/internal/security"
	"github.com/m1thrandir225/galore-services/pkg/shared"
)

func authMiddleware(tokenMaker security.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(security.AuthorizationHeaderKey)

		if len(authHeader) == 0 {
			err := errors.New("authorization Header not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, shared.ErrorResponse(err))
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			err := errors.New("invalid header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, shared.ErrorResponse(err))
			return
		}

		authorizationType := fields[0]
		if strings.ToLower(authorizationType) != security.AuthorizationTypeBearer {
			err := errors.New("invalid authorization type")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, shared.ErrorResponse(err))
			return
		}

		authToken := fields[1]
		payload, err := tokenMaker.VerifyToken(authToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, shared.ErrorResponse(err))
			return
		}

		ctx.Set(security.AuthorizationPayloadKey, payload)
		ctx.Next()
	}
}

func adminMiddleware(store db.Store) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		payload, err := security.GetPayloadFromContext(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, shared.ErrorResponse(err))
			return
		}

		userRole, err := store.GetUserRole(ctx, payload.UserId)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, shared.ErrorResponse(err))
			return
		}
		if userRole != db.UserRolesAdmin {
			ctx.AbortWithStatusJSON(http.StatusForbidden, shared.ErrorResponse(err))
			return
		}
		ctx.Next()

	}
}
