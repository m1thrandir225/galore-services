package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/m1thrandir225/galore-services/db/sqlc"
	"net/http"
)

func (server *Server) LikeFlavour(ctx *gin.Context) {
	var uriData UriId

	if err := ctx.ShouldBindUri(&uriData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	flavourId, err := uuid.Parse(uriData.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	payload, err := getPayloadFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.LikeFlavourParams{
		FlavourID: flavourId,
		UserID:    payload.UserId,
	}

	likedFlavour, err := server.store.LikeFlavour(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, likedFlavour)
}

func (server *Server) GetLikedFlavour(ctx *gin.Context) {}

func (server *Server) GetUserLikedFlavours(ctx *gin.Context) {}

func (server *Server) UnlikeFlavour(ctx *gin.Context) {}
