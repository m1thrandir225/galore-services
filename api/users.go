package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/m1thrandir225/galore-services/db/sqlc"
	"github.com/m1thrandir225/galore-services/util"
	"mime/multipart"
	"net/http"
)

type UpdateUserInformationRequest struct {
	Email     string                `form:"email" json:"email" binding:"required"`
	AvatarUrl *multipart.FileHeader `form:"avatar_url" json:"avatar_url" binding:"omitempty"`
	Name      string                `form:"name" json:"name" binding:"required"`
	Birthday  string                `form:"birthday" json:"birthday" binding:"required"`
}
type UpdateUserEmailNotificationsRequest struct {
	Enabled bool `form:"enabled" json:"enabled" binding:"required"`
}

type UpdateUserPushNotificationsRequest struct {
	Enabled bool `form:"enabled" json:"enabled" binding:"required"`
}

type UpdateUserPasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required"`
}

func (server *Server) getUserDetails(ctx *gin.Context) {
	var uriData UriId

	if err := ctx.ShouldBindUri(&uriData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userId, err := uuid.Parse(uriData.ID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = verifyUserIdWithToken(userId, ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	userDetails, err := server.store.GetUser(ctx, userId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, userDetails)

}

func (server *Server) deleteUser(ctx *gin.Context) {
	var uriData UriId

	if err := ctx.ShouldBindUri(&uriData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userId, err := uuid.Parse(uriData.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = verifyUserIdWithToken(userId, ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	err = server.store.DeleteUser(ctx, userId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Status(http.StatusOK)
}

func (server *Server) updateUserPassword(ctx *gin.Context) {
	var uriData UriId
	var requestData UpdateUserPasswordRequest

	if err := ctx.ShouldBindUri(&uriData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userId, err := uuid.Parse(uriData.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = verifyUserIdWithToken(userId, ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassowrd(requestData.NewPassword)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.UpdateUserPasswordParams{
		ID:       userId,
		Password: hashedPassword,
	}

	err = server.store.UpdateUserPassword(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Status(http.StatusOK)

}

func (server *Server) updateUserInformation(ctx *gin.Context) {
	var uriData UriId
	var requestData UpdateUserInformationRequest
	var avatarFilePath string
	if err := ctx.ShouldBindUri(&uriData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	userId, err := uuid.Parse(uriData.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = verifyUserIdWithToken(userId, ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	userInformation, err := server.store.GetUser(ctx, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	newBirthday, err := util.TimeToDbDate(requestData.Birthday)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if requestData.AvatarUrl == nil {
		avatarFilePath = userInformation.AvatarUrl
	} else {
		newAvatarData, err2 := util.BytesFromFile(requestData.AvatarUrl)
		if err2 != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		newFilePath, err2 := server.storage.ReplaceFile(userInformation.AvatarUrl, newAvatarData)
		if err2 != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		avatarFilePath = newFilePath
	}

	arg := db.UpdateUserInformationParams{
		ID:        userId,
		Birthday:  newBirthday,
		AvatarUrl: avatarFilePath,
		Name:      requestData.Name,
		Email:     requestData.Email,
	}

	updated, err := server.store.UpdateUserInformation(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, updated)
}

func (server *Server) updateUserPushNotifications(ctx *gin.Context) {
	var uriData UriId
	var requestData UpdateUserPushNotificationsRequest
	if err := ctx.ShouldBindUri(&uriData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	userId, err := uuid.Parse(uriData.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = verifyUserIdWithToken(userId, ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	arg := db.UpdateUserPushNotificationsParams{
		ID:                       userId,
		EnabledPushNotifications: requestData.Enabled,
	}

	updated, err := server.store.UpdateUserPushNotifications(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, updated)
}

func (server *Server) updateUserEmailNotifications(ctx *gin.Context) {
	var uriData UriId
	var requestData UpdateUserEmailNotificationsRequest

	if err := ctx.ShouldBindUri(&uriData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userId, err := uuid.Parse(uriData.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = verifyUserIdWithToken(userId, ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	arg := db.UpdateUserEmailNotificationsParams{
		ID:                        userId,
		EnabledEmailNotifications: requestData.Enabled,
	}

	updated, err := server.store.UpdateUserEmailNotifications(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, updated)
}
