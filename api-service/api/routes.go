package api

import (
	"github.com/gin-gonic/gin"
)

func (server *Server) setupRouter() {
	router := gin.Default()

	router.Static("/public", "./public")

	v1 := router.Group("/api/v1")
	/**
	Public Routes
	*/
	v1.POST("/register", server.registerUser)
	v1.POST("/login", server.loginUser)
	v1.POST("/refresh", server.refreshToken)
	v1.POST("/forgot-password", server.forgotPassword)
	v1.POST("/verify-otp", server.verifyOTP)
	v1.POST("/reset-password", server.resetPassword)
	v1.GET("/health", server.checkServiceHealth)

	/**
	Private Routes
	*/
	authRoutes := v1.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/logout", server.logout)
	/**
	User Routes
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
	authRoutes.GET("/users/:id/homescreen", server.getHomescreenForUser)
	authRoutes.GET("/users/:id/notifications", server.getUserNotifications)
	/**
	Cocktail Routes
	*/
	authRoutes.GET("/cocktails", server.getCocktails)
	authRoutes.POST("/cocktails", server.createCocktail)
	authRoutes.GET("/cocktails/:id", server.getCocktail)
	authRoutes.PUT("/cocktails/:id", server.updateCocktail)
	authRoutes.DELETE("/cocktails/:id", server.deleteCocktail)
	authRoutes.GET("/cocktails/:id/was_featured", server.checkWasCocktailFeatured)
	authRoutes.GET("/cocktails/:id/categories", server.getCategoriesForCocktail)
	authRoutes.GET("/cocktails/:id/simillar", server.getSimilarCocktails)

	/**
	Liked Cocktail Routes
	*/
	authRoutes.POST("/cocktails/:id/is_liked", server.getCocktailLikedStatus)
	authRoutes.POST("/cocktails/:id/like", server.likeCocktail)
	authRoutes.POST("/cocktails/:id/unlike", server.unlikeCocktail)

	/**
	Daily Featured Cocktail Routes
	*/
	authRoutes.GET("/cocktails/featured", server.getDailyFeatured)
	authRoutes.POST("/cocktails/featured", server.createDailyFeatured)
	authRoutes.DELETE("/cocktails/featured", server.deleteOlderFeatured)

	/**
	Category Routes
	*/
	authRoutes.GET("/categories", server.getAllCategories)
	authRoutes.POST("/categories", server.createCategory)
	authRoutes.GET("/categories/:id", server.getCategoryById)
	authRoutes.DELETE("/categories/:id", server.deleteCategory)
	authRoutes.PATCH("/categories/:id", server.updateCategory)
	authRoutes.GET("/categories/:id/cocktails", server.getCocktailsByCategory)

	/**
	Flavour Routes
	*/
	authRoutes.GET("/flavours", server.getAllFlavours)
	authRoutes.POST("/flavours", server.createFlavour)
	authRoutes.GET("/flavours/:id", server.getFlavourId)
	authRoutes.DELETE("/flavours/:id", server.deleteFlavour)
	authRoutes.PATCH("/flavours/:id", server.updateFlavour)
	/**
	Liked Flavour Routes
	*/
	authRoutes.GET("/flavours/:id/like", server.getLikedFlavour)
	authRoutes.POST("/flavours/like", server.likeFlavours)
	authRoutes.POST("/flavours/:id/like", server.likeFlavour)
	authRoutes.POST("/flavours/:id/unlike", server.unlikeFlavour)

	/**
	Category Flavour Routes
	*/
	authRoutes.POST("/category_flavour", server.createCategoryFlavour)
	authRoutes.DELETE("/category_flavour/:id", server.deleteCategoryFlavour)

	/**
	Notification Type Routes
	*/
	authRoutes.GET("/notification_types", server.getNotificationTypes)
	authRoutes.POST("/notification_types", server.createNotificationType)
	authRoutes.GET("/notification_types/:id", server.getNotificationType)
	authRoutes.DELETE("/notification_types/:id", server.deleteNotificationType)
	authRoutes.PUT("/notification_types/:id", server.updateNotificationType)
	/**
	Notifications Routes
	*/
	authRoutes.POST("/notifications", server.createNotification)
	authRoutes.PATCH("/notifications/:id", server.updateNotificationStatus)

	/**
	Generate Cocktail Routes
	*/
	v1.POST("/generate", server.generateCocktail)

	server.setupMigrationServiceRoutes(v1)
	server.setupCategorizerServiceRoutes(v1)

	server.router = router
}

func (server *Server) setupMigrationServiceRoutes(rg *gin.RouterGroup) {
	migrationRoutes := rg.Group("/migration").Use(microserviceAccessMiddleware(server.config.MigrationServiceKey))

	migrationRoutes.POST("/cocktails", server.createCocktailMigrate)
}

func (server *Server) setupCategorizerServiceRoutes(rg *gin.RouterGroup) {
	categorizerRoutes := rg.Group("/categorizer").Use(microserviceAccessMiddleware(server.config.CategorizerServiceKey))

	categorizerRoutes.GET("/category/:tag", server.getCategoryByTag)           // Check if category exists
	categorizerRoutes.POST("/category", server.createCategory)                 // Create a category
	categorizerRoutes.POST("/category_cocktail", server.addCocktailToCategory) // Add a cocktail to a given category
}
