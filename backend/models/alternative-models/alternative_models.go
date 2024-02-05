package models

import (
	"time"
)

// Department represents the departments table.
type Department struct {
	DepartmentID   int        `json:"department_id" gorm:"primary_key;auto_increment"`
	DepartmentName string     `json:"department_name" gorm:"type:varchar(255);unique;not_null"`
	CreatedAt      time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty" gorm:"default:NULL"`
}

// Position represents the positions table.
type Position struct {
	PositionID   int        `json:"position_id" gorm:"primary_key;auto_increment"`
	PositionName string     `json:"position_name" gorm:"type:varchar(255);unique;not_null"`
	CreatedAt    time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" gorm:"default:NULL"`
}

// Vendor represents the vendors table.
type Vendor struct {
	VendorID    int        `json:"vendor_id" gorm:"primary_key;auto_increment"`
	VendorName  string     `json:"vendor_name" gorm:"type:varchar(255);unique;not_null"`
	Description string     `json:"description,omitempty" gorm:"type:text"`
	ContactInfo string     `json:"contact_info,omitempty" gorm:"type:text"`
	Address     string     `json:"address,omitempty" gorm:"type:text"`
	CreatedAt   time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"default:NULL"`
}

// Category represents the categories table.
type Category struct {
	CategoryID   int        `json:"category_id" gorm:"primary_key;auto_increment"`
	CategoryName string     `json:"category_name" gorm:"type:varchar(255);not_null"`
	Icon         string     `json:"icon,omitempty" gorm:"type:varchar(255)"`
	Description  string     `json:"description,omitempty" gorm:"type:text"`
	CreatedAt    time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" gorm:"default:NULL"`
}

// Priority represents the priority table for SLA and tickets.
type Priority struct {
	PriorityID  int        `json:"priority_id" gorm:"primary_key;auto_increment"`
	Name        string     `json:"name" gorm:"type:varchar(255);not_null"`
	Description string     `json:"description,omitempty" gorm:"type:text"`
	Colour      string     `json:"colour" gorm:"type:varchar(6);default:'#FFFFFF'"`
	CreatedAt   time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"default:NULL"`
}

// Status represents the status table for tickets.
type Status struct {
	StatusID   int        `json:"status_id" gorm:"primary_key;auto_increment"`
	StatusName string     `json:"status_name" gorm:"type:varchar(255);not_null"`
	IsClosed   bool       `json:"is_closed" gorm:"default:false"`
	CreatedAt  time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty" gorm:"default:NULL"`
}

// Permission represents the permissions table.
type Permission struct {
	ID          int        `json:"id" gorm:"primary_key;auto_increment"`
	Name        string     `json:"name" gorm:"type:varchar(255);unique;not_null"`
	Description string     `json:"description,omitempty" gorm:"type:text"`
	CreatedAt   time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"default:NULL"`
}

// Role represents the roles table.
type Role struct {
	ID        int        `json:"id" gorm:"primary_key;auto_increment"`
	RoleName  string     `json:"role_name" gorm:"type:varchar(255);unique;not_null"`
	CreatedAt time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"default:NULL"`
}

// User represents the users table, including foreign keys to positions and departments.
type User struct {
	UserID       int        `json:"user_id" gorm:"primary_key;auto_increment"`
	FirstName    string     `json:"first_name" gorm:"type:varchar(255);not_null"`
	LastName     string     `json:"last_name" gorm:"type:varchar(255);not_null"`
	Email        string     `json:"email" gorm:"type:varchar(255);unique;not_null"`
	Phone        string     `json:"phone,omitempty" gorm:"type:varchar(255)"`
	PositionID   *int       `json:"position_id,omitempty" gorm:"default:NULL"`
	DepartmentID *int       `json:"department_id,omitempty" gorm:"default:NULL"`
	IsActive     bool       `json:"is_active" gorm:"default:true"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty" gorm:"default:NULL"`
	CreatedAt    time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" gorm:"default:NULL"`
}

// Agent represents the agents table, including references to roles, teams, units, and supervisors.
type Agent struct {
	ID           int        `json:"id" gorm:"primary_key;auto_increment"`
	FirstName    string     `json:"first_name" gorm:"type:varchar(255);not_null"`
	LastName     string     `json:"last_name" gorm:"type:varchar(255);not_null"`
	AgentEmail   string     `json:"agent_email" gorm:"type:varchar(255);unique;not_null"`
	Phone        string     `json:"phone,omitempty" gorm:"type:varchar(255);unique"`
	RoleID       *int       `json:"role_id,omitempty" gorm:"default:NULL"`
	TeamID       *int       `json:"team_id,omitempty" gorm:"default:NULL"`
	UnitID       *int       `json:"unit_id,omitempty" gorm:"default:NULL"`
	SupervisorID *int       `json:"supervisor_id,omitempty" gorm:"default:NULL"`
	CreatedAt    time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" gorm:"default:NULL"`
}

// Asset represents the assets table.
type Asset struct {
	ID                 int        `json:"id" gorm:"primary_key;auto_increment"`
	AssetTag           int        `json:"asset_tag" gorm:"type:int"`
	AssetName          string     `json:"asset_name" gorm:"type:varchar(255);not_null"`
	AssetType          string     `json:"asset_type" gorm:"type:varchar(255);not_null"`
	AssetTypeID        int        `json:"asset_type_id" gorm:"not_null"`
	Description        string     `json:"description,omitempty" gorm:"type:text"`
	VendorID           *int       `json:"vendor_id,omitempty" gorm:"default:NULL"`
	PurchaseDate       *time.Time `json:"purchase_date,omitempty" gorm:"type:date"`
	SerialNumber       string     `json:"serial_number" gorm:"type:varchar(255);unique"`
	Status             string     `json:"status" gorm:"type:varchar(50);default:'active'"`
	WarrantyExpiration *time.Time `json:"warranty_expiration,omitempty" gorm:"type:date"`
	Location           string     `json:"location,omitempty" gorm:"type:varchar(255)"`
	UserID             *int       `json:"user_id,omitempty" gorm:"default:NULL"`
	CreatedAt          time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt          time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty" gorm:"default:NULL"`
}

// ServiceRequest represents the service_requests table.
type ServiceRequest struct {
	RequestID   int        `json:"request_id" gorm:"primary_key;auto_increment"`
	Title       string     `json:"title" gorm:"type:varchar(255);not_null"`
	Description string     `json:"description" gorm:"type:text"`
	UserID      int        `json:"user_id" gorm:"not_null"`
	Status      string     `json:"status" gorm:"type:varchar(50);not_null"`
	CategoryID  int        `json:"category_id" gorm:"not_null"`
	Priority    string     `json:"priority,omitempty" gorm:"type:varchar(50)"`
	CreatedAt   time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"default:NULL"`
}

// Ticket represents the tickets table.
type Ticket struct {
	ID            int        `json:"id" gorm:"primary_key;auto_increment"`
	Subject       string     `json:"subject" gorm:"type:varchar(255);not_null"`
	Description   string     `json:"description" gorm:"type:text;not_null"`
	CategoryID    int        `json:"category_id" gorm:"not_null"`
	SubCategoryID *int       `json:"sub_category_id,omitempty" gorm:"default:NULL"`
	PriorityID    int        `json:"priority_id" gorm:"not_null"`
	SLAID         int        `json:"sla_id" gorm:"not_null"`
	UserID        int        `json:"user_id" gorm:"not_null"`
	AgentID       *int       `json:"agent_id,omitempty" gorm:"default:NULL"`
	AssignedAt    *time.Time `json:"assigned_at,omitempty" gorm:"default:NULL"`
	ClosedAt      *time.Time `json:"closed_at,omitempty" gorm:"default:NULL"`
	DueAt         *time.Time `json:"due_at,omitempty" gorm:"default:NULL"`
	Site          string     `json:"site,omitempty" gorm:"type:varchar(255)"`
	StatusID      int        `json:"status_id" gorm:"not_null"`
	CreatedAt     time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty" gorm:"default:NULL"`
}

// Feedback represents the feedback table.
type Feedback struct {
	FeedbackID int        `json:"feedback_id" gorm:"primary_key;auto_increment"`
	TicketID   int        `json:"ticket_id" gorm:"not_null"`
	UserID     int        `json:"user_id" gorm:"not_null"`
	Rating     int        `json:"rating" gorm:"not_null"`
	Comment    string     `json:"comment,omitempty" gorm:"type:text"`
	CreatedAt  time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty" gorm:"default:NULL"`
}

// UserProfile represents the user_profiles table.
type UserProfile struct {
	ProfileID      int        `json:"profile_id" gorm:"primary_key;auto_increment"`
	UserID         int        `json:"user_id" gorm:"unique;not_null"`
	Address        string     `json:"address,omitempty" gorm:"type:text"`
	ProfilePicture string     `json:"profile_picture,omitempty" gorm:"type:varchar(255)"`
	Bio            string     `json:"bio,omitempty" gorm:"type:text"`
	CreatedAt      time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty" gorm:"default:NULL"`
}

// SLA represents the sla table.
type SLA struct {
	SLAID          int        `json:"sla_id" gorm:"primary_key;auto_increment"`
	SLAName        string     `json:"sla_name" gorm:"type:varchar(255);not_null"`
	PriorityID     int        `json:"priority_id" gorm:"not_null"`
	ResponseTime   int        `json:"response_time" gorm:"not_null"`
	ResolutionTime int        `json:"resolution_time" gorm:"not_null"`
	CreatedAt      time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty" gorm:"default:NULL"`
}

// Incident represents the incidents table.
type Incident struct {
	IncidentID  int        `json:"incident_id" gorm:"primary_key;auto_increment"`
	Title       string     `json:"title" gorm:"type:varchar(255);not_null"`
	Description string     `json:"description" gorm:"type:text"`
	ReportedBy  int        `json:"reported_by" gorm:"not_null"`
	AssignedTo  *int       `json:"assigned_to,omitempty" gorm:"default:NULL"`
	Status      string     `json:"status" gorm:"type:varchar(50);not_null"`
	PriorityID  int        `json:"priority_id" gorm:"not_null"`
	CreatedAt   time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"default:NULL"`
}

// OperationalMetric represents the operational_metrics table.
type OperationalMetric struct {
	MetricID        int        `json:"metric_id" gorm:"primary_key;auto_increment"`
	Name            string     `json:"name" gorm:"type:varchar(255);not_null"`
	Value           float64    `json:"value" gorm:"type:decimal(10,2)"`
	MeasurementDate time.Time  `json:"measurement_date" gorm:"type:date;not_null"`
	CreatedAt       time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty" gorm:"default:NULL"`
}

// AuditLog represents the audit_logs table.
type AuditLog struct {
	AuditLogID    int       `json:"audit_log_id" gorm:"primary_key;auto_increment"`
	UserID        *int      `json:"user_id,omitempty" gorm:"default:NULL"`
	ActionType    string    `json:"action_type" gorm:"type:varchar(255);not_null"`
	Description   string    `json:"description" gorm:"type:text"`
	AffectedTable string    `json:"affected_table,omitempty" gorm:"type:varchar(255)"`
	AffectedRowID *int      `json:"affected_row_id,omitempty" gorm:"default:NULL"`
	ChangeDetails string    `json:"change_details,omitempty" gorm:"type:text"` // Consider JSON or similar structured data
	CreatedAt     time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}
