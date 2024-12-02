package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/m1thrandir225/galore-services/db/sqlc"
)

type createCategoryFlavourRequest struct {
	CategoryID string `json:"category_id" binding:"required,uuid"`
	FlavourID  string `json:"flavour_id" binding:"required,uuid"`
}

// create-category-flavour mapping
func (server *Server) CreateCategoryFlavour(ctx *gin.Context) {
	var requestData createCategoryFlavourRequest

	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	categoryId, err := uuid.Parse(requestData.CategoryID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	flavourId, err := uuid.Parse(requestData.FlavourID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateCategoryFlavourParams{
		CategoryID: categoryId,
		FlavourID:  flavourId,
	}

	categoryFlavourMapping, err := server.store.CreateCategoryFlavour(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, categoryFlavourMapping)
}

// get categories for user based on liked flavours
func (server *Server) GetCategoriesBasedOnLikedFlavours(ctx *gin.Context) {
	var uriId UriId
	if err := ctx.ShouldBindUri(&uriId); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userId, err := uuid.Parse(uriId.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	categories, err := server.store.GetCategoriesFromLikedFlavours(ctx, userId)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, categories)
}

// get a single category-flavour mapping (useless)

// delete category-flavour mapping
func (server *Server) DeleteCategoryFlavour(ctx *gin.Context) {
	var uriId UriId

	if err := ctx.ShouldBindUri(&uriId); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := uuid.Parse(uriId.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.store.DeleteCategoryFlavour(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Status(http.StatusOK)
}
