package models

import (
	"time"

	"gorm.io/gorm"
)

type Workflow struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	Name        string          `gorm:"type:varchar(255);not null" json:"name"`
	Description string          `gorm:"type:text" json:"description,omitempty"`
	IsActive    bool            `gorm:"default:true" json:"is_active"`
	Steps       []WorkflowStep  `gorm:"foreignKey:WorkflowID" json:"steps,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Workflow) TableName() string {
	return "workflows"
}

type WorkflowStep struct {
	ID            uint            `gorm:"primaryKey" json:"id"`
	WorkflowID    uint            `json:"workflow_id" gorm:"index;not null"`
	StepName      string          `gorm:"type:varchar(255);not null" json:"step_name"`
	StepOrder     int             `gorm:"not null" json:"step_order"`
	ActionType    string          `gorm:"type:varchar(100);not null" json:"action_type"`
	ActionDetails string          `gorm:"type:text" json:"action_details,omitempty"` // JSON detailing the action to be taken
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	DeletedAt     *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (WorkflowStep) TableName() string {
	return "workflow_steps"
}
