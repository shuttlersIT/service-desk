// backend/services/advertisement_service.go

package services

import (
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

// DefaultAdvertisementService is the default implementation of AdvertisementService
type DefaultTicketingService struct {
	DB            *gorm.DB
	TicketDBModel *models.TicketDBModel
	// Add any dependencies or data needed for the service
}

// NewDefaultAdvertisementService creates a new DefaultAdvertisementService.
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
