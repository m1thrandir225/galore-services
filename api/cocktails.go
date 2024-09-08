package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/m1thrandir225/galore-services/db/sqlc"
	"github.com/m1thrandir225/galore-services/dto"
	"io"
	"log"
	"mime/multipart"
	"net/http"
)

type CreateCocktailRequest struct {
	Name         string                `form:"name" binding:"required"`
	IsAlcoholic  bool                  `form:"isAlcoholic" binding:"required"`
	Glass        string                `form:"glass" binding:"required"`
	Image        *multipart.FileHeader `form:"file" binding:"required"`
	Instructions string                `form:"instructions" json:"instructions" binding:"required"`
	Ingredients  string                `form:"ingredients" json:"ingredients" binding:"required"`
}

func (server *Server) createCocktail(ctx *gin.Context) {
	var requestData CreateCocktailRequest

	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var isAlcoholic pgtype.Bool

	err := isAlcoholic.Scan(requestData.IsAlcoholic)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Unmarshal ingredients JSON string to IngredientDto
	var ingredientDto dto.IngredientDto
	if err = json.Unmarshal([]byte(requestData.Ingredients), &ingredientDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ingredients format"})
		return
	}

	var instructionDto dto.InstructionDto
	if err = json.Unmarshal([]byte(requestData.Instructions), &instructionDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instructions format"})
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	openedFile, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	defer openedFile.Close()

	fileData, err := io.ReadAll(openedFile)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	filePath, err := server.storage.UploadFile(fileData, file.Filename)

	log.Println(filePath)

	arg := db.CreateCocktailParams{
		Name:         requestData.Name,
		Image:        filePath,
		Glass:        requestData.Glass,
		Ingredients:  ingredientDto,
		Instructions: instructionDto,
		IsAlcoholic:  isAlcoholic,
	}

	cocktail, err := server.store.CreateCocktail(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, cocktail)
}
