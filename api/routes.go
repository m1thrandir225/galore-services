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

	/*
	 * User routes (Private)
	 *
	 */
	authRoutes.GET("/users/:id", server.getUserDetails)
	authRoutes.DELETE("/users/:id", server.deleteUser)
	authRoutes.POST("/users/:id", server.updateUserInformation)
	authRoutes.PUT("/users/:id/password", server.changeUserPassword)
	authRoutes.PUT("/users/:id/push-notifications", server.updateUserPushNotifications)
	authRoutes.PUT("/users/:id/email-notifications", server.updateUserEmailNotifications)

	/*
	 * Notification Type routes (Private)
	 */
	authRoutes.GET("/notification_types", server.getNotificationTypes)
	authRoutes.DELETE("/notification_types/:id", server.deleteNotificationType)
	authRoutes.PUT("/notification_types/:id", server.updateNotificationType)
	authRoutes.GET("/notification_types/:id", server.getNotificationType)
	authRoutes.POST("/notification_types", server.createNotificationType)

	/*
	 * Notifications routes (Private)
	 */

	authRoutes.GET("/users/:id/notifications", server.getUserNotifications)
	authRoutes.POST("/notifications", server.createNotification)
	authRoutes.PATCH("/notifications/:id", server.updateNotificationStatus)
	/*
	 * Cocktail routes (Private)
	 */

	/*
	 * AI Cocktail routes (Private)
	 */

	/*
	 * Flavour routes (Private)
	 */
	server.router = router
}
