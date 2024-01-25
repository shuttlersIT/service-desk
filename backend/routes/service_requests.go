package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shuttlersit/service-desk/backend/controllers"
)

func SetupServiceRequestRoutes(r *gin.Engine, serviceRequestController *controllers.ServiceRequestController) {
	serviceRequestRoutes := r.Group("/service-requests")
	{
		serviceRequestRoutes.POST("/", serviceRequestController.CreateServiceRequestHandler)
		serviceRequestRoutes.PUT("/:id", serviceRequestController.UpdateServiceRequestHandler)
		serviceRequestRoutes.GET("/user/:user_id", serviceRequestController.GetUserRequests)
		serviceRequestRoutes.PUT("/close/:id", serviceRequestController.CloseServiceRequestHandler)
		serviceRequestRoutes.GET("/:id/history", serviceRequestController.GetServiceRequestHistoryHandler)
		serviceRequestRoutes.POST("/:id/comments", serviceRequestController.AddCommentToServiceRequestHandler)
		serviceRequestRoutes.GET("/:id/comments", serviceRequestController.GetServiceRequestCommentsHandler)
	}
}
