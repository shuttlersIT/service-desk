// backend/routes/public.go

package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shuttlersit/service-desk/backend/controllers"
)

func SetOpenRoutes(r *gin.Engine, public *controllers.AuthController) {

	p := r.Group("/")
	//p.GET("/", public.index)
	p.GET("/login", public.Registration)
	p.GET("/login", public.Login)
	//p.POST("logout", public.Logout)
	//publics.PUT("/support", public.UpdateAdvertisement)
	//publics.DELETE("/shuttlers-admin", public.DeleteAdvertisement)

}
