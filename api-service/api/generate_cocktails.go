package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type GenerateCocktailRequest struct {
	ReferenceFlavours  []string `json:"reference_flavours" binding:"required"`
	ReferenceCocktails []string `json:"reference_cocktails" binding:"required"`
}

func (server *Server) generateCocktail(ctx *gin.Context) {
	var req GenerateCocktailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	prompt, err := server.cocktailGenerator.GenerateRecipe(req.ReferenceFlavours, req.ReferenceCocktails)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, prompt)
}
