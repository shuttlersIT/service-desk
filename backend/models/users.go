// backend/models/users.go

package models

import (
	"time"

	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	ID          int                   `gorm:"primaryKey" json:"user_id"`
	FirstName   string                `json:"first_name" binding:"required"`
	LastName    string                `json:"last_name" binding:"required"`
	Email       string                `json:"staff_email" binding:"required,email"`
	Credentials UsersLoginCredentials `json:"user_credentials" gorm:"foreignKey:UserID"`
	Phone       string                `json:"phoneNumber" binding:"required,e164"`
	Position    Position              `json:"position_id" gorm:"embedded"`
	Department  Department            `json:"department_id" gorm:"embedded"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
	Asset       []AssetAssignment     `json:"asset_assignment" gorm:"foreignKey:UserID"`
}

// TableName sets the table name for the Users model.
func (Users) TableName() string {
	return "users"
}

type UsersLoginCredentials struct {
	gorm.Model
	ID        int       `gorm:"primaryKey" json:"_"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName sets the table name for the UsersLoginCredentials model.
func (UsersLoginCredentials) TableName() string {
	return "usersLoginCredentials"
}

type Position struct {
	gorm.Model
	PositionID   int       `gorm:"primaryKey" json:"position_id"`
	PositionName string    `json:"position_name"`
	CadreName    string    `json:"cadre_name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName sets the table name for the Position model.
func (Position) TableName() string {
	return "position"
}

type Department struct {
	gorm.Model
	DepartmentID   int       `gorm:"primaryKey" json:"department_id"`
	DepartmentName string    `json:"department_name"`
	Emoji          string    `json:"emoji"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// TableName sets the table name for the Department model.
func (Department) TableName() string {
	return "department"
}

type UserStorage interface {
	CreateTicket(*Ticket) error
	DeleteTicket(int) error
	UpdateTicket(*Ticket) error
	GetTickets() ([]*Ticket, error)
	GetTicketByID(int) (*Ticket, error)
	GetTicketByNumber(int) (*Ticket, error)
}

type PositionStorage interface {
	CreateTicket(*Ticket) error
	DeleteTicket(int) error
	UpdateTicket(*Ticket) error
	GetTickets() ([]*Ticket, error)
	GetTicketByID(int) (*Ticket, error)
	GetTicketByNumber(int) (*Ticket, error)
}

type DepartmentStorage interface {
	CreateTicket(*Ticket) error
	DeleteTicket(int) error
	UpdateTicket(*Ticket) error
	GetTickets() ([]*Ticket, error)
	GetTicketByID(int) (*Ticket, error)
	GetTicketByNumber(int) (*Ticket, error)
}

type UserPasswordLoginStorage interface {
	CreateTicket(*Ticket) error
	DeleteTicket(int) error
	UpdateTicket(*Ticket) error
	GetTickets() ([]*Ticket, error)
	GetTicketByID(int) (*Ticket, error)
	GetTicketByNumber(int) (*Ticket, error)
}

// UserModel handles database operations for User
type UserDBModel struct {
	DB *gorm.DB
}

// NewUserModel creates a new instance of UserModel
func NewUserDBModel(db *gorm.DB) *UserDBModel {
	return &UserDBModel{
		DB: db,
	}
}
