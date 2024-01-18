package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shuttlersit/service-desk/backend/controllers"
)

func SetTicketRoutes(r *gin.Engine, tickets *controllers.TicketController) {

	t := r.Group("/tickets")
	t.GET("/", tickets.GetAllTickets)
	t.GET("/:id", tickets.GetTicketByID)
	t.POST("/", tickets.CreateTicket)
	t.PUT("/:id", tickets.UpdateTicket)
	t.DELETE("/:id", tickets.DeleteTicket)

}
