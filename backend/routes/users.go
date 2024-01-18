package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shuttlersit/service-desk/backend/controllers"
)

func SetUserRoutes(r *gin.Engine, users *controllers.UserController) {

	u := r.Group("/users")
	u.GET("/", users.GetAllUsers)
	u.GET("/:id", users.GetUserByID)
	u.POST("/", users.CreateUser)
	u.PUT("/:id", users.UpdateUser)
	u.DELETE("/:id", users.DeleteUser)

}
