// backend/models/assets.go

package models

import (
	"time"

	"gorm.io/gorm"
)

type Assets struct {
	gorm.Model
	ID            uint            `gorm:"primaryKey" json:"asset_id"`
	Tag           AssetTag        `json:"asset_tag" gorm:"foreignKey:AssetID"`
	AssetType     AssetType       `json:"asset_type" gorm:"embedded"`
	AssetName     string          `json:"asset_name"`
	Assignment    AssetAssignment `json:"asset_assignment" gorm:"foreignKey:AssetID"`
	Description   string          `json:"description"`
	Manufacturer  string          `json:"manufacturer"`
	Asset_Model   string          `json:"model"`
	SerialNumber  string          `json:"serial_number"`
	PurchaseDate  time.Time       `json:"purchase_date"`
	PurchasePrice string          `json:"purchase_price"`
	Vendor        string          `json:"vendor"`
	Site          string          `json:"site"`
	Status        string          `json:"status"`
	CreatedBy     uint            `json:"created_by"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

// TableName sets the table name for the Asset model.
func (Assets) TableName() string {
	return "assets"
}

// Hashtag represents a hashtag entity
type AssetTag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	AssetTag  string    `json:"tag"`
	Tags      []string  `json:"tags"` // Added Tags field
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	AssetID   uint      `json:"-"`
}

// TableName sets the table name for the AssetTags model.
func (AssetTag) TableName() string {
	return "asset_tag"
}

type AssetType struct {
	gorm.Model
	ID        int       `gorm:"primaryKey" json:"asset_type_id"`
	AssetType string    `json:"asset_type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName sets the table name for the AssetType model.
func (AssetType) TableName() string {
	return "assetType"
}

type AssetAssignment struct {
	gorm.Model
	ID             int       `gorm:"primaryKey" json:"assignment_id"`
	AssetID        int       `json:"_"`
	UserID         int       `json:"user_id"`
	AssignedBy     int       `json:"assigned_by"`
	AssignmentType string    `json:"assignment_type"`
	DueAt          time.Time `json:"due_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// TableName sets the table name for the Asset Assignment model.
func (AssetAssignment) TableName() string {
	return "asset_assignment"
}

type AssetsStorage interface {
	CreateAsset(*Assets) error
	DeleteAsset(int) error
	UpdateAsset(*Assets) error
	GetAssets() ([]*Assets, error)
	GetAssetByID(int) (*Assets, error)
	GetAssetByNumber(int) (*Assets, error)
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
	CreateAssetAssignment(*AssetAssignment) error
	DeleteAssetAssignment(int) error
	UpdateAssetAssignment(*AssetAssignment) error
	GetAssetAssignment() ([]*AssetAssignment, error)
	GetAssetAssignmentByID(int) (*AssetAssignment, error)
	GetAssetAssignmentByNumber(int) (*AssetAssignment, error)
}

// AssetModel handles database operations for Asset
type AssetDBModel struct {
	DB *gorm.DB
}

// NewAssetModel creates a new instance of TicketModel
func NewAssetDBModel(db *gorm.DB) *AssetDBModel {
	return &AssetDBModel{
		DB: db,
	}
}

// CreateAssets creates a new asset.
func (as *AssetDBModel) CreateAsset(asset *Assets) error {
	return as.DB.Create(asset).Error
}

// GetAssetsByID retrieves a user by its ID.
func (as *AssetDBModel) GetAssetByID(id uint) (*Assets, error) {
	var asset Assets
	err := as.DB.Where("id = ?", id).First(&asset).Error
	return &asset, err
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
func (as *AssetDBModel) GetAllAssets() (*[]Assets, error) {
	var assets []Assets
	err := as.DB.Find(&assets).Error
	return &assets, err
}
