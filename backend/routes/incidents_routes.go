package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shuttlersit/service-desk/backend/controllers"
)

func SetupIncidentRoutes(router *gin.Engine, incidentController *controllers.IncidentController) {
	incidentRoutes := router.Group("/incidents")
	{
		incidentRoutes.POST("/", incidentController.CreateIncidentHandler)
		incidentRoutes.PUT("/:id", incidentController.UpdateIncidentHandler)
		incidentRoutes.GET("/:id", incidentController.GetIncidentByIDHandler)
		incidentRoutes.DELETE("/:id", incidentController.DeleteIncidentHandler)
		incidentRoutes.GET("/", incidentController.GetAllIncidentsHandler)
		incidentRoutes.GET("/:id/comments", incidentController.GetIncidentCommentsHandler)
		incidentRoutes.POST("/:id/comments", incidentController.AddIncidentCommentHandler)
		incidentRoutes.GET("/:id/history", incidentController.GetIncidentHistoryHandler)
	}
}
