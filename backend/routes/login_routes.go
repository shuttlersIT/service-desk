// backend/routes/auth_routes.go

package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shuttlersit/service-desk/backend/controllers"
	"github.com/shuttlersit/service-desk/backend/middleware"
)

func SetupAuthRoutes(router *gin.Engine, auths *controllers.AuthController) {

	// Public routes
	router.GET("/auth/google", controllers.GoogleLogin)
	router.GET("/auth/google/callback", controllers.GoogleAuthCallback)

	// Protected routes
	api := router.Group("/api")
	{
		api.Use(middleware.AuthMiddleware())
		api.GET("/protected", func(c *gin.Context) {
			// Example protected route
			c.JSON(http.StatusOK, gin.H{"message": "You are authenticated"})
		})
	}
}
