package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/m1thrandir225/galore-services/db/sqlc"
)

type GetNotificationTypesResponse struct {
	Types []db.NotificationType `json:"notification_types"`
}

type NotificationTypeIDUri struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type UpdateNotificationTypeRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

type CreateNotificationTypeRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

func (server *Server) createNotificationType(ctx *gin.Context) {
	var requestData CreateNotificationTypeRequest

	if err := ctx.Bind(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateNotificationTypeParams{
		Tag:     requestData.Tag,
		Content: requestData.Content,
		Title:   requestData.Title,
	}

	notification_type, err := server.store.CreateNotificationType(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, notification_type)
}

func (server *Server) getNotificationTypes(ctx *gin.Context) {
	notification_types, err := server.store.GetAllTypes(ctx)

	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, GetNotificationTypesResponse{
		Types: notification_types,
	})
}

func (server *Server) getNotificationType(ctx *gin.Context) {
	var requestData NotificationTypeIDUri

	if err := ctx.ShouldBindUri(&requestData); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	id, err := uuid.Parse(requestData.ID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	notification_type, err := server.store.GetNotificationType(ctx, id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, notification_type)
}

func (server *Server) deleteNotificationType(ctx *gin.Context) {
	var requestData NotificationTypeIDUri

	if err := ctx.ShouldBindUri(&requestData); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}
	id, err := uuid.Parse(requestData.ID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.store.DeleteNotificationType(ctx, id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Status(http.StatusOK)
}

func (server *Server) updateNotificationType(ctx *gin.Context) {
	var uriData NotificationTypeIDUri
	var requestData UpdateNotificationTypeRequest
	if err := ctx.ShouldBindUri(&uriData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := uuid.Parse(uriData.ID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.UpdateNotificationTypeParams{
		ID:      id,
		Title:   requestData.Title,
		Content: requestData.Content,
		Tag:     requestData.Tag,
	}

	updated, err := server.store.UpdateNotificationType(ctx, args)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, updated)
}
