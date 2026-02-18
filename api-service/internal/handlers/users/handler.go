package users

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/m1thrandir225/galore-services/internal/security"
	"github.com/m1thrandir225/galore-services/internal/services"
	"github.com/m1thrandir225/galore-services/pkg/shared"
)

type Handler struct {
	userService services.UserService
}

func NewHandler(userService services.UserService) *Handler {
	return &Handler{
		userService: userService,
	}
}

func (h *Handler) GetUserDetails(ctx *gin.Context) {
	var uriData UserID

	if err := ctx.ShouldBindUri(&uriData); err != nil {
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(err))
		return
	}

	userId, err := uuid.Parse(uriData.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(err))
		return
	}

	err = security.VerifyUserIDWithToken(userId, ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, shared.ErrorResponse(err))
		return
	}

	userDetails, err := h.userService.GetUserDetails(ctx, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, shared.ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, shared.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, userDetails)
}

func (h *Handler) UpdateUserInformation(ctx *gin.Context) {
	var uriData UserID
	var requestData UpdateUserInformationRequest

	if err := ctx.ShouldBindUri(&uriData); err != nil {
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(err))
		return
	}

	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(err))
		return
	}

	userId, err := uuid.Parse(uriData.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(err))
		return
	}

	err = security.VerifyUserIDWithToken(userId, ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, shared.ErrorResponse(err))
		return
	}

	updatedData, err := h.userService.UpdateUserInformation(ctx, userId, requestData.Email, requestData.AvatarUrl, requestData.Name, requestData.Birthday)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, shared.ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, shared.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, updatedData)
}

func (h *Handler) DeleteUser(ctx *gin.Context) {
	var uriData UserID

	if err := ctx.ShouldBindUri(&uriData); err != nil {
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(err))
		return
	}

	userId, err := uuid.Parse(uriData.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(err))
		return
	}

	err = security.VerifyUserIDWithToken(userId, ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, shared.ErrorResponse(err))
		return
	}

	err = h.userService.DeleteUser(ctx, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, shared.ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, shared.ErrorResponse(err))
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Handler) UpdateUserPassword(ctx *gin.Context) {
	var uriData UserID
	var requestData UpdateUserPasswordRequest

	if err := ctx.ShouldBindUri(&uriData); err != nil {
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(err))
		return
	}

	userId, err := uuid.Parse(uriData.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(err))
		return
	}
	err = security.VerifyUserIDWithToken(userId, ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, shared.ErrorResponse(err))
		return
	}

	err = h.userService.UpdateUserPassword(ctx, userId, requestData.NewPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, shared.ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, shared.ErrorResponse(err))
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Handler) UpdateUserPushNotificationsSettings(ctx *gin.Context) {
	var uriData UserID
	var requestData UpdateUserPushNotificationsRequest

	if err := ctx.ShouldBindUri(&uriData); err != nil {
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(err))
		return
	}
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(err))
		return
	}

	userId, err := uuid.Parse(uriData.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(err))
		return
	}

	err = security.VerifyUserIDWithToken(userId, ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, shared.ErrorResponse(err))
		return
	}

	updatedData, err := h.userService.UpdateUserPushNotificationsSettings(ctx, userId, *requestData.Enabled)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, shared.ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, shared.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, updatedData)
}

func (h *Handler) UpdateUserEmailNotificationsSettings(ctx *gin.Context) {
	var uriData UserID
	var requestData UpdateUserEmailNotificationsRequest

	if err := ctx.ShouldBindUri(&uriData); err != nil {
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(err))
		return
	}
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(err))
		return
	}

	userId, err := uuid.Parse(uriData.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(err))
		return
	}

	err = security.VerifyUserIDWithToken(userId, ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, shared.ErrorResponse(err))
		return
	}

	updatedData, err := h.userService.UpdateUserEmailNotificationsSettings(ctx, userId, *requestData.Enabled)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, shared.ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, shared.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, updatedData)
}
