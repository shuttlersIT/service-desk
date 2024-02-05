// backend/routes/public.go

package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shuttlersit/service-desk/backend/controllers"
)

func SetupOpenRoutes(r *gin.Engine, public *controllers.AuthController) {

	p := r.Group("/")
	//p.GET("/", public.index)
	p.GET("/register", public.Register)
	p.GET("/login", public.Login)
	p.GET("/register/agent", public.RegisterAgent)
	p.GET("/login/agent", public.LoginAgent)
	//p.POST("logout", public.Logout)
	//publics.PUT("/support", public.UpdateAdvertisement)
	//publics.DELETE("/shuttlers-admin", public.DeleteAdvertisement)

}
