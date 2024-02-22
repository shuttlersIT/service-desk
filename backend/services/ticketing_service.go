// backend/services/ticketing_service.go

package services

import (
	"fmt"

	"github.com/shuttlersit/service-desk/models"
	"gorm.io/gorm"
)

// TicketServiceInterface provides methods for managing ticketss.
type TicketingServiceInterface interface {
	CreateTicket(ticket *models.Ticket) error
	UpdateTicket(ticket *models.Ticket) error
	GetTicketByID(id uint) (*models.Ticket, error)
	DeleteTicket(ticketID uint) error
	GetAllTickets() ([]*models.Ticket, error)
	AssignTicketToAgent(ticketID, agentID uint) error
	ChangeTicketStatus(ticketID uint, newStatus models.Status) error
	AddCommentToTicket(ticketID uint, c string) error
	GetTicketHistory(ticketID uint) ([]*models.TicketHistoryEntry, error)
	CreateCategory(category *models.Category) error
	UpdateCategory(category *models.Category) error
	DeleteCategory(categoryID uint) error
	CreateSubcategory(subcategory *models.SubCategory) error
	UpdateSubcategory(subcategory *models.SubCategory) error
	DeleteSubcategory(subcategoryID uint) error
	CreateTag(ticketID uint, tag string) (*models.Tag, error)
	AddTagToTicket(ticketID uint, tag string) error
	IndirectlyAddTagToTicket(ticketID uint, tag string) error
	RemoveTagFromTicket(ticketID uint, tag string) error
	IndirectlyRemoveTagFromTicket(ticketID uint, tag string) error
	CreateSLA(sla *models.SLA) error
	UpdateSLA(sla *models.SLA) error
	DeleteSLA(slaID uint) error
	CreatePriority(priority *models.Priority) error
	UpdatePriority(priority *models.Priority) error
	DeletePriority(priorityID uint) error
	CreateStatus(status *models.Status) error
	UpdateStatus(status *models.Status) error
	DeleteStatus(statusID uint) error
}

// DefaultUserService is the default implementation of UserService
type DefaultTicketingService struct {
	DB                 *gorm.DB
	TicketDBModel      *models.TicketDBModel
	TicketComment      *models.TicketCommentDBModel
	TicketHistoryEntry *models.TicketHistoryEntryDBModel
	UserDBModel        *models.UserDBModel
	AgentDBModel       *models.AgentDBModel
	EventPublisher     *models.EventPublisherImpl
	log                *models.PrintLogger
	// Add any dependencies or data needed for the service
}

// NewDefaultUserService creates a new DefaultUserService.
func NewDefaultTicketingService(db *gorm.DB, ticketDBModel *models.TicketDBModel, ticketComment *models.TicketCommentDBModel, history *models.TicketHistoryEntryDBModel, userDBModel *models.UserDBModel, agentDBModel *models.AgentDBModel, eventPublisher *models.EventPublisherImpl, log *models.PrintLogger) *DefaultTicketingService {
	return &DefaultTicketingService{
		DB:                 db,
		TicketDBModel:      ticketDBModel,
		TicketComment:      ticketComment,
		TicketHistoryEntry: history,
		UserDBModel:        userDBModel,
		AgentDBModel:       agentDBModel,
		EventPublisher:     eventPublisher,
		log:                log,
	}
}

// GetAllTickets retrieves all tickets.
func (ps *DefaultTicketingService) GetAllTickets() ([]*models.Ticket, error) {
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
	ticket.AgentID = &a.ID

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
	ticket.Status = newStatus.Name

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

	newCommentID, er := ts.TicketComment.CreateTicketComment(ticketID, c)
	if er != nil {
		return er
	}

	newComment, err := ts.TicketComment.GetCommentByID(newCommentID)
	if err != nil {
		return err
	}

	// Add the comment to the ticket
	ticket.Comments = append(ticket.Comments, *newComment)

	// Save the updated ticket
	err = ts.TicketDBModel.UpdateTicket(ticket)
	if err != nil {
		return err
	}

	return nil
}

func (ts *DefaultTicketingService) GetTicketHistory(ticketID uint) ([]*models.TicketHistoryEntry, error) {
	// Retrieve the ticket history entries for the given ticketID
	historyEntries, err := ts.TicketHistoryEntry.GetHistoryEntriesByTicketID(ticketID)
	if historyEntries == nil {

		return nil, fmt.Errorf("could not find ticket history")
	} else if err != nil {
		return nil, fmt.Errorf("unable to get ticket history from db")
	}

	return historyEntries, nil
}

///////////////////////////////////////////////////////////////////////////////////////
// backend/services/category_service.go

func (cs *DefaultTicketingService) CreateCategory(category *models.Category) error {
	// Create a new category
	err := cs.TicketDBModel.CreateCategory(category)
	if err != nil {
		return err
	}

	return nil
}

func (cs *DefaultTicketingService) UpdateCategory(category *models.Category) error {
	// Update an existing category
	err := cs.TicketDBModel.UpdateCategory(category)
	if err != nil {
		return err
	}

	return nil
}

func (cs *DefaultTicketingService) DeleteCategory(categoryID uint) error {
	// Delete a category by categoryID
	err := cs.TicketDBModel.DeleteCategory(categoryID)
	if err != nil {
		return err
	}

	return nil
}

func (cs *DefaultTicketingService) CreateSubcategory(subcategory *models.SubCategory) error {
	// Create a new subcategory
	err := cs.TicketDBModel.CreateSubCategory(subcategory)
	if err != nil {
		return err
	}

	return nil
}

func (cs *DefaultTicketingService) UpdateSubcategory(subcategory *models.SubCategory) error {
	// Update an existing subcategory
	err := cs.TicketDBModel.UpdateSubCategory(subcategory)
	if err != nil {
		return err
	}

	return nil
}

func (cs *DefaultTicketingService) DeleteSubcategory(subcategoryID uint) error {
	// Delete a subcategory by subcategoryID
	err := cs.TicketDBModel.DeleteSubcategory(subcategoryID)
	if err != nil {
		return err
	}

	return nil
}

// Handle Tags

func (ts *DefaultTicketingService) CreateTag(ticketID uint, tag string) (*models.Tag, error) {
	Newtag, status, err := ts.TicketDBModel.CreateTag(ticketID, tag)
	if !status {
		return nil, err
	}

	//ticket, _ := ts.GetTicketByID(ticketID)
	return Newtag, nil
}

func (ts *DefaultTicketingService) CreateTag2(ticketID uint, tag *models.Tag) (*models.Tag, error) {

	Newtag, status, err := ts.TicketDBModel.CreateTag(ticketID, tag.Name)
	if !status {
		return nil, err
	}

	//ticket, _ := ts.GetTicketByID(ticketID)
	return Newtag, nil
}

func (ts *DefaultTicketingService) AddTagToTicket(ticketID uint, tag string) error {
	err := ts.TicketDBModel.AddTagToTicket(ticketID, tag)
	if err != nil {
		return err
	}
	return nil
}

func (ts *DefaultTicketingService) IndirectlyAddTagToTicket(ticketID uint, tag string) error {
	// Retrieve the ticket by ticketID
	ticket, err := ts.TicketDBModel.GetTicketByID(ticketID)
	if err != nil {
		return err
	}
	t, e := ts.CreateTag(ticketID, tag)
	if e != nil {
		return e
	}

	// Add the tag to the ticket's tags
	ticket.Tags = append(ticket.Tags, *t)

	// Save the updated ticket
	err = ts.TicketDBModel.UpdateTicket(ticket)
	if err != nil {
		return err
	}

	return nil
}

func (ts *DefaultTicketingService) RemoveTagFromTicket(ticketID uint, tag string) error {
	err := ts.TicketDBModel.RemoveTagFromTicket(ticketID, tag)
	if err != nil {
		return err
	}
	return nil
}

func (ts *DefaultTicketingService) IndirectlyRemoveTagFromTicket(ticketID uint, tag string) error {
	// Retrieve the ticket by ticketID
	ticket, err := ts.TicketDBModel.GetTicketByID(ticketID)
	if err != nil {
		return err
	}

	// Remove the tag from the ticket's tags
	for i, t := range ticket.Tags {
		if t.Name == tag {
			ticket.Tags = append(ticket.Tags[:i], ticket.Tags[i+1:]...)
			break
		}
	}

	// Save the updated ticket
	err = ts.TicketDBModel.UpdateTicket(ticket)
	if err != nil {
		return err
	}

	return nil
}

//

// backend/services/sla_service.go

func (ss *DefaultTicketingService) CreateSLA(sla *models.SLA) error {
	// Create a new SLA
	err := ss.TicketDBModel.CreateSla(sla)
	if err != nil {
		return err
	}

	return nil
}

func (ss *DefaultTicketingService) UpdateSLA(sla *models.SLA) error {
	// Update an existing SLA
	err := ss.TicketDBModel.UpdateSla(sla)
	if err != nil {
		return err
	}

	return nil
}

func (ss *DefaultTicketingService) DeleteSLA(slaID uint) error {
	// Delete an SLA by slaID
	err := ss.TicketDBModel.DeleteSla(slaID)
	if err != nil {
		return err
	}

	return nil
}

// backend/services/priority_service.go

func (ps *DefaultTicketingService) CreatePriority(priority *models.Priority) error {
	// Create a new priority level
	err := ps.TicketDBModel.CreatePriority(priority)
	if err != nil {
		return err
	}

	return nil
}

func (ps *DefaultTicketingService) UpdatePriority(priority *models.Priority) error {
	// Update an existing priority level
	err := ps.TicketDBModel.UpdatePriority(priority)
	if err != nil {
		return err
	}

	return nil
}

func (ps *DefaultTicketingService) DeletePriority(priorityID uint) error {
	// Delete a priority level by priorityID
	err := ps.TicketDBModel.DeletePriority(priorityID)
	if err != nil {
		return err
	}

	return nil
}

// backend/services/status_service.go

func (ss *DefaultTicketingService) CreateStatus(status *models.Status) error {
	// Create a new ticket status
	err := ss.TicketDBModel.CreateStatus(status)
	if err != nil {
		return err
	}

	return nil
}

func (ss *DefaultTicketingService) UpdateStatus(status *models.Status) error {
	// Update an existing ticket status
	err := ss.TicketDBModel.UpdateStatus(status)
	if err != nil {
		return err
	}

	return nil
}

func (ss *DefaultTicketingService) DeleteStatus(statusID uint) error {
	// Delete a ticket status by statusID
	err := ss.TicketDBModel.DeleteStatus(statusID)
	if err != nil {
		return err
	}

	return nil
}
