// backend/controllers/ticket_controllers.go

package controllers

import (
	"net/http"
	"strconv"

	// "github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin"
	"github.com/shuttlersit/service-desk/backend/models"
	"github.com/shuttlersit/service-desk/backend/services"
)

type TicketController struct {
	TicketService *services.DefaultTicketingService
}

func NewTicketController() *TicketController {
	return &TicketController{}
}

// Implement controller methods like GetTickets, CreateTicket, GetTicket, UpdateTicket, DeleteTicket, GetAllTickets

func (tc *TicketController) CreateTicket(ctx *gin.Context) {
	var newTicket models.Ticket
	if err := ctx.ShouldBindJSON(&newTicket); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err := tc.TicketService.CreateTicket(&newTicket)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Ticket"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Ticket created successfully"})
}

// GetTicketByID handles the HTTP request to retrieve a user by ID.
func (pc *TicketController) GetTicketByID(ctx *gin.Context) {
	ticketID, _ := strconv.Atoi(ctx.Param("id"))
	ticket, err := pc.TicketService.GetTicketByID(uint(ticketID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}
	ctx.JSON(http.StatusOK, ticket)
}

// UpdateTicket handles PUT /ticket/:id route.
func (pc *TicketController) UpdateTicket(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var ad models.Ticket
	if err := ctx.ShouldBindJSON(&ad); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ad.ID = uint(id)

	updatedAd, err := pc.TicketService.UpdateTicket(&ad)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedAd)
}

// DeleteTicket handles DELETE /Ticket/:id route.
func (pc *TicketController) DeleteTicket(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	status, err := pc.TicketService.DeleteTicket(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, status)
}

// GetAllTickets handles the HTTP request to retrieve a agents by ID.
func (pc *TicketController) GetAllTickets(ctx *gin.Context) {
	tickets, err := pc.TicketService.GetAllTickets()
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "tickets not found"})
		return
	}
	ctx.JSON(http.StatusOK, tickets)
}
