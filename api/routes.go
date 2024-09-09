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
	* Flavour routes (Private)
	 */

	authRoutes.GET("/flavours/:id", server.getFlavourId)
	authRoutes.POST("/flavours", server.createFlavour)
	authRoutes.DELETE("/flavours/:id", server.deleteFlavour)
	authRoutes.PATCH("/flavours/:id", server.updateFlavour)
	authRoutes.GET("/flavours", server.getAllFlavours)

	/*
	* Categories routes (Private)
	 */
	authRoutes.GET("/categories", server.getAllCategories)
	authRoutes.GET("/categories/:id", server.getCategoryById)
	authRoutes.POST("/categories", server.createCategory)
	authRoutes.DELETE("/categories/:id", server.deleteCategory)
	authRoutes.PATCH("/categories/:id", server.updateCategory)

	/*
	 * Cocktail routes (Private)
	 */
	authRoutes.POST("/cocktails", server.createCocktail)
	authRoutes.DELETE("/cocktails/:id", server.deleteCocktail)
	authRoutes.GET("/cocktails/:id", server.getCocktail)
	authRoutes.PUT("/cocktails/:id", server.updateCocktail)

	/*
	 * AI Cocktail routes (Private)
	 */

	/*
	 * Flavour routes (Private)
	 */
	server.router = router
}
