package api

import "github.com/gin-gonic/gin"

func (server *Server) registerUser(ctx *gin.Context) {
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
