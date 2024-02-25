// backend/models/tickets.go

package models

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type Ticket struct {
	gorm.Model
	ID               uint                    `gorm:"primaryKey" json:"id"`
	Subject          string                  `gorm:"size:255;not null" json:"subject"`                   // Summary of the ticket issue
	Description      string                  `gorm:"type:text;not null" json:"description"`              // Detailed description of the issue
	CategoryID       uint                    `gorm:"index;not null" json:"category_id"`                  // Categorizes the ticket for routing or reporting
	SubCategoryID    uint                    `gorm:"index;not null" json:"sub_category_id"`              // Further refines the ticket category
	PriorityID       uint                    `gorm:"index;not null" json:"priority_id"`                  // Indicates the urgency of the ticket
	SLAID            uint                    `json:"sla_id"`                                             // Associates the ticket with a specific Service Level Agreement
	UserID           uint                    `gorm:"index;not null" json:"user_id"`                      // The user who submitted the ticket
	AgentID          *uint                   `gorm:"index" json:"agent_id,omitempty"`                    // Optionally assigns the ticket to a specific agent
	DueAt            *time.Time              `json:"due_at,omitempty"`                                   // Expected resolution time
	ClosedAt         *time.Time              `json:"closed_at,omitempty"`                                // Time when the ticket was closed
	Site             string                  `gorm:"size:255" json:"site"`                               // Location or site related to the ticket
	StatusID         uint                    `gorm:"index;not null" json:"status_id"`                    // Current status of the ticket
	Status           string                  `gorm:"size:100;not null" json:"status" binding:"required"` // Descriptive status
	StatusObject     Status                  `json:"status_details"`                                     // Embeds status details
	Priority         Priority                `json:"priority"`                                           // Embeds priority details
	Category         Category                `gorm:"foreignKey:CategoryID" json:"-"`                     // Links to the category entity
	SubCategory      SubCategory             `gorm:"foreignKey:SubCategoryID" json:"-"`                  // Links to the sub-category entity
	SLA              SLA                     `gorm:"foreignKey:sla_id" json:"-"`                         // Links to the SLA entity
	MediaAttachments []TicketMediaAttachment `gorm:"foreignKey:TicketID" json:"media_attachments"`       // Related media files
	Tags             []Tag                   `gorm:"many2many:ticket_tags;" json:"tags"`                 // Tags for categorization or filtering
	Comments         []Comment               `gorm:"foreignKey:TicketID" json:"comments"`                // User or agent comments on the ticket
	TicketHistory    []TicketHistoryEntry    `gorm:"foreignKey:TicketID" json:"ticket_history"`          // Audit trail of ticket changes
	User             Users                   `gorm:"foreignKey:UserID" json:"-"`                         // Links to the submitting user
	Agent            Agents                  `gorm:"foreignKey:AgentID" json:"-"`                        // Links to the assigned agent, if any
}

func (Ticket) TableName() string {
	return "tickets"
}

type Comment struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	TicketID  uint            `json:"ticket_id"`                         // Links comment to a specific ticket
	AuthorID  uint            `json:"author_id"`                         // Identifies the author of the comment
	Comment   string          `json:"comment" gorm:"type:text;not null"` // The content of the comment
	CreatedAt time.Time       `json:"created_at"`                        // Timestamp of comment creation
	UpdatedAt time.Time       `json:"updated_at"`                        // Timestamp of last update to comment
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // Soft delete flag
}

func (Comment) TableName() string {
	return "comments"
}

type TicketHistoryEntry struct {
	gorm.Model
	TicketID uint   `json:"ticket_id"`                       // Associates history entry with a ticket
	Action   string `gorm:"size:255;not null" json:"action"` // Describes the action taken
}

func (TicketHistoryEntry) TableName() string {
	return "ticket_history_entries"
}

// RelatedAd struct for storing related advertisements
type RelatedTicket struct {
	gorm.Model
	TicketID        uint `json:"ticket_id" gorm:"index;not null"`
	RelatedTicketID uint `json:"related_ticket_id" gorm:"index;not null"`
}

func (RelatedTicket) TableName() string {
	return "related_tickets"
}

// Hashtag represents a hashtag entity
type Tag struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	Name      string          `gorm:"size:255;not null;unique" json:"name"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Tag) TableName() string {
	return "tags"
}

type SLA struct {
	ID            uint            `gorm:"primaryKey" json:"id"`
	Name          string          `gorm:"size:255;not null" json:"name"`
	Description   string          `gorm:"type:text" json:"description,omitempty"`
	Target        int             `json:"target" gorm:"type:int;not null"`
	TimeToResolve int             `json:"time_to_resolve" gorm:"type:int;not null"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	DeletedAt     *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (SLA) TableName() string {
	return "slas"
}

type Priority struct {
	gorm.Model
	Name        string `gorm:"size:255;not null" json:"name"`
	Level       int    `json:"level" gorm:"type:int;not null"`
	Description string `gorm:"type:text" json:"description,omitempty"`
	Color       string `gorm:"size:7;default:'#FFFFFF'" json:"color"`
}

func (Priority) TableName() string {
	return "priorities"
}

type Satisfaction struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	Name      string          `gorm:"size:255;not null" json:"name"`
	Rank      int             `json:"rank" gorm:"type:int;not null"`
	Emoji     string          `json:"emoji,omitempty" gorm:"size:255"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName sets the table name for the Satisfaction model.
func (Satisfaction) TableName() string {
	return "satisfaction"
}

type Category struct {
	ID               uint            `gorm:"primaryKey" json:"id"`
	Name             string          `gorm:"size:255;not null;unique" json:"name"`
	Description      string          `gorm:"type:text" json:"description,omitempty"`
	SubCategories    []*Category     `gorm:"foreignKey:ParentCategoryID" json:"sub_categories,omitempty"`
	Icon             string          `gorm:"size:255" json:"icon,omitempty"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
	DeletedAt        *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	ParentCategoryID *uint           `json:"parent_category_id,omitempty" gorm:"index"`
}

func (Category) TableName() string {
	return "categories"
}

type SubCategory struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"size:255;not null" json:"name"`
	CategoryID  uint   `json:"category_id" gorm:"index;not null"`
	Description string `gorm:"type:text" json:"description,omitempty"`
	Icon        string `gorm:"size:255" json:"icon,omitempty"`
}

func (SubCategory) TableName() string {
	return "sub_categories"
}

type Status struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"size:255;not null" json:"name"`
	Description string `gorm:"type:text" json:"description,omitempty"`
	IsClosed    bool   `json:"is_closed" gorm:"not null;default:false"`
}

func (Status) TableName() string {
	return "statuses"
}

type Policies struct {
	ID           uint            `gorm:"primaryKey" json:"id"`
	Name         string          `gorm:"size:255;not null;unique" json:"name"`
	EmbeddedLink string          `json:"embedded_link,omitempty" gorm:"type:text"`
	PolicyUrl    string          `json:"policy_url,omitempty" gorm:"type:text"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	DeletedAt    *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName sets the table name for the Policies model.
func (Policies) TableName() string {
	return "policies"
}

// MediaAttachment struct for storing media attachments related to the Tickets
type TicketMediaAttachment struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	TicketID  uint            `json:"ticket_id" gorm:"index;not null"`
	FileName  string          `gorm:"size:255" json:"file_name"`
	FilePath  string          `gorm:"size:255" json:"file_path"`
	MimeType  string          `gorm:"size:50" json:"mime_type"`
	URL       string          `gorm:"size:255;not null" json:"url"`
	Type      string          `gorm:"size:255" json:"type"`
	Caption   string          `gorm:"size:255" json:"caption,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (TicketMediaAttachment) TableName() string {
	return "ticket_media_attachments"
}

type TicketUpdate struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	TicketID  uint            `json:"ticket_id" gorm:"index;not null"`
	UserID    uint            `json:"user_id" gorm:"index;not null"`
	Update    string          `gorm:"type:text" json:"update"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (TicketUpdate) TableName() string {
	return "ticket_updates"
}

type TicketAsset struct {
	ID       uint `gorm:"primaryKey" json:"id"`
	TicketID uint `json:"ticket_id" gorm:"index;not null"`
	AssetID  uint `json:"asset_id" gorm:"index;not null"`
}

func (TicketAsset) TableName() string {
	return "ticket_assets"
}

type SatisfactionSurvey struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TicketID  uint      `gorm:"index;not null" json:"ticket_id"`
	Rating    int       `json:"rating" gorm:"type:int;not null"`
	Comment   string    `gorm:"type:text" json:"comment,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func (SatisfactionSurvey) TableName() string {
	return "satisfaction_surveys"
}

type SupportResponse struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	TicketID    uint      `gorm:"index;not null" json:"ticket_id"`
	ResponderID uint      `gorm:"index;not null" json:"responder_id"`
	Message     string    `gorm:"type:text;not null" json:"message"`
	RespondedAt time.Time `json:"responded_at"`
}

func (SupportResponse) TableName() string {
	return "support_responses"
}

type TicketResolution struct {
	gorm.Model
	TicketID         uint      `json:"ticket_id" gorm:"index;not null"`
	ResolvedBy       uint      `json:"resolved_by" gorm:"index;not null"`
	Resolution       string    `json:"resolution" gorm:"type:text;not null"`
	ResolvedAt       time.Time `json:"resolved_at"`
	CustomerFeedback string    `json:"customer_feedback" gorm:"type:text"` // Optional feedback from the customer
}

func (TicketResolution) TableName() string {
	return "ticket_resolution"
}

type CustomerSatisfactionSurvey struct {
	gorm.Model
	TicketID    uint      `json:"ticket_id" gorm:"index;not null"`
	Rating      int       `json:"rating" gorm:"type:int;not null"` // E.g., 1-5 scale
	Comments    string    `json:"comments" gorm:"type:text"`
	SubmittedAt time.Time `json:"submitted_at"`
}

type SLAPolicy struct {
	gorm.Model
	Name             string `json:"name" gorm:"type:varchar(255);not null;unique"`
	Description      string `json:"description" gorm:"type:text"`
	ResponseTarget   int    `json:"response_target" gorm:"type:int;not null"`     // Target response time in minutes
	ResolutionTarget int    `json:"resolution_target" gorm:"type:int;not null"`   // Target resolution time in hours
	AppliesTo        string `json:"applies_to" gorm:"type:varchar(100);not null"` // E.g., "All Tickets", "VIP Customers"
}

func (SLAPolicy) TableName() string {
	return "sla_policy"
}

type TicketStorage interface {
	CreateTicket(*Ticket) error
	DeleteTicket(int) error
	UpdateTicket(*Ticket) error
	GetTickets() ([]*Ticket, error)
	GetTicketByID(int) (*Ticket, error)
	GetTicketByNumber(int) (*Ticket, error)
	ListTicketsByStatus(status string) ([]Ticket, error)
}

type CommentStorage interface {
	CreateTicketComment(*Comment) error
	//DeleteSubCategory(int) error
	//UpdateSubCategory(*SubCategory) error
	//GetAllSubCategories() ([]*SubCategory, error)
	//GetSubCategoryByID(int) (*SubCategory, error)
	//GetSubCategoryByNumber(int) (*SubCategory, error)
}

type TicketHistoryEntryStorage interface {
	CreateTicketHistoryEntry(*TicketHistoryEntry) error
	//DeleteSubCategory(int) error
	//UpdateSubCategory(*SubCategory) error
	//GetAllSubCategories() ([]*SubCategory, error)
	//GetSubCategoryByID(int) (*SubCategory, error)
	//GetSubCategoryByNumber(int) (*SubCategory, error)
}

type SlaStorage interface {
	CreateSla(*SLA) error
	DeleteSla(int) error
	UpdateSla(*SLA) error
	GetAllSla() ([]*SLA, error)
	GetSlaByID(int) (*SLA, error)
	GetSlaByNumber(int) (*SLA, error)
}

type PriorityStorage interface {
	CreatePriority(*Priority) error
	DeletePriority(int) error
	UpdatePriority(*Priority) error
	GetPriorities() ([]*Priority, error)
	GetPriorityByID(int) (*Priority, error)
	GetPriorityByNumber(int) (*Priority, error)
}

type SatisfactionStorage interface {
	CreateSatisfaction(*Satisfaction) error
	DeleteSatisfaction(int) error
	UpdateSatisfaction(*Satisfaction) error
	GetSatisfactions() ([]*Satisfaction, error)
	GetSatisfactionByID(int) (*Satisfaction, error)
	GetSatisfactionByNumber(int) (*Satisfaction, error)
}

type CategoryStorage interface {
	CreateCategory(*Category) error
	DeleteCategory(int) error
	UpdateCategory(*Category) error
	GetAllCategories() ([]*Category, error)
	GetCategoryByID(int) (*Category, error)
	GetCategoryByNumber(int) (*Category, error)
}

type SubCategoryStorage interface {
	CreateSubCategory(*SubCategory) error
	DeleteSubCategory(int) error
	UpdateSubCategory(*SubCategory) error
	GetAllSubCategories() ([]*SubCategory, error)
	GetSubCategoryByID(int) (*SubCategory, error)
	GetSubCategoryByNumber(int) (*SubCategory, error)
}

type StatusStorage interface {
	CreateStatus(*Status) error
	DeleteStatus(int) error
	UpdateStatus(*Status) error
	GetStatus() ([]*Status, error)
	GetStatusByID(int) (*Status, error)
	GetStatusByNumber(int) (*Status, error)
}

// TicketModel handles database operations for Ticket
type TicketDBModel struct {
	DB             *gorm.DB
	EventPublisher EventPublisherImpl
	log            Logger
}

// NewTicketModel creates a new instance of TicketModel
func NewTicketDBModel(db *gorm.DB, log Logger, eventPublisher EventPublisherImpl) *TicketDBModel {
	return &TicketDBModel{
		DB:             db,
		EventPublisher: eventPublisher,
		log:            log,
	}
}

// TicketModel handles database operations for Ticket
type TicketCommentDBModel struct {
	DB             *gorm.DB
	EventPublisher EventPublisherImpl
	log            Logger
}

// NewTicketModel creates a new instance of TicketModel
func NewTicketCommentDBModel(db *gorm.DB, log Logger, eventPublisher EventPublisherImpl) *TicketCommentDBModel {
	return &TicketCommentDBModel{
		DB:             db,
		EventPublisher: eventPublisher,
		log:            log,
	}
}

// TicketModel handles database operations for Ticket
type TicketHistoryEntryDBModel struct {
	DB             *gorm.DB
	log            Logger
	EventPublisher EventPublisherImpl
}

// NewTicketModel creates a new instance of TicketModel
func NewTicketHistoryEntryDBModel(db *gorm.DB, log Logger, eventPublisher EventPublisherImpl) *TicketHistoryEntryDBModel {
	return &TicketHistoryEntryDBModel{
		DB:             db,
		log:            log,
		EventPublisher: eventPublisher,
	}
}

// CreateTicket creates a new ticket in the database.
func (db *TicketDBModel) CreateTicket(ticket *Ticket) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(ticket).Error; err != nil {
			return err
		}
		return nil
	})
}

// GetTicketByID retrieves a ticket by its ID.
func (db *TicketDBModel) GetTicketByID(ticketID uint) (*Ticket, error) {
	var ticket Ticket
	result := db.DB.First(&ticket, ticketID)
	return &ticket, result.Error
}

// UpdateTicket updates an existing ticket's details.
func (db *TicketDBModel) UpdateTicket(ticket *Ticket) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(ticket).Error; err != nil {
			return err
		}
		return nil
	})
}

// DeleteTicket removes a ticket from the database.
func (db *TicketDBModel) DeleteTicket(ticketID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&Ticket{}, ticketID).Error; err != nil {
			return err
		}
		return nil
	})
}

// GetAllTickets retrieves all tickets from the database.
func (as *TicketDBModel) GetAllTickets() ([]*Ticket, error) {
	var tickets []*Ticket
	err := as.DB.Find(&tickets).Error
	return tickets, err
}

// ////////////////////////////////////////////////////////////////////////////////////
// CreateTicketComment creates a new TicketComment.
func (repo *TicketCommentDBModel) CreateTicketComment(ticketID uint, commentText string) (uint, error) {
	if err := validateCommentInput(commentText); err != nil {
		return 0, err
	}

	var newCommentID uint
	operation := func(tx *gorm.DB) error {
		newComment := Comment{
			TicketID: ticketID,
			Comment:  commentText,
		}

		if err := tx.Create(&newComment).Error; err != nil {
			return err
		}

		newCommentID = newComment.ID
		return nil
	}

	retryPolicy := backoff.NewExponentialBackOff()
	retryPolicy.InitialInterval = 1 * time.Second
	retryPolicy.MaxElapsedTime = 10 * time.Second
	err := backoff.Retry(func() error {
		return repo.DB.Transaction(operation)
	}, retryPolicy)

	if err != nil {
		log.Printf("Transaction failed after retries: %v", err)
		return 0, err
	}

	go repo.publishEventAsync(TicketCommentCreatedEvent{CommentID: newCommentID, TicketID: ticketID})

	return newCommentID, nil
}

// GetCommentByID retrieves a Comment by its ID.
func (as *TicketCommentDBModel) GetCommentByID(id uint) (*Comment, error) {
	var comment Comment
	err := as.DB.Where("id = ?", id).First(&comment).Error
	return &comment, err
}

////////////////////////////////////////////////////////////////////////////////////////////

// CreateTicketHistoryEntry creates a new TicketHistoryEntry.
func (as *TicketHistoryEntryDBModel) CreateTicketHistoryEntry(ticketHistory *TicketHistoryEntry, action string) error {
	ticketHistory.Action = action
	return as.DB.Create(ticketHistory).Error
}

// GetCommentByID retrieves a Comment by its ID.
func (thdm *TicketHistoryEntryDBModel) GetHistoryEntriesByTicketID(ticketID uint) ([]*TicketHistoryEntry, error) {
	var historyEntries []*TicketHistoryEntry
	err := thdm.DB.Where("ticket_id = ?", ticketID).Find(&historyEntries).Error
	return historyEntries, err
}

// ///////////////////////////////////////////////////////////////////////////////////////////
func (as *TicketDBModel) CreateTag(ticketID uint, tagName string) (*Tag, bool, error) {
	var tag Tag
	addTagStatus := false

	// Check if the tag already exists; if not, create a new one
	err := as.DB.FirstOrCreate(&tag, Tag{Name: tagName}).Error
	if err != nil {
		return nil, addTagStatus, fmt.Errorf("error handling the tag: %v", err)
	}

	// Fetch the ticket with its tags to avoid replacing existing associations
	var ticket Ticket
	if err := as.DB.Preload("Tags").First(&ticket, ticketID).Error; err != nil {
		return nil, addTagStatus, fmt.Errorf("ticket not found: %v", err)
	}

	// Associate the tag with the ticket
	// Check if the tag is already associated to avoid duplicates
	alreadyExists := false
	for _, existingTag := range ticket.Tags {
		if existingTag.ID == tag.ID {
			alreadyExists = true
			break
		}
	}

	if !alreadyExists {
		if err := as.DB.Model(&ticket).Association("Tags").Append(&tag); err != nil {
			return nil, addTagStatus, fmt.Errorf("failed to associate tag with ticket: %v", err)
		}
		addTagStatus = true // Indicate that the tag was successfully added
	}

	return &tag, addTagStatus, nil
}

func (tm *TicketDBModel) CreateTagAndAssociateWithTicket(ticketID uint, tagName string) error {
	return tm.DB.Transaction(func(tx *gorm.DB) error {
		var tag Tag
		if err := tx.FirstOrCreate(&tag, Tag{Name: tagName}).Error; err != nil {
			return err
		}

		if err := tx.Model(&Ticket{Model: gorm.Model{ID: ticketID}}).Association("Tags").Append(&tag); err != nil {
			return err
		}

		return nil
	})
}

///////////////////////////////////////////////////////////////////////////////////////////

// AssignAgentToTicket assigns an agent to a ticket and updates the ticket's status.
func (db *TicketDBModel) AssignTicketToAgent(ticketID, agentID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		ticket := &Ticket{}
		if err := tx.First(ticket, ticketID).Error; err != nil {
			return err
		}
		ticket.AgentID = &agentID
		if err := tx.Save(ticket).Error; err != nil {
			return err
		}
		return nil
	})
}

// ListTicketsByStatus lists all tickets with a specific status.
func (db *TicketDBModel) ListTicketsByStatus(status string) ([]Ticket, error) {
	var tickets []Ticket
	result := db.DB.Where("status = ?", status).Find(&tickets)
	return tickets, result.Error
}

func (tdb *TicketDBModel) ChangeTicketStatus(ticketID uint, newStatus *Status) error {
	// Change the status of a ticket
	ticket := &Ticket{}
	if err := tdb.DB.First(&ticket, ticketID).Error; err != nil {
		return err
	}

	ticket.Status = newStatus.Name
	ticket.StatusObject = *newStatus
	if err := tdb.DB.Save(ticket).Error; err != nil {
		return err
	}

	return nil
}

func (tdb *TicketDBModel) AddCommentToTicket(ticketID uint, comment string) error {
	// Add a comment to a ticket
	ticketComment := &Comment{
		TicketID: ticketID,
		Comment:  comment,
		//Comment:  comment,
	}
	if err := tdb.DB.Create(&ticketComment).Error; err != nil {
		return err
	}

	return nil
}

func (tdb *TicketDBModel) GetTicketHistory(ticketID uint) ([]*TicketHistoryEntry, error) {
	// Retrieve the history of a ticket, including comments and status changes
	var history []*TicketHistoryEntry

	// Get all comments for the ticket
	if err := tdb.DB.Where("ticket_id = ?", ticketID).Find(&history).Error; err != nil {
		return nil, err
	}

	// Include status changes in the history
	statusChanges := make([]*TicketHistoryEntry, 0)
	if err := tdb.DB.Where("ticket_id = ? AND field = ?", ticketID, "status").Find(&statusChanges).Error; err != nil {
		return nil, err
	}

	history = append(history, statusChanges...)

	// Sort the history by timestamp
	sort.Slice(history, func(i, j int) bool {
		return history[i].CreatedAt.Before(history[j].CreatedAt)
	})

	return history, nil
}

///////////////////////////////////////////////////////// SLA

func (tdb *TicketDBModel) CreateSla(sla *SLA) error {
	return tdb.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(sla).Error; err != nil {
			return err
		}
		return nil
	})
}

func (tdb *TicketDBModel) UpdateSla(sla *SLA) error {
	// Update an existing SLA
	if err := tdb.DB.Save(sla).Error; err != nil {
		return err
	}

	return nil
}

func (tdb *TicketDBModel) DeleteSla(slaID uint) error {
	// Delete an SLA by its ID
	if err := tdb.DB.Delete(&SLA{}, slaID).Error; err != nil {
		return err
	}

	return nil
}

// GetAllSLAs retrieves all SLAs from the database.
func (as *TicketDBModel) GetAllSLAs() ([]*SLA, error) {
	var slas []*SLA
	err := as.DB.Find(&slas).Error
	if err != nil {
		return nil, err
	}
	return slas, nil
}

// GetSLAByID retrieves an SLA by its ID.
func (as *TicketDBModel) GetSLAByID(id uint) (*SLA, error) {
	var sla SLA
	err := as.DB.Where("id = ?", id).First(&sla).Error
	if err != nil {
		return nil, err
	}
	return &sla, nil
}

// GetSLAByNumber retrieves an SLA by its SLA number.
func (as *TicketDBModel) GetSLAByNumber(slaNumber int) (*SLA, error) {
	var sla SLA
	err := as.DB.Where("sla_id = ?", slaNumber).First(&sla).Error
	if err != nil {
		return nil, err
	}
	return &sla, nil
}

///////////////////////////////////////////////////////// PRIORITY

func (tdb *TicketDBModel) CreatePriority(priority *Priority) error {
	return tdb.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(priority).Error; err != nil {
			return err
		}
		return nil
	})
}

func (tdb *TicketDBModel) UpdatePriority(priority *Priority) error {
	// Update an existing priority level
	if err := tdb.DB.Save(priority).Error; err != nil {
		return err
	}

	return nil
}

func (tdb *TicketDBModel) DeletePriority(priorityID uint) error {
	// Delete a priority level by its ID
	if err := tdb.DB.Delete(&Priority{}, priorityID).Error; err != nil {
		return err
	}

	return nil
}

// GetAllPriorities retrieves all priority levels from the database.
func (as *TicketDBModel) GetAllPriorities() ([]*Priority, error) {
	var priorities []*Priority
	err := as.DB.Find(&priorities).Error
	if err != nil {
		return nil, err
	}
	return priorities, nil
}

// GetPriorityByID retrieves a priority level by its ID.
func (as *TicketDBModel) GetPriorityByID(id uint) (*Priority, error) {
	var priority Priority
	err := as.DB.Where("id = ?", id).First(&priority).Error
	if err != nil {
		return nil, err
	}
	return &priority, nil
}

// GetPriorityByNumber retrieves a priority level by its priority number.
func (as *TicketDBModel) GetPriorityByNumber(priorityNumber int) (*Priority, error) {
	var priority Priority
	err := as.DB.Where("priority_id = ?", priorityNumber).First(&priority).Error
	if err != nil {
		return nil, err
	}
	return &priority, nil
}

///////////////////////////////////////////////////////// TAGS

func (tdb *TicketDBModel) AddTagToTicket(ticketID uint, tagName string) error {
	return tdb.DB.Transaction(func(tx *gorm.DB) error {
		var tag Tag
		// Ensure the tag exists or create a new one
		if err := tx.FirstOrCreate(&tag, Tag{Name: tagName}).Error; err != nil {
			return fmt.Errorf("error finding or creating tag: %w", err)
		}

		// Associate the tag with the ticket
		if err := tx.Model(&Ticket{Model: gorm.Model{ID: ticketID}}).Association("Tags").Append(&tag); err != nil {
			return fmt.Errorf("error appending tag to ticket: %w", err)
		}

		return nil
	})
}

func (tdb *TicketDBModel) RemoveTagFromTicket(ticketID uint, tagName string) error {
	return tdb.DB.Transaction(func(tx *gorm.DB) error {
		var tag Tag
		if err := tx.Where("name = ?", tagName).First(&tag).Error; err != nil {
			return fmt.Errorf("error finding tag: %w", err)
		}

		// Remove the association between the tag and the ticket
		if err := tx.Model(&Ticket{Model: gorm.Model{ID: ticketID}}).Association("Tags").Delete(&tag); err != nil {
			return fmt.Errorf("error removing tag from ticket: %w", err)
		}

		return nil
	})
}

///////////////////////////////////////////////////////// STATUS

func (tdb *TicketDBModel) CreateStatus(status *Status) error {
	// Create a new ticket status
	if err := tdb.DB.Create(&status).Error; err != nil {
		return err
	}

	return nil
}

func (tdb *TicketDBModel) UpdateStatus(status *Status) error {
	// Update an existing ticket status
	if err := tdb.DB.Save(&status).Error; err != nil {
		return err
	}

	return nil
}

func (tdb *TicketDBModel) DeleteStatus(statusID uint) error {
	// Delete a ticket status by its ID
	if err := tdb.DB.Delete(&Status{}, statusID).Error; err != nil {
		return err
	}

	return nil
}

// GetAllStatuses retrieves all ticket statuses from the database.
func (as *TicketDBModel) GetAllStatuses() ([]*Status, error) {
	var statuses []*Status
	err := as.DB.Find(&statuses).Error
	if err != nil {
		return nil, err
	}
	return statuses, nil
}

// GetStatusByID retrieves a ticket status by its ID.
func (as *TicketDBModel) GetStatusByID(id uint) (*Status, error) {
	var status Status
	err := as.DB.Where("id = ?", id).First(&status).Error
	if err != nil {
		return nil, err
	}
	return &status, nil
}

// GetStatusByNumber retrieves a ticket status by its status number.
func (as *TicketDBModel) GetStatusByNumber(statusNumber int) (*Status, error) {
	var status Status
	err := as.DB.Where("status_id = ?", statusNumber).First(&status).Error
	if err != nil {
		return nil, err
	}
	return &status, nil
}

////////////////////////////////////////////////////////// CATEGORY

func (tdb *TicketDBModel) CreateCategory(category *Category) error {
	// Create a new category
	if err := tdb.DB.Create(&category).Error; err != nil {
		return err
	}

	return nil
}

func (tdb *TicketDBModel) UpdateCategory(category *Category) error {
	// Update an existing category
	if err := tdb.DB.Save(&category).Error; err != nil {
		return err
	}

	return nil
}

func (tdb *TicketDBModel) DeleteCategory(categoryID uint) error {
	// Delete a category by its ID
	if err := tdb.DB.Delete(&Category{}, categoryID).Error; err != nil {
		return err
	}

	return nil
}

// GetAllCategories retrieves all categories from the database.
func (as *TicketDBModel) GetAllCategories() ([]*Category, error) {
	var categories []*Category
	err := as.DB.Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// GetCategoryByID retrieves a category by its ID.
func (as *TicketDBModel) GetCategoryByID(id uint) (*Category, error) {
	var category Category
	err := as.DB.Where("id = ?", id).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// GetCategoryByNumber retrieves a category by its category number.
func (as *TicketDBModel) GetCategoryByNumber(categoryNumber int) (*Category, error) {
	var category Category
	err := as.DB.Where("category_id = ?", categoryNumber).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

//////////////////////////////////////////////////////// SUBCATEGORY

// backend/models/ticket_db_model.go

func (tdb *TicketDBModel) CreateSubCategory(subCategory *SubCategory) error {
	return tdb.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(subCategory).Error; err != nil {
			return err
		}
		return nil
	})
}

func (tdb *TicketDBModel) UpdateSubCategory(subCategory *SubCategory) error {
	return tdb.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(subCategory).Error; err != nil {
			return err
		}
		return nil
	})
}

func (tdb *TicketDBModel) DeleteSubcategory(subcategoryID uint) error {
	// Delete a subcategory by its ID
	if err := tdb.DB.Delete(&SubCategory{}, subcategoryID).Error; err != nil {
		return err
	}

	return nil
}

// GetAllSubcategories retrieves all subcategories from the database.
func (as *TicketDBModel) GetAllSubcategories() ([]*SubCategory, error) {
	var subcategories []*SubCategory
	err := as.DB.Find(&subcategories).Error
	if err != nil {
		return nil, err
	}
	return subcategories, nil
}

// GetSubcategoryByID retrieves a subcategory by its ID.
func (as *TicketDBModel) GetSubcategoryByID(id uint) (*SubCategory, error) {
	var subcategory SubCategory
	err := as.DB.Where("id = ?", id).First(&subcategory).Error
	if err != nil {
		return nil, err
	}
	return &subcategory, nil
}

// GetSubcategoryByNumber retrieves a subcategory by its subcategory number.
func (as *TicketDBModel) GetSubcategoryByNumber(subcategoryNumber int) (*SubCategory, error) {
	var subcategory SubCategory
	err := as.DB.Where("subcategory_id = ?", subcategoryNumber).First(&subcategory).Error
	if err != nil {
		return nil, err
	}
	return &subcategory, nil
}

// //////////////////////////////////////////////////////////////
func (db *TicketDBModel) GetTicketByIDWithAssociations(ticketID uint) (*Ticket, error) {
	var ticket Ticket
	err := db.DB.Preload("Tags").Preload("Comments").Preload("TicketHistory").First(&ticket, ticketID).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

// Soft delete a ticket
func (db *TicketDBModel) SoftDeleteTicket(ticketID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&Ticket{}, ticketID).Error; err != nil {
			return err
		}
		return nil
	})
}

func (db *TicketDBModel) GetOverdueHighPriorityTickets() ([]*Ticket, error) {
	var tickets []*Ticket
	err := db.DB.Joins("Priority").Where("DueAt < ? AND priorities.level >= ?", time.Now(), 10).Find(&tickets).Error
	if err != nil {
		return nil, err
	}
	return tickets, nil
}

func (db *TicketDBModel) BatchInsertComments(comments []*Comment) error {
	return db.DB.CreateInBatches(comments, 100).Error // Inserts in batches of 100
}

func (t *Ticket) BeforeCreate(tx *gorm.DB) (err error) {
	if t.Status == "" {
		t.Status = "open" // Default status if not provided
	}
	return nil
}

/*
func (db *TicketDBModel) GetTicketStatusesCached() ([]*Status, error) {
	// Attempt to fetch from cache first
	if cachedStatuses, found := cache.Get("ticket_statuses"); found {
		return cachedStatuses.([]*Status), nil
	}

	// Fallback to database and update cache
	statuses, err := db.GetAllStatuses()
	if err != nil {
		return nil, err
	}
	cache.Set("ticket_statuses", statuses, cache.DefaultExpiration)
	return statuses, nil
}

func (db *TicketDBModel) GetHighPriorityTicketsOptimized() ([]*Ticket, error) {
	var tickets []*Ticket
	// Use raw SQL or GORM's advanced features for complex queries
	err := db.DB.Raw("SELECT * FROM tickets WHERE priority_id IN (SELECT id FROM priorities WHERE level > ?)", 10).Scan(&tickets).Error
	if err != nil {
		return nil, err
	}
	return tickets, nil
}

func (db *TicketDBModel) CreateTicketSecure(ticket *Ticket) error {
	p := bluemonday.UGCPolicy()
	safeHTML := p.Sanitize(unsafeHTML)
	// Sanitize input data
	ticket.Subject = sanitize(ticket.Subject)
	ticket.Description = sanitize(ticket.Description)

	// Validate data
	if err := validateTicket(ticket); err != nil {
		return err
	}

	return db.CreateTicket(ticket)
}

func (db *TicketDBModel) UpdateTicketAsUser(ticket *Ticket, userID uint) error {
	// Check if the user has the right to update the ticket
	if !db.HasUpdatePermission(ticket.ID, userID) {
		return fmt.Errorf("unauthorized")
	}

	return db.UpdateTicket(ticket)
}
*/

func (repo *TicketCommentDBModel) publishEventAsync(event interface{}) {
	if err := repo.EventPublisher.Publish(event); err != nil {
		log.Printf("Failed to publish event: %v", err)
	}
}

func validateCommentInput(comment string) error {
	if comment == "" {
		return errors.New("comment text cannot be empty")
	}
	return nil
}

func (repo *TicketCommentDBModel) ExecuteWithRetry(operation func(tx *gorm.DB) error) error {
	retryBackoff := backoff.NewExponentialBackOff()
	retryBackoff.InitialInterval = 1 * time.Second
	retryBackoff.MaxElapsedTime = 10 * time.Second

	operationWithBackoff := func() error {
		return repo.DB.Transaction(operation)
	}

	err := backoff.Retry(operationWithBackoff, retryBackoff)
	if err != nil {
		log.Printf("Transaction failed after retries: %v", err)
		return err
	}
	return nil
}

func IsTransientError(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1205 || mysqlErr.Number == 1213
	}
	return false
}
