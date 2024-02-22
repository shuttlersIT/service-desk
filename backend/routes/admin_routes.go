// backend/routes/auth_routes.go

package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shuttlersit/service-desk/backend/controllers"
	"github.com/shuttlersit/service-desk/backend/middleware"
)

func SetupAdminRoutes(router *gin.Engine, auths *controllers.AuthController) {

	// Protected routes
	admin := router.Group("/admin")
	admin.Use(middleware.RequireRole("admin")) // Only users with the "admin" role can access routes in this group
	{
		admin.Use(middleware.AuthMiddleware())
		admin.GET("/protected", func(c *gin.Context) {
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
