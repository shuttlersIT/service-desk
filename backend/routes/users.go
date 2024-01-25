// backend/routes/users.go

package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shuttlersit/service-desk/backend/controllers"
)

func SetUserRoutes(r *gin.Engine, user *controllers.UserController) {

	u := r.Group("/users")
	u.GET("/", user.GetAllUsersHandler)
	u.GET("/:id", user.GetUserByIDHandler)
	u.POST("/", user.CreateUserHandler)
	u.PUT("/:id", user.UpdateUserHandler)
	u.DELETE("/:id", user.DeleteUserHandler)

}
