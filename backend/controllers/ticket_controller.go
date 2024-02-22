// backend/controllers/ticket_controllers.go

package controllers

import (
	"net/http"
	"strconv"

	// "github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin"
	"github.com/shuttlersit/service-desk/models"
	"github.com/shuttlersit/service-desk/services"
)

type TicketController struct {
	TicketService *services.DefaultTicketingService
}

func NewTicketController(ticketService *services.DefaultTicketingService) *TicketController {
	return &TicketController{
		TicketService: ticketService,
	}
}

// Implement controller methods like GetTickets, CreateTicket, GetTicket, UpdateTicket, DeleteTicket, GetAllTickets

// CreateTicket handles the HTTP request to create a new ticket.
func (tc *TicketController) CreateTicket(c *gin.Context) {
	var newTicket models.Ticket
	if err := c.ShouldBindJSON(&newTicket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err := tc.TicketService.CreateTicket(&newTicket)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ticket"})
		return
	}

	c.JSON(http.StatusCreated, newTicket)
}

// UpdateTicket handles the HTTP request to update an existing ticket.
func (tc *TicketController) UpdateTicket(c *gin.Context) {
	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	var updatedTicket models.Ticket
	if err := c.ShouldBindJSON(&updatedTicket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	updatedTicket.ID = uint(ticketID)

	uTicket, err := tc.TicketService.UpdateTicket(&updatedTicket)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ticket"})
		return
	}

	c.JSON(http.StatusOK, uTicket)
}

// GetTicketByID handles the HTTP request to retrieve a ticket by its ID.
func (tc *TicketController) GetTicketByID(c *gin.Context) {
	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	ticket, err := tc.TicketService.GetTicketByID(uint(ticketID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

// DeleteTicket handles the HTTP request to delete a ticket by its ID.
func (tc *TicketController) DeleteTicket(c *gin.Context) {
	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	_, err = tc.TicketService.DeleteTicket(uint(ticketID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete ticket"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket deleted successfully"})
}

// GetAllTickets handles the HTTP request to retrieve all tickets.
func (tc *TicketController) GetAllTickets(c *gin.Context) {
	tickets, err := tc.TicketService.GetAllTickets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tickets"})
		return
	}

	c.JSON(http.StatusOK, tickets)
}

// AssignTicketToAgent handles the HTTP request to assign a ticket to an agent.
func (tc *TicketController) AssignTicketToAgent(c *gin.Context) {
	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	agentID, err := strconv.ParseUint(c.Param("agent_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid agent ID"})
		return
	}

	err = tc.TicketService.AssignTicketToAgent(uint(ticketID), uint(agentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign ticket to agent"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket assigned to agent successfully"})
}

// ChangeTicketStatus handles the HTTP request to change the status of a ticket.
func (tc *TicketController) ChangeTicketStatus(c *gin.Context) {
	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	var status models.Status
	if err := c.ShouldBindJSON(&status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err = tc.TicketService.ChangeTicketStatus(uint(ticketID), status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to change ticket status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket status changed successfully"})
}

// AddCommentToTicket handles the HTTP request to add a comment to a ticket.
func (tc *TicketController) AddCommentToTicket(c *gin.Context) {
	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err = tc.TicketService.AddCommentToTicket(uint(ticketID), comment.Comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add comment to ticket"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Comment added to ticket successfully"})
}

// GetTicketHistory handles the HTTP request to retrieve the history of a ticket.
func (tc *TicketController) GetTicketHistory(c *gin.Context) {
	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	history, err := tc.TicketService.GetTicketHistory(uint(ticketID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve ticket history"})
		return
	}

	c.JSON(http.StatusOK, history)
}

// CreateTag handles the HTTP request to create a new ticket tag.
func (tc *TicketController) CreateTag2(c *gin.Context) {
	ticketID, _ := strconv.Atoi(c.Param("ticket_id"))
	var newTag models.Tag
	if err := c.ShouldBindJSON(&newTag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	tag, err := tc.TicketService.CreateTag(uint(ticketID), newTag.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tag"})
		return
	}

	c.JSON(http.StatusCreated, tag)
}

// CreateTag handles the HTTP request to create a new ticket tag.
func (tc *TicketController) CreateTag1(c *gin.Context) {
	ticketID, _ := strconv.Atoi(c.Param("ticket_id"))
	tag := c.Param("ticket_id")
	//var newTag models.Tags
	//if err := c.ShouldBindJSON(&newTag); err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
	//	return
	//}

	nTag, err := tc.TicketService.CreateTag(uint(ticketID), tag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tag"})
		return
	}

	c.JSON(http.StatusCreated, nTag)
}

// AddTagToTicket handles the HTTP request to add a tag to a ticket.
func (tc *TicketController) AddTagToTicket(c *gin.Context) {
	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	tag := c.Param("tag")

	err = tc.TicketService.AddTagToTicket(uint(ticketID), tag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add tag to ticket"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag added to ticket successfully"})
}

// IndirectlyAddTagToTicket handles the HTTP request to indirectly add a tag to a ticket.
func (tc *TicketController) IndirectlyAddTagToTicket(c *gin.Context) {
	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	tag := c.Param("tag")

	err = tc.TicketService.IndirectlyAddTagToTicket(uint(ticketID), tag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to indirectly add tag to ticket"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag indirectly added to ticket successfully"})
}

// RemoveTagFromTicket handles the HTTP request to remove a tag from a ticket.
func (tc *TicketController) RemoveTagFromTicket(c *gin.Context) {
	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	tag := c.Param("tag")

	err = tc.TicketService.RemoveTagFromTicket(uint(ticketID), tag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove tag from ticket"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag removed from ticket successfully"})
}

// IndirectlyRemoveTagFromTicket handles the HTTP request to indirectly remove a tag from a ticket.
func (tc *TicketController) IndirectlyRemoveTagFromTicket(c *gin.Context) {
	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	tag := c.Param("tag")

	err = tc.TicketService.IndirectlyRemoveTagFromTicket(uint(ticketID), tag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to indirectly remove tag from ticket"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag indirectly removed from ticket successfully"})
}

// CreateSLA handles the HTTP request to create a new Service Level Agreement (SLA).
func (tc *TicketController) CreateSLA(c *gin.Context) {
	var newSLA models.SLA
	if err := c.ShouldBindJSON(&newSLA); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err := tc.TicketService.CreateSLA(&newSLA)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create SLA"})
		return
	}

	c.JSON(http.StatusCreated, newSLA)
}

// UpdateSLA handles the HTTP request to update an existing Service Level Agreement (SLA).
func (tc *TicketController) UpdateSLA(c *gin.Context) {
	slaID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SLA ID"})
		return
	}

	var updatedSLA models.SLA
	if err := c.ShouldBindJSON(&updatedSLA); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	updatedSLA.ID = uint(slaID)

	err = tc.TicketService.UpdateSLA(&updatedSLA)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update SLA"})
		return
	}

	c.JSON(http.StatusOK, updatedSLA)
}

// DeleteSLA handles the HTTP request to delete a Service Level Agreement (SLA) by its ID.
func (tc *TicketController) DeleteSLA(c *gin.Context) {
	slaID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SLA ID"})
		return
	}

	err = tc.TicketService.DeleteSLA(uint(slaID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete SLA"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SLA deleted successfully"})
}

// CreatePriority handles the HTTP request to create a new ticket priority.
func (tc *TicketController) CreatePriority(c *gin.Context) {
	var newPriority models.Priority
	if err := c.ShouldBindJSON(&newPriority); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err := tc.TicketService.CreatePriority(&newPriority)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create priority"})
		return
	}

	c.JSON(http.StatusCreated, newPriority)
}

// UpdatePriority handles the HTTP request to update an existing ticket priority.
func (tc *TicketController) UpdatePriority(c *gin.Context) {
	priorityID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid priority ID"})
		return
	}

	var updatedPriority models.Priority
	if err := c.ShouldBindJSON(&updatedPriority); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	updatedPriority.ID = uint(priorityID)

	err = tc.TicketService.UpdatePriority(&updatedPriority)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update priority"})
		return
	}

	c.JSON(http.StatusOK, updatedPriority)
}

// DeletePriority handles the HTTP request to delete a ticket priority by its ID.
func (tc *TicketController) DeletePriority(c *gin.Context) {
	priorityID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid priority ID"})
		return
	}

	err = tc.TicketService.DeletePriority(uint(priorityID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete priority"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Priority deleted successfully"})
}

// CreateStatus handles the HTTP request to create a new ticket status.
func (tc *TicketController) CreateStatus(c *gin.Context) {
	var newStatus models.Status
	if err := c.ShouldBindJSON(&newStatus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err := tc.TicketService.CreateStatus(&newStatus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create status"})
		return
	}

	c.JSON(http.StatusCreated, newStatus)
}

// UpdateStatus handles the HTTP request to update an existing ticket status.
func (tc *TicketController) UpdateStatus(c *gin.Context) {
	statusID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status ID"})
		return
	}

	var updatedStatus models.Status
	if err := c.ShouldBindJSON(&updatedStatus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	updatedStatus.ID = uint(statusID)

	err = tc.TicketService.UpdateStatus(&updatedStatus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}

	c.JSON(http.StatusOK, updatedStatus)
}

// DeleteStatus handles the HTTP request to delete a ticket status by its ID.
func (tc *TicketController) DeleteStatus(c *gin.Context) {
	statusID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status ID"})
		return
	}

	err = tc.TicketService.DeleteStatus(uint(statusID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status deleted successfully"})
}

// CreateCategory handles the HTTP request to create a new ticket category.
func (tc *TicketController) CreateCategory(c *gin.Context) {
	var newCategory models.Category
	if err := c.ShouldBindJSON(&newCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err := tc.TicketService.CreateCategory(&newCategory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, newCategory)
}

// UpdateCategory handles the HTTP request to update an existing ticket category.
func (tc *TicketController) UpdateCategory(c *gin.Context) {
	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var updatedCategory models.Category
	if err := c.ShouldBindJSON(&updatedCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	updatedCategory.ID = uint(categoryID)

	err = tc.TicketService.UpdateCategory(&updatedCategory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	c.JSON(http.StatusOK, updatedCategory)
}

// DeleteCategory handles the HTTP request to delete a ticket category by its ID.
func (tc *TicketController) DeleteCategory(c *gin.Context) {
	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	err = tc.TicketService.DeleteCategory(uint(categoryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}

// CreateSubcategory handles the HTTP request to create a new ticket subcategory.
func (tc *TicketController) CreateSubcategory(c *gin.Context) {
	var newSubcategory models.SubCategory
	if err := c.ShouldBindJSON(&newSubcategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err := tc.TicketService.CreateSubcategory(&newSubcategory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subcategory"})
		return
	}

	c.JSON(http.StatusCreated, newSubcategory)
}

// UpdateSubcategory handles the HTTP request to update an existing ticket subcategory.
func (tc *TicketController) UpdateSubcategory(c *gin.Context) {
	subcategoryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subcategory ID"})
		return
	}

	var updatedSubcategory models.SubCategory
	if err := c.ShouldBindJSON(&updatedSubcategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	updatedSubcategory.ID = uint(subcategoryID)

	err = tc.TicketService.UpdateSubcategory(&updatedSubcategory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update subcategory"})
		return
	}

	c.JSON(http.StatusOK, updatedSubcategory)
}

// DeleteSubcategory handles the HTTP request to delete a ticket subcategory by its ID.
func (tc *TicketController) DeleteSubcategory(c *gin.Context) {
	subcategoryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subcategory ID"})
		return
	}

	err = tc.TicketService.DeleteSubcategory(uint(subcategoryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete subcategory"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subcategory deleted successfully"})
}

// CreateTag handles the HTTP request to create a new ticket tag.
func (tc *TicketController) CreateTag(c *gin.Context) {
	ticketID, err := strconv.ParseUint(c.Param("ticketID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	var newTag models.Tag
	if err := c.ShouldBindJSON(&newTag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Assuming you have a ticket service method to create a tag for a specific ticket.
	tag, err := tc.TicketService.CreateTag(uint(ticketID), newTag.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tag"})
		return
	}

	c.JSON(http.StatusCreated, tag)
}

// GetTicketByIDHandler retrieves a ticket by its ID.
func (tc *TicketController) GetTicketByIDHandler(c *gin.Context) {
	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	// Assuming you have a ticket service method to get a ticket by its ID.
	ticket, err := tc.TicketService.GetTicketByID(uint(ticketID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

// GetTicketHistoryHandler retrieves the history of a ticket.
func (tc *TicketController) GetTicketHistoryHandler(c *gin.Context) {
	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	// Assuming you have a ticket service method to get the history of a ticket.
	history, err := tc.TicketService.GetTicketHistory(uint(ticketID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket history not found"})
		return
	}

	c.JSON(http.StatusOK, history)
}

// GetAllTicketsHandler retrieves all tickets.
func (tc *TicketController) GetAllTicketsHandler(c *gin.Context) {
	// Assuming you have a ticket service method to get all tickets.
	tickets, err := tc.TicketService.GetAllTickets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tickets"})
		return
	}

	c.JSON(http.StatusOK, tickets)
}

// DeleteTicketHandler deletes a ticket by its ID.
func (tc *TicketController) DeleteTicketHandler(c *gin.Context) {
	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	status, err := tc.TicketService.DeleteTicket(uint(ticketID))
	if err != nil || !status {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete ticket"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket deleted successfully"})
}

// AssignTicketToAgentHandler assigns a ticket to a specific agent.
func (tc *TicketController) AssignTicketToAgentHandler(c *gin.Context) {
	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	agentID, err := strconv.ParseUint(c.Param("agent_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid agent ID"})
		return
	}

	// Assuming you have a ticket service method to assign a ticket to an agent.
	err = tc.TicketService.AssignTicketToAgent(uint(ticketID), uint(agentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign ticket to agent"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket assigned to agent successfully"})
}

// ChangeTicketStatusHandler changes the status of a ticket.
func (tc *TicketController) ChangeTicketStatusHandler(c *gin.Context) {
	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	var newStatus models.Status
	if err := c.ShouldBindJSON(&newStatus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Assuming you have a ticket service method to change the status of a ticket.
	err = tc.TicketService.ChangeTicketStatus(uint(ticketID), newStatus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to change ticket status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket status changed successfully"})
}

// AddCommentToTicketHandler adds a comment to a ticket.
func (tc *TicketController) AddCommentToTicketHandler(c *gin.Context) {
	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	var comment string
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Assuming you have a ticket service method to add a comment to a ticket.
	err = tc.TicketService.AddCommentToTicket(uint(ticketID), comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add comment to ticket"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment added to ticket successfully"})
}

// GetTicketByIDHandler retrieves a ticket by its ID.
func (tc *TicketController) GetTicketByIDHandler2(c *gin.Context) {
	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	ticket, err := tc.TicketService.GetTicketByID(uint(ticketID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

// DeleteTicketHandler deletes a ticket by its ID.
func (tc *TicketController) DeleteTicketHandler2(c *gin.Context) {
	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	status, err := tc.TicketService.DeleteTicket(uint(ticketID))
	if err != nil || !status {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket deleted successfully"})
}

// GetAllTicketsHandler retrieves all tickets.
func (tc *TicketController) GetAllTicketsHandler2(c *gin.Context) {
	tickets, err := tc.TicketService.GetAllTickets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tickets"})
		return
	}

	c.JSON(http.StatusOK, tickets)
}
