// backend/models/assets.go

package models

import (
	"time"

	"gorm.io/gorm"
)

type Assets struct {
	gorm.Model
	ID            int             `gorm:"primaryKey" json:"asset_id"`
	Tag           AssetTag        `json:"asset_tag" gorm:"foreignKey:AssetID"`
	AssetType     AssetType       `json:"asset_type" gorm:"embedded"`
	AssetName     string          `json:"asset_name"`
	Assignment    AssetAssignment `json:"asset_assignment" gorm:"foreignKey:AssetID"`
	Description   string          `json:"description"`
	Manufacturer  string          `json:"manufacturer"`
	Asset_Model   string          `json:"model"`
	SerialNumber  string          `json:"serial_number"`
	PurchaseDate  time.Time       `json:"purchase_date"`
	PurchasePrice string          `json:"purchase_price"`
	Vendor        string          `json:"vendor"`
	Site          string          `json:"site"`
	Status        string          `json:"status"`
	CreatedBy     int             `json:"created_by"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

// TableName sets the table name for the Asset model.
func (Assets) TableName() string {
	return "assets"
}

// Hashtag represents a hashtag entity
type AssetTag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	AssetTag  string    `json:"tag"`
	Tags      []string  `json:"tags"` // Added Tags field
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	AssetID   uint      `json:"-"`
}

// TableName sets the table name for the AssetTags model.
func (AssetTag) TableName() string {
	return "asset_tag"
}

type AssetType struct {
	gorm.Model
	ID        int       `gorm:"primaryKey" json:"asset_type_id"`
	AssetType string    `json:"asset_type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName sets the table name for the AssetType model.
func (AssetType) TableName() string {
	return "assetType"
}

type AssetAssignment struct {
	gorm.Model
	ID             int       `gorm:"primaryKey" json:"assignment_id"`
	AssetID        int       `json:"_"`
	UserID         int       `json:"user_id"`
	AssignedBy     int       `json:"assigned_by"`
	AssignmentType string    `json:"assignment_type"`
	DueAt          time.Time `json:"due_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// TableName sets the table name for the Asset Assignment model.
func (AssetAssignment) TableName() string {
	return "asset_assignment"
}

type AssetStorage interface {
	CreateTicket(*Ticket) error
	DeleteTicket(int) error
	UpdateTicket(*Ticket) error
	GetTickets() ([]*Ticket, error)
	GetTicketByID(int) (*Ticket, error)
	GetTicketByNumber(int) (*Ticket, error)
}

type AssetTypeStorage interface {
	CreateTicket(*Ticket) error
	DeleteTicket(int) error
	UpdateTicket(*Ticket) error
	GetTickets() ([]*Ticket, error)
	GetTicketByID(int) (*Ticket, error)
	GetTicketByNumber(int) (*Ticket, error)
}

type AssetAssignmentStorage interface {
	CreateTicket(*Ticket) error
	DeleteTicket(int) error
	UpdateTicket(*Ticket) error
	GetTickets() ([]*Ticket, error)
	GetTicketByID(int) (*Ticket, error)
	GetTicketByNumber(int) (*Ticket, error)
}

// AssetModel handles database operations for Asset
type AssetDBModel struct {
	DB *gorm.DB
}

// NewAssetModel creates a new instance of TicketModel
func NewAssetDBModel(db *gorm.DB) *AssetDBModel {
	return &AssetDBModel{
		DB: db,
	}
}
