package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/m1thrandir225/galore-services/db/sqlc"
	"github.com/m1thrandir225/galore-services/dto"
	"mime/multipart"
	"net/http"
)

type CreateCocktailRequest struct {
	Name         string                `form:"name" binding:"required"`
	IsAlcoholic  bool                  `form:"isAlcoholic" binding:"required"`
	Glass        string                `form:"glass" binding:"required"`
	Image        *multipart.FileHeader `form:"file" binding:"required"`
	Instructions string                `form:"instructions" binding:"required"`
	Ingredients  string                `form:"ingredients" binding:"required"`
}

func (server *Server) createCocktail(ctx *gin.Context) {
	var requestData CreateCocktailRequest

	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fileBytes, err := extractFileBytes(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	fileHeader, err := ctx.FormFile("file")

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	fileUploadedPath, err := server.storage.UploadFile(fileBytes, fileHeader.Filename)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var isAlcoholic pgtype.Bool
	var ingredients dto.IngredientDto

	var instructions dto.InstructionDto

	err = isAlcoholic.Scan(requestData.IsAlcoholic)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = json.Unmarshal([]byte(requestData.Ingredients), &ingredients)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	err = json.Unmarshal([]byte(requestData.Instructions), &instructions)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateCocktailParams{
		Name:         requestData.Name,
		Image:        fileUploadedPath,
		Glass:        requestData.Glass,
		Ingredients:  ingredients,
		Instructions: instructions,
		IsAlcoholic:  isAlcoholic,
	}

	cocktail, err := server.store.CreateCocktail(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, cocktail)
}
