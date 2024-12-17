package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/m1thrandir225/galore-services/db/sqlc"
)

type CreateDailyFeaturedItemRequest struct {
	CocktailId string `json:"cocktail_id" binding:"required,uuid"`
}

type CheckWasCocktailFeaturedRequest struct {
	CocktailId string `json:"cocktail_id" binding:"required,uuid"`
}

type GetDailyFeatured struct {
	Timezone string `form:"timezone" binding:"required"`
}

func (server *Server) setupDailyFeaturedRoutes(cocktailRoutes *gin.RouterGroup) {
	cocktailRoutes.GET("/featured", server.getDailyFeatured)
	cocktailRoutes.POST("/featured", server.createDailyFeatured)
	cocktailRoutes.DELETE("/featured", server.deleteOlderFeatured)
}

func (server *Server) createDailyFeatured(ctx *gin.Context) {
	var requestData CreateDailyFeaturedItemRequest

	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	cocktailId, err := uuid.Parse(requestData.CocktailId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	featuredItem, err := server.store.CreateDailyFeatured(ctx, cocktailId)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, featuredItem)
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

func (server *Server) checkWasCocktailFeatured(ctx *gin.Context) {
	/*
	* Check if the current cocktail was featured recently (7 days)
	 */
	var uriId UriId
	if err := ctx.ShouldBindUri(&uriId); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	cocktailId, err := uuid.Parse(uriId.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	_, err = server.store.CheckWasCocktailFeatured(ctx, cocktailId)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusOK, gin.H{
				"was_featured": false,
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, gin.H{
		"was_featrued": true,
	})
}

func (server *Server) deleteOlderFeatured(ctx *gin.Context) {
	err := server.store.DeleteOlderFeatured(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Status(http.StatusOK)
}
