package api

import "github.com/gin-gonic/gin"

func (server *Server) setupRouter() {
	router := gin.Default()

	server.router = router
}
