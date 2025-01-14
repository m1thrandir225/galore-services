package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/m1thrandir225/galore-services/db/sqlc"
)

type GetUserCocktailLikedStatus struct {
	UserID string `json:"user_id" binding:"required,uuid"`
}

func (server *Server) likeCocktail(ctx *gin.Context) {
	var uriData UriId

	if err := ctx.ShouldBindUri(&uriData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	cocktailId, err := uuid.Parse(uriData.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	payload, err := getPayloadFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.LikeCocktailParams{
		CocktailID: cocktailId,
		UserID:     payload.UserId,
	}

	_, err = server.store.LikeCocktail(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Status(http.StatusOK)
}

func (server *Server) unlikeCocktail(ctx *gin.Context) {
	var uriData UriId

	if err := ctx.ShouldBindUri(&uriData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	cocktailId, err := uuid.Parse(uriData.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	payload, err := getPayloadFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.UnlikeCocktailParams{
		CocktailID: cocktailId,
		UserID:     payload.UserId,
	}

	err = server.store.UnlikeCocktail(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Status(http.StatusOK)
}

func (server *Server) getCocktailLikedStatus(ctx *gin.Context) {
	var uriData UriId
	var reqData GetUserCocktailLikedStatus

	if err := ctx.ShouldBindUri(&uriData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&reqData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	cocktailID, err := uuid.Parse(uriData.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	userId, err := uuid.Parse(reqData.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	status, err := server.store.IsCocktailLiked(ctx, db.IsCocktailLikedParams{
		CocktailID: cocktailID,
		UserID:     userId,
	})
	log.Println(status)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, status)
}

func (server *Server) getUserLikedCocktails(ctx *gin.Context) {
	var uriData UriId
	if err := ctx.ShouldBindUri(&uriData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userId, err := uuid.Parse(uriData.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	cocktails, err := server.store.GetLikedCocktails(ctx, userId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, cocktails)
}
