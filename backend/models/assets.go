// backend/models/assets.go

package models

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// DeviceRegistration tracks devices registered for push notifications or other services.
type DeviceRegistration struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"index;notNull" json:"user_id"`
	DeviceID  string         `gorm:"type:varchar(255);notNull;unique" json:"device_id"`
	Platform  string         `gorm:"type:varchar(100);notNull" json:"platform"`
	Token     string         `gorm:"type:text;notNull" json:"token"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (DeviceRegistration) TableName() string {
	return "device_registrations"
}

// Assets represent details about various assets.
type Assets struct {
	gorm.Model
	AssetTag          string     `gorm:"size:100;uniqueIndex;not null" json:"assetTag"`
	Name              string     `gorm:"size:255;not null" json:"name"`
	Description       string     `gorm:"size:500" json:"description,omitempty"`
	Status            string     `gorm:"size:100;not null" json:"status"`
	CategoryID        uint       `gorm:"not null" json:"categoryId"`
	Location          string     `gorm:"size:255" json:"location,omitempty"`
	PurchaseDate      *time.Time `json:"purchaseDate,omitempty"`
	WarrantyExpiry    *time.Time `json:"warrantyExpiry,omitempty"`
	Manufacturer      string     `gorm:"size:255" json:"manufacturer,omitempty"`
	SerialNumber      string     `gorm:"size:255;uniqueIndex" json:"serialNumber,omitempty"`
	PurchaseCost      float64    `json:"purchaseCost,omitempty"`
	CurrentValue      float64    `json:"currentValue,omitempty"`
	Depreciation      float64    `json:"depreciation,omitempty"`
	AssetTypeID       uint       `gorm:"index;not null" json:"asset_type_id"`
	AssetName         string     `gorm:"size:255;not null" json:"asset_name"`
	AssetModel        string     `gorm:"size:255" json:"model,omitempty"`
	WarrantyExpire    *time.Time `json:"warranty_expire,omitempty"`
	PurchasePrice     float64    `gorm:"type:decimal(10,2)" json:"purchase_price,omitempty"`
	Vendor            Vendor     `gorm:"foreignKey:VendorID" json:"-"`
	User              Users      `gorm:"foreignKey:UserID" json:"-"`
	AssigneeID        *uint      `gorm:"index" json:"user_id,omitempty"`
	SiteID            uint       `gorm:"index" json:"site_id"`
	CreatedByID       uint       `gorm:"index" json:"created_by_id"`
	AssetAssignment   *uint      `json:"assetAssignment" gorm:"foreignKey:AssetAssignmentID"`
	DecommissionDate  *time.Time `json:"decommission_date,omitempty"`
	Tag               string     `gorm:"type:varchar(100);uniqueIndex" json:"tag"`
	DepreciationValue float64    `gorm:"type:decimal(10,2)" json:"depreciationValue"`
	UsefulLife        int        `gorm:"type:int" json:"usefulLife"`
	SalvageValue      float64    `gorm:"type:decimal(10,2)" json:"salvageValue"`
}

func (Assets) TableName() string {
	return "assets"
}

// InventoryItem represents an item in the inventory.
type InventoryItem struct {
	gorm.Model
	ProductID   uint   `gorm:"uniqueIndex;not null" json:"product_id"`
	Name        string `gorm:"type:varchar(255);not null" json:"name"`
	Description string `gorm:"type:text;not null" json:"description"`
	StockLevel  int    `gorm:"not null" json:"stock_level"`
}

func (InventoryItem) TableName() string {
	return "inventory_items"
}

// AssetTag represents tags associated with an asset.
type AssetTag struct {
	gorm.Model
	AssetTag string   `json:"tag"`
	Tags     []string `gorm:"type:text[]" json:"tags"`
	AssetID  uint     `gorm:"index" json:"asset_id"`
}

func (AssetTag) TableName() string {
	return "asset_tag"
}

// AssetType represents the type of an asset.
type AssetType struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	AssetType   string `json:"asset_type"`
	Description string `gorm:"type:text" json:"description"`
}

func (AssetType) TableName() string {
	return "assetType"
}

// AssetCategory represents the category of an asset.
type AssetCategory struct {
	gorm.Model
	Name        string `gorm:"size:255;uniqueIndex;not null" json:"name"`
	Description string `gorm:"size:500" json:"description,omitempty"`
}

func (AssetCategory) TableName() string {
	return "asset_category"
}

// AssetAssignment tracks the assignment of assets to users or locations.
type AssetAssignment struct {
	gorm.Model
	AssetID        uint       `gorm:"not null" json:"assetId"`
	AssignedTo     uint       `gorm:"not null" json:"assignedTo"`
	AssignmentDate time.Time  `json:"assignmentDate"`
	ReturnDate     *time.Time `json:"returnDate,omitempty"`
	AssignmentType string     `gorm:"size:255" json:"assignment_type"`
	Status         string     `gorm:"size:100;not null" json:"status"`
	AssignedBy     uint       `gorm:"index;not null" json:"assigned_by"`
	DueAt          *time.Time `json:"due_at,omitempty"`
}

func (AssetAssignment) TableName() string {
	return "asset_assignments"
}

// AssetLog records actions taken on assets for auditing purposes.
type AssetLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	AssetID   uint      `gorm:"not null" json:"asset_id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	Action    string    `gorm:"type:varchar(255)" json:"action"`
	Timestamp time.Time `json:"timestamp"`
	Details   string    `gorm:"type:text" json:"details"`
}

func (AssetLog) TableName() string {
	return "asset_log"
}

// Resource represents a resource.
type Resource struct {
	gorm.Model
	Name               string    `json:"name" binding:"required"`
	Description        string    `json:"description,omitempty"`
	Type               string    `json:"type"`
	Location           string    `json:"location,omitempty"`
	AvailabilityStatus string    `json:"availability_status"`
	Metadata           string    `gorm:"type:text" json:"metadata,omitempty"`
	Status             string    `gorm:"type:varchar(100);not null" json:"status"`
	Bookings           []Booking `gorm:"foreignKey:ResourceID" json:"bookings,omitempty"`
}

func (Resource) TableName() string {
	return "resources"
}

// ResourceAllocation represents the allocation of a resource.
type ResourceAllocation struct {
	gorm.Model
	ResourceID        uint      `gorm:"index;not null" json:"resource_id"`
	UserID            uint      `gorm:"index" json:"user_id,omitempty"`
	AllocationDetails string    `gorm:"type:text" json:"allocation_details"`
	StartTime         time.Time `json:"start_time"`
	EndTime           time.Time `json:"end_time"`
}

func (ResourceAllocation) TableName() string {
	return "resource_allocations"
}

// AgentResourceAllocation represents the allocation of a resource to an agent.
type AgentResourceAllocation struct {
	ID                uint            `gorm:"primaryKey" json:"id"`
	AgentID           uint            `json:"agent_id" gorm:"index;not null"`
	ResourceID        uint            `json:"resource_id" gorm:"index;not null"`
	AllocationDetails string          `gorm:"type:text" json:"allocation_details,omitempty"`
	StartTime         time.Time       `json:"start_time"`
	EndTime           time.Time       `json:"end_time"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
	DeletedAt         *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (AgentResourceAllocation) TableName() string {
	return "agent_resource_allocations"
}

// Booking represents a booking of a resource.
type Booking struct {
	gorm.Model
	ResourceID uint      `json:"resource_id"`
	UserID     uint      `json:"user_id"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Status     string    `json:"status"`
}

func (Booking) TableName() string {
	return "bookings"
}

// AccessControlList represents an access control list entry.
type AccessControlList struct {
	gorm.Model
	Resource   string `gorm:"type:varchar(255);not null" json:"resource"`
	Action     string `gorm:"type:varchar(100);not null" json:"action"`
	RoleID     uint   `gorm:"index" json:"role_id"`
	Permission bool   `json:"permission"`
}

func (AccessControlList) TableName() string {
	return "access_control_lists"
}

// ComplianceAuditLog represents an entry in the compliance audit log.
type ComplianceAuditLog struct {
	gorm.Model
	Action      string `gorm:"type:varchar(255);not null"`
	UserID      uint   `gorm:"index"`
	Description string `gorm:"type:text;not null"`
	Details     string `gorm:"type:text"`
}

func (ComplianceAuditLog) TableName() string {
	return "compliance_audit_logs"
}

// Vendor represents a vendor.
type Vendor struct {
	gorm.Model
	VendorName    string `gorm:"size:255;not null;unique" json:"vendor_name"`
	Description   string `gorm:"type:text" json:"description,omitempty"`
	ContactInfo   string `gorm:"type:text" json:"contact_info,omitempty"`
	ContactPerson string `gorm:"size:255" json:"contact_person,omitempty"`
	Address       string `gorm:"type:text" json:"address,omitempty"`
}

func (Vendor) TableName() string {
	return "vendors"
}

// LifecycleEvent documents significant events in an asset's lifecycle.
type LifecycleEvent struct {
	gorm.Model
	AssetID uint      `gorm:"not null" json:"assetId"`
	Date    time.Time `json:"date"`
	Type    string    `gorm:"size:255;not null" json:"type"`
	Notes   string    `gorm:"size:500" json:"notes,omitempty"`
}

// AssetMaintenance represents maintenance activities for an asset.
type AssetMaintenance struct {
	gorm.Model
	AssetID       uint       `gorm:"index;not null" json:"assetId"`
	ScheduledDate *time.Time `json:"scheduledDate"`
	CompletedDate *time.Time `json:"completedDate,omitempty"`
	Description   string     `gorm:"type:text" json:"description,omitempty"`
	Status        string     `gorm:"size:100;not null" json:"status"`
}

// AssetPerformanceAnalysis documents the performance analysis of an asset.
type AssetPerformanceAnalysis struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	AssetID           uint      `gorm:"index;not null" json:"asset_id"`
	Summary           string    `gorm:"type:text" json:"summary"`
	DowntimeEvents    int       `gorm:"default:0" json:"downtime_events"`
	AverageDowntime   float64   `gorm:"type:decimal(10,2);default:0.0" json:"average_downtime"`
	AnalysisStartDate time.Time `gorm:"type:datetime" json:"analysis_start_date"`
	AnalysisEndDate   time.Time `gorm:"type:datetime" json:"analysis_end_date"`
	CreatedAt         time.Time `gorm:"<-:create" json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// AssetPerformanceLog records performance metrics for an asset over time.
type AssetPerformanceLog struct {
	gorm.Model
	AssetID     uint      `gorm:"index;not null" json:"assetId"`
	LogDate     time.Time `json:"logDate"`
	Performance string    `gorm:"type:text" json:"performance"`
	Notes       string    `gorm:"type:text" json:"notes"`
}

// AssetPerformanceMetrics represents metrics for an asset's performance.
type AssetPerformanceMetrics struct {
	gorm.Model
	AssetID    uint      `gorm:"index;not null" json:"assetId"`
	MetricDate time.Time `json:"metricDate"`
	MetricType string    `gorm:"type:varchar(100);not null" json:"metricType"`
	Value      float64   `json:"value"`
}

// AssetUtilizationReport provides a summary of how an asset is being utilized.
type AssetUtilizationReport struct {
	AssetID         uint    `json:"assetId"`
	Utilization     float64 `json:"utilization"`
	ReportingPeriod string  `json:"reportingPeriod"`
}

// AssetAuditLog records changes and access to an asset for auditing purposes.
type AssetAuditLog struct {
	gorm.Model
	AssetID     uint      `gorm:"index;not null" json:"assetId"`
	Action      string    `gorm:"type:varchar(255);not null" json:"action"`
	ActionDate  time.Time `json:"actionDate"`
	PerformedBy uint      `gorm:"index" json:"performedBy"`
	Notes       string    `gorm:"type:text" json:"notes"`
}

// AssetIssue tracks issues reported for specific assets.
type AssetIssue struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	AssetID     uint           `gorm:"index;notNull" json:"asset_id"`
	ReportedBy  uint           `gorm:"index" json:"reported_by"`
	ReportedAt  time.Time      `json:"reported_at"`
	IssueType   string         `gorm:"size:255" json:"issue_type"`
	Description string         `gorm:"type:text" json:"description"`
	Status      string         `gorm:"size:100;notNull" json:"status"`
	ResolvedBy  *uint          `gorm:"index" json:"resolved_by,omitempty"`
	ResolvedAt  *time.Time     `json:"resolved_at,omitempty"`
	Resolution  string         `gorm:"type:text" json:"resolution,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// AssetCostAnalysis details the cost analysis of an asset.
type AssetCostAnalysis struct {
	AssetID        uint               `json:"assetId"`
	TotalCost      float64            `json:"totalCost"`
	CostComponents map[string]float64 `json:"costComponents"`
	AnalysisDate   time.Time          `json:"analysis_date"`
}

// AssetAccessLog records access or modifications to asset data.
type AssetAccessLog struct {
	gorm.Model
	AssetID     uint      `gorm:"index;not null" json:"assetId"`
	AccessedBy  uint      `gorm:"index" json:"accessedBy"`
	AccessDate  time.Time `json:"accessDate"`
	ActionType  string    `gorm:"type:varchar(100);not null" json:"actionType"`
	Description string    `gorm:"type:text" json:"description"`
}

// AssetPermission defines permissions for accessing or modifying asset data.
type AssetPermission struct {
	gorm.Model
	AssetID    uint       `gorm:"index;not null" json:"assetId"`
	GrantedTo  uint       `gorm:"index" json:"grantedTo"`
	Permission string     `gorm:"type:varchar(100);not null" json:"permission"`
	ValidUntil *time.Time `json:"validUntil"`
}

// AssetReport compiles comprehensive details about an asset into a report.
type AssetReport struct {
	AssetID    uint      `json:"assetId"`
	ReportDate time.Time `json:"reportDate"`
	Content    string    `gorm:"type:text" json:"content"`
}

// AssetLifecycleEvent records events in the lifecycle of an asset.
type AssetLifecycleEvent struct {
	gorm.Model
	AssetID   uint      `gorm:"index;not null" json:"assetId"`
	EventType string    `gorm:"type:varchar(100);not null" json:"eventType"`
	EventDate time.Time `json:"eventDate"`
	Details   string    `gorm:"type:text" json:"details"`
}

// AssetRepairLog records repairs performed on assets.
type AssetRepairLog struct {
	gorm.Model
	AssetID           uint              `gorm:"index;not null" json:"assetId"`
	Description       string            `gorm:"type:text" json:"description"`
	Status            string            `gorm:"type:varchar(100)" json:"status"`
	StartDate         time.Time         `json:"startDate"`
	CompletionDetails CompletionDetails `gorm:"foreignKey:UserID" json:"-"`
}

// CompletionDetails contains details of the completion of an asset repair.
type CompletionDetails struct {
	gorm.Model
	RepairLogID    uint      `gorm:"index;not null" json:"repairLogId"`
	Cost           float64   `gorm:"type:decimal(10,2)" json:"cost"`
	RepairProvider string    `gorm:"type:varchar(255)" json:"repairProvider"`
	RepairDate     time.Time `json:"repairDate"`
	StartDate      time.Time `json:"startDate"`
	ReturnDate     time.Time `json:"returnDate"`
}

// AssetDepreciationRecord records depreciation details for an asset.
type AssetDepreciationRecord struct {
	gorm.Model
	AssetID                 uint      `gorm:"index;not null" json:"assetId"`
	DepreciationDate        time.Time `json:"depreciationDate"`
	DepreciatedValue        float64   `gorm:"type:decimal(10,2)" json:"depreciatedValue"`
	AccumulatedDepreciation float64   `gorm:"type:decimal(10,2)" json:"accumulatedDepreciation"`
}

// ExternalServiceData stores data from external services related to assets.
type ExternalServiceData struct {
	ID                    uint      `gorm:"primaryKey" json:"id"`
	AssetID               uint      `gorm:"index;not null" json:"asset_id"`
	ExternalServiceID     string    `gorm:"type:varchar(255);not null" json:"external_service_id"`
	LastSynced            time.Time `gorm:"type:datetime" json:"last_synced"`
	ExternalServiceStatus string    `gorm:"type:varchar(255)" json:"external_service_status"`
	AdditionalInformation string    `gorm:"type:text" json:"additional_information"`
	CreatedAt             time.Time `gorm:"<-:create" json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// AssetDecommission records the decommissioning process of assets.
type AssetDecommission struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	AssetID          uint      `gorm:"not null" json:"asset_id"`
	DecommissionDate time.Time `json:"decommission_date"`
	Reason           string    `gorm:"type:text" json:"reason"`
	Status           string    `gorm:"type:varchar(100)" json:"status"`
}

// AssetUtilization captures utilization data for assets.
type AssetUtilization struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	AssetID      uint      `gorm:"not null" json:"asset_id"`
	Utilization  float64   `json:"utilization"`
	ReportedDate time.Time `json:"reported_date"`
}

// AssetHealthReport represents the health status of an asset.
type AssetHealthReport struct {
	gorm.Model
	AssetID      uint      `gorm:"index;notNull" json:"asset_id"`
	Report       string    `json:"report"`
	ReportDate   time.Time `json:"report_date"`
	HealthStatus string    `gorm:"type:varchar(100);notNull" json:"health_status"`
	Notes        string    `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
}

func (AssetHealthReport) TableName() string {
	return "asset_health_reports"
}

// AssetMaintenanceSchedule represents the scheduled maintenance for an asset.
type AssetMaintenanceSchedule struct {
	gorm.Model
	AssetID       uint       `gorm:"index;notNull" json:"asset_id"`
	ScheduledDate time.Time  `json:"scheduled_date"`
	CompletedDate *time.Time `json:"completedDate,omitempty"`
	Description   string     `gorm:"type:text" json:"description"`
	Status        string     `gorm:"type:varchar(100);notNull" json:"status"`
}

func (AssetMaintenanceSchedule) TableName() string {
	return "asset_maintenance_schedules"
}

// AssetInspectionRecord documents the inspections performed on assets.
type AssetInspectionRecord struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	AssetID        uint           `gorm:"index;notNull" json:"asset_id"`
	InspectionDate time.Time      `json:"inspection_date"`
	InspectedBy    uint           `gorm:"index;notNull" json:"inspected_by"`
	Outcome        string         `gorm:"size:255" json:"outcome"`
	Notes          string         `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (AssetInspectionRecord) TableName() string {
	return "asset_inspection_record"
}

type AssetStorage interface {
	CreateAsset(asset *Assets) error
	DeleteAsset(assetID uint) error
	UpdateAsset(assetID uint, updates map[string]interface{}) (*Assets, error)
	GetAssetByID(assetID uint) (*Assets, error)
	ListAllAssets() ([]*Assets, error)
	AssignAssetToUser(assetID, userID uint) error
	ReturnAssetFromUser(assignmentID uint) error
	ListAssetsByStatus(status string) ([]*Assets, error)
	CalculateAssetDepreciation(assetID uint) error
	RegisterAsset(asset *Assets) error
	GetAllAssets() ([]*Assets, error)
	GetAssetByNumber(assetNumber int) (*Assets, error)
	CreateAssetType(assetType *AssetType) error
	UpdateAssetType(typeID uint, updates map[string]interface{}) error
	DeleteAssetType(typeID uint) error
	ListAssetTypes() ([]*AssetType, error)
	ScheduleAssetMaintenance(maintenance *AssetMaintenance) error
	CompleteAssetMaintenance(maintenanceID uint) error
	ListAssetMaintenanceSchedules(assetID uint) ([]*AssetMaintenance, error)
	PerformCostAnalysis(assetID uint) (*AssetCostAnalysis, error)
	GenerateAssetReport(assetID uint) (*AssetReport, error)
	AuditAssetAccess(accessLog *AssetAccessLog) error
	ListAssetAudits(assetID uint) ([]*AssetAuditLog, error)
	SyncAssetWithExternalService(assetID uint, data *ExternalServiceData) error
	InitiateDecommissioning(decommission *AssetDecommission) error
	CompleteDecommissioning(assetID uint) error
	GrantAssetPermission(permission *AssetPermission) error
	RevokeAssetPermission(permissionID uint) error
	ListAssetPermissions(assetID uint) ([]*AssetPermission, error)
	CreateAssetAssignment(assetAssignment *AssetAssignment) error
	UpdateAssetAssignment(*AssetAssignment) error
	DeleteAssetAssignment(int) error
	GetAssetAssignmentsByUser(userID uint) ([]*AssetAssignment, error)
	GetAssetAssignmentsByAsset(assetID uint) ([]*AssetAssignment, error)
	UnassignAsset(assetID uint) error
	ReassignAsset(assetID, fromUserID, toUserID uint) error
	LogIncident(incident *Incident) error
	ResolveIncident(incidentID uint, resolution string) error
	ScheduleBulkAssetMaintenance(assetIDs []uint, maintenanceDetails *AssetMaintenanceSchedule) error
	DecommissionAssets(assetIDs []uint) error
	ReactivateAssets(assetIDs []uint) error
	AddLifecycleEvent(assetID uint, event *LifecycleEvent) error
	ListLifecycleEvents(assetID uint) ([]*LifecycleEvent, error)
	ListOverdueMaintenances() ([]*AssetMaintenanceSchedule, error)
	RecordAssetPerformance(assetID uint, performance *AssetPerformanceLog) error
	RetrievePerformanceHistory(assetID uint) ([]*AssetPerformanceLog, error)
	AnalyzeAssetUtilization(assetID uint) (*AssetUtilizationReport, error)
	ListDepreciatedAssets(valueThreshold float64) ([]*Assets, error)
	LogAssetRepair(repair *AssetRepairLog) error
	ListAssetRepairLogs(assetID uint) ([]*AssetRepairLog, error)
	RecordAssetIssue(assetID uint, issue *AssetIssue) error
	ResolveAssetIssue(issueID uint, resolutionDetails string) error
	ListAssetIssues(assetID uint) ([]*AssetIssue, error)
}

type AssetAssignmentStorage interface {
	CreateAssetAssignment(assetAssignment *AssetAssignment) error
	DeleteAssetAssignment(int) error
	UpdateAssetAssignment(*AssetAssignment) error
	GetAssetAssignment() ([]*AssetAssignment, error)
	GetAssetAssignmentByID(int) (*AssetAssignment, error)
	GetAssetAssignmentByNumber(int) (*AssetAssignment, error)
	AssignAsset(asset *Assets, userID uint) (*AssetAssignment, error)
	UnassignAsset(assetAssignment *AssetAssignment, agentID uint) (*AssetAssignment, error)
}

// AssetAssignmentModel handles database operations for Asset
type AssetAssignmentDBModel struct {
	DB           *gorm.DB
	AssetDBModel *AssetDBModel
	log          Logger
}

// NewAssetAssignmentModel creates a new instance of TicketModel
func NewAssetAssignmentDBModel(db *gorm.DB, logger Logger) *AssetAssignmentDBModel {
	return &AssetAssignmentDBModel{
		DB: db,
		//AssetDBModel: assetDBModel,
		log: logger,
	}
}

// AssetModel handles database operations for Asset
type AssetDBModel struct {
	DB              *gorm.DB
	AssetAssignment *AssetAssignmentDBModel
	log             Logger
	EventPublisher  EventPublisherImpl
}

// NewAssetModel creates a new instance of TicketModel
func NewAssetDBModel(db *gorm.DB, assetAssignment *AssetAssignmentDBModel, log Logger, eventPublisher EventPublisherImpl) *AssetDBModel {
	return &AssetDBModel{
		DB:              db,
		AssetAssignment: assetAssignment,
		log:             log,
		EventPublisher:  eventPublisher,
	}
}

// CreateAsset adds a new asset to the database.
func (db *AssetDBModel) CreateAsset(asset *Assets) error {
	tx := db.DB.Begin()
	if tx.Error != nil {
		db.log.Error("Failed to begin transaction:", tx.Error)
		return tx.Error
	}

	if err := tx.Create(asset).Error; err != nil {
		tx.Rollback()
		db.log.Error("Failed to create asset:", err)
		return fmt.Errorf("creating asset: %w", err)
	}

	return tx.Commit().Error
}

// GetAssetByID fetches the details of a specific asset by its ID.
func (model *AssetDBModel) GetAssetByID(id uint) (*Assets, error) {
	var asset Assets
	if err := model.DB.Where("id = ?", id).First(&asset).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("asset with ID %d not found: %w", id, err)
		}
		return nil, fmt.Errorf("failed to retrieve asset with ID %d: %w", id, err)
	}
	return &asset, nil
}

// UpdateAsset modifies an existing asset identified by its ID.
func (model *AssetDBModel) UpdateAsset(assetID uint, updates *Assets) error {
	tx := model.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Model(&Assets{}).Where("id = ?", assetID).Updates(updates).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update asset with ID %d: %w", assetID, err)
	}

	return tx.Commit().Error
}

// DeleteAsset removes an asset from the database by its ID.
func (model *AssetDBModel) DeleteAsset(id uint) error {
	tx := model.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Delete(&Assets{}, id).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete asset with ID %d: %w", id, err)
	}

	return tx.Commit().Error
}

// RegisterAsset adds a new asset to the database.
func (db *AssetDBModel) RegisterAsset(asset *Assets) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&asset).Error; err != nil {
			return err
		}
		return nil
	})
}

// GetAllAssets retrieves a list of all assets in the database.
func (model *AssetDBModel) GetAllAssets() ([]*Assets, error) {
	var assets []*Assets
	if err := model.DB.Find(&assets).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve assets: %w", err)
	}
	return assets, nil
}

// ListAvailableAssets lists all assets that are not currently assigned.
func (model *AssetDBModel) ListAvailableAssetsA() ([]*Assets, error) {
	var assets []*Assets
	err := model.DB.Where("assigned_to_user_id IS NULL").Find(&assets).Error
	if err != nil {
		return nil, fmt.Errorf("error listing available assets: %w", err)
	}
	return assets, nil
}

// ListAvailableAssets lists all assets that are not currently assigned.
func (model *AssetDBModel) ListAvailableAssetsB() ([]*Assets, error) {
	var assets []*Assets
	err := model.DB.Where("status = ?", "available").Find(&assets).Error
	if err != nil {
		return nil, fmt.Errorf("error listing available assets: %w", err)
	}
	return assets, nil
}

// ReassignAsset transfers an asset from one user to another.
func (model *AssetDBModel) ReassignAsset(assetID, fromUserID, toUserID uint) error {
	return model.DB.Transaction(func(tx *gorm.DB) error {
		// Update asset assignment to the new user.
		err := tx.Model(&AssetAssignment{}).Where("asset_id = ? AND assigned_to = ?", assetID, fromUserID).
			Update("assigned_to", toUserID).Error
		if err != nil {
			return fmt.Errorf("error reassigning asset %d: %w", assetID, err)
		}
		return nil
	})
}

// GetAssetAssignmentHistory retrieves the history of assignments for a given asset.
func (model *AssetDBModel) GetAssetAssignmentHistory(assetID uint) ([]*AssetAssignment, error) {
	var history []*AssetAssignment
	err := model.DB.Where("asset_id = ?", assetID).Find(&history).Error
	if err != nil {
		return nil, fmt.Errorf("error getting asset assignment history for asset %d: %w", assetID, err)
	}
	return history, nil
}

// ValidateAssetOwnership checks if a given asset is assigned to a specified user.
func (model *AssetDBModel) ValidateAssetOwnership(assetID, userID uint) (bool, error) {
	var count int64
	err := model.DB.Model(&AssetAssignment{}).Where("asset_id = ? AND assigned_to = ?", assetID, userID).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("error validating ownership for asset %d: %w", assetID, err)
	}
	return count > 0, nil
}

// ListAssetsByType lists all assets of a specified type.
func (model *AssetDBModel) ListAssetsByType(assetType string) ([]*Assets, error) {
	var assets []*Assets
	err := model.DB.Joins("JOIN asset_types on assets.asset_type_id = asset_types.id").
		Where("asset_types.type = ?", assetType).Find(&assets).Error
	if err != nil {
		return nil, fmt.Errorf("error listing assets by type '%s': %w", assetType, err)
	}
	return assets, nil
}

// AssetDecommission marks an asset as decommissioned and out of service.
func (s *AssetDBModel) AssetDecommission(assetID uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Assets{}).Where("id = ?", assetID).Update("status", "Decommissioned").Error; err != nil {
			s.log.Error(fmt.Sprintf("Failed to decommission asset %d: %v", assetID, err))
			return err
		}
		return nil
	})
}

// RetrieveAssetMaintenanceRecords fetches all maintenance records associated with an asset.
func (s *AssetDBModel) RetrieveAssetMaintenanceRecords(assetID uint) ([]*AssetMaintenance, error) {
	var records []*AssetMaintenance
	if err := s.DB.Where("asset_id = ?", assetID).Find(&records).Error; err != nil {
		s.log.Error("Failed to retrieve maintenance records: ", err)
		return nil, err
	}
	return records, nil
}

// GetAssetByNumber retrieves an asset by its asset number.
func (as *AssetDBModel) GetAssetByNumber(assetNumber int) (*Assets, error) {
	var asset Assets
	err := as.DB.Where("asset_id = ?", assetNumber).First(&asset).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

// CreateAssetType adds a new type of asset to the system.
func (model *AssetDBModel) CreateAssetType(assetType *AssetType) error {
	if err := model.DB.Create(assetType).Error; err != nil {
		return fmt.Errorf("failed to create asset type: %w", err)
	}
	return nil
}

// UpdateAssetType modifies the details of an existing asset type.
func (model *AssetDBModel) UpdateAssetType(id uint, updates map[string]interface{}) error {
	if err := model.DB.Model(&AssetType{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update asset type with ID %d: %w", id, err)
	}
	return nil
}

// CompleteAssetRepair marks a scheduled repair as completed.
func (s *AssetDBModel) CompleteAssetRepair(repairID uint, completionDetails CompletionDetails) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&AssetRepairLog{}).Where("id = ?", repairID).Updates(AssetRepairLog{Status: "Completed", CompletionDetails: completionDetails}).Error; err != nil {
			s.log.Error(fmt.Sprintf("Failed to complete repair %d: %v", repairID, err))
			return err
		}
		s.log.Info(fmt.Sprintf("Repair %d completed", repairID))
		return nil
	})
}

///////////////////////////////////////////////////////////////////////////////////////////

// GetAssetAssignmentByID retrieves a AssetAssignment by its ID.
func (as *AssetAssignmentDBModel) GetAssetAssignmentByID(id uint) (*AssetAssignment, error) {
	var asset AssetAssignment
	err := as.DB.Where("id = ?", id).First(&asset).Error
	return &asset, err
}

// UpdateAssetAssignment updates the assignment of an existing asset.
func (as *AssetAssignmentDBModel) AssignAsset(asset *Assets, userID uint) (*AssetAssignment, error) {
	var assetAssignment AssetAssignment
	assetAssignment.AssetID = asset.ID
	assetAssignment.AssignmentType = "permanent"
	assetAssignment.CreatedAt = time.Now()

	id := as.DB.Save(&assetAssignment).RowsAffected
	assignment, err := as.GetAssetAssignmentByID(uint(id))
	if err != nil {
		return nil, err
	}

	return assignment, nil
}

// InitiateDecommissioning marks an asset for decommissioning.
func (db *AssetDBModel) InitiateDecommissioning(decommission *AssetDecommission) error {
	tx := db.DB.Begin()
	if tx.Error != nil {
		db.log.Error("Failed to initiate transaction for asset decommissioning", tx.Error)
		return tx.Error
	}

	if err := tx.Create(decommission).Error; err != nil {
		tx.Rollback()
		db.log.Error("Failed to initiate decommissioning for asset", err)
		return fmt.Errorf("initiating decommissioning: %w", err)
	}

	return tx.Commit().Error
}

// CompleteDecommissioning finalizes the decommissioning process for an asset.
func (db *AssetDBModel) CompleteDecommissioning(assetID uint) error {
	tx := db.DB.Begin()
	if tx.Error != nil {
		db.log.Error("Failed to start transaction for completing decommissioning", tx.Error)
		return tx.Error
	}

	if err := tx.Model(&AssetDecommission{}).Where("asset_id = ?", assetID).Update("status", "Completed").Error; err != nil {
		tx.Rollback()
		db.log.Error("Failed to complete decommissioning for asset", err)
		return fmt.Errorf("completing decommissioning: %w", err)
	}

	return tx.Commit().Error
}

// RecordAssetUsage logs the usage details of an asset.
func (s *AssetDBModel) RecordAssetUsage(assetID uint, usageDetails *AssetUtilization) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(usageDetails).Error; err != nil {
			s.log.Error(fmt.Sprintf("Failed to record usage for asset %d: %v", assetID, err))
			return err
		}
		return nil
	})
}

// RetrieveAssetUsageHistory fetches the usage history of an asset.
func (s *AssetDBModel) RetrieveAssetUsageHistory(assetID uint) ([]*AssetUtilization, error) {
	var usages []*AssetUtilization
	if err := s.DB.Where("asset_id = ?", assetID).Find(&usages).Error; err != nil {
		s.log.Error("Failed to retrieve asset usage history: ", err)
		return nil, err
	}
	return usages, nil
}

// GrantAssetPermission grants a user permission to access or modify an asset.
func (db *AssetDBModel) GrantAssetPermission(permission *AssetPermission) error {
	tx := db.DB.Begin()
	if tx.Error != nil {
		db.log.Error("Failed to start transaction for granting asset permission", tx.Error)
		return tx.Error
	}

	if err := tx.Create(permission).Error; err != nil {
		tx.Rollback()
		db.log.Error("Failed to grant asset permission", err)
		return fmt.Errorf("granting asset permission: %w", err)
	}

	return tx.Commit().Error
}

// RevokeAssetPermission revokes a previously granted permission from a user.
func (db *AssetDBModel) RevokeAssetPermission(permissionID uint) error {
	tx := db.DB.Begin()
	if tx.Error != nil {
		db.log.Error("Failed to start transaction for revoking asset permission", tx.Error)
		return tx.Error
	}

	if err := tx.Delete(&AssetPermission{}, permissionID).Error; err != nil {
		tx.Rollback()
		db.log.Error("Failed to revoke asset permission", err)
		return fmt.Errorf("revoking asset permission: %w", err)
	}

	return tx.Commit().Error
}

// ListAssetPermissions lists all permissions granted for an asset.
func (db *AssetDBModel) ListAssetPermissions(assetID uint) ([]*AssetPermission, error) {
	var permissions []*AssetPermission
	if err := db.DB.Where("asset_id = ?", assetID).Find(&permissions).Error; err != nil {
		db.log.Error("Failed to list asset permissions", err)
		return nil, fmt.Errorf("listing asset permissions: %w", err)
	}
	return permissions, nil
}

// PerformCostAnalysis conducts a detailed analysis of the costs associated with maintaining an asset.
// PerformCostAnalysis conducts a detailed analysis of the costs associated with maintaining an asset.
func (db *AssetDBModel) PerformCostAnalysis(assetID uint) (*AssetCostAnalysis, error) {
	var analysis AssetCostAnalysis
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		db.log.Info("Performing cost analysis for asset %d", assetID)
		// Example logic to aggregate costs from maintenance and repair records
		// Replace with actual queries to fetch and sum costs
		maintenanceCosts := 0.0
		repairCosts := 0.0
		tx.Table("asset_maintenance_logs").Select("SUM(cost)").Where("asset_id = ?", assetID).Row().Scan(&maintenanceCosts)
		tx.Table("asset_repair_logs").Select("SUM(cost)").Where("asset_id = ?", assetID).Row().Scan(&repairCosts)
		analysis.AssetID = assetID
		analysis.TotalCost = maintenanceCosts + repairCosts
		analysis.AnalysisDate = time.Now()
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error performing cost analysis for asset %d: %w", assetID, err)
	}

	return &analysis, nil
}

// GenerateAssetReport compiles a comprehensive report detailing all aspects of an asset.
func (db *AssetDBModel) GenerateAssetReport(assetID uint) (*AssetReport, error) {
	// Placeholder for report generation logic
	db.log.Info("Generating report for asset", assetID)
	return &AssetReport{
		AssetID:    assetID,
		ReportDate: time.Now(),
		Content:    "Comprehensive asset report content",
	}, nil
}

// AuditAssetAccess records every access or interaction with an asset for security auditing.
func (db *AssetDBModel) AuditAssetAccess(accessLog *AssetAccessLog) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(accessLog).Error; err != nil {
			db.log.Error("Failed to audit asset access", err)
			return fmt.Errorf("auditing asset access: %w", err)
		}
		return nil
	})
}

// ListAssetAudits retrieves all audit logs for a given asset.
func (db *AssetDBModel) ListAssetAudits(assetID uint) ([]*AssetAuditLog, error) {
	var audits []*AssetAuditLog
	if err := db.DB.Where("asset_id = ?", assetID).Find(&audits).Error; err != nil {
		db.log.Error("Failed to list asset audits", err)
		return nil, fmt.Errorf("listing asset audits: %w", err)
	}
	return audits, nil
}

// SyncAssetWithExternalService updates asset details based on external service data.
func (db *AssetDBModel) SyncAssetWithExternalService(assetID uint, data *ExternalServiceData) error {
	// Placeholder for synchronization logic with external services
	db.log.Info("Synchronizing asset with external service data", assetID)
	return nil // Implement actual synchronization logic
}

// UpdateAssetAssignment updates the assignment of an existing asset.
func (as *AssetAssignmentDBModel) UnassignAsset2(assetAssignment *AssetAssignment, agentID uint) (*AssetAssignment, error) {
	due_at := time.Now().AddDate(1, 0, 0)
	assetAssignment, err := as.GetAssetAssignmentByID(assetAssignment.ID)
	if err != nil {
		return nil, err
	}
	newAssetAssignment := &AssetAssignment{
		AssetID:        assetAssignment.AssetID,
		AssignedTo:     0,
		AssignedBy:     agentID,      // Assuming the same user assigns the asset
		AssignmentType: "unassigned", // Update as needed
		Status:         "unassigned",
		DueAt:          &due_at, // Due date example
	}

	erro := as.AssetDBModel.CreateAssetAssignment(assetAssignment)
	if erro != nil {
		return nil, fmt.Errorf("unable to assign asset")
	}

	return newAssetAssignment, nil
}

// AssignAssetToUser associates an asset with a user, effectively assigning it.
func (model *AssetDBModel) AssignAssetToUser(assetID, userID uint) error {
	tx := model.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Create or update asset assignment
	var assignment AssetAssignment
	result := tx.Where("asset_id = ?", assetID).First(&assignment)
	if result.RowsAffected == 0 {
		// No existing assignment, create a new one
		assignment = AssetAssignment{
			AssetID:        assetID,
			AssignedTo:     userID,
			AssignmentType: "assigned",
			Status:         "active",
		}
	} else {
		// Update existing assignment
		assignment.AssignedTo = userID
		assignment.Status = "active"
	}

	if err := tx.Save(&assignment).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to assign asset %d to user %d: %w", assetID, userID, err)
	}

	return tx.Commit().Error
}

// UnassignAssetFromUser removes the association between an asset and a user.
func (model *AssetDBModel) UnassignAssetFromUser(assetID uint) error {
	tx := model.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Update the asset assignment to mark it as unassigned
	if err := tx.Model(&AssetAssignment{}).Where("asset_id = ?", assetID).Updates(map[string]interface{}{"user_id": gorm.Expr("NULL"), "assignment_status": "inactive"}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to unassign asset %d: %w", assetID, err)
	}

	return tx.Commit().Error
}

func (as *AssetAssignmentDBModel) UnassignAsset(assignmentID, agentID uint) error {
	// Fetch the existing asset assignment
	var assetAssignment AssetAssignment
	if err := as.DB.First(&assetAssignment, assignmentID).Error; err != nil {
		return err // Handle error if assignment is not found
	}

	// Update fields to reflect unassignment
	assetAssignment.AssignedTo = 0 // Assuming this indicates unassignment
	assetAssignment.AssignedBy = agentID
	assetAssignment.AssignmentType = "unassigned"
	assetAssignment.Status = "unassigned"
	assetAssignment.DueAt = nil // Clear due date if appropriate

	// Save changes
	return as.DB.Save(&assetAssignment).Error
}

// backend/models/asset_db_model.go

func (adb *AssetDBModel) AssignAssetToUser2(assetID, userID uint) error {
	// Assign an asset to a user by updating the user_id field in the asset record
	if err := adb.DB.Model(&Assets{}).Where("id = ?", assetID).Update("user_id", userID).Error; err != nil {
		return err
	}

	return nil
}

func (adb *AssetDBModel) UnassignAssetUser2(assetID uint) error {
	// Unassign an asset from a user by setting the user_id field to null
	if err := adb.DB.Model(&Assets{}).Where("id = ?", assetID).Update("user_id", nil).Error; err != nil {
		return err
	}

	return nil
}

// AddAsset records a new asset in the database.
func (db *AssetDBModel) AddAsset(asset *Assets) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(asset).Error; err != nil {
			return err
		}
		return nil
	})
}

// UpdateAssetDetails updates specific details of an asset.
func (model *AssetDBModel) UpdateAssetDetails(assetID uint, newStatus string) error {
	return model.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&Assets{}).Where("id = ?", assetID).Update("status", newStatus).Error
		if err != nil {
			return fmt.Errorf("error updating asset details for asset %d: %w", assetID, err)
		}
		return nil
	})
}

// BulkUpdateAssetsConditionally performs bulk updates on assets based on a condition.
func (db *AssetDBModel) BulkUpdateAssetsConditionally(condition map[string]interface{}, updates map[string]interface{}) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Assets{}).Where(condition).Updates(updates).Error; err != nil {
			return err
		}
		return nil
	})
}

// LogIncident records a new incident, capturing all relevant details.
func (db *IncidentDBModel) LogIncident(incident *Incident) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(incident).Error; err != nil {
			return err
		}
		return nil
	})
}

// ResolveIncident marks an incident as resolved with a resolution summary.
func (db *IncidentDBModel) ResolveIncident(incidentID uint, resolution string) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Incident{}).Where("id = ?", incidentID).Update("resolution", resolution).Update("status", "Resolved").Error; err != nil {
			return err
		}
		return nil
	})
}

// CreateAssetAssignment assigns an asset to a user or updates an existing assignment.
func (model *AssetDBModel) CreateAssetAssignment(assetAssignment *AssetAssignment) error {
	if err := model.DB.Create(assetAssignment).Error; err != nil {
		return fmt.Errorf("failed to create asset assignment: %w", err)
	}
	return nil
}

// UpdateAssetAssignment updates the details of an existing asset assignment with transactional integrity.
func (db *AssetDBModel) UpdateAssetAssignment(assetAssignment *AssetAssignment) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(assetAssignment).Error; err != nil {
			return err
		}
		return nil
	})
}

// DeleteAssetAssignment removes an existing assignment of an asset to a user.
func (model *AssetDBModel) DeleteAssetAssignment(assignmentID uint) error {
	if err := model.DB.Delete(&AssetAssignment{}, assignmentID).Error; err != nil {
		return fmt.Errorf("failed to delete asset assignment with ID %d: %w", assignmentID, err)
	}
	return nil
}

// GetAssetAssignmentsByUser retrieves asset assignments for a user by their user ID.
func (as *AssetDBModel) GetAssetAssignmentsByUser(userID uint) ([]*AssetAssignment, error) {
	var assetAssignments []*AssetAssignment
	err := as.DB.Where("user_id = ?", userID).Find(&assetAssignments).Error
	if err != nil {
		return nil, err
	}
	return assetAssignments, nil
}

// GetAssetAssignmentsByAsset retrieves asset assignments for an asset by its asset ID.
func (as *AssetDBModel) GetAssetAssignmentsByAsset(assetID uint) ([]*AssetAssignment, error) {
	var assetAssignments []*AssetAssignment
	err := as.DB.Where("asset_id = ?", assetID).Find(&assetAssignments).Error
	if err != nil {
		return nil, err
	}
	return assetAssignments, nil
}

// GetAssetAssignmentByID retrieves an asset assignment by its ID.
func (as *AssetDBModel) GetAssetAssignmentByID(id uint) (*AssetAssignment, error) {
	var assetAssignment AssetAssignment
	err := as.DB.Where("id = ?", id).First(&assetAssignment).Error
	if err != nil {
		return nil, err
	}
	return &assetAssignment, nil
}

// GetAssetAssignmentByNumber retrieves an asset assignment by its asset assignment number.
func (as *AssetDBModel) GetAssetAssignmentByNumber(assetAssignmentNumber int) (*AssetAssignment, error) {
	var assetAssignment AssetAssignment
	err := as.DB.Where("asset_assignment_id = ?", assetAssignmentNumber).First(&assetAssignment).Error
	if err != nil {
		return nil, err
	}
	return &assetAssignment, nil
}

// GetAssetByNumber retrieves an asset by its asset number.
func (as *UserDBModel) GetAssetByNumber(assetNumber int) (*Assets, error) {
	var asset Assets
	err := as.DB.Where("asset_id = ?", assetNumber).First(&asset).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

// UpdateAssetStatusAndLocation updates both the status and location of an asset.
func (s *AssetDBModel) UpdateAssetStatusAndLocation(assetID uint, newStatus, newLocation string) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Assets{}).Where("id = ?", assetID).Updates(Assets{Status: newStatus, Location: newLocation}).Error; err != nil {
			s.log.Error(fmt.Sprintf("Failed to update status and location for asset %d: %v", assetID, err))
			return err
		}
		return nil
	})
}

// TransferAsset between users with transactional integrity.
func (db *AssetDBModel) TransferAsset(assetID, fromUserID, toUserID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		// Step 1: Unassign the asset from the current user.
		if err := tx.Model(&AssetAssignment{}).Where("asset_id = ? AND user_id = ?", assetID, fromUserID).Delete(&AssetAssignment{}).Error; err != nil {
			return err
		}

		// Step 2: Assign the asset to the new user.
		newAssignment := AssetAssignment{AssetID: assetID, AssignedTo: toUserID}
		if err := tx.Create(&newAssignment).Error; err != nil {
			return err
		}

		return nil
	})
}

// UpdateAssetStatus updates an asset's status with transactional integrity.
func (db *AssetDBModel) UpdateAssetStatus(assetID uint, newStatus string) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Assets{}).Where("id = ?", assetID).Update("status", newStatus).Error; err != nil {
			return err
		}
		return nil
	})
}

// ProcessEndOfDayUpdates performs end-of-day processing, potentially affecting multiple tables.
func (db *AssetDBModel) ProcessEndOfDayUpdates() error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		// Example: Archive old asset assignments.
		if err := tx.Where("end_date < ?", time.Now()).Model(&AssetAssignment{}).Update("status", "archived").Error; err != nil {
			return err
		}

		// Example: Update user statuses based on some criteria.
		if err := tx.Model(&Users{}).Where("last_login < ?", time.Now().AddDate(0, -1, 0)).Update("status", "inactive").Error; err != nil {
			return err
		}

		// Add more operations as needed.

		return nil
	})
}

func (db *AssetDBModel) ScheduleAssetMaintenance(schedule *AssetMaintenance) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(schedule).Error; err != nil {
			db.log.Error("Failed to schedule maintenance for asset %d: %v", schedule.AssetID, err)
			return fmt.Errorf("scheduling asset maintenance: %w", err)
		}
		db.log.Info("Maintenance scheduled for asset %d on %s", schedule.AssetID, schedule.ScheduledDate.Format("2006-01-02"))
		return nil
	})
}

func (db *AssetDBModel) RecordAssetPerformanceMetrics(metrics *AssetPerformanceMetrics) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(metrics).Error; err != nil {
			db.log.Error("Failed to record performance metrics for asset %d: %v", metrics.AssetID, err)
			return fmt.Errorf("recording asset performance metrics: %w", err)
		}
		db.log.Info("Performance metrics recorded for asset %d", metrics.AssetID)
		return nil
	})
}

func (db *AssetDBModel) CompleteAssetMaintenance(maintenanceID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&AssetMaintenance{}).Where("id = ?", maintenanceID).Update("status", "Completed").Error; err != nil {
			db.log.Error("Failed to complete maintenance for maintenance ID %d: %v", maintenanceID, err)
			return fmt.Errorf("completing asset maintenance: %w", err)
		}
		db.log.Info("Maintenance completed for maintenance ID %d", maintenanceID)
		return nil
	})
}

func (db *AssetDBModel) ListAssetMaintenanceSchedules(assetID uint) ([]*AssetMaintenance, error) {
	var schedules []*AssetMaintenance
	if err := db.DB.Where("asset_id = ?", assetID).Find(&schedules).Error; err != nil {
		db.log.Error("Failed to list maintenance schedules for asset %d: %v", assetID, err)
		return nil, fmt.Errorf("listing asset maintenance schedules: %w", err)
	}
	return schedules, nil
}

func (db *AssetDBModel) AnalyzeAssetPerformanceOverTime(assetID uint, start, end time.Time) (*AssetPerformanceAnalysis, error) {
	var analysis AssetPerformanceAnalysis
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		db.log.Info("Analyzing performance for asset %d between %s and %s", assetID, start.Format("2006-01-02"), end.Format("2006-01-02"))
		// Example logic to simulate performance analysis
		// Replace with actual data fetching and analysis
		result := tx.Table("asset_performance_records").Select("COUNT(*) as downtime_events, AVG(downtime_duration) as average_downtime").Where("asset_id = ? AND event_date BETWEEN ? AND ?", assetID, start, end).Scan(&analysis)
		if result.Error != nil {
			db.log.Error("Error analyzing performance for asset %d: %v", assetID, result.Error)
			return result.Error
		}
		analysis.AssetID = assetID
		analysis.Summary = "Performance met expected benchmarks with no significant downtime."
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error analyzing performance for asset %d: %w", assetID, err)
	}

	return &analysis, nil
}

// ScheduleMaintenance plans maintenance activities for an asset.
func (db *AssetDBModel) ScheduleMaintenance(maintenance *AssetMaintenance) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(maintenance).Error; err != nil {
			db.log.Error("Failed to schedule maintenance", err)
			return fmt.Errorf("scheduling maintenance: %w", err)
		}
		return nil
	})
}

// LogPerformanceMetric records performance metrics for an asset.
func (db *AssetDBModel) LogPerformanceMetric(performanceLog *AssetPerformanceLog) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(performanceLog).Error; err != nil {
			db.log.Error("Failed to log performance metric", err)
			return fmt.Errorf("logging performance metric: %w", err)
		}
		return nil
	})
}

// ListUpcomingMaintenance retrieves all upcoming maintenance activities for assets.
func (model *AssetDBModel) ListUpcomingMaintenance() ([]*AssetMaintenance, error) {
	var maintenances []*AssetMaintenance
	if err := model.DB.Where("scheduled_date > ?", time.Now()).Find(&maintenances).Error; err != nil {
		return nil, err // Failed to retrieve maintenance schedules
	}
	return maintenances, nil
}

// UpdateMaintenanceRecord updates details of a scheduled maintenance activity.
func (model *AssetDBModel) UpdateMaintenanceRecord(maintenanceID uint, updates map[string]interface{}) error {
	return model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&AssetMaintenance{}).Where("id = ?", maintenanceID).Updates(updates).Error; err != nil {
			return err // Failed to update maintenance record
		}
		return nil
	})
}

// RecordAssetPerformance logs performance metrics for an asset.
func (model *AssetDBModel) RecordAssetPerformance(performance *AssetPerformanceLog) error {
	return model.DB.Create(performance).Error // Assumes error handling outside
}

// RetrievePerformanceHistory fetches the performance log history of an asset.
func (model *AssetDBModel) RetrievePerformanceHistory(assetID uint) ([]*AssetPerformanceLog, error) {
	var logs []*AssetPerformanceLog
	if err := model.DB.Where("asset_id = ?", assetID).Find(&logs).Error; err != nil {
		return nil, err // Failed to retrieve performance logs
	}
	return logs, nil
}

// AnalyzeAssetUtilization provides insights into how effectively assets are being used.
func (model *AssetDBModel) AnalyzeAssetUtilization(assetID uint) (*AssetUtilizationReport, error) {
	// Placeholder for complex analysis logic, potentially involving machine learning or statistical analysis
	return &AssetUtilizationReport{}, nil
}

// ListUnderutilizedAssets identifies assets that are not being utilized to their full potential.
func (db *AssetDBModel) ListUnderutilizedAssets(threshold float64) ([]*Assets, error) {
	var underutilizedAssets []*Assets
	if err := db.DB.Where("utilization_rate < ?", threshold).Find(&underutilizedAssets).Error; err != nil {
		db.log.Error("Failed to list underutilized assets: %v", err)
		return nil, fmt.Errorf("listing underutilized assets: %w", err)
	}
	return underutilizedAssets, nil
}

// ListUserAssets retrieves all assets assigned to a specific user.
func (model *AssetDBModel) ListUserAssets(userID uint) ([]*Assets, error) {
	var assets []*Assets
	if err := model.DB.Joins("JOIN asset_user_assignments on asset_user_assignments.asset_id = assets.id").
		Where("asset_user_assignments.user_id = ?", userID).Find(&assets).Error; err != nil {
		return nil, fmt.Errorf("failed to list assets for user ID %d: %w", userID, err)
	}
	return assets, nil
}

// CalculateNewDepreciationValue updates the depreciation value of an asset using the Straight-Line Depreciation method.
func (model *AssetDBModel) CalculateNewDepreciationValue(assetID uint, yearsSincePurchase int) error {
	return model.DB.Transaction(func(tx *gorm.DB) error {
		var asset Assets
		err := tx.Find(&asset, assetID).Error
		if err != nil {
			return fmt.Errorf("asset not found for ID %d: %w", assetID, err)
		}

		// Assuming useful life and salvage value are part of the asset struct.
		// These values should ideally be set based on asset type or specific asset details.
		usefulLife := asset.UsefulLife // Assuming this is in years
		salvageValue := asset.SalvageValue

		if yearsSincePurchase > int(usefulLife) {
			// Asset has exceeded its useful life; depreciation value is the difference between purchase price and salvage value.
			newDepreciationValue := asset.PurchasePrice - salvageValue
			err = tx.Model(&asset).Update("depreciation_value", newDepreciationValue).Error
		} else {
			// Calculate annual depreciation expense
			annualDepreciationExpense := (asset.PurchasePrice - salvageValue) / float64(usefulLife)

			// Calculate total depreciation until now
			totalDepreciation := float64(yearsSincePurchase) * annualDepreciationExpense
			newDepreciationValue := asset.PurchasePrice - totalDepreciation

			// Update asset with new depreciation value
			err = tx.Model(&asset).Update("depreciation_value", newDepreciationValue).Error
		}

		if err != nil {
			return fmt.Errorf("updating depreciation value for asset %d: %w", assetID, err)
		}
		return nil
	})
}

// CalculateDepreciation updates the depreciation value of an asset over time.
func (model *AssetDBModel) CalculateDepreciation(assetID uint, depreciationRate float64) error {
	return model.DB.Transaction(func(tx *gorm.DB) error {
		var asset Assets
		if err := tx.First(&asset, assetID).Error; err != nil {
			return err // Asset not found
		}

		// Calculate new depreciation value based on some logic
		newValue := asset.PurchasePrice * (1 - depreciationRate)
		if err := tx.Model(&asset).Update("current_value", newValue).Error; err != nil {
			return err // Failed to update asset value
		}

		return nil
	})
}

// ListDepreciatedAssets retrieves assets that have depreciated below a certain threshold.
func (model *AssetDBModel) ListDepreciatedAssets(valueThreshold float64) ([]*Assets, error) {
	var assets []*Assets
	if err := model.DB.Where("current_value < ?", valueThreshold).Find(&assets).Error; err != nil {
		return nil, err // Failed to retrieve assets
	}
	return assets, nil
}

// LogAssetPerformance records performance metrics for an asset.
func (db *AssetDBModel) LogAssetPerformance(performanceLog *AssetPerformanceLog) error {
	tx := db.DB.Begin()
	if tx.Error != nil {
		db.log.Error("Failed to start transaction for logging asset performance", tx.Error)
		return tx.Error
	}

	if err := tx.Create(performanceLog).Error; err != nil {
		tx.Rollback()
		db.log.Error("Failed to log asset performance", err)
		return fmt.Errorf("logging asset performance: %w", err)
	}

	return tx.Commit().Error
}

// LogAuditEvent records an audit event related to asset management activities.
func (db *AssetDBModel) LogAuditEvent(auditLog *AssetAuditLog) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(auditLog).Error; err != nil {
			db.log.Error("Failed to log audit event", err)
			return fmt.Errorf("logging audit event: %w", err)
		}
		return nil
	})
}

// DetectAnomalousAccess uses patterns and heuristics to identify unusual or unauthorized asset access.
func (model *AssetDBModel) DetectAnomalousAccess() ([]*AssetAccessLog, error) {
	// Placeholder for logic to detect anomalies in asset access patterns
	return nil, nil
}

// GetAssetAudits retrieves all audit logs for a given asset.
func (model *AssetDBModel) GetAssetAudits(assetID uint) ([]*AssetAuditLog, error) {
	var audits []*AssetAuditLog
	if err := model.DB.Where("asset_id = ?", assetID).Find(&audits).Error; err != nil {
		return nil, err // Failed to retrieve audits
	}
	return audits, nil
}

// ListMaintenanceActivities retrieves scheduled maintenance activities for an asset.
func (db *AssetDBModel) ListMaintenanceActivities(assetID uint) ([]*AssetMaintenance, error) {
	var activities []*AssetMaintenance
	if err := db.DB.Where("asset_id = ?", assetID).Find(&activities).Error; err != nil {
		db.log.Error("Failed to list maintenance activities", err)
		return nil, fmt.Errorf("listing maintenance activities: %w", err)
	}
	return activities, nil
}

// ListLifecycleEvents retrieves all lifecycle events for a given asset.
func (db *AssetDBModel) ListLifecycleEvents(assetID uint) ([]*AssetLifecycleEvent, error) {
	var events []*AssetLifecycleEvent
	if err := db.DB.Where("asset_id = ?", assetID).Find(&events).Error; err != nil {
		db.log.Error("Failed to list lifecycle events", err)
		return nil, fmt.Errorf("listing lifecycle events: %w", err)
	}
	return events, nil
}

// ListPerformanceLogs retrieves all performance logs for a given asset.
func (db *AssetDBModel) ListPerformanceLogs(assetID uint) ([]*AssetPerformanceLog, error) {
	var logs []*AssetPerformanceLog
	if err := db.DB.Where("asset_id = ?", assetID).Find(&logs).Error; err != nil {
		db.log.Error("Failed to list performance logs", err)
		return nil, fmt.Errorf("listing performance logs: %w", err)
	}
	return logs, nil
}

// CompleteMaintenance marks a scheduled maintenance activity as completed.
func (db *AssetDBModel) CompleteMaintenance(maintenanceID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&AssetMaintenance{}).Where("id = ?", maintenanceID).Update("status", "Completed").Error; err != nil {
			db.log.Error("Failed to mark maintenance as completed", err)
			return fmt.Errorf("completing maintenance: %w", err)
		}
		return nil
	})
}

// UpdateLifecycleEvent modifies details of a recorded lifecycle event.
func (db *AssetDBModel) UpdateLifecycleEvent(eventID uint, updates map[string]interface{}) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&AssetLifecycleEvent{}).Where("id = ?", eventID).Updates(updates).Error; err != nil {
			db.log.Error("Failed to update lifecycle event", err)
			return fmt.Errorf("updating lifecycle event: %w", err)
		}
		return nil
	})
}

// RecordLifecycleEvent logs a significant event in the lifecycle of an asset.
func (db *AssetDBModel) RecordLifecycleEvent(event *AssetLifecycleEvent) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(event).Error; err != nil {
			db.log.Error("Failed to record lifecycle event", err)
			return fmt.Errorf("recording lifecycle event: %w", err)
		}
		return nil
	})
}

// GenerateUtilizationReport provides insights into the usage patterns of assets.
func (db *AssetDBModel) GenerateUtilizationReport(assetID uint) (*AssetUtilizationReport, error) {
	var report AssetUtilizationReport
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		db.log.Info("Generating utilization report for asset %d", assetID)
		// Simulate fetching utilization data from a database
		// Example logic, replace with actual query to calculate utilization based on asset logs
		result := tx.Table("asset_utilization_logs").Select("AVG(utilization_rate) as utilization_rate").Where("asset_id = ?", assetID).Scan(&report)
		if result.Error != nil {
			db.log.Error("Error generating utilization report for asset %d: %v", assetID, result.Error)
			return result.Error
		}
		report.AssetID = assetID
		report.ReportingPeriod = time.Now().String()
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error generating utilization report for asset %d: %w", assetID, err)
	}

	return &report, nil
}

func (db *AssetDBModel) InitiateAssetDecommissioning(decommission *AssetDecommission) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(decommission).Error; err != nil {
			db.log.Error("Failed to initiate decommissioning for asset %d: %v", decommission.AssetID, err)
			return fmt.Errorf("initiating asset decommissioning: %w", err)
		}
		return nil
	})
}

func (db *AssetDBModel) CompleteAssetDecommissioning(decommissionID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		update := map[string]interface{}{"status": "Completed"}
		if err := tx.Model(&AssetDecommission{}).Where("id = ?", decommissionID).Updates(update).Error; err != nil {
			db.log.Error("Failed to complete decommissioning for decommission ID %d: %v", decommissionID, err)
			return fmt.Errorf("completing asset decommissioning: %w", err)
		}
		return nil
	})
}

func (db *AssetDBModel) LogAssetRepair(log *AssetRepairLog) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(log).Error; err != nil {
			db.log.Error("Failed to log repair for asset %d: %v", log.AssetID, err)
			return fmt.Errorf("logging asset repair: %w", err)
		}
		db.log.Info("Repair logged for asset %d", log.AssetID)
		return nil
	})
}

func (db *AssetDBModel) CalculateAssetDepreciation(record *AssetDepreciationRecord) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(record).Error; err != nil {
			db.log.Error("Failed to calculate depreciation for asset %d: %v", record.AssetID, err)
			return fmt.Errorf("calculating asset depreciation: %w", err)
		}
		db.log.Info("Depreciation calculated for asset %d", record.AssetID)
		return nil
	})
}

func (db *AssetDBModel) RetrieveAssetRepairHistory(assetID uint) ([]*AssetRepairLog, error) {
	var logs []*AssetRepairLog
	if err := db.DB.Where("asset_id = ?", assetID).Find(&logs).Error; err != nil {
		db.log.Error("Failed to retrieve repair history for asset %d: %v", assetID, err)
		return nil, fmt.Errorf("retrieving asset repair history: %w", err)
	}
	return logs, nil
}

func (db *AssetDBModel) ListAssetDepreciationRecords(assetID uint) ([]*AssetDepreciationRecord, error) {
	var records []*AssetDepreciationRecord
	if err := db.DB.Where("asset_id = ?", assetID).Find(&records).Error; err != nil {
		db.log.Error("Failed to list depreciation records for asset %d: %v", assetID, err)
		return nil, fmt.Errorf("listing asset depreciation records: %w", err)
	}
	return records, nil
}

// UpdateAssetCategory changes the category of an asset.
func (s *AssetDBModel) UpdateAssetCategory(assetID uint, newCategoryID uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Assets{}).Where("id = ?", assetID).Update("category_id", newCategoryID).Error; err != nil {
			s.log.Error(fmt.Sprintf("Failed to update category for asset %d: %v", assetID, err))
			return err
		}
		return nil
	})
}

// CalculateNewDepreciationValue updates the depreciation value of an asset using the Straight-Line Depreciation method.
// This version ensures flexibility for different asset types, logical depreciation calculations, and includes detailed logging.
func (model *AssetDBModel) CalculateNewDepreciationValue1(assetID uint, yearsSincePurchase int) error {
	return model.DB.Transaction(func(tx *gorm.DB) error {
		var asset Assets
		// Fetch the asset to determine its purchase price, useful life, and salvage value.
		if err := tx.First(&asset, assetID).Error; err != nil {
			model.log.Error("Asset not found", "assetID", assetID, "error", err)
			return fmt.Errorf("asset not found for ID %d: %w", assetID, err)
		}

		// Validate input and asset's parameters to ensure they are within logical and valid ranges.
		if yearsSincePurchase < 0 {
			model.log.Info("Invalid years since purchase", "assetID", assetID, "yearsSincePurchase", yearsSincePurchase)
			return fmt.Errorf("invalid years since purchase for asset ID %d", assetID)
		}
		if asset.UsefulLife <= 0 || asset.SalvageValue < 0 {
			model.log.Info("Invalid asset parameters", "assetID", assetID, "UsefulLife", asset.UsefulLife, "SalvageValue", asset.SalvageValue)
			return fmt.Errorf("asset with ID %d has invalid parameters", assetID)
		}

		// Calculate depreciation using the Straight-Line method, adjusting for the asset's actual useful life.
		annualDepreciation := (asset.PurchasePrice - asset.SalvageValue) / float64(asset.UsefulLife)
		depreciationAccumulated := float64(min(yearsSincePurchase, int(asset.UsefulLife))) * annualDepreciation

		// Cap the depreciation at the maximum possible value to prevent it from exceeding logical limits.
		maxDepreciation := asset.PurchasePrice - asset.SalvageValue
		if depreciationAccumulated > maxDepreciation {
			depreciationAccumulated = maxDepreciation
		}

		// Update the asset with the newly calculated depreciation value.
		if err := tx.Model(&asset).Update("depreciation_value", depreciationAccumulated).Error; err != nil {
			model.log.Error("Failed to update depreciation value", "assetID", assetID, "error", err)
			return fmt.Errorf("updating depreciation value for asset %d: %w", assetID, err)
		}

		model.log.Info("Depreciation value updated successfully", "assetID", assetID, "NewDepreciationValue", depreciationAccumulated)
		return nil
	})
}

//  ............................. .

// CalculateNewDepreciationValue updates the depreciation value of an asset using the Straight-Line Depreciation method,
// accounting for various asset types and ensuring logical depreciation calculations.
func (model *AssetDBModel) CalculateNewDepreciationValue1Extra(assetID uint, yearsSincePurchase int) error {
	return model.DB.Transaction(func(tx *gorm.DB) error {
		var asset Assets
		// Fetch the asset to get its purchase price, useful life, and salvage value.
		if err := tx.First(&asset, assetID).Error; err != nil {
			model.log.Error("Asset not found", "assetID", assetID, "error", err)
			return fmt.Errorf("asset not found for ID %d: %w", assetID, err)
		}

		// Validate input to ensure valid calculation parameters.
		if yearsSincePurchase < 0 || asset.UsefulLife <= 0 || asset.SalvageValue < 0 || asset.SalvageValue >= asset.PurchasePrice {
			model.log.Info("Invalid calculation parameters", "assetID", assetID, "yearsSincePurchase", yearsSincePurchase, "UsefulLife", asset.UsefulLife, "SalvageValue", asset.SalvageValue)
			return fmt.Errorf("invalid calculation parameters for asset ID %d", assetID)
		}

		// Calculate depreciation using the Straight-Line method.
		annualDepreciation := (asset.PurchasePrice - asset.SalvageValue) / float64(asset.UsefulLife)
		depreciationAccumulated := float64(min(yearsSincePurchase, int(asset.UsefulLife))) * annualDepreciation

		// Cap the depreciation at the maximum possible value.
		maxDepreciation := asset.PurchasePrice - asset.SalvageValue
		if depreciationAccumulated > maxDepreciation {
			depreciationAccumulated = maxDepreciation
		}

		// Update the asset with the new depreciation value.
		if err := tx.Model(&asset).Update("depreciation_value", depreciationAccumulated).Error; err != nil {
			model.log.Error("Failed to update depreciation value", "assetID", assetID, "error", err)
			return fmt.Errorf("updating depreciation value for asset %d: %w", assetID, err)
		}

		model.log.Info("Depreciation value updated successfully", "assetID", assetID, "NewDepreciationValue", depreciationAccumulated)
		return nil
	})
}

// min is a helper function to find the minimum of two integers.
func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// CalculateNewDepreciationValue updates the depreciation value of an asset using the Straight-Line Depreciation method.
func (model *AssetDBModel) CalculateNewDepreciationValue2(assetID uint, yearsSincePurchase int) error {
	return model.DB.Transaction(func(tx *gorm.DB) error {
		var asset Assets
		err := tx.Find(&asset, assetID).Error
		if err != nil {
			return fmt.Errorf("asset not found for ID %d: %w", assetID, err)
		}

		// Assuming useful life and salvage value are part of the asset struct.
		// These values should ideally be set based on asset type or specific asset details.
		usefulLife := asset.UsefulLife // Assuming this is in years
		salvageValue := asset.SalvageValue

		if yearsSincePurchase > int(usefulLife) {
			// Asset has exceeded its useful life; depreciation value is the difference between purchase price and salvage value.
			newDepreciationValue := asset.PurchasePrice - salvageValue
			err = tx.Model(&asset).Update("depreciation_value", newDepreciationValue).Error
		} else {
			// Calculate annual depreciation expense
			annualDepreciationExpense := (asset.PurchasePrice - salvageValue) / float64(usefulLife)

			// Calculate total depreciation until now
			totalDepreciation := float64(yearsSincePurchase) * annualDepreciationExpense
			newDepreciationValue := asset.PurchasePrice - totalDepreciation

			// Update asset with new depreciation value
			err = tx.Model(&asset).Update("depreciation_value", newDepreciationValue).Error
		}

		if err != nil {
			return fmt.Errorf("updating depreciation value for asset %d: %w", assetID, err)
		}
		return nil
	})
}

func (s *AssetDBModel) ScheduleAssetInspection(inspection *AssetInspectionRecord) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(inspection).Error; err != nil {
			return fmt.Errorf("scheduling asset inspection: %w", err)
		}
		return nil
	})
}
