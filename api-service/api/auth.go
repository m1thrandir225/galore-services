package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/m1thrandir225/galore-services/db/sqlc"
	"github.com/m1thrandir225/galore-services/util"
)

type registerUserRequest struct {
	Name      string                `form:"name" json:"name" binding:"required"`
	Email     string                `form:"email" json:"email" binding:"required,email"`
	AvatarUrl *multipart.FileHeader `form:"avatar_url" json:"avatar_url" binding:"required"`
	Password  string                `form:"password" json:"password" binding:"required"`
	Birthday  string                `form:"birthday" json:"birthday" binding:"required"`
}

type registerUserResponse struct {
	User                  db.CreateUserRow `json:"user"`
	AccessTokenExpiresAt  time.Time        `json:"access_token_expires_at"`
	AccessToken           string           `json:"access_token"`
	RefreshToken          string           `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time        `json:"refresh_token_expires_at"`
	SessionID             uuid.UUID        `json:"session_id"`
}

type loginUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
	SessionID    string `json:"session_id" binding:"required,uuid"`
}

type refreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type loginUserResponse struct {
	User                  db.CreateUserRow `json:"user"`
	RefreshToken          string           `json:"refresh_token"`
	AccessToken           string           `json:"access_token"`
	RefreshTokenExpiresAt time.Time        `json:"refresh_token_expires_at"`
	AccessTokenExpiresAt  time.Time        `json:"access_token_expires_at"`
	SessionID             uuid.UUID        `json:"session_id"`
}

type logoutRequest struct {
	SessionID uuid.UUID `json:"session_id" binding:"required"`
}

type logoutResponse struct {
	Message string `json:"message"`
}

func (server *Server) registerUser(ctx *gin.Context) {
	var requestData registerUserRequest

	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}
	dbDate, err := util.TimeToDbDate(requestData.Birthday)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassowrd(requestData.Password)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	imageData, err := util.BytesFromFile(requestData.AvatarUrl)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	userId := uuid.New()

	avatarUrl, err := server.storage.UploadFile(imageData, userId.String(), requestData.AvatarUrl.Filename)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	args := db.CreateUserParams{
		ID:        userId,
		Email:     requestData.Email,
		Birthday:  dbDate,
		Name:      requestData.Name,
		AvatarUrl: avatarUrl,
		Password:  hashedPassword,
	}

	newEntry, err := server.store.CreateUser(ctx, args)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	accessToken, accessTokenPayload, err := server.tokenMaker.CreateToken(newEntry.ID, server.config.AccessTokenDuration)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshTokenPayload, err := server.tokenMaker.CreateToken(newEntry.ID, server.config.RefreshTokenDuration)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshTokenPayload.ID,
		Email:        args.Email,
		RefreshToken: refreshToken,
		ClientIp:     ctx.ClientIP(),
		UserAgent:    ctx.Request.UserAgent(),
		IsBlocked:    false,
		ExpiresAt:    refreshTokenPayload.ExpiredAt,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	b, err := json.Marshal(session)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	server.cache.SaveItem(ctx, args.Email, string(b))

	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := registerUserResponse{
		User:                  newEntry,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessTokenPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshTokenPayload.ExpiredAt,
		SessionID:             session.ID,
	}

	ctx.JSON(http.StatusOK, response)
}

func (server *Server) loginUser(ctx *gin.Context) {
	var requestData loginUserRequest

	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByEmail(ctx, requestData.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.ComparePassword(user.Password, requestData.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessTokenPayload, err := server.tokenMaker.CreateToken(user.ID, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshTokenPayload, err := server.tokenMaker.CreateToken(user.ID, server.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshTokenPayload.ID,
		Email:        requestData.Email,
		RefreshToken: refreshToken,
		ClientIp:     ctx.ClientIP(),
		UserAgent:    ctx.Request.UserAgent(),
		ExpiresAt:    time.Now().Add(server.config.RefreshTokenDuration),
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := loginUserResponse{
		SessionID: session.ID,
		User: db.CreateUserRow{
			ID:                        user.ID,
			Name:                      user.Name,
			Email:                     user.Email,
			AvatarUrl:                 user.AvatarUrl,
			Birthday:                  user.Birthday,
			EnabledPushNotifications:  user.EnabledPushNotifications,
			EnabledEmailNotifications: user.EnabledPushNotifications,
			CreatedAt:                 user.CreatedAt,
		},
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessTokenPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshTokenPayload.ExpiredAt,
	}

	ctx.JSON(http.StatusOK, response)
}

func (server *Server) logout(ctx *gin.Context) {
	// This will be an auth only route so you would need to send a acces token, and as a specific meter need to send the session id
	// so when you logout the session can be deleted or maybe set it to blocked, let it be for the time being just be deleted
	// TODO: set the session to blocked instead of deleting it so the user can see his own previous sessions.

	var requestData logoutRequest

	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	_, err := server.store.InvalidateSession(ctx, requestData.SessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// No need to return anything
	ctx.JSON(http.StatusOK, logoutResponse{
		Message: "You have been logged out",
	})
}

func (server *Server) refreshToken(ctx *gin.Context) {
	var requestData refreshTokenRequest

	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	sessionId, err := uuid.Parse(requestData.SessionID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	payload, err := server.tokenMaker.VerifyToken(requestData.RefreshToken)
	if err != nil {
		_, _ = server.store.InvalidateSession(ctx, sessionId)
		//No need to handle the error as if there is an error it is still unauthorized.

		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	newToken, _, err := server.tokenMaker.CreateToken(payload.UserId, server.config.AccessTokenDuration)

	ctx.JSON(http.StatusOK, refreshTokenResponse{
		AccessToken: newToken,
	})

}

func (server *Server) verifyAuthCookie(cookie string) bool {
	payload, err := server.tokenMaker.VerifyToken(cookie)
	if err != nil {
		return false
	}
	err = payload.Valid()

	return err == nil
}
