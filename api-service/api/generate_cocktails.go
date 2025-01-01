package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/galore-services/cocktail_gen"
	db "github.com/m1thrandir225/galore-services/db/sqlc"
	"log"
	"net/http"
)

type GenerateCocktailRequest struct {
	ReferenceFlavours  []string `json:"reference_flavours" binding:"required"`
	ReferenceCocktails []string `json:"reference_cocktails" binding:"required"`
}

type GenerateImageRequest struct {
	Prompt string `json:"prompt" binding:"required"`
	Model  string `json:"model" binding:"required"`
}

type CreateGenerateCocktailRequest struct {
	ReferenceFlavours  []string `json:"reference_flavours" binding:"required"`
	ReferenceCocktails []string `json:"reference_cocktails" binding:"required"`
}

func (server *Server) createGenerateCocktailRequest(ctx *gin.Context) {
	var reqData CreateGenerateCocktailRequest
	if err := ctx.ShouldBindJSON(&reqData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	payload, err := getPayloadFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	log.Println(payload.UserId)
	cocktailPrompt := cocktail_gen.GeneratePrompt(reqData.ReferenceFlavours, reqData.ReferenceCocktails)

	arg := db.CreateGenerateCocktailRequestParams{
		UserID: payload.UserId,
		Prompt: cocktailPrompt,
		Status: db.GenerationStatusGeneratingCocktail,
	}

	cocktailRequest, err := server.store.CreateGenerateCocktailRequest(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	server.scheduler.EnqueueJob("generate_cocktail_draft", map[string]interface{}{
		"user_id":             payload.UserId,
		"prompt":              cocktailPrompt,
		"cocktail_request_id": cocktailRequest.ID,
	})

	ctx.JSON(http.StatusOK, cocktailRequest)

}

func (server *Server) generateCocktail(ctx *gin.Context) {
	var req GenerateCocktailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	cocktailPrompt := cocktail_gen.GeneratePrompt(req.ReferenceFlavours, req.ReferenceCocktails)

	prompt, err := server.cocktailGenerator.GenerateRecipe(cocktailPrompt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, prompt)
}

func (server *Server) generateImage(ctx *gin.Context) {
	var req GenerateImageRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	httpClient := &http.Client{}

	generatedImage, err := server.imageGenerator.GenerateImage(req.Prompt, httpClient, req.Model)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	filePath, err := server.storage.UploadFile(generatedImage.Content, "generated-images", fmt.Sprintf("%s%s", generatedImage.FileName, generatedImage.FileExt))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	log.Println(filePath)

	ctx.JSON(http.StatusOK, generatedImage)
}
