// backend/models/events.go

package models

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Title       string    `gorm:"size:255;not null" json:"title"` // Title of the event
	Description string    `gorm:"type:text" json:"description"`   // Detailed description of the event
	StartTime   time.Time `json:"start_time" gorm:"not null"`     // Start time of the event
	EndTime     time.Time `json:"end_time" gorm:"not null"`       // End time of the event
	AllDay      bool      `json:"all_day"`                        // Indicates if the event lasts all day
	Location    string    `gorm:"size:255" json:"location"`       // Location of the event
	UserID      uint      `gorm:"not null;index" json:"user_id"`  // ID of the user who created the event
	AgentID     uint      `gorm:"not null;index" json:"agent_id"` // ID of the agent associated with the event, if any
	OrganizerID uint      `gorm:"index" json:"organizer_id"`      // ID of the organizer of the event, if different from the user
}

func (Event) TableName() string {
	return "event"
}

type LogEntry struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Level     string    `gorm:"size:50;not null" json:"level"`     // Severity level of the log entry (INFO, WARN, ERROR)
	Message   string    `gorm:"type:text;not null" json:"message"` // The log message
	Context   string    `gorm:"type:text" json:"context"`          // Additional context for the log entry, possibly in JSON format
}

func (LogEntry) TableName() string {
	return "log_entry"
}

type UserSetting struct {
	gorm.Model
	UserID       uint   `gorm:"not null;uniqueIndex:idx_user_setting" json:"user_id"`              // ID of the user to whom the setting applies
	SettingKey   string `gorm:"size:255;not null;uniqueIndex:idx_user_setting" json:"setting_key"` // The key or name of the setting
	SettingValue string `gorm:"type:text;not null" json:"setting_value"`                           // The value of the setting, stored as a string for flexibility
}

func (UserSetting) TableName() string {
	return "user_setting"
}

type ApplicationSetting struct {
	gorm.Model
	Key         string `gorm:"size:255;not null;unique" json:"key"` // The key or name of the application-wide setting
	Value       string `gorm:"type:text;not null" json:"value"`     // The value of the setting, allowing for complex configurations stored as strings
	Description string `gorm:"size:255" json:"description"`         // Optional description of what the setting controls or affects
}

func (ApplicationSetting) TableName() string {
	return "application_setting"
}

type Project struct {
	gorm.Model
	Name        string  `gorm:"size:255;not null" json:"name"`                    // Name of the project
	Description string  `gorm:"type:text" json:"description"`                     // Detailed description of the project
	OwnerID     uint    `gorm:"not null" json:"owner_id"`                         // ID of the user who owns or manages the project
	Status      string  `gorm:"size:100;not null;default:'active'" json:"status"` // Current status of the project (e.g., active, completed, archived)
	TeamMembers []Users `gorm:"many2many:project_members;" json:"team_members"`   // Users assigned as team members of the project
}

func (Project) TableName() string {
	return "project"
}

type Task struct {
	gorm.Model
	ProjectID   uint      `json:"project_id"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	AssigneeID  uint      `json:"assignee_id"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	Deadline    time.Time `json:"deadline"`
}

func (Task) TableName() string {
	return "tasks"
}

type Notification struct {
	gorm.Model
	UserID  uint   `gorm:"not null;index" json:"user_id"`     // The user who will receive the notification
	Type    string `gorm:"size:100;not null" json:"type"`     // The type of notification (e.g., email, in-app)
	Title   string `gorm:"size:255;not null" json:"title"`    // The title of the notification
	Content string `gorm:"type:text;not null" json:"content"` // The main content of the notification
	Message string `gorm:"type:text;not null" json:"message"` // An additional message or information related to the notification
	IsRead  bool   `gorm:"default:false" json:"is_read"`      // Flag indicating whether the notification has been read
	Seen    bool   `gorm:"default:false" json:"seen"`         // Flag indicating whether the notification has been seen (used for in-app notifications)
}

func (Notification) TableName() string {
	return "notification"
}

type NotificationSetting struct {
	gorm.Model
	UserID               uint   `gorm:"index;not null" json:"user_id"`          // ID of the user to whom the settings apply
	EmailNotifications   bool   `json:"email_notifications"`                    // Whether the user opts in for email notifications
	SMSSettings          bool   `json:"sms_notifications"`                      // Whether the user opts in for SMS notifications
	PushNotifications    bool   `json:"push_notifications"`                     // Whether the user opts in for push notifications
	NotificationChannels string `gorm:"type:text" json:"notification_channels"` // JSON-encoded string detailing other notification channels the user has opted into
}

func (NotificationSetting) TableName() string {
	return "notification_setting"
}

type FileUpload struct {
	gorm.Model
	UserID     uint   `gorm:"not null;index" json:"user_id"`      // ID of the user who uploaded the file
	FileName   string `gorm:"size:255;not null" json:"file_name"` // Original name of the file
	FileType   string `gorm:"size:100;not null" json:"file_type"` // MIME type of the file
	FileSize   int64  `gorm:"not null" json:"file_size"`          // Size of the file in bytes
	URL        string `gorm:"size:255;not null" json:"url"`       // URL where the file can be accessed
	Associated string `gorm:"size:255" json:"associated"`         // Optionally associates the file with an entity (e.g., "Project", "Task")
	EntityID   uint   `json:"entity_id"`                          // ID of the associated entity, if any
}

func (FileUpload) TableName() string {
	return "file_upload"
}

type UserFeedback struct {
	gorm.Model
	UserID    uint   `gorm:"not null;index" json:"user_id"`      // ID of the user providing feedback
	Feedback  string `gorm:"type:text;not null" json:"feedback"` // The content of the feedback
	ContactMe bool   `gorm:"default:false" json:"contact_me"`    // Whether the user consents to be contacted for further information
	Content   string `gorm:"type:text;not null"`                 // Additional content or context related to the feedback
	Response  string `gorm:"type:text" json:"response"`          // An optional response to the feedback, possibly from an admin or support team
}

func (UserFeedback) TableName() string {
	return "user_feedback"
}

type EventLog struct {
	gorm.Model
	UserID    uint   `gorm:"index" json:"user_id"`            // Optional: ID of the user involved in the event
	Event     string `gorm:"type:text;not null" json:"event"` // Detailed description of the event
	Level     string `gorm:"size:50;not null" json:"level"`   // Severity level of the event (e.g., INFO, WARNING, ERROR)
	Metadata  string `gorm:"type:text" json:"metadata"`       // JSON-encoded metadata providing additional context
	EventType string `gorm:"type:varchar(100);not null"`      // Categorization of the event type
	Payload   string `gorm:"type:text;not null"`              // JSON payload detailing the event specifics
	Status    string `gorm:"type:varchar(100);not null"`      // Current processing status of the event (e.g., New, Processed, Error)
}

func (EventLog) TableName() string {
	return "event_log"
}

type AppConfig struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Key         string `gorm:"size:255;unique;not null" json:"key"` // Unique key for the configuration setting
	Value       string `gorm:"type:text;not null" json:"value"`     // Value of the configuration setting
	Description string `gorm:"size:255" json:"description"`         // Optional description of the configuration setting
}

func (AppConfig) TableName() string {
	return "app_config"
}

type AuditLog struct {
	gorm.Model
	UserID    uint   `gorm:"index" json:"user_id"`             // ID of the user performing the action
	Action    string `gorm:"type:text;not null" json:"action"` // Description of the action performed
	Entity    string `gorm:"size:255;not null" json:"entity"`  // The entity affected by the action
	EntityID  uint   `json:"entity_id"`                        // ID of the affected entity
	Details   string `gorm:"type:text" json:"details"`         // Detailed information about the action
	IP        string `gorm:"size:45" json:"ip"`                // IP address from which the action was performed
	UserAgent string `gorm:"type:text" json:"user_agent"`      // User agent of the browser/device used
}

func (AuditLog) TableName() string {
	return "audit_log"
}

type UserPreference struct {
	gorm.Model
	UserID          uint   `json:"user_id"`
	PreferenceKey   string `json:"preference_key"`
	PreferenceValue string `json:"preference_value"`
	Language        string `gorm:"type:varchar(10);not null" json:"language"` // ISO language code
	TimeZone        string `gorm:"type:varchar(100);not null" json:"time_zone"`
}

func (UserPreference) TableName() string {
	return "user_preferences"
}

type Feedback struct {
	gorm.Model
	UserID     uint      `gorm:"index;not null" json:"user_id"`      // User ID associated with the feedback
	Feedback   string    `gorm:"type:text;not null" json:"feedback"` // The content of the feedback
	Category   string    `gorm:"size:255" json:"category"`           // Type of feedback (e.g., Bug Report, Suggestion)
	Status     string    `gorm:"size:100" json:"status"`             // Status of the feedback (e.g., New, Reviewed, Resolved)
	ResolvedAt time.Time `json:"resolved_at,omitempty"`              // Date and time when the feedback was resolved
}

func (Feedback) TableName() string {
	return "feedback"
}

// GenerateUserActivityReport aggregates user activities over a specified period.
func (db *UserDBModel) GenerateUserActivityReport(userID uint, startDate, endDate time.Time) ([]UserActivityLog, error) {
	var reports []UserActivityLog
	return reports, db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Raw(`
            SELECT activity_type, COUNT(*) as count, DATE(activity_date) as date 
            FROM activities 
            WHERE user_id = ? AND (activity_date BETWEEN ? AND ?) 
            GROUP BY activity_type, DATE(activity_date)
        `, userID, startDate, endDate).Scan(&reports).Error; err != nil {
			return err
		}
		return nil
	})
}

// AllocateResources dynamically allocates resources based on current demand.
func (db *AssetDBModel) AllocateResources(requestID uint, resourcesNeeded int) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		// Example logic for resource allocation
		if err := tx.Model(&Resource{}).Where("available > ?", resourcesNeeded).Update("available", gorm.Expr("available - ?", resourcesNeeded)).Error; err != nil {
			return err
		}
		return nil
	})
}

type ProductReview struct {
	ID        uint   `gorm:"primaryKey"`
	ProductID uint   `gorm:"index;not null" json:"product_id"`  // ID of the product being reviewed
	UserID    uint   `gorm:"index;not null" json:"user_id"`     // ID of the user submitting the review
	Review    string `gorm:"type:text;not null" json:"review"`  // Content of the review
	Rating    int    `gorm:"check:rating >= 1 AND rating <= 5"` // Rating given by the user, on a scale of 1 to 5
	CreatedAt time.Time
}

func (ProductReview) TableName() string {
	return "product_reviews"
}

type ScheduledTask struct {
	ID       uint       `gorm:"primaryKey" json:"id"`
	Name     string     `gorm:"type:varchar(255);not null" json:"name"`     // Unique name of the scheduled task
	Type     string     `gorm:"type:varchar(100);not null" json:"type"`     // Type of task, e.g., "email", "cleanup"
	Schedule string     `gorm:"type:varchar(100);not null" json:"schedule"` // Scheduling format, e.g., cron expression
	LastRun  *time.Time `json:"last_run,omitempty"`                         // Timestamp of the last run
	NextRun  time.Time  `json:"next_run" gorm:"not null"`                   // Timestamp of the next scheduled run
	IsActive bool       `gorm:"default:true" json:"is_active"`              // Indicates if the task is active
}

func (ScheduledTask) TableName() string {
	return "scheduled_tasks"
}

type SecurityEventLog struct {
	gorm.Model
	EventType string    `gorm:"type:varchar(100);not null"` // E.g., "login_attempt", "password_change"
	UserID    uint      `gorm:"index"`                      // Optional, not all events may be user-specific
	Details   string    `gorm:"type:text"`                  // JSON or similar structured format recommended
	LogTime   time.Time `json:"log_time,omitempty"`
}

func (SecurityEventLog) TableName() string {
	return "security_event_logs"
}

type Poll struct {
	gorm.Model
	Question string     `gorm:"type:varchar(255);not null" json:"question"` // The question or statement of the poll
	Options  string     `gorm:"type:text;not null" json:"options"`          // JSON-encoded array of options for the poll
	IsActive bool       `gorm:"default:true" json:"is_active"`              // Whether the poll is currently active
	EndTime  *time.Time `json:"end_time,omitempty"`                         // Optional end time for when the poll closes

}

func (Poll) TableName() string {
	return "polls"
}

type PollVote struct {
	ID      uint      `gorm:"primaryKey" json:"id"`
	PollID  uint      `gorm:"index;not null" json:"poll_id"`            // The poll to which the vote belongs
	UserID  uint      `gorm:"index;not null" json:"user_id"`            // The user who cast the vote
	Option  string    `gorm:"type:varchar(255);not null" json:"option"` // The chosen option
	VotedAt time.Time `json:"voted_at"`                                 // Timestamp of when the vote was cast
}

func (PollVote) TableName() string {
	return "poll_votes"
}

type ServiceDeskEvent struct {
	gorm.Model
	EventName       string    `json:"event_name" gorm:"type:varchar(255);not null"`
	Description     string    `json:"description" gorm:"type:text"`
	EventDate       time.Time `json:"event_date"`
	Severity        string    `json:"severity" gorm:"type:varchar(100);not null"` // E.g., "Low", "Medium", "High"
	AffectedSystems string    `json:"affected_systems" gorm:"type:text"`          // JSON array of affected systems
}

type EventsDBModel struct {
	DB  *gorm.DB
	log Logger
}

func NewEventsDBModel(db *gorm.DB, log Logger) *EventsDBModel {
	return &EventsDBModel{
		DB:  db,
		log: log,
	}
}

type TicketCommentCreatedEvent struct {
	CommentID   uint      `gorm:"primaryKey" json:"commentId"`
	TicketID    uint      `gorm:"not null" json:"ticketId"`
	CommenterID uint      `gorm:"not null" json:"commenterId"`
	CommentText string    `gorm:"type:text;not null" json:"commentText"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

// UserActivityLog records activities performed by users within the system.
type UserActivityLog struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	UserID       uint           `gorm:"index;notNull" json:"user_id"`
	Action       string         `gorm:"type:varchar(255);notNull" json:"action"`
	Description  string         `gorm:"type:text" json:"description,omitempty"`
	ActivityType string         `json:"activity_type" gorm:"type:varchar(255);not null"`
	Details      string         `json:"details" gorm:"type:text;not null"`
	Timestamp    time.Time      `json:"timestamp"` // Timestamp of when the activity occurred
	OccurredAt   time.Time      `json:"occurred_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (UserActivityLog) TableName() string {
	return "user_activity_logs"
}

// AssetLocation defines physical or logical locations where assets are stored or used.
type AssetLocation struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:255;notNull;unique" json:"name"`
	Description string         `gorm:"type:text" json:"description,omitempty"`
	Address     string         `gorm:"type:text" json:"address,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// AssetReservation allows users to reserve assets for specific periods.
type AssetReservation struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	AssetID   uint           `gorm:"index;notNull" json:"asset_id"`
	UserID    uint           `gorm:"index;notNull" json:"user_id"`
	StartTime time.Time      `json:"start_time"`
	EndTime   time.Time      `json:"end_time"`
	Status    string         `gorm:"size:100;notNull" json:"status"` // e.g., Reserved, Cancelled
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// AssetCheckInOut records check-in and check-out activities for assets.
type AssetCheckInOut struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	AssetID    uint           `gorm:"index;notNull" json:"asset_id"`
	UserID     uint           `gorm:"index;notNull" json:"user_id"`
	Action     string         `gorm:"type:varchar(100);notNull" json:"action"` // CheckIn, CheckOut
	OccurredAt time.Time      `json:"occurred_at"`
	Notes      string         `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// SystemSetting represents configurable settings within the asset management system.
type SystemSetting struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Key       string         `gorm:"size:255;notNull;unique" json:"key"`
	Value     string         `gorm:"type:text;notNull" json:"value"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// AssetTransferLog records the transfer of assets between different locations or users.
type AssetTransferLog struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	AssetID       uint           `gorm:"index;notNull" json:"asset_id"`
	FromLocation  uint           `gorm:"index" json:"from_location"`
	ToLocation    uint           `gorm:"index" json:"to_location"`
	TransferDate  time.Time      `json:"transfer_date"`
	TransferredBy uint           `gorm:"index;notNull" json:"transferred_by"`
	Notes         string         `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// UserNotification stores notifications sent to users within the system.
type UserNotification struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UserID     uint           `gorm:"index;notNull" json:"user_id"`
	Message    string         `gorm:"type:text;notNull" json:"message"`
	Read       bool           `gorm:"notNull;default:false" json:"read"`
	NotifiedAt time.Time      `json:"notified_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// MaintenanceRequest captures requests for maintenance of assets, including details of the request and status.
type MaintenanceRequest struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	AssetID         uint           `gorm:"index;notNull" json:"asset_id"`
	RequestedBy     uint           `gorm:"index;notNull" json:"requested_by"`
	Description     string         `gorm:"type:text;notNull" json:"description"`
	Status          string         `gorm:"size:100;notNull" json:"status"` // e.g., Pending, InProgress, Completed
	MaintenanceDate *time.Time     `json:"maintenance_date,omitempty"`
	CompletedDate   *time.Time     `json:"completed_date,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// AssetWarrantyDetails holds warranty information for assets, including start and end dates.
type AssetWarrantyDetails struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	AssetID       uint           `gorm:"index;notNull" json:"asset_id"`
	WarrantyStart time.Time      `json:"warranty_start"`
	WarrantyEnd   time.Time      `json:"warranty_end"`
	Provider      string         `gorm:"size:255" json:"provider"`
	WarrantyTerms string         `gorm:"type:text" json:"warranty_terms"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
