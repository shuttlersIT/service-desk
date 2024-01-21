package models

import (
	"time"
)

type Incident struct {
	ID          uint      `gorm:"primaryKey" json:"incident_id"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Severity    string    `json:"severity" binding:"required"`
	TeamID      uint      `json:"team_id"`
	Status      string    `json:"status" binding:"required"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Additional fields and methods can be added as needed.
