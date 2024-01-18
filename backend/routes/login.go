// backend/routes/auth_routes.go

package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shuttlersit/service-desk/backend/controllers"
)

func SetAuthRoutes(r *gin.Engine, auths *controllers.AuthController) {

	//auth := r.Group("/auths")
	//auth.GET("/", auth.GetAllAdvertisements)
	//auth.GET("/:id", auth.GetAdvertisementByID)
	//auth.POST("/", auth.CreateAdvertisement)
	//auth.PUT("/:id", auth.UpdateAdvertisement)
	//auth.DELETE("/:id", auth.DeleteAdvertisement)

}
