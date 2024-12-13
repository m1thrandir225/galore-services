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
	/*
	 * Public routes
	 */
	v1.POST("/register", server.registerUser)
	v1.POST("/login", server.loginUser)
	v1.POST("/refresh", server.refreshToken)
	v1.GET("/health", server.checkService)

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
	authRoutes.GET("/users/:id/flavours", server.getUserLikedFlavours)
	authRoutes.GET("/users/:id/cocktails", server.getUserLikedCocktails)
	authRoutes.GET("/users/:id/categories", server.getCategoriesBasedOnLikedFlavours)

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
	categorizerServerRoutes.GET("/category/:tag", server.getCategoryByTag)           // Check if category exists
	categorizerServerRoutes.POST("/category", server.createCategory)                 // Create a category
	categorizerServerRoutes.POST("/category_cocktail", server.addCocktailToCategory) // Add a cocktail to a given category

	authRoutes.GET("/categories", server.getAllCategories)
	authRoutes.GET("/categories/:id", server.getCategoryById)
	authRoutes.POST("/categories", server.createCategory)
	authRoutes.DELETE("/categories/:id", server.deleteCategory)
	authRoutes.PATCH("/categories/:id", server.updateCategory)
	authRoutes.PUT("/categories/:id/cocktails", server.getCocktailsByCategory)

	/*
	 * Cocktail routes (Private)
	 */
	migrationServerRoutes.POST("/cocktails", server.createCocktailMigrate)

	authRoutes.GET("/cocktails", server.getCocktails)
	authRoutes.POST("/cocktails", server.createCocktail)
	authRoutes.GET("/cocktails/:id", server.getCocktail)
	authRoutes.PUT("/cocktails/:id", server.updateCocktail)
	authRoutes.DELETE("/cocktails/:id", server.deleteCocktail)
	authRoutes.GET("/cocktails/:id/was_featured", server.checkWasCocktailFeatured)
	authRoutes.GET("/cocktails/:id/categories", server.getCategoriesForCocktail)

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
	authRoutes.POST("/flavours/:id/like", server.likeFlavour)
	authRoutes.POST("/flavours/:id/unlike", server.unlikeFlavour)
	authRoutes.GET("/flavours/:id/like", server.getLikedFlavour)

	/**
	* Category Flavours (Private)
	 */
	authRoutes.POST("/category-flavour", server.createCategoryFlavour)
	authRoutes.DELETE("/category-flavour/:id", server.deleteCategoryFlavour)

	/*
	* Daily Featured Cocktails
	 */
	authRoutes.GET("/cocktails/featured", server.getDailyFeatured)
	authRoutes.POST("/cocktails/featured", server.createDailyFeatured)
	authRoutes.DELETE("/cocktails/featured", server.deleteOlderFeatured)

	server.router = router
}
