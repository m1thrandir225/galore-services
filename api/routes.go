package api

import (
	"github.com/gin-gonic/gin"
)

func (server *Server) setupRouter() {
	router := gin.Default()

	v1 := router.Group("/api/v1")

	authRoutes := v1.Group("/").Use(authMiddleware(server.tokenMaker))

	/*
	 * Public routes
	 */
	v1.POST("/register", server.registerUser)
	v1.POST("/login", server.loginUser)

	/*
	 * Private routes
	 */
	authRoutes.POST("/logout", server.logout)
	authRoutes.GET("/users/:id/profile", server.getUserDetails)
	authRoutes.DELETE("/users/:id/profile", server.deleteUser)
	authRoutes.POST("/users/:id/profile", server.editUser)
	authRoutes.POST("/users/:id/change-password", server.changeUserPassword)

	server.router = router
}
