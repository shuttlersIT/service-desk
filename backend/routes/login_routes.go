// backend/routes/auth_routes.go

package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shuttlersit/service-desk/controllers"
	"github.com/shuttlersit/service-desk/middleware"
)

func SetupAuthRoutes(router *gin.Engine, auths *controllers.AuthController) {

	// Public routes
	router.GET("/auth/google", controllers.GoogleLogin)
	//router.GET("/auth/google/callback", auths.GoogleAuthCallback)

	// Protected routes
	api := router.Group("/api")
	{
		api.Use(middleware.AuthMiddleware())
		api.GET("/protected", func(c *gin.Context) {
			// Example protected route
			c.JSON(http.StatusOK, gin.H{"message": "You are authenticated"})
		})
		//api.GET("/profile", controllers.UserProfile)
	}

	// Protect an admin route
	//adminRoutes := router.Group("/admin").Use(middleware.RequireRole("admin"))
	{
		//adminRoutes.GET("/dashboard", controllers.index)
	}

}

func SetupGoogleAuthRoutes(router *gin.Engine, auths *controllers.GoogleAuthMainController) {

	// Public routes
	router.GET("/auth/google", controllers.GoogleLogin)
	router.GET("/auth/google/callback", auths.CustomAuthHandler)

	// Protected routes
	api := router.Group("/api")
	{
		api.Use(middleware.AuthMiddleware())
		api.GET("/protected", func(c *gin.Context) {
			// Example protected route
			c.JSON(http.StatusOK, gin.H{"message": "You are authenticated"})
		})
		//api.GET("/profile", controllers.UserProfile)
	}

	// Protect an admin route
	//adminRoutes := router.Group("/admin").Use(middleware.RequireRole("admin"))
	{
		//adminRoutes.GET("/dashboard", controllers.index)
	}

}
