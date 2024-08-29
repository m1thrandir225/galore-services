package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
	User        db.CreateUserRow `json:"user"`
	AccessToken string           `json:"access_token"`
}

type LoginUserRequest struct{}

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

	accessToken, err := server.tokenMaker.CreateToken(args.Email, server.config.AccessTokenDuration)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//TODO: implement adding file to storage

	response := registerUserResponse{
		User:        newEntry,
		AccessToken: accessToken,
	}

	ctx.JSON(http.StatusOK, response)
}

func (server *Server) loginUser(ctx *gin.Context) {
}

func (server *Server) logout(ctx *gin.Context) {}

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
