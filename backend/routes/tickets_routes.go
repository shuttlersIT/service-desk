// backend/routes/tickets.go

package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shuttlersit/service-desk/backend/controllers"
)

func SetupTicketRoutes(router *gin.Engine, ticketController *controllers.TicketController) {

	// Define Ticket routes
	ticketRoutes := router.Group("/tickets")
	{
		ticketRoutes.POST("/", ticketController.CreateTicket)
		ticketRoutes.PUT("/:id", ticketController.UpdateTicket)
		ticketRoutes.GET("/:id", ticketController.GetTicketByIDHandler)
		ticketRoutes.DELETE("/:id", ticketController.DeleteTicketHandler)
		ticketRoutes.GET("/", ticketController.GetAllTicketsHandler)
		ticketRoutes.PUT("/:id/assign/:agentID", ticketController.AssignTicketToAgentHandler)
		ticketRoutes.PUT("/:id/change-status", ticketController.ChangeTicketStatusHandler)
		ticketRoutes.POST("/:id/comment", ticketController.AddCommentToTicketHandler)
		ticketRoutes.GET("/:id/history", ticketController.GetTicketHistoryHandler)
		ticketRoutes.POST("/categories", ticketController.CreateCategory)
		ticketRoutes.PUT("/categories/:id", ticketController.UpdateCategory)
		ticketRoutes.DELETE("/categories/:id", ticketController.DeleteCategory)
		ticketRoutes.POST("/subcategories", ticketController.CreateSubcategory)
		ticketRoutes.PUT("/subcategories/:id", ticketController.UpdateSubcategory)
		ticketRoutes.DELETE("/subcategories/:id", ticketController.DeleteSubcategory)
		ticketRoutes.POST("/:id/tags", ticketController.AddTagToTicket)
		ticketRoutes.DELETE("/:id/tags/:tag", ticketController.RemoveTagFromTicket)
		ticketRoutes.POST("/sla", ticketController.CreateSLA)
		ticketRoutes.PUT("/sla/:id", ticketController.UpdateSLA)
		ticketRoutes.DELETE("/sla/:id", ticketController.DeleteSLA)
		ticketRoutes.POST("/priorities", ticketController.CreatePriority)
		ticketRoutes.PUT("/priorities/:id", ticketController.UpdatePriority)
		ticketRoutes.DELETE("/priorities/:id", ticketController.DeletePriority)
		ticketRoutes.POST("/statuses", ticketController.CreateStatus)
		ticketRoutes.PUT("/statuses/:id", ticketController.UpdateStatus)
		ticketRoutes.DELETE("/statuses/:id", ticketController.DeleteStatus)
	}

}
