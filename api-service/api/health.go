package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthResponse struct {
	Status string `json:"status"`
}

func (server *Server) checkServiceHealth(ctx *gin.Context) {
	message := HealthResponse{
		Status: "health",
	}
	ctx.JSON(http.StatusOK, message)
}
