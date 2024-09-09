package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/m1thrandir225/galore-services/db/sqlc"
	"github.com/m1thrandir225/galore-services/dto"
	"github.com/m1thrandir225/galore-services/token"
	"github.com/m1thrandir225/galore-services/util"
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
	var isAlcoholic pgtype.Bool
	var ingredientDto dto.IngredientDto
	var instructionDto dto.InstructionDto

	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := isAlcoholic.Scan(requestData.IsAlcoholic)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Unmarshal ingredients and instructions to dto objects
	if err = json.Unmarshal([]byte(requestData.Ingredients), &ingredientDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ingredients format"})
		return
	}
	if err = json.Unmarshal([]byte(requestData.Instructions), &instructionDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instructions format"})
		return
	}

	//Unload the bytes from the uploaded file
	fileData, err := util.BytesFromFile(requestData.Image)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	/*
	* Get the id of the currently logged-in user to use as a name for the folder that the uploaded file will be placed in.
	 */
	data, exists := ctx.Get(authorizationPayloadKey)
	payload := data.(*token.Payload)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	//Upload the file to the public/user_id/file
	filePath, err := server.storage.UploadFile(fileData, payload.ID.String(), requestData.Image.Filename)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
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
