package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/m1thrandir225/galore-services/db/sqlc"
	"github.com/m1thrandir225/galore-services/util"
)

type registerUserRequest struct {
	Name      string `json:"name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	AvatarUrl string `json:"avatar_url" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Birthday  string `json:"birthday" binding:"required"`
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

type loginUserResponse struct {
	SessionID             uuid.UUID `json:"session_id"`
	User                  db.User   `json:"user"`
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

func (server *Server) registerUser(ctx *gin.Context) {
	var requestData registerUserRequest

	if err := ctx.Bind(&requestData); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}
	format := "2006-01-02"

	birthdayDate, err := time.Parse(format, requestData.Birthday)

	var dbDate pgtype.Date

	err = dbDate.Scan(birthdayDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassowrd(requestData.Password)

	args := db.CreateUserParams{
		Email:     requestData.Email,
		Birthday:  dbDate,
		Name:      requestData.Name,
		AvatarUrl: requestData.AvatarUrl,
		Password:  hashedPassword,
	}

	newEntry, err := server.store.CreateUser(ctx, args)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	accessToken, accessTokenPayload, err := server.tokenMaker.CreateToken(args.Email, server.config.AccessTokenDuration)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshTokenPayload, err := server.tokenMaker.CreateToken(args.Email, server.config.RefreshTokenDuration)

	if err != nil {
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

	if err := ctx.Bind(&requestData); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByEmail(ctx, requestData.Email)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.ComparePassword(user.Password, requestData.Password)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessTokenPayload, err := server.tokenMaker.CreateToken(requestData.Email, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshTokenPayload, err := server.tokenMaker.CreateToken(requestData.Email, server.config.RefreshTokenDuration)

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
		SessionID:             session.ID,
		User:                  user,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessTokenPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshTokenPayload.ExpiredAt,
	}

	ctx.JSON(http.StatusOK, response)

}

func (server *Server) logout(ctx *gin.Context) {
	// probably send a session_id and an email and logout??
	// this is mainly to set the session to is_blocked
}

func (server *Server) verifyAuthCookie(cookie string) bool {
	payload, err := server.tokenMaker.VerifyToken(cookie)
	if err != nil {
		return false
	}
	err = payload.Valid()

	if err != nil {
		return false
	}

	return true
}
