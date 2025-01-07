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
	authRoutes := v1.Group("/")
	authRoutes.Use(authMiddleware(server.tokenMaker))
	{
		/**
		Admin Routes
		*/
		adminRoutes := authRoutes.Group("/")
		adminRoutes.Use(adminMiddleware(server.store))
		{
			/**
			Admin User Routes
			*/

			/**
			Admin Cocktail Routes
			*/
			cocktailAdminRoutes := adminRoutes.Group("/cocktails")

			cocktailAdminRoutes.POST("/", server.createCocktail)
			cocktailAdminRoutes.DELETE("/:id", server.deleteCocktail)
			cocktailAdminRoutes.PUT("/:id", server.updateCocktail)

			/**
			Admin Category Routes
			*/
			categoryAdminRoutes := adminRoutes.Group("/categories")

			categoryAdminRoutes.POST("/", server.createCategory)
			categoryAdminRoutes.PATCH("/:id", server.updateCategory)
			categoryAdminRoutes.DELETE("/:id", server.deleteCategory)

			/**
			Admin Flavour Routes
			*/
			flavourAdminRoutes := adminRoutes.Group("/flavours")

			flavourAdminRoutes.POST("/", server.createFlavour)
			flavourAdminRoutes.DELETE("/:id", server.deleteFlavour)
			flavourAdminRoutes.PATCH("/:id", server.updateFlavour)

			/**
			Category Flavour Routes
			*/
			adminRoutes.POST("/category_flavour", server.createCategoryFlavour)
			adminRoutes.DELETE("/category_flavour/:id", server.deleteCategoryFlavour)

			/**
			Notification Type Routes
			*/
			notificationTypesRoutes := adminRoutes.Group("/notification_types")

			notificationTypesRoutes.GET("/", server.getNotificationTypes)
			notificationTypesRoutes.POST("/", server.createNotificationType)
			notificationTypesRoutes.GET("/:id", server.getNotificationType)
			notificationTypesRoutes.DELETE("/:id", server.deleteNotificationType)
			notificationTypesRoutes.PUT("/:id", server.updateNotificationType)
		}

		/**
		User Routes
		*/
		userRoutes := authRoutes.Group("/users")

		//TODO: add verification to check if the user-id is the one in the context
		userRoutes.GET("/:id", server.getUserDetails)
		userRoutes.DELETE("/:id", server.deleteUser)
		userRoutes.POST("/:id", server.updateUserInformation)
		userRoutes.PUT("/:id/password", server.updateUserPassword)
		userRoutes.PUT("/:id/push-notifications", server.updateUserPushNotifications)
		userRoutes.PUT("/:id/email-notifications", server.updateUserEmailNotifications)
		userRoutes.GET("/:id/flavours", server.getUserLikedFlavours)
		userRoutes.GET("/:id/cocktails", server.getUserLikedCocktails)
		userRoutes.GET("/:id/generated-cocktails", server.getUserGeneratedCocktails)
		userRoutes.GET("/:id/generate-requests", server.getIncompleteUserCocktailGenerationRequests)
		userRoutes.GET("/:id/categories", server.getCategoriesBasedOnLikedFlavours)
		userRoutes.GET("/:id/homescreen", server.getHomescreenForUser)
		userRoutes.GET("/:id/notifications", server.getUserNotifications)

		/**
		Cocktail Routes
		*/
		cocktailRoutes := authRoutes.Group("/cocktails")

		cocktailRoutes.GET("/", server.getCocktails)
		cocktailRoutes.GET("/:id", server.getCocktail)
		cocktailRoutes.GET("/:id/categories", server.getCategoriesForCocktail)
		cocktailRoutes.GET("/:id/simillar", server.getSimilarCocktails)

		/**
		Liked Cocktail Routes
		*/
		cocktailRoutes.POST("/:id/is_liked", server.getCocktailLikedStatus)
		cocktailRoutes.POST("/:id/like", server.likeCocktail)
		cocktailRoutes.POST("/:id/unlike", server.unlikeCocktail)

		/**
		Daily Featured Cocktail Routes
		*/
		cocktailRoutes.GET("/featured", server.getDailyFeatured)

		/**
		Category Routes
		*/
		categoryRoutes := authRoutes.Group("/categories")

		categoryRoutes.GET("/", server.getAllCategories)
		categoryRoutes.GET("/:id", server.getCategoryById)
		categoryRoutes.GET("/:id/cocktails", server.getCocktailsByCategory)

		/**
		Flavour Routes
		*/

		flavourRoutes := authRoutes.Group("/flavours")

		flavourRoutes.GET("/", server.getAllFlavours)
		flavourRoutes.GET("/:id", server.getFlavourId)
		/**
		Liked Flavour Routes
		*/
		flavourRoutes.GET("/:id/like", server.getLikedFlavour)
		flavourRoutes.POST("/like", server.likeFlavours)
		flavourRoutes.POST("/:id/like", server.likeFlavour)
		flavourRoutes.POST("/:id/unlike", server.unlikeFlavour)

		/**
		Notifications Routes
		*/
		notificationRoutes := authRoutes.Group("/notifications")

		notificationRoutes.POST("/", server.createNotification)
		notificationRoutes.PATCH("/:id", server.updateNotificationStatus)

		/**
		Generate Cocktail Routes
		*/
		authRoutes.GET("/generated/:id", server.getGeneratedCocktail)
		authRoutes.POST("/generate-cocktail", server.createGenerateCocktailRequest)

	}

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
