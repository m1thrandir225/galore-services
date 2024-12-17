package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/m1thrandir225/galore-services/db/sqlc"
)

type GetNotificationTypesResponse struct {
	Types []db.NotificationType `json:"notification_types"`
}

type UpdateNotificationTypeRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	Tag     string `json:"tag" binding:"required"`
}

type CreateNotificationTypeRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	Tag     string `json:"tag" binding:"required"`
}

func (server *Server) setupNotificationTypesRoutes(authRoutes *gin.RouterGroup) {
	notificationTypes := authRoutes.Group("/notification_types")

	notificationTypes.GET("/", server.getNotificationTypes)
	notificationTypes.DELETE("/:id", server.deleteNotificationType)
	notificationTypes.PUT("/:id", server.updateNotificationType)
	notificationTypes.GET("/:id", server.getNotificationType)
	notificationTypes.POST("/", server.createNotificationType)
}

func (server *Server) createNotificationType(ctx *gin.Context) {
	var requestData CreateNotificationTypeRequest

	if err := ctx.BindJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateNotificationTypeParams{
		Tag:     requestData.Tag,
		Content: requestData.Content,
		Title:   requestData.Title,
	}
	notificationType, err := server.store.CreateNotificationType(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, notificationType)
}

func (server *Server) getNotificationTypes(ctx *gin.Context) {
	notificationTypes, err := server.store.GetAllTypes(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, GetNotificationTypesResponse{
		Types: notificationTypes,
	})
}

func (server *Server) getNotificationType(ctx *gin.Context) {
	var requestData UriId

	if err := ctx.ShouldBindUri(&requestData); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	id, err := uuid.Parse(requestData.ID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	notificationType, err := server.store.GetNotificationType(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, notificationType)
}

func (server *Server) deleteNotificationType(ctx *gin.Context) {
	var requestData UriId

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
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Status(http.StatusOK)
}

func (server *Server) updateNotificationType(ctx *gin.Context) {
	var uriData UriId
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
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, updated)
}
