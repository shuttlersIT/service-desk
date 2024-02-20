// backend/models/incident_db_model.go

package models

import (
	"time"

	"gorm.io/gorm"
)

// Incident represents an incident report.
type Analytics struct {
	gorm.Model
	Name           string    `gorm:"size:255;not null" json:"name"`
	Value          float64   `json:"value" gorm:"type:decimal(10,2)"`
	ReportedAt     time.Time `json:"reported_at" gorm:"not null"`
	UserID         uint      `json:"user_id" gorm:"not null;index"`
	Title          string    `json:"title" gorm:"size:255;not null"`
	Description    string    `json:"description" gorm:"type:text;not null"`
	Category       string    `json:"category" gorm:"size:100;not null;index"`
	Priority       string    `json:"priority" gorm:"size:50;not null"`
	Tags           []Tag     `json:"tags" gorm:"many2many:analytics_tags;"`
	AttachmentURL  string    `json:"attachment_url" gorm:"size:255"`
	HasAttachments bool      `json:"has_attachments"`
	Severity       string    `json:"severity" gorm:"size:50;not null"`
}

func (Analytics) TableName() string {
	return "engagement_metrics"
}

// EngagementMetrics handles database operations for incidents.
type AnalyticsDBModel struct {
	DB  *gorm.DB
	log Logger
}

// NewIncidentDBModel creates a new instance of EngagementMetricsDBModel.
func NewEngagementMetricsDBModel(db *gorm.DB, log Logger) *AnalyticsDBModel {
	return &AnalyticsDBModel{
		DB:  db,
		log: log,
	}
}

// GetUserEngagementMetrics provides analytics on user engagement and system interaction.
func (db *AnalyticsDBModel) GetUserEngagementMetrics(userID uint) (*AnalyticsDBModel, error) {
	var metrics AnalyticsDBModel
	// Example SQL might involve complex JOINs and aggregation functions
	err := db.DB.Raw(`
        SELECT COUNT(*) AS total_actions, AVG(response_time) AS average_response_time
        FROM user_actions
        WHERE user_id = ?
    `, userID).Scan(&metrics).Error
	if err != nil {
		return nil, err
	}
	return &metrics, nil
}

// OptimizeSystemResources suggests optimizations based on usage patterns and resource consumption.
func (db *OptimizationDBModel) OptimizeSystemResources(threshold float64) ([]OptimizationSuggestion, error) {
	var suggestions []OptimizationSuggestion
	// This query identifies resources whose average usage is below a certain threshold,
	// indicating they might be under-utilized.
	err := db.DB.Raw(`
        SELECT resource_id, SUM(usage) AS total_usage, AVG(usage) AS avg_usage
        FROM resource_usage
        GROUP BY resource_id
        HAVING AVG(usage) < ?
    `, threshold).Scan(&suggestions).Error
	if err != nil {
		return nil, err
	}
	return suggestions, nil
}

type PerformanceMetric struct {
	gorm.Model
	MetricName string  `gorm:"type:varchar(255);not null"`
	Value      float64 `gorm:"not null"`
	RecordedAt time.Time
}

type AgentPerformanceMetric struct {
	gorm.Model
	AgentID     uint    `json:"agent_id" gorm:"index;not null"`
	MetricName  string  `json:"metric_name" gorm:"type:varchar(255);not null"`
	Value       float64 `json:"value" gorm:"type:decimal(10,2);not null"`
	TargetValue float64 `json:"target_value" gorm:"type:decimal(10,2)"`
	Period      string  `json:"period" gorm:"type:varchar(100);not null"` // E.g., "Monthly", "Quarterly"
}

func (PerformanceMetric) TableName() string {
	return "performance_metrics"
}

type OptimizationSuggestion struct {
	ResourceID   uint    `json:"resource_id"`
	TotalUsage   float64 `json:"total_usage"`
	AverageUsage float64 `json:"average_usage"`
}

type OptimizationDBModel struct {
	DB *gorm.DB
}

func NewOptimizationDBModel(db *gorm.DB) *OptimizationDBModel {
	return &OptimizationDBModel{DB: db}
}
