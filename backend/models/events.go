// backend/models/events.go

package models

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Title       string    `gorm:"size:255;not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	AllDay      bool      `json:"all_day"`
	Location    string    `gorm:"size:255" json:"location"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	User        Users     `gorm:"foreignKey:UserID" json:"-"`
}

func (Event) TableName() string {
	return "event"
}

type LogEntry struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Level     string    `gorm:"size:50;not null" json:"level"` // Log level (e.g., INFO, WARN, ERROR)
	Message   string    `gorm:"type:text;not null" json:"message"`
	Context   string    `gorm:"type:text" json:"context"` // Additional context (JSON, XML, etc.)
}

func (LogEntry) TableName() string {
	return "log_entry"
}

type UserSetting struct {
	gorm.Model
	UserID       uint   `gorm:"not null;uniqueIndex:idx_user_setting" json:"user_id"`
	User         Users  `gorm:"foreignKey:UserID" json:"-"`
	SettingKey   string `gorm:"size:255;not null;uniqueIndex:idx_user_setting" json:"setting_key"`
	SettingValue string `gorm:"type:text;not null" json:"setting_value"` // Stored as JSON string for flexibility
}

func (UserSetting) TableName() string {
	return "user_setting"
}

type ApplicationSetting struct {
	gorm.Model
	Key         string `gorm:"size:255;not null;unique" json:"key"`
	Value       string `gorm:"type:text;not null" json:"value"` // Stored as JSON string for flexibility
	Description string `gorm:"size:255" json:"description"`     // Optional description of the setting
}

func (ApplicationSetting) TableName() string {
	return "application_setting"
}

type Project struct {
	gorm.Model
	Name        string  `gorm:"size:255;not null" json:"name"`
	Description string  `gorm:"type:text" json:"description"`
	OwnerID     uint    `gorm:"not null" json:"owner_id"`
	Owner       Users   `gorm:"foreignKey:OwnerID" json:"-"`
	Status      string  `gorm:"size:100;not null;default:'active'" json:"status"` // Example statuses: active, completed, archived
	TeamMembers []Users `gorm:"many2many:project_members;" json:"team_members"`
}

func (Project) TableName() string {
	return "project"
}

type Task struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	ProjectID   uint           `json:"project_id"`
	Title       string         `json:"title" binding:"required"`
	Description string         `json:"description"`
	AssigneeID  uint           `json:"assignee_id"`
	Status      string         `json:"status"`
	Priority    string         `json:"priority"`
	Deadline    time.Time      `json:"deadline"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Task) TableName() string {
	return "tasks"
}

type Notification struct {
	gorm.Model
	UserID  uint   `gorm:"not null;index" json:"user_id"`
	User    Users  `gorm:"foreignKey:UserID" json:"-"`
	Type    string `gorm:"size:100;not null" json:"type"` // e.g., email, in-app
	Title   string `gorm:"size:255;not null" json:"title"`
	Content string `gorm:"type:text;not null" json:"content"`
	Message string `gorm:"type:text;not null" json:"message"`
	IsRead  bool   `gorm:"default:false" json:"is_read"`
}

func (Notification) TableName() string {
	return "notification"
}

type NotificationSetting struct {
	gorm.Model
	UserID               uint   `gorm:"index;not null" json:"user_id"`
	EmailNotifications   bool   `json:"email_notifications"`                    // Whether the user opts in for email notifications
	SMSSettings          bool   `json:"sms_notifications"`                      // Whether the user opts in for SMS notifications
	PushNotifications    bool   `json:"push_notifications"`                     // Whether the user opts in for push notifications on devices
	NotificationChannels string `gorm:"type:text" json:"notification_channels"` // JSON-encoded string of custom notification channels
}

func (NotificationSetting) TableName() string {
	return "notification_setting"
}

type FileUpload struct {
	gorm.Model
	UserID     uint   `gorm:"not null;index" json:"user_id"`
	User       Users  `gorm:"foreignKey:UserID" json:"-"`
	FileName   string `gorm:"size:255;not null" json:"file_name"`
	FileType   string `gorm:"size:100;not null" json:"file_type"` // e.g., image/png, application/pdf
	FileSize   int64  `gorm:"not null" json:"file_size"`          // in bytes
	URL        string `gorm:"size:255;not null" json:"url"`       // URL to access the file
	Associated string `gorm:"size:255" json:"associated"`         // Optional: associated entity (e.g., "Project", "Task")
	EntityID   uint   `json:"entity_id"`                          // Optional: ID of the associated entity
}

func (FileUpload) TableName() string {
	return "file_upload"
}

type UserFeedback struct {
	gorm.Model
	UserID    uint   `gorm:"not null;index" json:"user_id"`
	User      Users  `gorm:"foreignKey:UserID" json:"-"`
	Feedback  string `gorm:"type:text;not null" json:"feedback"`
	ContactMe bool   `gorm:"default:false" json:"contact_me"` // Whether the user is open to being contacted for further discussion
}

func (UserFeedback) TableName() string {
	return "user_feedback"
}

type EventLog struct {
	gorm.Model
	UserID   uint   `gorm:"index" json:"user_id"` // Optional: associated user, if applicable
	User     Users  `gorm:"foreignKey:UserID" json:"-"`
	Event    string `gorm:"type:text;not null" json:"event"` // Description of the event
	Level    string `gorm:"size:50;not null" json:"level"`   // e.g., INFO, WARNING, ERROR
	Metadata string `gorm:"type:text" json:"metadata"`       // JSON-encoded metadata for additional context
}

func (EventLog) TableName() string {
	return "event_log"
}

type AppConfig struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Key   string `gorm:"size:255;unique;not null" json:"key"` // Configuration key
	Value string `gorm:"type:text;not null" json:"value"`     // Configuration value, stored as a string for flexibility
}

func (AppConfig) TableName() string {
	return "app_config"
}

type AuditLog struct {
	gorm.Model
	UserID    uint   `gorm:"index" json:"user_id"`
	Action    string `gorm:"type:text;not null" json:"action"` // Description of the action performed
	Entity    string `gorm:"size:255;not null" json:"entity"`  // Entity affected by the action
	EntityID  uint   `json:"entity_id"`                        // ID of the affected entity
	Details   string `gorm:"type:text" json:"details"`         // Detailed information about the action
	IP        string `gorm:"size:45" json:"ip"`                // IP address from which the action was performed
	UserAgent string `gorm:"type:text" json:"user_agent"`      // User agent of the browser/device used for the action
}

func (AuditLog) TableName() string {
	return "audit_log"
}

type UserPreference struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	UserID          uint           `json:"user_id"`
	PreferenceKey   string         `json:"preference_key"`
	PreferenceValue string         `json:"preference_value"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (UserPreference) TableName() string {
	return "user_preferences"
}

type Feedback struct {
	gorm.Model
	UserID     uint      `gorm:"index;not null" json:"user_id"`
	Feedback   string    `gorm:"type:text;not null" json:"feedback"`
	Category   string    `gorm:"size:255" json:"category"` // Example: 'Bug Report', 'Suggestion', 'Praise'
	Status     string    `gorm:"size:100" json:"status"`   // Example: 'New', 'Reviewed', 'Resolved'
	ResolvedAt time.Time `json:"resolved_at,omitempty"`
}

func (Feedback) TableName() string {
	return "feedback"
}

func (em *EventDBModel) CollectEventFeedback(eventID uint) (*EventFeedbackSummary, error) {
	// Aggregate feedback for a given event
}

func (em *EventDBModel) SendEventReminders(eventID uint) error {
	// Implementation to send reminders
}

func LogAuditEvent(entity string, action string, details string, performedBy uint) error {
    // Record an audit event in the system
}

func ApplyRateLimiting(agentID uint, resource string) error {
    // Implement rate limiting
