package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/m1thrandir225/galore-services/db/sqlc"
)

type CreateNotificationRequest struct {
	UserID             string `json:"user_id" binding:"required,uuid"`
	NotificationTypeId string `json:"notification_type_id" binding:"required,uuid"`
}

type UpdateNotificationStatusRequest struct {
	Opened *bool `json:"opened" binding:"required"`
}

func (server *Server) createNotification(ctx *gin.Context) {
	var requestData CreateNotificationRequest

	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userId, err := uuid.Parse(requestData.UserID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	notificationTypeId, err := uuid.Parse(requestData.NotificationTypeId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateNotificationParams{
		UserID:             userId,
		NotificationTypeID: notificationTypeId,
	}

	notification, err := server.store.CreateNotification(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//TODO: IMPLEMENT SENDING THE NOTIFICATION
	err = server.scheduler.EnqueueJob("send_notification", map[string]interface{}{
		"notification_type_id": notificationTypeId,
		"user_id":              userId,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, notification)
}

func (server *Server) updateNotificationStatus(ctx *gin.Context) {
	var uriData UriId
	var requestData UpdateNotificationStatusRequest

	if err := ctx.ShouldBindUri(&uriData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	notificationId, err := uuid.Parse(uriData.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.UpdateUserNotificationParams{
		ID:     notificationId,
		Opened: *requestData.Opened,
	}

	updated, err := server.store.UpdateUserNotification(ctx, arg)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, updated)
}

func (server *Server) getUserNotifications(ctx *gin.Context) {
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

	notifications, err := server.store.GetUserNotifications(ctx, userId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, notifications)
}
