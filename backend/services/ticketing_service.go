// backend/services/ticketing_service.go

package services

import (
	"fmt"

	"github.com/shuttlersit/service-desk/backend/models"
	"gorm.io/gorm"
)

// TicketServiceInterface provides methods for managing ticketss.
type TicketingServiceInterface interface {
	CreateTicket(ticket *models.Ticket) (*models.Ticket, error)
	UpdateTicket(ticket *models.Ticket) (*models.Ticket, error)
	GetTicketByID(id uint) (*models.Ticket, error)
	DeleteTicket(ticketID uint) (bool, error)
	GetAllTickets() *[]models.Ticket
}

// DefaultUserService is the default implementation of UserService
type DefaultTicketingService struct {
	DB                 *gorm.DB
	TicketDBModel      *models.TicketDBModel
	TicketComment      *models.TicketCommentDBModel
	TicketHistoryEntry *models.TicketHistoryEntryDBModel
	UserDBModel        *models.UserDBModel
	AgentDBModel       *models.AgentDBModel
	// Add any dependencies or data needed for the service
}

// NewDefaultUserService creates a new DefaultUserService.
func NewDefaultTicketingService(ticketDBModel *models.TicketDBModel) *DefaultTicketingService {
	return &DefaultTicketingService{
		TicketDBModel: ticketDBModel,
	}
}

// GetAllTickets retrieves all tickets.
func (ps *DefaultTicketingService) GetAllTickets() (*[]models.Ticket, error) {
	tickets, err := ps.TicketDBModel.GetAllTickets()
	if err != nil {
		return nil, err
	}
	return tickets, nil
}

// CreateTicket creates a new Ticket.
func (ps *DefaultTicketingService) CreateTicket(ticket *models.Ticket) error {
	err := ps.TicketDBModel.CreateTicket(ticket)
	if err != nil {
		return err
	}
	return nil
}

// CreateUser creates a new Ticket.
func (ps *DefaultTicketingService) GetTicketByID(id uint) (*models.Ticket, error) {
	ticket, err := ps.TicketDBModel.GetTicketByID(id)
	if err != nil {
		return nil, err
	}
	return ticket, nil
}

// UpdateTicket updates an existing Ticket.
func (ps *DefaultTicketingService) UpdateTicket(ticket *models.Ticket) (*models.Ticket, error) {
	err := ps.TicketDBModel.UpdateTicket(ticket)
	if err != nil {
		return nil, err
	}
	return ticket, nil
}

// DeleteTicket deletes an ticket by ID.
func (ps *DefaultTicketingService) DeleteTicket(ticketID uint) (bool, error) {
	status := false
	err := ps.TicketDBModel.DeleteTicket(ticketID)
	if err != nil {
		return status, err
	}
	status = true
	return status, nil
}

// /////////////////////////////////////////////////////////////////////////////
func (ts *DefaultTicketingService) AssignTicketToAgent(ticketID, agentID uint) error {
	// Retrieve the ticket by ticketID
	ticket, err := ts.TicketDBModel.GetTicketByID(ticketID)
	if err != nil {
		return err
	}

	// Update the assigned agent ID
	a, _ := ts.AgentDBModel.GetAgentByID(agentID)
	ticket.AgentID = *a

	// Save the updated ticket
	err = ts.TicketDBModel.UpdateTicket(ticket)
	if err != nil {
		return err
	}

	return nil
}

func (ts *DefaultTicketingService) ChangeTicketStatus(ticketID uint, newStatus models.Status) error {
	// Retrieve the ticket by ticketID
	ticket, err := ts.TicketDBModel.GetTicketByID(ticketID)
	if err != nil {
		return err
	}

	// Update the ticket status
	ticket.Status = newStatus

	// Save the updated ticket
	err = ts.TicketDBModel.UpdateTicket(ticket)
	if err != nil {
		return err
	}

	return nil
}

func (ts *DefaultTicketingService) AddCommentToTicket(ticketID uint, c string) error {
	// Retrieve the ticket by ticketID
	ticket, err := ts.TicketDBModel.GetTicketByID(ticketID)
	if err != nil {
		return err
	}
	comment, er := ts.TicketComment.CreateTicketComment(ticketID, c)
	if er != nil {
		return er
	}

	// Add the comment to the ticket
	ticket.Comments = append(ticket.Comments, *comment)

	// Save the updated ticket
	err = ts.TicketDBModel.UpdateTicket(ticket)
	if err != nil {
		return err
	}

	return nil
}

func (ts *DefaultTicketingService) GetTicketHistory(ticketID uint) (*[]models.TicketHistoryEntry, error) {
	// Retrieve the ticket history entries for the given ticketID
	historyEntries := ts.TicketHistoryEntry.GetHistoryEntriesByTicketID(ticketID)
	if historyEntries == nil {

		return nil, fmt.Errorf("could not find ticket history")
	}

	return historyEntries, nil
}
