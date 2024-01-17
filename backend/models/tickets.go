// backend/models/tickets.go

package models

import (
	"time"

	"gorm.io/gorm"
)

type Ticket struct {
	gorm.Model
	ID               uint                    `gorm:"primaryKey" json:"ticket_id"`
	Subject          string                  `json:"subject"`
	Description      string                  `json:"description"`
	Category         Category                `json:"category" gorm:"embedded"`
	SubCategory      SubCategory             `json:"sub_category" gorm:"embedded"`
	Priority         Priority                `json:"priority" gorm:"embedded"`
	SLA              Sla                     `json:"sla" gorm:"embedded"`
	UserID           Users                   `json:"user" gorm:"embedded"`
	AgentID          Agents                  `json:"agent" gorm:"embedded"`
	CreatedAt        time.Time               `json:"created_at"`
	UpdatedAt        time.Time               `json:"updated_at"`
	DueAt            time.Time               `json:"due_at"`
	AssetID          []Assets                `json:"asset_id" gorm:"embedded"`
	RelatedTickets   []RelatedTicket         `json:"related_ticket_id" gorm:"foreignKey:TicketID"`
	MediaAttachments []TicketMediaAttachment `json:"mediaAttachments" gorm:"foreignKey:TicketID"`
	Tags             []Tags                  `json:"hashtags" gorm:"foreignKey:TicketID"`
	Site             string                  `json:"site"`
	Status           Status                  `json:"status" gorm:"embedded"`
}

// TableName sets the table name for the Ticket model.
func (Ticket) TableName() string {
	return "tickets"
}

// RelatedAd struct for storing related advertisements
type RelatedTicket struct {
	gorm.Model
	TicketID        uint `json:"-"`
	RelatedTicketID uint `json:"relatedTicketID" gorm:"foreignKey:TicketID"`
	Order           int  `json:"order" gorm:"default:0"`
}

// TableName sets the table name for the Advertisement model.
func (RelatedTicket) TableName() string {
	return "related_tickets"
}

// Hashtag represents a hashtag entity
type Tags struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TagName   string    `json:"tag"`
	Tags      []string  `json:"tags"` // Added Tags field
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TicketID  uint      `json:"-"`
}

// TableName sets the table name for the Tags model.
func (Tags) TableName() string {
	return "tags"
}

type Sla struct {
	gorm.Model
	SlaID          int       `gorm:"primaryKey" json:"sla_id"`
	SlaName        string    `json:"sla_name"`
	PriorityID     int       `json:"priority_id"`
	SatisfactionID int       `json:"satisfaction_id"`
	PolicyID       int       `json:"policy_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// TableName sets the table name for the Sla model.
func (Sla) TableName() string {
	return "sla"
}

type Priority struct {
	gorm.Model
	PriorityID    int       `gorm:"primaryKey" json:"priority_id"`
	Name          string    `json:"priority_name"`
	FirstResponse int       `json:"first_response"`
	Colour        string    `json:"red"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// TableName sets the table name for the priority model.
func (Priority) TableName() string {
	return "priority"
}

type Satisfaction struct {
	gorm.Model
	SatisfactionID int       `gorm:"primaryKey" json:"satisfaction_id"`
	Name           string    `json:"satisfaction_name"`
	Rank           int       `json:"rank"`
	Emoji          string    `json:"emoji"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// TableName sets the table name for the Satisfaction model.
func (Satisfaction) TableName() string {
	return "satisfaction"
}

type Category struct {
	gorm.Model
	ID           int       `gorm:"primaryKey" json:"category_id"`
	CategoryName string    `json:"category_name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName sets the table name for the Category model.
func (Category) TableName() string {
	return "category"
}

type SubCategory struct {
	gorm.Model
	SubCategoryID   int       `gorm:"primaryKey" json:"sub_category_id"`
	SubCategoryName string    `json:"sub_category_name"`
	CategoryID      int       `json:"category_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// TableName sets the table name for the SubCategory model.
func (SubCategory) TableName() string {
	return "subCategory"
}

type Status struct {
	gorm.Model
	StatusID   int       `gorm:"primaryKey" json:"status_id"`
	StatusName string    `json:"status_name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// TableName sets the table name for the Status model.
func (Status) TableName() string {
	return "status"
}

type Policies struct {
	gorm.Model
	PolicyID     int       `gorm:"primaryKey" json:"policy_id"`
	PolicyName   string    `json:"policy_name"`
	EmbeddedLink string    `json:"policy_embed"`
	PolicyUrl    string    `json:"policy_url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName sets the table name for the Policies model.
func (Policies) TableName() string {
	return "policies"
}

// MediaAttachment struct for storing media attachments related to the Tickets
type TicketMediaAttachment struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey" json:"-"`
	URL       string `json:"url"`
	Type      string `json:"type"`
	Caption   string `json:"caption"`
	AltText   string `json:"altText"`
	IsPrimary bool   `json:"isPrimary" gorm:"default:false"`
	Order     int    `json:"order" gorm:"default:0"`
	TicketID  uint   `json:"-"`
	// Add more fields as needed
}

// TableName sets the table name for the Ticket model.
func (TicketMediaAttachment) TableName() string {
	return "ticket_media_attachment"
}

type TicketStorage interface {
	CreateTicket(*Ticket) error
	DeleteTicket(int) error
	UpdateTicket(*Ticket) error
	GetTickets() ([]*Ticket, error)
	GetTicketByID(int) (*Ticket, error)
	GetTicketByNumber(int) (*Ticket, error)
}

type SlaStorage interface {
	CreateTicket(*Ticket) error
	DeleteTicket(int) error
	UpdateTicket(*Ticket) error
	GetTickets() ([]*Ticket, error)
	GetTicketByID(int) (*Ticket, error)
	GetTicketByNumber(int) (*Ticket, error)
}

type PriorityStorage interface {
	CreateTicket(*Ticket) error
	DeleteTicket(int) error
	UpdateTicket(*Ticket) error
	GetTickets() ([]*Ticket, error)
	GetTicketByID(int) (*Ticket, error)
	GetTicketByNumber(int) (*Ticket, error)
}

type SatisfactionStorage interface {
	CreateTicket(*Ticket) error
	DeleteTicket(int) error
	UpdateTicket(*Ticket) error
	GetTickets() ([]*Ticket, error)
	GetTicketByID(int) (*Ticket, error)
	GetTicketByNumber(int) (*Ticket, error)
}

type CategoryStorage interface {
	CreateTicket(*Ticket) error
	DeleteTicket(int) error
	UpdateTicket(*Ticket) error
	GetTickets() ([]*Ticket, error)
	GetTicketByID(int) (*Ticket, error)
	GetTicketByNumber(int) (*Ticket, error)
}

type SubCategoryStorage interface {
	CreateTicket(*Ticket) error
	DeleteTicket(int) error
	UpdateTicket(*Ticket) error
	GetTickets() ([]*Ticket, error)
	GetTicketByID(int) (*Ticket, error)
	GetTicketByNumber(int) (*Ticket, error)
}

type StatusStorage interface {
	CreateTicket(*Ticket) error
	DeleteTicket(int) error
	UpdateTicket(*Ticket) error
	GetTickets() ([]*Ticket, error)
	GetTicketByID(int) (*Ticket, error)
	GetTicketByNumber(int) (*Ticket, error)
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
