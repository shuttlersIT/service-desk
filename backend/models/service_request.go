package models

import (
	"time"
)

type ServiceRequest struct {
	ID          uint      `gorm:"primaryKey" json:"service_request_id"`
	UserID      uint      `json:"user_id" binding:"required"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Status      string    `json:"status" binding:"required"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Additional fields and methods can be added as needed.
