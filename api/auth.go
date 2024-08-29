package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/m1thrandir225/galore-services/db/sqlc"
)

type RegisterUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required, email"`
	Password string `json:"password" binding:"required"`
	Birthday string `json:"birthday" binding:"required"`
}

type LoginUserRequest struct{}

func (server *Server) registerUser(ctx *gin.Context) {
	var requestData RegisterUserRequest

	if err := ctx.Bind(&requestData); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	dateFormat := "2022-03-01"

	birthdayDate, err := time.Parse(dateFormat, requestData.Birthday)

	dbDate := pgtype.Date{
		Time: birthdayDate,
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//TODO: hash password
	args := db.CreateUserParams{
		Email:     requestData.Email,
		Birthday:  dbDate,
		Name:      requestData.Name,
		AvatarUrl: "",
		Password:  "",
	}

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
