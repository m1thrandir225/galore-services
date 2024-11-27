package api

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/m1thrandir225/galore-services/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateCocktailCategoryRequest struct {
	CocktailId string `json:"cocktail_id" binding:"required"`
	CategoryId string `json:"category_id" binding:"required"`
}

func (server *Server) addCocktailToCategory(ctx *gin.Context) {
	var requestData CreateCocktailCategoryRequest

	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	categoryId, err := uuid.Parse(requestData.CategoryId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	cocktailId, err := uuid.Parse(requestData.CocktailId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	arg := db.CreateCocktailCategoryParams{
		CocktailID: cocktailId,
		CategoryID: categoryId,
	}

	cocktailCategory, err := server.store.CreateCocktailCategory(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, cocktailCategory)
}

func (server *Server) getCocktailCategoryById(ctx *gin.Context) {
	var uriData UriId

	if err := ctx.ShouldBindUri(&uriData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := uuid.Parse(uriData.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	cocktailCategory, err := server.store.GetCocktailCategory(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, cocktailCategory)
}

func (server *Server) deleteCocktailCategory(ctx *gin.Context) {
	var uriData UriId

	if err := ctx.ShouldBindUri(&uriData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	id, err := uuid.Parse(uriData.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.store.DeleteCocktailCategory(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Status(http.StatusOK)
}
