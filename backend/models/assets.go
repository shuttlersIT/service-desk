// backend/models/assets.go

package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Assets struct {
	gorm.Model
	AssetTag      string     `gorm:"size:255;not null;unique" json:"asset_tag"`
	AssetTypeID   uint       `gorm:"index;not null" json:"asset_type_id"`
	Name          string     `gorm:"size:255;not null" json:"name"`
	Description   string     `gorm:"type:text" json:"description"`
	Manufacturer  string     `gorm:"size:255" json:"manufacturer"`
	AssetModel    string     `gorm:"size:255" json:"model"`
	SerialNumber  string     `gorm:"size:255;unique" json:"serial_number"`
	PurchaseDate  *time.Time `json:"purchase_date,omitempty"`
	PurchasePrice float64    `gorm:"type:decimal(10,2)" json:"purchase_price"`
	Vendor        string     `gorm:"size:255" json:"vendor"`
	Status        string     `gorm:"size:100;not null" json:"status"`
	Location      string     `gorm:"size:255" json:"location"`
	// Assuming UserID is the identifier for the user to whom the asset is currently assigned.
	UserID              uint       `gorm:"index" json:"user_id,omitempty"`
	SiteID              uint       `gorm:"index" json:"site_id"`
	CreatedByID         uint       `gorm:"index" json:"created_by_id"`
	AssetAssignment     *uint      `json:"assetAssignment" gorm:"foreignKey:AssetAssignmentID"`
	LastMaintenanceDate *time.Time `json:"last_maintenance_date,omitempty"`
	NextMaintenanceDate *time.Time `json:"next_maintenance_date,omitempty"`
}

func (Assets) TableName() string {
	return "assets"
}

type EventParticipant struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	EventID  uint      `json:"event_id" gorm:"index;not null"`
	AgentID  uint      `json:"agent_id" gorm:"index;not null"`
	JoinedAt time.Time `json:"joined_at"`
}

func (EventParticipant) TableName() string {
	return "event_participant"
}

// Hashtag represents a hashtag entity
type AssetTag struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	AssetTag  string         `json:"tag"`
	Tags      []string       `json:"tags"` // Added Tags field
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	AssetID   uint           `json:"-"`
}

// TableName sets the table name for the AssetTags model.
func (AssetTag) TableName() string {
	return "asset_tag"
}

type AssetType struct {
	ID        uint           `gorm:"primaryKey" json:"asset_type_id"`
	AssetType string         `json:"asset_type"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName sets the table name for the AssetType model.
func (AssetType) TableName() string {
	return "assetType"
}

type AssetAssignment struct {
	ID               uint           `gorm:"primaryKey" json:"assignment_id"`
	AssetID          uint           `gorm:"index;not null" json:"asset_id"`
	UserID           uint           `gorm:"index;not null" json:"user_id"`
	AssignedBy       uint           `gorm:"index;not null" json:"assigned_by"`
	AssignmentType   string         `gorm:"size:255" json:"assignment_type"`
	AssignmentStatus string         `gorm:"size:255" json:"assignment_status"`
	DueAt            *time.Time     `json:"due_at,omitempty"` // Made optional
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	// DeletedAt removed to align with gorm.Model inclusion
}

func (AssetAssignment) TableName() string {
	return "asset_assignments"
}

type Resource struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	Name               string         `json:"name" binding:"required"`
	Description        string         `json:"description"`
	Type               string         `json:"type"`
	Location           string         `json:"location"`
	AvailabilityStatus string         `json:"availability_status"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	Bookings           []Booking      `json:"bookings" gorm:"foreignKey:ResourceID"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Resource) TableName() string {
	return "resources"
}

type Booking struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	ResourceID uint           `json:"resource_id"`
	UserID     uint           `json:"user_id"`
	StartTime  time.Time      `json:"start_time"`
	EndTime    time.Time      `json:"end_time"`
	Status     string         `json:"status"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Booking) TableName() string {
	return "bookings"
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

// UpdateAssets updates the details of an existing asset.
func (as *AssetDBModel) UpdateAsset(asset *Assets) error {
	if err := as.DB.Save(asset).Error; err != nil {
		return err
	}
	return nil
}

// DeleteAssets deletes a asset from the database.
func (as *AssetDBModel) DeleteAsset(id uint) error {
	if err := as.DB.Delete(&Assets{}, id).Error; err != nil {
		return err
	}
	return nil
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

type AssetDepreciation struct {
	AssetID      uint      `json:"asset_id" gorm:"index;not null"`
	Depreciation float64   `json:"depreciation"`
	RecordedAt   time.Time `json:"recorded_at"`
}

func (am *AssetDBModel) ScheduleAssetAudit(auditDate time.Time) error {
	// Schedule an audit for the specified date
}
