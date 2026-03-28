package api

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db2 "github.com/m1thrandir225/galore-services/internal/db/sqlc"
	"github.com/m1thrandir225/galore-services/internal/recipe"
)

type CreateGenerateCocktailRequest struct {
	ReferenceFlavours  []string `json:"reference_flavours" binding:"required"`
	ReferenceCocktails []string `json:"reference_cocktails" binding:"required"`
}

func (server *Server) getUserGeneratedCocktails(ctx *gin.Context) {
	var uriId UriId

	if err := ctx.ShouldBindUri(&uriId); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userId, err := uuid.Parse(uriId.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	generatedCocktails, err := server.store.GetUserGeneratedCocktails(ctx, userId)
	if err != nil {
		if errors.Is(err, db2.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, generatedCocktails)
}

func (server *Server) getGeneratedCocktail(ctx *gin.Context) {
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

	cocktail, err := server.store.GetGeneratedCocktail(ctx, cocktailId)
	if err != nil {
		if errors.Is(err, db2.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, cocktail)
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
	cocktailPrompt := recipe.GeneratePrompt(reqData.ReferenceFlavours, reqData.ReferenceCocktails)

	arg := db2.CreateGenerateCocktailRequestParams{
		UserID: payload.UserId,
		Prompt: cocktailPrompt,
		Status: db2.GenerationStatusGeneratingCocktail,
	}

	cocktailRequest, err := server.store.CreateGenerateCocktailRequest(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_ = server.scheduler.EnqueueJob("generate_cocktail_draft", map[string]interface{}{
		"user_id":             payload.UserId,
		"prompt":              cocktailPrompt,
		"cocktail_request_id": cocktailRequest.ID,
	})

	ctx.JSON(http.StatusOK, cocktailRequest)

}
