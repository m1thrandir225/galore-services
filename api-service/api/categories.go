package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/m1thrandir225/galore-services/db/sqlc"
	"net/http"
)

type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
	Tag  string `json:"tag" binding:"required"`
}
type GetCategoryByTagRequest struct {
	Name string `uri:"name" binding:"required"`
}

type UpdateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
	Tag  string `json:"tag" binding:"required"`
}

func (server *Server) getAllCategories(ctx *gin.Context) {
	categories, err := server.store.GetAllCategories(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, categories)
}

func (server *Server) createCategory(ctx *gin.Context) {
	var requestData CreateCategoryRequest

	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateCategoryParams{
		Name: requestData.Name,
		Tag:  requestData.Tag,
	}

	category, err := server.store.CreateCategory(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, category)
}

// TODO: deprecate
func (server *Server) getCategoryByTag(ctx *gin.Context) {
	var uriData GetCategoryByTagRequest

	if err := ctx.ShouldBindUri(&uriData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	category, err := server.store.GetCategoryByTag(ctx, uriData.Name)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, category)
}

func (server *Server) getCategoryById(ctx *gin.Context) {
	var uriData UriId

	if err := ctx.ShouldBindUri(&uriData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	categoryId, err := uuid.Parse(uriData.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	category, err := server.store.GetCategoryById(ctx, categoryId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, category)
}

func (server *Server) updateCategory(ctx *gin.Context) {
	var uriData UriId
	var requestData UpdateCategoryRequest

	if err := ctx.ShouldBindUri(&uriData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	categoryId, err := uuid.Parse(uriData.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.UpdateCategoryParams{
		ID:   categoryId,
		Name: requestData.Name,
		Tag:  requestData.Tag,
	}

	updated, err := server.store.UpdateCategory(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, updated)

}

func (server *Server) deleteCategory(ctx *gin.Context) {
	var uriData UriId

	if err := ctx.ShouldBindUri(&uriData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	categoryId, err := uuid.Parse(uriData.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.store.DeleteCategory(ctx, categoryId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Status(http.StatusOK)
}