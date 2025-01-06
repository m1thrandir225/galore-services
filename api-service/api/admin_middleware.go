package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/m1thrandir225/galore-services/db/sqlc"
	"net/http"
)

func adminMiddleware(store db.Store) gin.HandlerFunc {
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
		if userRole != db.UserRolesAdmin {
			ctx.AbortWithStatusJSON(http.StatusForbidden, errorResponse(err))
			return
		}
		ctx.Next()

	}
}
