// backend/services/advertisement_service.go

package services

import (
	"time"

	"github.com/shuttlersit/service-desk/backend/models"
	"gorm.io/gorm"
)

// AssetServiceInterface provides methods for managing assets.
type AssetServiceInterface interface {
	CreateAsset(asset *models.Assets) error
	DeleteAsset(assetID uint) error
	UpdateAsset(asset *models.Assets) error
	GetAssetByID(assetID uint) (*models.Assets, error)
	GetAllAssets() ([]*models.Assets, error)
	GetAssetByNumber(assetNumber int) (*models.Assets, error)
	AssignAsset(assetID, userID uint, agentID uint) error
	AssignAssetToUser(assetID, userID uint, agentID uint) (*models.AssetAssignment, error)
	UnassignAsset(assetID uint, agentID uint) error
	UnassignAssetFromUser(assetID uint, agentID uint) error
}

// DefaultAssetService is the default implementation of AssettService
type DefaultAssetService struct {
	DB                     *gorm.DB
	AssetDBModel           *models.AssetDBModel
	AssetAssignmentDBModel *models.AssetAssignmentDBModel
	// Add any dependencies or data needed for the service
}

// NewDefaultAssetService creates a new DefaultAssetService.
func NewDefaultAssetService(assetDBModel *models.AssetDBModel, assetAssignmentDBModel *models.AssetAssignmentDBModel) *DefaultAssetService {
	return &DefaultAssetService{
		AssetDBModel:           assetDBModel,
		AssetAssignmentDBModel: assetAssignmentDBModel,
	}
}

// GetAllAssets retrieves all assets.
func (ps *DefaultAssetService) GetAllAssets() ([]*models.Assets, error) {
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

// backend/services/asset_service.go

func (as *DefaultAssetService) AssignAsset(assetID, userID uint, agentID uint) error {
	// Retrieve the asset by assetID
	// Update the asset user ID
	asset, err := as.GetAssetByID(assetID)
	if err != nil {
		return err
	}

	assetAssignment, err := as.AssetAssignmentDBModel.AssignAsset(asset, userID)
	if err != nil {
		return err
	}

	asset.Status = "assigned"
	asset.Assignment = *assetAssignment

	// Save the updated asset
	err = as.AssetDBModel.UpdateAsset(asset)
	if err != nil {
		return err
	}

	return nil
}

func (as *DefaultAssetService) UnassignAsset(assetID uint, agentID uint) error {
	// Retrieve the asset by assetID
	// Update the asset user ID
	asset, err := as.GetAssetByID(assetID)
	if err != nil {
		return err
	}

	assetAssignment, err := as.AssetAssignmentDBModel.UnassignAsset(&asset.Assignment, agentID)
	if err != nil {
		return err
	}

	asset.Status = "unassigned"
	asset.Assignment = *assetAssignment

	// Save the updated asset
	err = as.AssetDBModel.UpdateAsset(asset)
	if err != nil {
		return err
	}

	return nil
}

// GetAssetByNumber retrieves an Asset by its asset number.
func (as *DefaultAssetService) GetAssetByNumber(assetNumber int) (*models.Assets, error) {
	asset, err := as.AssetDBModel.GetAssetByNumber(assetNumber)
	if err != nil {
		return nil, err
	}
	return asset, nil
}

// AssignAsset assigns an Asset to a user.
func (as *DefaultAssetService) AssignAssetToUser(assetID, userID uint, agentID uint) (*models.AssetAssignment, error) {
	// Retrieve the asset by assetID
	asset, err := as.AssetDBModel.GetAssetByID(assetID)
	if err != nil {
		return nil, err
	}

	// Create a new asset assignment record
	assetAssignment := &models.AssetAssignment{
		AssetID:          asset.ID,
		UserID:           userID,
		AssignedBy:       agentID,     // Assuming the same user assigns the asset
		AssignmentType:   "permanent", // Update as needed
		AssignmentStatus: "assigned",
		DueAt:            time.Now().AddDate(1, 0, 0), // Due date example
		CreatedAt:        time.Now(),
	}

	erro := as.AssetAssignmentDBModel.CreateAssetAssignment(assetAssignment)
	if erro != nil {
		return nil, err
	}

	// Update the asset's assignment and status
	asset.Assignment = *assetAssignment
	asset.Status = "assigned" // Update as needed

	if err := as.AssetDBModel.UpdateAsset(asset); err != nil {
		return nil, err
	}

	return assetAssignment, nil
}

// UnassignAsset unassigns an Asset from a user.
func (as *DefaultAssetService) UnassignAssetFromUser(assetID uint, agentID uint) error {
	// Retrieve the asset by assetID
	asset, err := as.AssetDBModel.GetAssetByID(assetID)
	if err != nil {
		return err
	}

	// Retrieve the asset assignment
	assetAssignment := &asset.Assignment

	// Update the asset assignment
	newAssetAssignment, err := as.AssetAssignmentDBModel.UnassignAsset(assetAssignment, agentID)
	if err != nil {
		return err
	}

	// Update the asset's assignment and status
	asset.Assignment = *newAssetAssignment
	asset.Status = "unassigned" // Update as needed

	if err := as.AssetDBModel.UpdateAsset(asset); err != nil {
		return err
	}

	return nil
}
