package api

import (
	"github.com/gin-gonic/gin"
)

func (server *Server) setupRouter() {
	router := gin.Default()

	router.Static("/public", "./public")

	v1 := router.Group("/api/v1")

	authRoutes := v1.Group("/").Use(authMiddleware(server.tokenMaker))
	migrationServerRoutes := v1.Group("/migration").Use(microserviceAccessMiddleware(server.config.MigrationServiceKey))
	categorizerServerRoutes := v1.Group("/categorizer").Use(microserviceAccessMiddleware(server.config.CategorizerServiceKey))

	// embeddingServerRoutes := v1.Group("/embedding").Use(microserviceAccessMiddleware(server.config.EmbeddingServiceKey))
	/*
	 * Public routes
	 */
	v1.POST("/register", server.registerUser)
	v1.POST("/login", server.loginUser)
	v1.POST("/refresh", server.refreshToken)

	/*
	 * Private routes (user routes)
	 */
	authRoutes.POST("/logout", server.logout)

	/*
	 * User routes (Private)
	 *
	 */
	authRoutes.GET("/users/:id", server.getUserDetails)
	authRoutes.DELETE("/users/:id", server.deleteUser)
	authRoutes.POST("/users/:id", server.updateUserInformation)
	authRoutes.PUT("/users/:id/password", server.updateUserPassword)
	authRoutes.PUT("/users/:id/push-notifications", server.updateUserPushNotifications)
	authRoutes.PUT("/users/:id/email-notifications", server.updateUserEmailNotifications)
	authRoutes.GET("/users/:id/flavours", server.GetUserLikedFlavours)
	authRoutes.GET("/users/:id/cocktails", server.getUserLikedCocktails)

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
	categorizerServerRoutes.GET("/category/:tag", server.getCategoryByTag)           //Check if category exists
	categorizerServerRoutes.POST("/category", server.createCategory)                 //Create a category
	categorizerServerRoutes.POST("/category_cocktail", server.addCocktailToCategory) //Add a cocktail to a given category

	authRoutes.GET("/categories", server.getAllCategories)
	authRoutes.GET("/categories/:id", server.getCategoryById)
	authRoutes.POST("/categories", server.createCategory)
	authRoutes.DELETE("/categories/:id", server.deleteCategory)
	authRoutes.PATCH("/categories/:id", server.updateCategory)

	/*
	 * Cocktail routes (Private)
	 */
	migrationServerRoutes.POST("/cocktails", server.createCocktailMigrate)

	authRoutes.POST("/cocktails", server.createCocktail)
	authRoutes.DELETE("/cocktails/:id", server.deleteCocktail)
	authRoutes.GET("/cocktails/:id", server.getCocktail)
	authRoutes.GET("/cocktails", server.getCocktails)
	authRoutes.PUT("/cocktails/:id", server.updateCocktail)

	/*
	* Liked Cocktail Routes (Private)
	 */
	authRoutes.POST("/cocktails/:id/like", server.likeCocktail)
	authRoutes.POST("/cocktails/:id/unlike", server.unlikeCocktail)

	/*
	 * AI Cocktail routes (Private)
	 */

	/*
	* Liked Flavours (Private)
	 */
	authRoutes.POST("/flavours/:id/like", server.LikeFlavour)
	authRoutes.POST("/flavours/:id/unlike", server.UnlikeFlavour)
	authRoutes.GET("/flavours/:id/like", server.GetLikedFlavour)

	server.router = router
}
