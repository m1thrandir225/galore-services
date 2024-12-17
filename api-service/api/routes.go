package api

import (
	"github.com/gin-gonic/gin"
)

func (server *Server) setupRouter() {
	router := gin.Default()

	router.Static("/public", "./public")

	v1 := router.Group("/api/v1")

	server.setupPublicRoutes(v1)
	server.setupAuthRoutes(v1)
	server.setupMigrationServiceRoutes(v1)
	server.setupCategorizerServiceRoutes(v1)

	server.router = router
}

func (server *Server) setupPublicRoutes(rg *gin.RouterGroup) {
	rg.POST("/register", server.registerUser)
	rg.POST("/login", server.loginUser)
	rg.POST("/refresh", server.refreshToken)
	rg.GET("/health", server.checkServiceHealth)
}

func (server *Server) setupMigrationServiceRoutes(rg *gin.RouterGroup) {
	migrationRoutes := rg.Group("/migration")
	migrationRoutes.Use(microserviceAccessMiddleware(server.config.MigrationServiceKey))

	migrationRoutes.POST("/cocktails", server.createCocktailMigrate)
}

func (server *Server) setupCategorizerServiceRoutes(rg *gin.RouterGroup) {
	categorizerRoutes := rg.Group("/categorizer")

	categorizerRoutes.Use(microserviceAccessMiddleware(server.config.CategorizerServiceKey))

	categorizerRoutes.GET("/category/:tag", server.getCategoryByTag)           // Check if category exists
	categorizerRoutes.POST("/category", server.createCategory)                 // Create a category
	categorizerRoutes.POST("/category_cocktail", server.addCocktailToCategory) // Add a cocktail to a given category
}

func (server *Server) setupAuthRoutes(rg *gin.RouterGroup) {
	authRoutes := rg.Group("/")
	authRoutes.Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/logout", server.logout)

	server.setupUserRoutes(authRoutes)
	server.setupCocktailRoutes(authRoutes)
	server.setupCategoryRoutes(authRoutes)
	server.setupFlavourRoutes(authRoutes)
	server.setupCategoryFlavourRoutes(authRoutes)
	server.setupNotificationTypesRoutes(authRoutes)
	server.setupNotificationRoutes(authRoutes)

}
