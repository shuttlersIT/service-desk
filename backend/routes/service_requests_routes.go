package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shuttlersit/service-desk/backend/controllers"
)

func SetupServiceRequestRoutes(router *gin.Engine, serviceRequestController *controllers.ServiceRequestController) {
	// Define Service Request routes
	serviceRequestRoutes := router.Group("/service-requests")
	{
		serviceRequestRoutes.POST("/", serviceRequestController.CreateServiceRequestHandler)
		serviceRequestRoutes.PUT("/:id", serviceRequestController.UpdateServiceRequestHandler)
		serviceRequestRoutes.GET("/:id", serviceRequestController.GetServiceRequestByIDHandler)
		serviceRequestRoutes.DELETE("/:id", serviceRequestController.DeleteServiceRequestHandler)
		serviceRequestRoutes.GET("/", serviceRequestController.GetAllServiceRequestsHandler)

	}
}
