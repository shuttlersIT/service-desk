// backend/models/agents.go

package models

import (
	"time"

	"gorm.io/gorm"
)

type Agents struct {
	gorm.Model
	ID           int                   `gorm:"primaryKey" json:"agent_id"`
	FirstName    string                `json:"first_name" binding:"required"`
	LastName     string                `json:"last_name" binding:"required"`
	AgentEmail   string                `json:"agent_email" binding:"required,email"`
	Credentials  AgentLoginCredentials `json:"agent_credentials" gorm:"foreignKey:AgentID"`
	Phone        string                `json:"phoneNumber" binding:"required,e164"`
	RoleID       Role                  `json:"role_id" gorm:"embedded"`
	Unit         Unit                  `json:"unit" gorm:"embedded"`
	SupervisorID int                   `json:"supervisor_id"`
	CreatedAt    time.Time             `json:"created_at"`
	UpdatedAt    time.Time             `json:"updated_at"`
}

// TableName sets the table name for the Agent model.
func (Agents) TableName() string {
	return "agents"
}

type AgentLoginCredentials struct {
	gorm.Model
	ID        int       `gorm:"primaryKey" json:"_"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	AgentID   int       `json:"agent_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName sets the table name for the Agent model.
func (AgentLoginCredentials) TableName() string {
	return "agentLoginDetails"
}

type Unit struct {
	gorm.Model
	ID        int       `gorm:"primaryKey" json:"unit_id"`
	UnitName  string    `json:"unit_name"`
	Emoji     string    `json:"emoji"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName sets the table name for the Agent model.
func (Unit) TableName() string {
	return "unit"
}

type Role struct {
	gorm.Model
	ID        int       `gorm:"primaryKey" json:"role_id"`
	RoleName  string    `json:"role_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName sets the table name for the Agent model.
func (Role) TableName() string {
	return "role"
}

type AgentStorage interface {
	CreateTicket(*Ticket) error
	DeleteTicket(int) error
	UpdateTicket(*Ticket) error
	GetTickets() ([]*Ticket, error)
	GetTicketByID(int) (*Ticket, error)
	GetTicketByNumber(int) (*Ticket, error)
}

type AgentPasswordLoginStorage interface {
	CreateTicketOperation(*Ticket) error
	DeleteTicket(int) error
	UpdateTicket(*Ticket) error
	GetTickets() ([]*Ticket, error)
	GetTicketByID(int) (*Ticket, error)
	GetTicketByNumber(int) (*Ticket, error)
}

type UnitStorage interface {
	CreateTicket(*Ticket) error
	DeleteTicket(int) error
	UpdateTicket(*Ticket) error
	GetTickets() ([]*Ticket, error)
	GetTicketByID(int) (*Ticket, error)
	GetTicketByNumber(int) (*Ticket, error)
}

type RoleStorage interface {
	CreateTicket(*Ticket) error
	DeleteTicket(int) error
	UpdateTicket(*Ticket) error
	GetTickets() ([]*Ticket, error)
	GetTicketByID(int) (*Ticket, error)
	GetTicketByNumber(int) (*Ticket, error)
}

// AgentModel handles database operations for Agent
type AgentDBModel struct {
	DB *gorm.DB
}

// NewAgentModel creates a new instance of Agent
func NewAgentDBModel(db *gorm.DB) *AgentDBModel {
	return &AgentDBModel{
		DB: db,
	}
}
