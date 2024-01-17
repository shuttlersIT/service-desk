//

package models

import (
	"gorm.io/gorm"
)

type GoogleCredentials struct {
	Cid     string `json:"cid"`
	Csecret string `json:"csecret"`
}

// TableName sets the table name for the Credentials model.
func (GoogleCredentials) TableName() string {
	return "googleCredentials"
}

type GoogeAuthStorage interface {
	CreateTicket(*Ticket) error
	DeleteTicket(int) error
	UpdateTicket(*Ticket) error
	GetTickets() ([]*Ticket, error)
	GetTicketByID(int) (*Ticket, error)
	GetTicketByNumber(int) (*Ticket, error)
}

// UserModel handles database operations for User
type GoogleCredentialsDBModel struct {
	DB *gorm.DB
}

// NewUserModel creates a new instance of UserModel
func NewGoogleCredentialsDBModel(db *gorm.DB) *GoogleCredentialsDBModel {
	return &GoogleCredentialsDBModel{
		DB: db,
	}
}
