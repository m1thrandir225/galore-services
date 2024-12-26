package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/m1thrandir225/galore-services/db/sqlc"
)

type LikeFlavoursRequest struct {
	FlavourIds []string `json:"flavour_ids" binding:"required"`
	UserId     string   `json:"user_id" binding:"required"`
}

func (server *Server) setupLikedFlavourRoutes(authRoutes *gin.RouterGroup) {

}

func (server *Server) likeFlavours(ctx *gin.Context) {
	var requestData LikeFlavoursRequest

	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userId := uuid.MustParse(requestData.UserId)

	var flavourIds []uuid.UUID
	for _, flavourId := range requestData.FlavourIds {
		flavourIds = append(flavourIds, uuid.MustParse(flavourId))
	}

	arg := db.LikeFlavoursParams{
		Flavourids: flavourIds,
		Userid:     userId,
	}

	flavours, err := server.store.LikeFlavours(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if len(flavours) != len(requestData.FlavourIds) {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Status(http.StatusOK)
}

func (server *Server) likeFlavour(ctx *gin.Context) {
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

func (server *Server) getLikedFlavour(ctx *gin.Context) {
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

	arg := db.GetLikedFlavourParams{
		UserID:    payload.UserId,
		FlavourID: flavourId,
	}

	flavour, err := server.store.GetLikedFlavour(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, flavour)
}

func (server *Server) getUserLikedFlavours(ctx *gin.Context) {
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

	flavours, err := server.store.GetUserLikedFlavours(ctx, userId)

	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusOK, []db.Flavour{})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, flavours)
}

func (server *Server) unlikeFlavour(ctx *gin.Context) {
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

	arg := db.UnlikeFlavourParams{
		FlavourID: flavourId,
		UserID:    payload.UserId,
	}

	err = server.store.UnlikeFlavour(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Status(http.StatusOK)
}
