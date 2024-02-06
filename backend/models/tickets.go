// backend/models/tickets.go

package models

import (
	"fmt"
	"sort"
	"time"

	"gorm.io/gorm"
)

type Ticket struct {
	gorm.Model
	Subject          string                  `gorm:"size:255;not null" json:"subject"`
	Description      string                  `gorm:"type:text;not null" json:"description"`
	CategoryID       uint                    `json:"category_id"`
	SubCategoryID    uint                    `json:"sub_category_id"`
	PriorityID       uint                    `json:"priority_id"`
	SLAID            uint                    `json:"sla_id"`
	UserID           uint                    `json:"user_id"`
	AgentID          uint                    `json:"agent_id"`
	DueAt            *time.Time              `json:"due_at,omitempty"`
	Site             string                  `gorm:"size:255" json:"site"`
	StatusID         uint                    `json:"status_id"`
	Status           Status                  `json:"status"`
	Priority         Priority                `json:"priority"`
	Category         SubCategory             `json:"category"`
	MediaAttachments []TicketMediaAttachment `gorm:"foreignKey:TicketID" json:"media_attachments"`
	Tags             []Tag                   `gorm:"many2many:ticket_tags;" json:"tags"`
	Comments         []Comment               `gorm:"foreignKey:TicketID" json:"comments"`
	TicketHistory    []TicketHistoryEntry    `gorm:"foreignKey:TicketID" json:"ticket_history"`
	DeletedAt        gorm.DeletedAt          `gorm:"index" json:"deleted_at,omitempty"`
}

func (Ticket) TableName() string {
	return "tickets"
}

type Comment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	TicketID  uint           `json:"ticket_id"`
	AuthorID  uint           `json:"author_id"`
	Comment   string         `json:"comment"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Comment) TableName() string {
	return "comments"
}

type TicketHistoryEntry struct {
	gorm.Model
	TicketID uint   `json:"ticket_id"`
	Action   string `gorm:"size:255;not null" json:"action"`
}

func (TicketHistoryEntry) TableName() string {
	return "ticket_history_entries"
}

// RelatedAd struct for storing related advertisements
type RelatedTicket struct {
	gorm.Model
	TicketID        uint `json:"ticket_id"`
	RelatedTicketID uint `json:"related_ticket_id"`
}

func (RelatedTicket) TableName() string {
	return "related_tickets"
}

// Hashtag represents a hashtag entity
type Tag struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:255;not null;unique" json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Tag) TableName() string {
	return "tags"
}

type SLA struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `gorm:"size:255;not null" json:"name"`
	Description   string         `gorm:"type:text" json:"description,omitempty"`
	Target        int            `json:"target"`          // Target resolution time in hours
	TimeToResolve int            `json:"time_to_resolve"` // In hours
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (SLA) TableName() string {
	return "slas"
}

type Priority struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:255;not null" json:"name"`
	Level     int            `json:"level"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Priority) TableName() string {
	return "priorities"
}

type Satisfaction struct {
	gorm.Model
	SatisfactionID uint      `gorm:"primaryKey" json:"satisfaction_id"`
	Name           string    `json:"satisfaction_name"`
	Rank           int       `json:"rank"`
	Emoji          string    `json:"emoji"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	DeletedAt      time.Time `json:"deleted_at"`
}

// TableName sets the table name for the Satisfaction model.
func (Satisfaction) TableName() string {
	return "satisfaction"
}

type Category struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:255;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Category) TableName() string {
	return "categories"
}

type SubCategory struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:255;not null" json:"name"`
	CategoryID  uint           `json:"category_id"`
	Description string         `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (SubCategory) TableName() string {
	return "sub_categories"
}

type Status struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:255;not null" json:"name"`
	IsClosed  bool           `json:"is_closed"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Status) TableName() string {
	return "statuses"
}

type Policies struct {
	gorm.Model
	PolicyID     uint      `gorm:"primaryKey" json:"policy_id"`
	PolicyName   string    `json:"policy_name"`
	EmbeddedLink string    `json:"policy_embed"`
	PolicyUrl    string    `json:"policy_url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at"`
}

// TableName sets the table name for the Policies model.
func (Policies) TableName() string {
	return "policies"
}

// MediaAttachment struct for storing media attachments related to the Tickets
type TicketMediaAttachment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	TicketID  uint           `json:"ticket_id"`
	FileName  string         `gorm:"size:255" json:"file_name"`
	FilePath  string         `gorm:"size:255" json:"file_path"`
	MimeType  string         `gorm:"size:50" json:"mime_type"`
	URL       string         `gorm:"size:255;not null" json:"url"`
	Type      string         `gorm:"size:255" json:"type"` // E.g., image, document
	Caption   string         `gorm:"size:255" json:"caption,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (TicketMediaAttachment) TableName() string {
	return "ticket_media_attachments"
}

type TicketUpdate struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	TicketID  uint           `json:"ticket_id"`
	UserID    uint           `json:"user_id"`
	Update    string         `gorm:"type:text" json:"update"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (TicketUpdate) TableName() string {
	return "ticket_updates"
}

type TicketAsset struct {
	ID       uint `gorm:"primaryKey" json:"id"`
	TicketID uint `json:"ticket_id"`
	AssetID  uint `json:"asset_id"`
}

func (TicketAsset) TableName() string {
	return "ticket_assets"
}

type SatisfactionSurvey struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TicketID  uint      `json:"ticket_id"`
	Rating    int       `json:"rating"` // For example, 1 to 5
	Comment   string    `gorm:"type:text" json:"comment,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func (SatisfactionSurvey) TableName() string {
	return "satisfaction_surveys"
}

type TicketStorage interface {
	CreateTicket(*Ticket) error
	DeleteTicket(int) error
	UpdateTicket(*Ticket) error
	GetTickets() ([]*Ticket, error)
	GetTicketByID(int) (*Ticket, error)
	GetTicketByNumber(int) (*Ticket, error)
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
	DB *gorm.DB
}

// NewTicketModel creates a new instance of TicketModel
func NewTicketDBModel(db *gorm.DB) *TicketDBModel {
	return &TicketDBModel{
		DB: db,
	}
}

// TicketModel handles database operations for Ticket
type TicketCommentDBModel struct {
	DB *gorm.DB
}

// NewTicketModel creates a new instance of TicketModel
func NewTicketCommentDBModel(db *gorm.DB) *TicketCommentDBModel {
	return &TicketCommentDBModel{
		DB: db,
	}
}

// TicketModel handles database operations for Ticket
type TicketHistoryEntryDBModel struct {
	DB *gorm.DB
}

// NewTicketModel creates a new instance of TicketModel
func NewTicketHistoryEntryDBModel(db *gorm.DB) *TicketHistoryEntryDBModel {
	return &TicketHistoryEntryDBModel{
		DB: db,
	}
}

// CreateTicket creates a new Ticket.
func (as *TicketDBModel) CreateTicket(ticket *Ticket) error {
	return as.DB.Create(ticket).Error
}

// GetTicketByID retrieves a Ticket by its ID.
func (as *TicketDBModel) GetTicketByID(id uint) (*Ticket, error) {
	var ticket Ticket
	err := as.DB.Where("id = ?", id).First(&ticket).Error
	return &ticket, err
}

// UpdateTicket updates the details of an existing Ticket.
func (as *TicketDBModel) UpdateTicket(ticket *Ticket) error {
	if err := as.DB.Save(ticket).Error; err != nil {
		return err
	}
	return nil
}

// DeleteTicket deletes a ticket from the database.
func (as *TicketDBModel) DeleteTicket(id uint) error {
	if err := as.DB.Delete(&Ticket{}, id).Error; err != nil {
		return err
	}
	return nil
}

// GetAllTickets retrieves all tickets from the database.
func (as *TicketDBModel) GetAllTickets() ([]*Ticket, error) {
	var tickets []*Ticket
	err := as.DB.Find(&tickets).Error
	return tickets, err
}

// ////////////////////////////////////////////////////////////////////////////////////
// CreateTicketComment creates a new TicketComment.
func (as *TicketCommentDBModel) CreateTicketComment(ticketID uint, c string) (*Comment, error) {
	var comment Comment
	comment.TicketID = ticketID
	//comment.Author = getCurrentUser()
	comment.Comment = c
	id := as.DB.Create(comment).RowsAffected
	return as.GetCommentByID(uint(id))
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
func (as *TicketHistoryEntryDBModel) GetHistoryEntriesByTicketID(ticketID uint) []*TicketHistoryEntry {
	var ticketHistory []*TicketHistoryEntry
	err := as.DB.Find(&ticketHistory).Error
	if err != nil {
		return nil
	}
	return ticketHistory
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// CreateTicket creates a new Tag.
func (as *TicketDBModel) CreateTag(ticketID uint, tag string) (*Tag, bool, error) {
	addTagStatus := false
	ticket, err := as.GetTicketByID(ticketID)
	if err != nil {
		return nil, addTagStatus, fmt.Errorf("ticket not found")
	}
	ticket.Tags = append(ticket.Tags, tag)
	erro := as.UpdateTicket(ticket)
	if erro != nil {
		return nil, addTagStatus, fmt.Errorf("ticket not found")
	}
	addTagStatus = true
	return ticket.Tags, addTagStatus, nil
}

///////////////////////////////////////////////////////////////////////////////////////////

func (tdb *TicketDBModel) AssignTicketToAgent(ticketID uint, agent *Agents) error {
	// Assign a ticket to an agent by updating the AgentID field in the ticket
	ticket := &Ticket{}
	if err := tdb.DB.First(&ticket, ticketID).Error; err != nil {
		return err
	}

	ticket.AgentID = *&agent.ID
	if err := tdb.DB.Save(&ticket).Error; err != nil {
		return err
	}

	return nil
}

func (tdb *TicketDBModel) ChangeTicketStatus(ticketID uint, newStatus *Status) error {
	// Change the status of a ticket
	ticket := &Ticket{}
	if err := tdb.DB.First(&ticket, ticketID).Error; err != nil {
		return err
	}

	ticket.Status = *newStatus
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
	// Create a new SLA
	if err := tdb.DB.Create(sla).Error; err != nil {
		return err
	}

	return nil
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
	// Create a new priority level
	if err := tdb.DB.Create(priority).Error; err != nil {
		return err
	}

	return nil
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

func (tdb *TicketDBModel) AddTagToTicket(ticketID uint, tag string) error {
	//var t Tag
	t, _, _ := tdb.CreateTag(ticketID, tag)
	// Add a tag to a ticket
	ticket := &Ticket{}
	if err := tdb.DB.First(&ticket, ticketID).Error; err != nil {
		return err
	}

	// Append the tag to the ticket's Tags field
	ticket.Tags = append(ticket.Tags, t)

	if err := tdb.DB.Save(&ticket).Error; err != nil {
		return err
	}

	return nil
}

func (tdb *TicketDBModel) RemoveTagFromTicket(ticketID uint, tag string) error {
	// Remove a tag from a ticket
	ticket := &Ticket{}
	if err := tdb.DB.First(&ticket, ticketID).Error; err != nil {
		return err
	}

	// Remove the tag from the ticket's Tags field
	for i, existingTag := range ticket.Tags.Tags {
		if existingTag == tag {
			ticket.Tags.Tags = append(ticket.Tags.Tags[:i], ticket.Tags.Tags[i+1:]...)
			break
		}
	}

	if err := tdb.DB.Save(&ticket).Error; err != nil {
		return err
	}

	return nil
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

func (tdb *TicketDBModel) CreateSubcategory(subcategory *SubCategory) error {
	// Create a new subcategory
	if err := tdb.DB.Create(&subcategory).Error; err != nil {
		return err
	}

	return nil
}

func (tdb *TicketDBModel) UpdateSubcategory(subcategory *SubCategory) error {
	// Update an existing subcategory
	if err := tdb.DB.Save(&subcategory).Error; err != nil {
		return err
	}

	return nil
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
