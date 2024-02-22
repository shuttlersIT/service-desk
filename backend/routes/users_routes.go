// backend/routes/users.go

package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shuttlersit/service-desk/controllers"
)

func SetupUserRoutes(router *gin.Engine, userController *controllers.UserController) {

	userRoutes := router.Group("/users")
	{
		userRoutes.POST("/", userController.CreateUserHandler)
		userRoutes.PUT("/:id", userController.UpdateUserHandler)
		userRoutes.GET("/:id", userController.GetUserByIDHandler)
		userRoutes.DELETE("/:id", userController.DeleteUserHandler)
		userRoutes.GET("/", userController.GetAllUsersHandler)
		userRoutes.POST("/positions", userController.CreatePositionHandler)
		userRoutes.PUT("/positions/:id", userController.UpdatePositionHandler)
		userRoutes.DELETE("/positions/:id", userController.DeletePositionHandler)
		userRoutes.GET("/positions", userController.GetPositionsHandler)
		userRoutes.GET("/positions/:number", userController.GetPositionByNumberHandler)
		userRoutes.POST("/departments", userController.CreateDepartmentHandler)
		userRoutes.PUT("/departments/:id", userController.UpdateDepartmentHandler)
		userRoutes.DELETE("/departments/:id", userController.DeleteDepartmentHandler)
		userRoutes.GET("/departments", userController.GetDepartmentsHandler)
		userRoutes.GET("/departments/:number", userController.GetDepartmentByNumberHandler)
	}

	// Define Position routes
	positionRoutes := router.Group("/positions")
	{
		positionRoutes.POST("/", userController.CreatePositionHandler)
		positionRoutes.PUT("/:id", userController.UpdatePositionHandler)
		positionRoutes.GET("/:id", userController.GetPositionByIDHandler)
		positionRoutes.DELETE("/:id", userController.DeletePositionHandler)
		positionRoutes.GET("/", userController.GetPositionsHandler)
	}

}
