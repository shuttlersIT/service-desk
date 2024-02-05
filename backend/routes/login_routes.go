// backend/routes/auth_routes.go

package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shuttlersit/service-desk/backend/controllers"
)

func SetupAuthRoutes(r *gin.Engine, auths *controllers.AuthController) {

	//authG := r.Group("/auth")
	//authG.POST("/auth/register", auths.GoogleRegister)
	//authG.GET("/auth/login", auths.GoogleLogin)
	//authG.GET("/auth/callback", auths.GoogleCallback)
	//auth.PUT("/:id", auth.UpdateAdvertisement)
	//auth.DELETE("/:id", auth.DeleteAdvertisement)

}
