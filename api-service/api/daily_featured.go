package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type GetDailyFeatured struct {
	Timezone string `form:"timezone"`
}

func (server *Server) getDailyFeatured(ctx *gin.Context) {
	var requestData GetDailyFeatured

	if err := ctx.ShouldBindQuery(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	featuredItems, err := server.store.GetDailyFeatured(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, featuredItems)
}
