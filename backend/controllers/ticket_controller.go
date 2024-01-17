package controllers

import (

	// "github.com/gin-gonic/gin"
	"github.com/shuttlersit/service-desk/backend/models"
)

type TicketController struct {
	Ticket *models.TicketDBModel
}

func NewTicketController() *TicketController {
	return &TicketController{}
}

// Implement controller methods like GetTickets, CreateTicket, GetTicket, UpdateTicket, DeleteTicket
