package auth

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/m1thrandir225/galore-services/internal/dto"
	"github.com/m1thrandir225/galore-services/internal/services"
	"github.com/m1thrandir225/galore-services/pkg/shared"
)

type Handler struct {
	authService          services.AuthService
	sessionService       services.SessionService
	userService          services.UserService
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func NewHandler(
	authService services.AuthService,
	sessionService services.SessionService,
	userService services.UserService,
	accessTokenDuration time.Duration,
	refreshTokenDuration time.Duration,
) *Handler {
	return &Handler{
		authService:          authService,
		sessionService:       sessionService,
		userService:          userService,
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}
}

func (h *Handler) RegisterUser(ctx *gin.Context) {
	var requestData registerUserRequest

	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(err))
		return
	}

	userData, err := h.authService.RegisterUser(
		ctx,
		requestData.Name,
		requestData.Email,
		requestData.Password,
		requestData.Birthday,
		requestData.AvatarUrl,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, shared.ErrorResponse(err))
		return
	}

	userDTO := dto.UserDTOFromDB(*userData)

	accessToken, accessTokenPayload, err := h.authService.CreateToken(
		ctx,
		userData.ID,
		h.accessTokenDuration,
		services.AccessToken,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, shared.ErrorResponse(err))
		return
	}

	refreshToken, refreshTokenPayload, err := h.authService.CreateToken(
		ctx,
		userData.ID,
		h.refreshTokenDuration,
		services.RefreshToken,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, shared.ErrorResponse(err))
		return
	}

	session, err := h.sessionService.CreateSession(
		ctx,
		userData.ID,
		userData.Email,
		refreshToken,
		ctx.ClientIP(),
		ctx.Request.UserAgent(),
		refreshTokenPayload.ExpiredAt,
		false,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, shared.ErrorResponse(err))
		return
	}

	response := authenticatedUserResponse{
		User:                  userDTO,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessTokenPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshTokenPayload.ExpiredAt,
		SessionID:             session.ID,
	}
	//TODO: Send welcome email

	ctx.JSON(http.StatusOK, response)
}

func (h *Handler) LoginUser(ctx *gin.Context) {
	var requestData loginUserRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(err))
		return
	}

	user, err := h.authService.LoginUser(ctx, requestData.Email, requestData.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, shared.ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(err)) //incorrect password
		return
	}

	userDTO := dto.UserDTOFromDB(*user)

	accessToken, accessTokenPayload, err := h.authService.CreateToken(
		ctx,
		user.ID,
		h.accessTokenDuration,
		services.AccessToken,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, shared.ErrorResponse(err))
		return
	}

	refreshToken, refreshTokenPayload, err := h.authService.CreateToken(
		ctx,
		user.ID,
		h.refreshTokenDuration,
		services.RefreshToken,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, shared.ErrorResponse(err))
		return
	}

	session, err := h.sessionService.CreateSession(
		ctx,
		user.ID,
		user.Email,
		refreshToken,
		ctx.ClientIP(),
		ctx.Request.UserAgent(),
		refreshTokenPayload.ExpiredAt,
		false,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, shared.ErrorResponse(err))
		return
	}

	response := authenticatedUserResponse{
		User:                  userDTO,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessTokenPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshTokenPayload.ExpiredAt,
		SessionID:             session.ID,
	}
	ctx.JSON(http.StatusOK, response)
}

func (h *Handler) LogoutUser(ctx *gin.Context) {
	var requestData logoutRequest

	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(err))
		return
	}

	_, err := h.sessionService.InvalidateSession(ctx, requestData.SessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, shared.ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, shared.ErrorResponse(err))
		return
	}

	response := shared.MessageResponse("You have been logged out")

	ctx.JSON(http.StatusOK, response)
}

func (h *Handler) RefreshToken(ctx *gin.Context) {
	var requestData refreshTokenRequest

	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(err))
		return
	}

	sessionId, err := uuid.Parse(requestData.SessionID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, shared.ErrorResponse(err))
		return
	}

	payload, err := h.authService.VerifyToken(ctx, requestData.RefreshToken)

	if err != nil {
		_, _ = h.sessionService.InvalidateSession(ctx, sessionId)

		ctx.JSON(http.StatusUnauthorized, shared.ErrorResponse(err))
		return
	}

	newToken, _, err := h.authService.CreateToken(
		ctx,
		payload.UserId,
		h.accessTokenDuration,
		services.AccessToken,
	)

	response := refreshTokenResponse{
		AccessToken: newToken,
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *Handler) ForgotPasswordRequest(ctx *gin.Context) {
	var requestData forgotPasswordRequest

	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

}

func (h *Handler) VerifyForgotPasswordRequest(ctx *gin.Context) {}

func (h *Handler) ResetPassword(ctx *gin.Context) {}
