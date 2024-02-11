// backend/models/assets.go

package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type DeviceRegistration struct {
	gorm.Model
	UserID   uint   `gorm:"index;not null" json:"user_id"`
	DeviceID string `gorm:"type:varchar(255);not null" json:"device_id"` // Unique identifier for the device
	Platform string `gorm:"type:varchar(100);not null" json:"platform"`  // E.g., "iOS", "Android"
	Token    string `gorm:"type:text;not null" json:"token"`             // Token for push notifications
}

func (DeviceRegistration) TableName() string {
	return "device_registrations"
}

type Assets struct {
	gorm.Model
	AssetTag        string     `gorm:"size:100;not null;unique" json:"asset_tag"`
	AssetTypeID     uint       `gorm:"index;not null" json:"asset_type_id"`
	AssetName       string     `gorm:"size:255;not null" json:"asset_name"`
	Description     string     `gorm:"type:text" json:"description,omitempty"`
	Manufacturer    string     `gorm:"size:255" json:"manufacturer,omitempty"`
	AssetModel      string     `gorm:"size:255" json:"model,omitempty"`
	SerialNumber    string     `gorm:"size:255;unique" json:"serial_number"`
	PurchaseDate    *time.Time `json:"purchase_date,omitempty"`
	WarrantyExpire  *time.Time `json:"warranty_expire,omitempty"`
	PurchasePrice   float64    `gorm:"type:decimal(10,2)" json:"purchase_price,omitempty"`
	Status          string     `gorm:"size:100;not null" json:"status"`
	Vendor          Vendor     `gorm:"foreignKey:VendorID" json:"-"`
	User            Users      `gorm:"foreignKey:UserID" json:"-"`
	Location        string     `gorm:"size:255" json:"location,omitempty"`
	UserID          *uint      `gorm:"index" json:"user_id,omitempty"`
	SiteID          uint       `gorm:"index" json:"site_id"`
	CreatedByID     uint       `gorm:"index" json:"created_by_id"`
	AssetAssignment *uint      `json:"assetAssignment" gorm:"foreignKey:AssetAssignmentID"`
}

func (Assets) TableName() string {
	return "assets"
}

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

// Hashtag represents a hashtag entity
type AssetTag struct {
	gorm.Model
	AssetTag string   `json:"tag"`
	Tags     []string `gorm:"type:text[]" json:"tags"` // Use pq.StringArray for PostgreSQL compatibility
	AssetID  uint     `gorm:"index" json:"asset_id"`
}

// TableName sets the table name for the AssetTags model.
func (AssetTag) TableName() string {
	return "asset_tag"
}

type AssetType struct {
	gorm.Model
	AssetType string `json:"asset_type"`
}

// TableName sets the table name for the AssetType model.
func (AssetType) TableName() string {
	return "assetType"
}

type AssetAssignment struct {
	gorm.Model
	AssetID          uint       `gorm:"index;not null" json:"asset_id"`
	UserID           uint       `gorm:"index;not null" json:"user_id"`
	AssignedBy       uint       `gorm:"index;not null" json:"assigned_by"`
	AssignmentType   string     `gorm:"size:255" json:"assignment_type"`
	AssignmentStatus string     `gorm:"size:255" json:"assignment_status"`
	DueAt            *time.Time `json:"due_at,omitempty"`
}

func (AssetAssignment) TableName() string {
	return "asset_assignments"
}

type Resource struct {
	gorm.Model
	Name               string    `json:"name" binding:"required"`
	Description        string    `json:"description,omitempty"`
	Type               string    `json:"type"`
	Location           string    `json:"location,omitempty"`
	AvailabilityStatus string    `json:"availability_status"`
	Metadata           string    `gorm:"type:text" json:"metadata,omitempty"`      // JSON for storing additional information
	Status             string    `gorm:"type:varchar(100);not null" json:"status"` // E.g., "active", "maintenance"
	Bookings           []Booking `gorm:"foreignKey:ResourceID" json:"bookings,omitempty"`
}

func (Resource) TableName() string {
	return "resources"
}

type ResourceAllocation struct {
	gorm.Model
	ResourceID        uint      `gorm:"index;not null" json:"resource_id"`
	UserID            uint      `gorm:"index" json:"user_id,omitempty"`      // Optional, if allocated to a user
	AllocationDetails string    `gorm:"type:text" json:"allocation_details"` // JSON detailing allocation
	StartTime         time.Time `json:"start_time"`
	EndTime           time.Time `json:"end_time"`
}

func (ResourceAllocation) TableName() string {
	return "resource_allocations"
}

type AgentResourceAllocation struct {
	ID                uint            `gorm:"primaryKey" json:"id"`
	AgentID           uint            `json:"agent_id" gorm:"index;not null"`
	ResourceID        uint            `json:"resource_id" gorm:"index;not null"`
	AllocationDetails string          `gorm:"type:text" json:"allocation_details,omitempty"` // JSON detailing allocation
	StartTime         time.Time       `json:"start_time"`
	EndTime           time.Time       `json:"end_time"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
	DeletedAt         *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (AgentResourceAllocation) TableName() string {
	return "agent_resource_allocations"
}

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

type AccessControlList struct {
	gorm.Model
	Resource   string `gorm:"type:varchar(255);not null" json:"resource"` // E.g., "article", "user_profile"
	Action     string `gorm:"type:varchar(100);not null" json:"action"`   // E.g., "read", "write"
	RoleID     uint   `gorm:"index" json:"role_id"`
	Permission bool   `json:"permission"`
}

func (AccessControlList) TableName() string {
	return "access_control_lists"
}

type ComplianceAuditLog struct {
	gorm.Model
	Action      string `gorm:"type:varchar(255);not null"`
	UserID      uint   `gorm:"index"`
	Description string `gorm:"type:text;not null"`
	Details     string `gorm:"type:text"` // JSON format recommended
}

func (ComplianceAuditLog) TableName() string {
	return "compliance_audit_logs"
}

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

type AssetsStorage interface {
	CreateAsset(*Assets) error
	DeleteAsset(int) error
	UpdateAsset(*Assets) error
	GetAssets() ([]*Assets, error)
	GetAssetByID(int) (*Assets, error)
	GetAssetByNumber(int) (*Assets, error)
	UnassignAsset(uint) error
	AssignAssetToUser2(assetID, userID uint) error
	UnassignAssetUser2(assetID uint) error
}

type AssetTypeStorage interface {
	CreateAssetType(*AssetType) error
	DeleteAssetType(int) error
	UpdateAssetType(*AssetType) error
	GetAssetType() ([]*AssetType, error)
	GetAssetTypeByID(int) (*AssetType, error)
	GetAssetTypeByNumber(int) (*AssetType, error)
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
	DB *gorm.DB
}

// NewAssetAssignmentModel creates a new instance of TicketModel
func NewAssetAssignmentDBModel(db *gorm.DB) *AssetAssignmentDBModel {
	return &AssetAssignmentDBModel{
		DB: db,
	}
}

// AssetModel handles database operations for Asset
type AssetDBModel struct {
	DB              *gorm.DB
	AssetAssignment *AssetAssignmentDBModel
}

// NewAssetModel creates a new instance of TicketModel
func NewAssetDBModel(db *gorm.DB, assetAssignment *AssetAssignmentDBModel) *AssetDBModel {
	return &AssetDBModel{
		DB:              db,
		AssetAssignment: assetAssignment,
	}
}

// CreateAssets creates a new asset.
func (as *AssetDBModel) CreateAsset(asset *Assets) error {
	return as.DB.Create(asset).Error
}

// GetAssetByID retrieves an asset by its ID.
func (as *AssetDBModel) GetAssetByID(id uint) (*Assets, error) {
	var asset Assets
	err := as.DB.Where("id = ?", id).First(&asset).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

// UpdateAsset updates an existing asset's details.
func (db *AssetDBModel) UpdateAsset(assetID uint, updates map[string]interface{}) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Assets{}).Where("id = ?", assetID).Updates(updates).Error; err != nil {
			return err
		}
		return nil
	})
}

// DeleteAssets deletes a asset from the database.
func (as *AssetDBModel) DeleteAsset(id uint) error {
	if err := as.DB.Delete(&Assets{}, id).Error; err != nil {
		return err
	}
	return nil
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

// GetAllAssets retrieves all Assets from the database.
func (as *AssetDBModel) GetAllAssets() ([]*Assets, error) {
	var assets []*Assets
	err := as.DB.Find(&assets).Error
	if err != nil {
		return nil, err
	}
	return assets, err
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

// UpdateAssetAssignment updates the assignment of an existing asset.
func (as *AssetAssignmentDBModel) UnassignAsset(assetAssignment *AssetAssignment, agentID uint) (*AssetAssignment, error) {
	due_at := time.Now().AddDate(1, 0, 0)
	assetAssignment, err := as.GetAssetAssignmentByID(assetAssignment.ID)
	if err != nil {
		return nil, err
	}
	newAssetAssignment := &AssetAssignment{
		AssetID:          assetAssignment.AssetID,
		UserID:           0,
		AssignedBy:       agentID,      // Assuming the same user assigns the asset
		AssignmentType:   "unassigned", // Update as needed
		AssignmentStatus: "unassigned",
		DueAt:            &due_at, // Due date example
		CreatedAt:        time.Now(),
	}

	erro := as.CreateAssetAssignment(newAssetAssignment)
	if erro != nil {
		return nil, fmt.Errorf("unable to assign asset")
	}

	return newAssetAssignment, nil
}

func (tdb *AssetDBModel) AssignAssetToUser(assetID, userID uint) error {
	// Assign an asset to a user
	asset := &Assets{}
	if err := tdb.DB.First(asset, assetID).Error; err != nil {
		return err
	}

	// Update the asset's UserID
	asset.AssetAssignment = &userID
	if err := tdb.DB.Save(asset).Error; err != nil {
		return err
	}

	return nil
}

func (ass *AssetAssignmentDBModel) CreateAssetAssignment(assetAssignment *AssetAssignment) error {
	return ass.DB.Create(assetAssignment).Error
}

func (tdb *AssetDBModel) UnassignAsset(assetID uint) error {
	due_at := time.Now().AddDate(1, 0, 0)
	// Unassign an asset from a user
	asset := &Assets{}
	if err := tdb.DB.First(asset, assetID).Error; err != nil {
		return err
	}

	newAssetAssignment := &AssetAssignment{
		AssetID:          assetID,
		UserID:           0,
		AssignedBy:       0,            // Assuming the same user assigns the asset
		AssignmentType:   "unassigned", // Update as needed
		AssignmentStatus: "unassigned",
		DueAt:            &due_at, // Due date example
		CreatedAt:        time.Now(),
	}

	// Clear the asset's UserID
	asset.AssetAssignment = &newAssetAssignment.ID
	if err := tdb.DB.Save(asset).Error; err != nil {
		return err
	}

	return nil
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

// UpdateAssetDetails allows for updating specific attributes of an asset.
func (db *AssetDBModel) UpdateAssetDetails(assetID uint, updates map[string]interface{}) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Assets{}).Where("id = ?", assetID).Updates(updates).Error; err != nil {
			return err
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

// CreateAssetAssignment creates a new asset assignment for a user with transactional integrity.
func (db *UserDBModel) CreateAssetAssignment(assetAssignment *AssetAssignment) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(assetAssignment).Error; err != nil {
			return err
		}
		return nil
	})
}

// UpdateAssetAssignment updates the details of an existing asset assignment with transactional integrity.
func (db *UserDBModel) UpdateAssetAssignment(assetAssignment *AssetAssignment) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(assetAssignment).Error; err != nil {
			return err
		}
		return nil
	})
}

// DeleteAssetAssignment deletes an asset assignment from the database with transactional integrity.
func (db *UserDBModel) DeleteAssetAssignment(id uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&AssetAssignment{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}

// GetAssetAssignmentsByUser retrieves asset assignments for a user by their user ID.
func (as *UserDBModel) GetAssetAssignmentsByUser(userID uint) ([]*AssetAssignment, error) {
	var assetAssignments []*AssetAssignment
	err := as.DB.Where("user_id = ?", userID).Find(&assetAssignments).Error
	if err != nil {
		return nil, err
	}
	return assetAssignments, nil
}

// GetAssetAssignmentsByAsset retrieves asset assignments for an asset by its asset ID.
func (as *UserDBModel) GetAssetAssignmentsByAsset(assetID uint) ([]*AssetAssignment, error) {
	var assetAssignments []*AssetAssignment
	err := as.DB.Where("asset_id = ?", assetID).Find(&assetAssignments).Error
	if err != nil {
		return nil, err
	}
	return assetAssignments, nil
}

// GetAssetAssignmentByID retrieves an asset assignment by its ID.
func (as *UserDBModel) GetAssetAssignmentByID(id uint) (*AssetAssignment, error) {
	var assetAssignment AssetAssignment
	err := as.DB.Where("id = ?", id).First(&assetAssignment).Error
	if err != nil {
		return nil, err
	}
	return &assetAssignment, nil
}

// GetAssetAssignmentByNumber retrieves an asset assignment by its asset assignment number.
func (as *UserDBModel) GetAssetAssignmentByNumber(assetAssignmentNumber int) (*AssetAssignment, error) {
	var assetAssignment AssetAssignment
	err := as.DB.Where("asset_assignment_id = ?", assetAssignmentNumber).First(&assetAssignment).Error
	if err != nil {
		return nil, err
	}
	return &assetAssignment, nil
}

// GetAssetByID retrieves an asset by its ID.
func (as *UserDBModel) GetAssetByID(id uint) (*Assets, error) {
	var asset Assets
	err := as.DB.Where("id = ?", id).First(&asset).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
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

// TransferAsset between users with transactional integrity.
func (db *UserDBModel) TransferAsset(assetID, fromUserID, toUserID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		// Step 1: Unassign the asset from the current user.
		if err := tx.Model(&AssetAssignment{}).Where("asset_id = ? AND user_id = ?", assetID, fromUserID).Delete(&AssetAssignment{}).Error; err != nil {
			return err
		}

		// Step 2: Assign the asset to the new user.
		newAssignment := AssetAssignment{AssetID: assetID, UserID: toUserID}
		if err := tx.Create(&newAssignment).Error; err != nil {
			return err
		}

		return nil
	})
}

// UpdateAssetStatus updates an asset's status with transactional integrity.
func (db *UserDBModel) UpdateAssetStatus(assetID uint, newStatus string) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Assets{}).Where("id = ?", assetID).Update("status", newStatus).Error; err != nil {
			return err
		}
		return nil
	})
}

// ActivateUsers activates a list of users by their IDs with transactional integrity.
func (db *UserDBModel) ActivateUsers(userIDs []uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Users{}).Where("id IN ?", userIDs).Update("active", true).Error; err != nil {
			return err
		}
		return nil
	})
}

// DeactivateUsers deactivates a list of users by their IDs with transactional integrity.
func (db *UserDBModel) DeactivateUsers(userIDs []uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Users{}).Where("id IN ?", userIDs).Update("active", false).Error; err != nil {
			return err
		}
		return nil
	})
}

// ProcessEndOfDayUpdates performs end-of-day processing, potentially affecting multiple tables.
func (db *UserDBModel) ProcessEndOfDayUpdates() error {
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
