package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shuttlersit/service-desk/backend/controllers"
)

func SetupServiceRequestRoutes(r *gin.Engine, serviceRequestController *controllers.ServiceRequestController) {
	serviceRequestRoutes := r.Group("/service-requests")
	{
		serviceRequestRoutes.POST("/", serviceRequestController.CreateServiceRequest)
		serviceRequestRoutes.PUT("/:id", serviceRequestController.UpdateServiceRequest)
		serviceRequestRoutes.GET("/user/:user_id", serviceRequestController.GetUserRequests)
		serviceRequestRoutes.PUT("/close/:id", serviceRequestController.CloseServiceRequest)
		serviceRequestRoutes.GET("/:id/history", serviceRequestController.GetServiceRequestHistory)
		serviceRequestRoutes.POST("/:id/comments", serviceRequestController.AddCommentToServiceRequest)
		serviceRequestRoutes.GET("/:id/comments", serviceRequestController.GetServiceRequestComments)
	}
}
