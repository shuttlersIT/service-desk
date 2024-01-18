// backend/services/advertisement_service.go

package services

import (
	"github.com/shuttlersit/service-desk/backend/models"
	"gorm.io/gorm"
)

// AssetServiceInterface provides methods for managing assets.
type AssetServiceInterface interface {
	CreateAsset(asset *models.Assets) (*models.Assets, error)
	UpdateAsset(asset *models.Assets) (*models.Assets, error)
	GetAssetByID(id uint) (*models.Assets, error)
	DeleteAsset(assetID uint) (bool, error)
	GetAllAssets() *[]models.Assets
}

// DefaultAssetService is the default implementation of AssettService
type DefaultAssetService struct {
	DB           *gorm.DB
	AssetDBModel *models.AssetDBModel
	// Add any dependencies or data needed for the service
}

// NewDefaultAssetService creates a new DefaultAssetService.
func NewDefaultAssetService(assetDBModel *models.AssetDBModel) *DefaultAssetService {
	return &DefaultAssetService{
		AssetDBModel: assetDBModel,
	}
}

// GetAllAssets retrieves all assets.
func (ps *DefaultAssetService) GetAllAssets() (*[]models.Assets, error) {
	assets, err := ps.AssetDBModel.GetAllAssets()
	if err != nil {
		return nil, err
	}
	return assets, nil
}

// CreateAsset creates a new Asset.
func (ps *DefaultAssetService) CreateAsset(asset *models.Assets) error {
	err := ps.AssetDBModel.CreateAsset(asset)
	if err != nil {
		return err
	}
	return nil
}

// CreateAsset creates a new asset.
func (ps *DefaultAssetService) GetAssetByID(id uint) (*models.Assets, error) {
	asset, err := ps.AssetDBModel.GetAssetByID(id)
	if err != nil {
		return nil, err
	}
	return asset, nil
}

// UpdateAsset updates an existing asset.
func (ps *DefaultAssetService) UpdateAsset(asset *models.Assets) (*models.Assets, error) {
	err := ps.AssetDBModel.UpdateAsset(asset)
	if err != nil {
		return nil, err
	}
	return asset, nil
}

// DeleteAsset deletes an asset by ID.
func (ps *DefaultAssetService) DeleteAsset(id uint) (bool, error) {
	status := false
	err := ps.AssetDBModel.DeleteAsset(id)
	if err != nil {
		return status, err
	}
	status = true
	return status, nil
}
