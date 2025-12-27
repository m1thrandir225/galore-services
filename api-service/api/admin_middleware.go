package api

import (
	"github.com/gin-gonic/gin"
	db2 "github.com/m1thrandir225/galore-services/internal/db/sqlc"

	"net/http"
)

func adminMiddleware(store db2.Store) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		payload, err := getPayloadFromContext(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		userRole, err := store.GetUserRole(ctx, payload.UserId)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		if userRole != db2.UserRolesAdmin {
			ctx.AbortWithStatusJSON(http.StatusForbidden, errorResponse(err))
			return
		}
		ctx.Next()

	}
}
