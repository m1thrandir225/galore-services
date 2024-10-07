package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	apiHeaderKey = "x-api-key"
)

func microserviceAccessMiddleware(accessKey string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apiKey := ctx.GetHeader(apiHeaderKey)

		if len(apiKey) == 0 {
			err := errors.New("Missing service key in header")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		if apiKey != accessKey {
			err := errors.New("Invalid api key")
			ctx.AbortWithStatusJSON(http.StatusForbidden, errorResponse(err))
			return
		}
		ctx.Next()
	}

}
