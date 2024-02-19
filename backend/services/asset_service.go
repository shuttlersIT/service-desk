// backend/services/advertisement_service.go

package services

import (
	"fmt"
	"log"
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
	log                    Logger
	// Add any dependencies or data needed for the service
}

// NewDefaultAssetService creates a new DefaultAssetService.
func NewDefaultAssetService(assetDBModel *models.AssetDBModel, assetAssignmentDBModel *models.AssetAssignmentDBModel, log Logger) *DefaultAssetService {
	return &DefaultAssetService{
		AssetDBModel:           assetDBModel,
		AssetAssignmentDBModel: assetAssignmentDBModel,
		log:                    log,
	}
}

// CreateAsset handles the creation of a new asset.
func (service *DefaultAssetService) CreateAsset(asset *models.Assets) error {
	err := service.DB.Transaction(func(tx *gorm.DB) error {
		if err := service.AssetDBModel.CreateAsset(asset); err != nil {
			service.log.Error(fmt.Sprintf("Error creating asset: %v", err))
			return err
		}
		service.log.Error(fmt.Sprintf("Asset created successfully: %v", asset))
		return nil
	})

	if err != nil {
		service.log.Error(fmt.Sprintf("Transaction failed: %v", err))
	}
	return err
}

// DeleteAsset handles the deletion of an asset by its ID.
func (service *DefaultAssetService) DeleteAsset(assetID uint) error {
	//var status bool
	err := service.DB.Transaction(func(tx *gorm.DB) error {
		if err := service.AssetDBModel.DeleteAsset(assetID); err != nil {
			service.log.Error(fmt.Sprintf("Error deleting asset: %v", err))
			return err
		}
		service.log.Error(fmt.Sprintf("Asset %d deleted successfully", assetID))
		return nil
	})

	if err != nil {
		service.log.Error(fmt.Sprintf("Transaction failed: %v", err))
		return err
	}
	return err
}

// UpdateAsset handles updates to an existing asset.
func (s *DefaultAssetService) UpdateAsset(assetID uint, updates *models.Assets) (*models.Assets, error) {
	var asset models.Assets
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", assetID).First(&asset).Error; err != nil {
			return fmt.Errorf("asset not found: %w", err)
		}

		if err := tx.Model(&asset).Updates(updates).Error; err != nil {
			return fmt.Errorf("updating asset: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &asset, nil
}

func (s *DefaultAssetService) ListAllAssets() ([]*models.Assets, error) {
	var assets []*models.Assets
	if err := s.DB.Find(&assets).Error; err != nil {
		return nil, fmt.Errorf("listing assets: %w", err)
	}
	return assets, nil
}

// AssignAssetToUser assigns an asset to a user, updating the AssetAssignment table.
func (s *DefaultAssetService) AssignAssetToUser(assetID, userID, assignedBy uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		assignment := models.AssetAssignment{
			AssetID:        assetID,
			AssignedTo:     userID,
			AssignmentDate: time.Now(),
			Status:         "assigned",
			AssignedBy:     assignedBy,
		}

		if err := tx.Create(&assignment).Error; err != nil {
			s.log.Error(fmt.Sprintf("Failed to assign asset %d to user %d: %v", assetID, userID, err))
			return err
		}

		s.log.Info(fmt.Sprintf("Asset %d successfully assigned to user %d", assetID, userID))
		return nil
	})
}

func (s *DefaultAssetService) ReturnAssetFromUser(assignmentID uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		var assignment models.AssetAssignment
		if err := tx.First(&assignment, assignmentID).Error; err != nil {
			return fmt.Errorf("finding assignment: %w", err)
		}

		t := time.Now()
		assignment.Status = "Returned"
		assignment.ReturnDate = &t

		if err := tx.Save(&assignment).Error; err != nil {
			return fmt.Errorf("updating assignment: %w", err)
		}
		return nil
	})
}

func (s *DefaultAssetService) GetAssetTypeByID(typeID uint) (*models.AssetType, error) {
	var assetType models.AssetType
	if err := s.DB.First(&assetType, typeID).Error; err != nil {
		return nil, fmt.Errorf("asset type not found: %w", err)
	}
	return &assetType, nil
}

func (s *DefaultAssetService) CreateAssetType(assetType *models.AssetType) (*models.AssetType, error) {
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(assetType).Error; err != nil {
			return fmt.Errorf("creating asset type: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return assetType, nil
}

func (s *DefaultAssetService) UpdateAssetType(typeID uint, updates map[string]interface{}) (*models.AssetType, error) {
	var assetType models.AssetType
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", typeID).First(&assetType).Error; err != nil {
			return fmt.Errorf("asset type not found: %w", err)
		}

		if err := tx.Model(&assetType).Updates(updates).Error; err != nil {
			return fmt.Errorf("updating asset type: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &assetType, nil
}

func (s *DefaultAssetService) DeleteAssetType(typeID uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.AssetType{}, typeID).Error; err != nil {
			return fmt.Errorf("deleting asset type: %w", err)
		}
		return nil
	})
}

// ListAssetTypes retrieves all asset types from the database.
func (s *DefaultAssetService) ListAssetTypes() ([]*models.AssetType, error) {
	var assetTypes []*models.AssetType
	if err := s.DB.Find(&assetTypes).Error; err != nil {
		s.log.Error(fmt.Sprintf("Failed to retrieve asset types: ", err))
		return nil, err
	}
	return assetTypes, nil
}

// GetAssetByID retrieves an asset by its ID, without involving transactions.
func (service *DefaultAssetService) GetAssetByID(assetID uint) (*models.Assets, error) {
	asset, err := service.AssetDBModel.GetAssetByID(assetID)
	if err != nil {
		service.log.Error(fmt.Sprintf("Error retrieving asset by ID %d: %v", assetID, err))
		return nil, err
	}
	return asset, nil
}

// GetAllAssets retrieves all assets, without involving transactions.
func (service *DefaultAssetService) GetAllAssets() ([]*models.Assets, error) {
	assets, err := service.AssetDBModel.GetAllAssets()
	if err != nil {
		service.log.Error(fmt.Sprintf("Error retrieving all assets: %v", err))
		return nil, err
	}
	return assets, nil
}

// GetAssetByNumber retrieves an asset by its number, without involving transactions.
func (service *DefaultAssetService) GetAssetByNumber(assetNumber int) (*models.Assets, error) {
	asset, err := service.AssetDBModel.GetAssetByNumber(assetNumber)
	if err != nil {
		service.log.Error(fmt.Sprintf("Error retrieving asset by number %d: %v", assetNumber, err))
		return nil, err
	}
	return asset, nil
}

// backend/services/asset_service.go

// UnassignAsset unassigns an asset from any user, demonstrating transactional integrity.
func (service *DefaultAssetService) UnassignAsset(assetID uint, agentID uint) error {
	return service.DB.Transaction(func(tx *gorm.DB) error {
		if err := service.AssetAssignmentDBModel.UnassignAsset(assetID, agentID); err != nil {
			log.Printf("Error unassigning asset %d: %v", assetID, err)
			return err
		}
		log.Printf("Asset %d successfully unassigned", assetID)
		return nil
	})
}

// UnassignAssetFromUser specifically unassigns an asset from a user, showing transactional integrity.
func (service *DefaultAssetService) UnassignAssetFromUser(assetID uint, userID uint, agentID uint) error {
	return service.DB.Transaction(func(tx *gorm.DB) error {
		if err := service.AssetAssignmentDBModel.AssetDBModel.UnassignAssetFromUser(assetID); err != nil {
			log.Printf("Error unassigning asset %d from user %d: %v", assetID, userID, err)
			return err
		}
		log.Printf("Asset %d successfully unassigned from user %d", assetID, userID)
		return nil
	})
}

// RetrieveAssetMaintenanceRecords fetches all maintenance records associated with an asset.
func (s *DefaultAssetService) RetrieveAssetMaintenanceRecords(assetID uint) ([]*models.AssetMaintenance, error) {
	var records []*models.AssetMaintenance
	if err := s.DB.Where("asset_id = ?", assetID).Find(&records).Error; err != nil {
		s.log.Error(fmt.Sprintf("Failed to retrieve maintenance records: ", err))
		return nil, err
	}
	return records, nil
}

// UpdateAssetCategory changes the category of an asset.
func (s *DefaultAssetService) UpdateAssetCategory(assetID uint, newCategoryID uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Assets{}).Where("id = ?", assetID).Update("category_id", newCategoryID).Error; err != nil {
			s.log.Error(fmt.Sprintf("Failed to update category for asset %d: %v", assetID, err))
			return err
		}
		return nil
	})
}

//////////////

// ListAssetsByUser lists all assets assigned to a specific user.
func (service *DefaultAssetService) ListAssetsByUser(userID uint) ([]*models.Assets, error) {
	assets, err := service.AssetDBModel.ListUserAssets(userID)
	if err != nil {
		log.Printf("Error listing assets for user %d: %v", userID, err)
		return nil, err
	}
	return assets, nil
}

// ListAvailableAssets lists all assets that are not currently assigned to any user.
func (service *DefaultAssetService) ListAvailableAssetsByEmptyAssignee() ([]*models.Assets, error) {
	assets, err := service.AssetDBModel.ListAvailableAssetsA()
	if err != nil {
		log.Printf("Error listing available assets: %v", err)
		return nil, err
	}
	return assets, nil
}

// ListAvailableAssets lists all assets that are not currently assigned to any user.
func (service *DefaultAssetService) ListAvailableAssetsByStatus() ([]*models.Assets, error) {
	assets, err := service.AssetDBModel.ListAvailableAssetsB()
	if err != nil {
		log.Printf("Error listing available assets: %v", err)
		return nil, err
	}
	return assets, nil
}

// ListAssetsByCategory lists all assets under a specific category.
func (s *DefaultAssetService) ListAssetsByCategory(categoryID uint) ([]*models.Assets, error) {
	var assets []*models.Assets
	if err := s.DB.Where("category_id = ?", categoryID).Find(&assets).Error; err != nil {
		s.log.Error(fmt.Sprintf("Failed to list assets by category %d: %v", categoryID, err))
		return nil, err
	}
	return assets, nil
}

// AssetCheckIn updates an asset's status to 'Available' when it is returned or checked in.
func (s *DefaultAssetService) AssetCheckIn(assetID uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Assets{}).Where("id = ?", assetID).Update("status", "Available").Error; err != nil {
			s.log.Error(fmt.Sprintf("Failed to check in asset %d: %v", assetID, err))
			return err
		}
		s.log.Info(fmt.Sprintf("Asset %d checked in and marked as available", assetID))
		return nil
	})
}

// AssetCheckOut updates an asset's status to 'In Use' when it is checked out.
func (s *DefaultAssetService) AssetCheckOut(assetID, userID uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		assignment := models.AssetAssignment{
			AssetID:        assetID,
			AssignedTo:     userID,
			AssignmentDate: time.Now(),
			Status:         "In Use",
		}
		if err := tx.Create(&assignment).Error; err != nil {
			s.log.Error(fmt.Sprintf("Failed to check out asset %d to user %d: %v", assetID, userID, err))
			return err
		}
		s.log.Info(fmt.Sprintf("Asset %d checked out to user %d", assetID, userID))
		return nil
	})
}

// ReassignAsset transfers an asset from one user to another.
func (service *DefaultAssetService) ReassignAsset(assetID, fromUserID, toUserID uint) error {
	return service.DB.Transaction(func(tx *gorm.DB) error {
		if err := service.AssetDBModel.ReassignAsset(assetID, fromUserID, toUserID); err != nil {
			log.Printf("Error reassigning asset %d from user %d to user %d: %v", assetID, fromUserID, toUserID, err)
			return err
		}
		log.Printf("Asset %d successfully reassigned from user %d to user %d", assetID, fromUserID, toUserID)
		return nil
	})
}

// ListAssetHistory provides a history of assignments for a given asset.
func (service *DefaultAssetService) ListAssetHistory(assetID uint) ([]*models.AssetAssignment, error) {
	history, err := service.AssetDBModel.GetAssetAssignmentHistory(assetID)
	if err != nil {
		log.Printf("Error listing history for asset %d: %v", assetID, err)
		return nil, err
	}
	return history, nil
}

// UpdateAssetStatus updates the status of an asset.
func (service *DefaultAssetService) UpdateAssetStatus(assetID uint, newStatus string) error {
	return service.DB.Transaction(func(tx *gorm.DB) error {
		if err := service.AssetDBModel.UpdateAssetDetails(assetID, newStatus); err != nil {
			log.Printf("Error updating status for asset %d to '%s': %v", assetID, newStatus, err)
			return err
		}
		log.Printf("Asset %d status updated to '%s'", assetID, newStatus)
		return nil
	})
}

// ValidateAssetOwnership checks if a given asset is assigned to a specified user.
func (service *DefaultAssetService) ValidateAssetOwnership(assetID, userID uint) (bool, error) {
	isOwner, err := service.AssetDBModel.ValidateAssetOwnership(assetID, userID)
	if err != nil {
		log.Printf("Error validating ownership of asset %d by user %d: %v", assetID, userID, err)
		return false, err
	}
	return isOwner, nil
}

// ListAssetsByType lists all assets of a specified type.
func (service *DefaultAssetService) ListAssetsByType(assetType string) ([]*models.Assets, error) {
	assets, err := service.AssetDBModel.ListAssetsByType(assetType)
	if err != nil {
		log.Printf("Error listing assets of type '%s': %v", assetType, err)
		return nil, err
	}
	return assets, nil
}

// LogAssetActivity records an action taken on an asset in the AssetLog table.
func (s *DefaultAssetService) LogAssetActivity(assetID, userID uint, action, details string) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		logEntry := models.AssetLog{
			AssetID:   assetID,
			UserID:    userID,
			Action:    action,
			Timestamp: time.Now(),
			Details:   details,
		}

		if err := tx.Create(&logEntry).Error; err != nil {
			s.log.Error(fmt.Sprintf("Failed to log activity for asset %d: %v", assetID, err))
			return err
		}

		s.log.Info(fmt.Sprintf("Activity for asset %d logged: %s", assetID, action))
		return nil
	})
}

// More methods can follow this pattern, adapting to the specific needs of each operation while ensuring
// transactional integrity, robust error handling, and detailed logging.

func (s *DefaultAssetService) RegisterAssetMaintenance(maintenance *models.AssetMaintenance) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(maintenance).Error; err != nil {
			return fmt.Errorf("registering asset maintenance: %w", err)
		}
		return nil
	})
}

func (s *DefaultAssetService) UpdateAssetMaintenance(maintenanceID uint, updates map[string]interface{}) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		var maintenance models.AssetMaintenance
		if err := tx.First(&maintenance, maintenanceID).Error; err != nil {
			return fmt.Errorf("asset maintenance not found: %w", err)
		}

		if err := tx.Model(&maintenance).Updates(updates).Error; err != nil {
			return fmt.Errorf("updating asset maintenance: %w", err)
		}
		return nil
	})
}

func (s *DefaultAssetService) ListAssetMaintenance(assetID uint) ([]*models.AssetMaintenance, error) {
	var maintenances []*models.AssetMaintenance
	if err := s.DB.Where("asset_id = ?", assetID).Find(&maintenances).Error; err != nil {
		return nil, fmt.Errorf("listing asset maintenance: %w", err)
	}
	return maintenances, nil
}

// ScheduleAssetMaintenance schedules maintenance for an asset.
func (s *DefaultAssetService) ScheduleAssetMaintenance(assetID uint, date *time.Time, description string) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		maintenance := models.AssetMaintenance{
			AssetID:       assetID,
			ScheduledDate: date,
			Description:   description,
			Status:        "scheduled",
		}

		if err := tx.Create(&maintenance).Error; err != nil {
			s.log.Error(fmt.Sprintf("Failed to schedule maintenance for asset %d: %v", assetID, err))
			return err
		}

		s.log.Info(fmt.Sprintf("Maintenance for asset %d scheduled on %s", assetID, date))
		return nil
	})
}

func (service *DefaultAssetService) ListAssetMaintenanceSchedules(assetID uint) ([]*models.AssetMaintenance, error) {
	var maintenanceSchedules []*models.AssetMaintenance
	if err := service.DB.Where("asset_id = ?", assetID).Find(&maintenanceSchedules).Error; err != nil {
		service.log.Error(fmt.Sprintf("Error listing maintenance schedules for asset ID %d: %v", assetID, err))
		return nil, fmt.Errorf("listing maintenance schedules for asset ID %d: %w", assetID, err)
	}
	return maintenanceSchedules, nil
}

// ScheduleAssetRepair schedules repair for an asset, updating its status accordingly.
func (s *DefaultAssetService) ScheduleAssetRepair(assetID uint, repairDetails string, scheduledDate time.Time) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		repairLog := models.AssetRepairLog{
			AssetID:     assetID,
			StartDate:   scheduledDate,
			Description: repairDetails,
		}
		if err := tx.Create(&repairLog).Error; err != nil {
			s.log.Error(fmt.Sprintf("Failed to schedule repair for asset %d: %v", assetID, err))
			return err
		}
		s.log.Info(fmt.Sprintf("Repair scheduled for asset %d on %s", assetID, scheduledDate.Format("2006-01-02")))
		return nil
	})
}

func (service *DefaultAssetService) ScheduleAssetInspection(assetID uint, inspectionDetails *models.AssetInspectionRecord) error {
	return service.DB.Transaction(func(tx *gorm.DB) error {
		inspection := models.AssetInspectionRecord{
			AssetID:        assetID,
			InspectionDate: inspectionDetails.InspectionDate,
			InspectedBy:    inspectionDetails.InspectedBy,
			Notes:          inspectionDetails.Notes,
		}
		if err := tx.Create(&inspection).Error; err != nil {
			service.log.Error(fmt.Sprintf("Error scheduling inspection for asset ID %d: %v", assetID, err))
			return fmt.Errorf("scheduling inspection for asset ID %d: %w", assetID, err)
		}
		service.log.Info(fmt.Sprintf("Inspection for asset ID %d scheduled", assetID))
		return nil
	})
}

// CompleteAssetRepair marks a scheduled repair as completed.
func (s *DefaultAssetService) CompleteAssetRepair(repairID uint, completionDetails models.CompletionDetails) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.AssetRepairLog{}).Where("id = ?", repairID).Updates(models.AssetRepairLog{Status: "Completed", CompletionDetails: completionDetails}).Error; err != nil {
			s.log.Error(fmt.Sprintf("Failed to complete repair %d: %v", repairID, err))
			return err
		}
		s.log.Info(fmt.Sprintf("Repair %d completed", repairID))
		return nil
	})
}

func (s *DefaultAssetService) ScheduleAssetInspection2(inspection *models.AssetInspectionRecord) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(inspection).Error; err != nil {
			return fmt.Errorf("scheduling asset inspection: %w", err)
		}
		return nil
	})
}

func (s *DefaultAssetService) CompleteAssetInspection(inspectionID uint, updates map[string]interface{}) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		var inspection models.AssetInspectionRecord
		if err := tx.First(&inspection, inspectionID).Error; err != nil {
			return fmt.Errorf("asset inspection not found: %w", err)
		}

		if err := tx.Model(&inspection).Updates(updates).Error; err != nil {
			return fmt.Errorf("completing asset inspection: %w", err)
		}
		return nil
	})
}

func (s *DefaultAssetService) GetAssetInspections(assetID uint) ([]*models.AssetInspectionRecord, error) {
	var inspections []*models.AssetInspectionRecord
	if err := s.DB.Where("asset_id = ?", assetID).Find(&inspections).Error; err != nil {
		return nil, fmt.Errorf("retrieving asset inspections: %w", err)
	}
	return inspections, nil
}

func (service *DefaultAssetService) CompleteAssetDecommissioning(assetID uint, decommissionDetails *models.AssetDecommission) error {
	return service.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Assets{}).Where("id = ?", assetID).Updates(map[string]interface{}{
			"status":              "Decommissioned",
			"decommission_date":   decommissionDetails.DecommissionDate,
			"decommission_reason": decommissionDetails.Reason,
		}).Error; err != nil {
			service.log.Error(fmt.Sprintf("Error completing decommissioning for asset ID %d: %v", assetID, err))
			return fmt.Errorf("completing decommissioning for asset ID %d: %w", assetID, err)
		}
		service.log.Info(fmt.Sprintf("Asset ID %d decommissioned", assetID))
		return nil
	})
}

// DecommissionAsset initiates the decommissioning process for an asset.
func (s *DefaultAssetService) DecommissionAsset(assetID uint, reason string) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		decommission := models.AssetDecommission{
			AssetID:          assetID,
			DecommissionDate: time.Now(),
			Reason:           reason,
			Status:           "in_progress",
		}

		if err := tx.Create(&decommission).Error; err != nil {
			s.log.Error(fmt.Sprintf("Failed to decommission asset %d: %v", assetID, err))
			return err
		}

		s.log.Info(fmt.Sprintf("Asset %d decommissioning initiated", assetID))
		return nil
	})
}

func (s *DefaultAssetService) ReactivateAsset(assetID uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		var asset models.Assets
		if err := tx.Where("id = ?", assetID).First(&asset).Error; err != nil {
			return fmt.Errorf("asset not found: %w", err)
		}

		asset.Status = "Active"
		if err := tx.Save(&asset).Error; err != nil {
			return fmt.Errorf("reactivating asset: %w", err)
		}
		return nil
	})
}

func (s *DefaultAssetService) ListAssetsByStatus(status string) ([]*models.Assets, error) {
	var assets []*models.Assets
	if err := s.DB.Where("status = ?", status).Find(&assets).Error; err != nil {
		return nil, fmt.Errorf("listing assets by status: %w", err)
	}
	return assets, nil
}

// ListAssetsByLocation retrieves all assets located in a specific location.
func (s *DefaultAssetService) ListAssetsByLocation(location string) ([]*models.Assets, error) {
	var assets []*models.Assets
	if err := s.DB.Where("location = ?", location).Find(&assets).Error; err != nil {
		s.log.Error(fmt.Sprintf("Failed to retrieve assets by location: ", err))
		return nil, err
	}
	return assets, nil
}

func (s *DefaultAssetService) RecordAssetPerformance(assetID uint, performanceData *models.AssetPerformanceLog) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(performanceData).Error; err != nil {
			return fmt.Errorf("recording asset performance for asset ID %d: %w", assetID, err)
		}
		return nil
	})
}

// AssetDecommission marks an asset as decommissioned and out of service.
func (s *DefaultAssetService) AssetDecommission(assetID uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Assets{}).Where("id = ?", assetID).Update("status", "Decommissioned").Error; err != nil {
			s.log.Error(fmt.Sprintf("Failed to decommission asset %d: %v", assetID, err))
			return err
		}
		return nil
	})
}

// ListAssetsForAudit selects assets that meet certain criteria for auditing.
func (s *DefaultAssetService) ListAssetsForAudit(criteria map[string]interface{}) ([]*models.Assets, error) {
	var assets []*models.Assets
	if err := s.DB.Where(criteria).Find(&assets).Error; err != nil {
		s.log.Error(fmt.Sprintf("Failed to list assets for audit: ", err))
		return nil, err
	}
	return assets, nil
}

func (s *DefaultAssetService) CalculateAssetDepreciation(assetID uint, years int) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		var asset models.Assets
		if err := tx.Where("id = ?", assetID).First(&asset).Error; err != nil {
			return fmt.Errorf("asset not found for depreciation calculation: %w", err)
		}

		// Placeholder: Depreciation calculation logic goes here.
		// For example, straight-line depreciation over the 'years' period.
		newDepreciationValue := asset.PurchasePrice - (asset.PurchasePrice / float64(asset.UsefulLife) * float64(years))

		if err := tx.Model(&asset).Update("depreciation_value", newDepreciationValue).Error; err != nil {
			return fmt.Errorf("updating depreciation value for asset ID %d: %w", assetID, err)
		}
		return nil
	})
}

func (s *DefaultAssetService) ListAssetByCondition(condition string) ([]*models.Assets, error) {
	var assets []*models.Assets
	if err := s.DB.Where("condition = ?", condition).Find(&assets).Error; err != nil {
		return nil, fmt.Errorf("listing assets by condition '%s': %w", condition, err)
	}
	return assets, nil
}

// UpdateAssetCondition updates the condition of an asset, such as 'New', 'Good', 'Needs Repair'.
func (s *DefaultAssetService) UpdateAssetCondition(assetID uint, newCondition string) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Assets{}).Where("id = ?", assetID).Update("condition", newCondition).Error; err != nil {
			s.log.Error(fmt.Sprintf("Failed to update condition for asset %d: %v", assetID, err))
			return err
		}
		s.log.Info(fmt.Sprintf("Asset %d condition updated to '%s'", assetID, newCondition))
		return nil
	})
}

func (s *DefaultAssetService) UpdateAssetLocation(assetID uint, newLocation *models.AssetLocation) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Assets{}).Where("id = ?", assetID).Update("location", newLocation).Error; err != nil {
			return fmt.Errorf("updating location for asset ID %d: %w", assetID, err)
		}
		return nil
	})
}

// PerformAssetHealthCheck simulates the process of performing a health check on an asset.
func (s *DefaultAssetService) PerformAssetHealthCheck(assetID uint) (*models.AssetHealthReport, error) {
	// Simulate fetching health check data and determine the health status
	healthStatus := "Good" // Placeholder value
	notes := "Asset is functioning within expected parameters."

	report := models.AssetHealthReport{
		AssetID:      assetID,
		ReportDate:   time.Now(),
		HealthStatus: healthStatus,
		Notes:        notes,
	}

	if err := s.DB.Create(&report).Error; err != nil {
		s.log.Error(fmt.Sprintf("Error performing health check for asset %d: %v", assetID, err))
		return nil, err
	}

	return &report, nil
}

func (s *DefaultAssetService) PerformAssetHealthCheck2(assetID uint) (*models.AssetHealthReport, error) {
	// Example placeholder logic. In a real scenario, this would involve more complex operations,
	// potentially including integration with external systems for health monitoring.
	var report models.AssetHealthReport
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		// Simulate fetching health check data
		report.AssetID = assetID
		report.Report = "Asset health check performed successfully. No issues found."
		report.CreatedAt = time.Now()

		if err := tx.Create(&report).Error; err != nil {
			s.log.Error(fmt.Sprintf("Error performing health check for asset ID %d: %v", assetID, err))
			return fmt.Errorf("performing health check for asset ID %d: %w", assetID, err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &report, nil
}

func (s *DefaultAssetService) ListAssetsNeedingMaintenance() ([]*models.Assets, error) {
	var assets []*models.Assets
	// Placeholder: Logic to list assets based on maintenance schedule and last maintenance date.
	return assets, nil
}

// ScheduleBulkAssetMaintenance schedules maintenance for multiple assets.
func (s *DefaultAssetService) ScheduleBulkAssetMaintenance(assetIDs []uint, maintenanceDetails *models.AssetMaintenanceSchedule) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		for _, assetID := range assetIDs {
			maintenanceRecord := models.AssetMaintenanceSchedule{
				AssetID:       assetID,
				ScheduledDate: maintenanceDetails.ScheduledDate,
				Description:   maintenanceDetails.Description,
				Status:        "Scheduled",
			}
			if err := tx.Create(&maintenanceRecord).Error; err != nil {
				return fmt.Errorf("scheduling maintenance for asset ID %d: %w", assetID, err)
			}
		}
		return nil
	})
}

func (s *DefaultAssetService) DecommissionAssets(assetIDs []uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// Batch update status to 'Decommissioned' for all assetIDs.
		if err := tx.Model(&models.Assets{}).Where("id IN ?", assetIDs).Update("status", "Decommissioned").Error; err != nil {
			return fmt.Errorf("decommissioning assets: %w", err)
		}
		return nil
	})
}

func (s *DefaultAssetService) ReassignAssets(fromUserID, toUserID uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// Update all assets assigned to fromUserID to now be assigned to toUserID.
		if err := tx.Model(&models.AssetAssignment{}).Where("user_id = ?", fromUserID).Update("user_id", toUserID).Error; err != nil {
			return fmt.Errorf("reassigning assets from user ID %d to user ID %d: %w", fromUserID, toUserID, err)
		}
		return nil
	})
}

func (s *DefaultAssetService) AddLifecycleEvent(assetID uint, event *models.LifecycleEvent) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.LifecycleEvent{}).Create(event).Error; err != nil {
			return fmt.Errorf("adding lifecycle event for asset ID %d: %w", assetID, err)
		}
		return nil
	})
}

func (s *DefaultAssetService) ListLifecycleEvents(assetID uint) ([]*models.LifecycleEvent, error) {
	var events []*models.LifecycleEvent
	if err := s.DB.Where("asset_id = ?", assetID).Find(&events).Error; err != nil {
		return nil, fmt.Errorf("listing lifecycle events for asset ID %d: %w", assetID, err)
	}
	return events, nil
}

func (s *DefaultAssetService) ScheduleMaintenance(assetID uint, maintenance *models.AssetMaintenanceSchedule) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(maintenance).Error; err != nil {
			return fmt.Errorf("scheduling maintenance for asset ID %d: %w", assetID, err)
		}
		return nil
	})
}

func (s *DefaultAssetService) ScheduleMaintenance2(assetID uint, maintenance *models.AssetMaintenanceSchedule) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(maintenance).Error; err != nil {
			return fmt.Errorf("scheduling maintenance for asset ID %d: %w", assetID, err)
		}
		return nil
	})
}

// CompleteAssetMaintenance marks a maintenance record as completed.
func (s *DefaultAssetService) CompleteAssetMaintenance(maintenanceID uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		if err := tx.Model(&models.AssetMaintenance{}).Where("id = ?", maintenanceID).
			Updates(map[string]interface{}{"status": "completed", "maintenance_date": now}).Error; err != nil {
			s.log.Error(fmt.Sprintf("Failed to complete maintenance for record %d: %v", maintenanceID, err))
			return err
		}

		s.log.Info(fmt.Sprintf("Maintenance record %d marked as completed", maintenanceID))
		return nil
	})
}

func (s *DefaultAssetService) RecordAssetIssue(assetID uint, issue *models.AssetIssue) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(issue).Error; err != nil {
			return fmt.Errorf("recording issue for asset ID %d: %w", assetID, err)
		}
		return nil
	})
}

func (s *DefaultAssetService) ResolveAssetIssue(issueID uint, resolutionDetails string) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.AssetIssue{}).Where("id = ?", issueID).Updates(map[string]interface{}{"status": "Resolved", "resolution_details": resolutionDetails}).Error; err != nil {
			return fmt.Errorf("resolving issue ID %d: %w", issueID, err)
		}
		return nil
	})
}

func (service *DefaultAssetService) AnalyzeAssetPerformance(assetID uint, startDate, endDate time.Time) (*models.AssetPerformanceAnalysis, error) {
	var analysis models.AssetPerformanceAnalysis
	// Placeholder logic for performance analysis; actual implementation would depend on available data and business requirements.
	analysis.AssetID = assetID
	analysis.AnalysisStartDate = startDate
	analysis.AnalysisEndDate = endDate
	analysis.Summary = "Asset performance within expected parameters."

	service.log.Info(fmt.Sprintf("Performance analysis for asset ID %d completed", assetID))
	return &analysis, nil
}

// GenerateAssetUtilizationReport compiles utilization data for an asset over a specified period.
func (s *DefaultAssetService) GenerateAssetUtilizationReport(assetID uint, startDate, endDate time.Time) (*models.AssetUtilizationReport, error) {
	var report models.AssetUtilizationReport
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		// Placeholder for actual utilization calculation logic
		utilization := 75.0 // Simulated utilization percentage
		report = models.AssetUtilizationReport{
			AssetID:         assetID,
			Utilization:     utilization,
			ReportingPeriod: fmt.Sprintf("%s to %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")),
		}

		s.log.Info(fmt.Sprintf("Utilization report generated for asset %d: %f%%", assetID, utilization))
		return nil
	})

	if err != nil {
		s.log.Error(fmt.Sprintf("Failed to generate utilization report for asset %d: %v", assetID, err))
		return nil, err
	}

	return &report, nil
}

func (service *DefaultAssetService) GenerateAssetPerformanceReport(assetID uint) (*models.AssetPerformanceAnalysis, error) {
	// Placeholder for performance report generation logic
	var report models.AssetPerformanceAnalysis
	report.AssetID = assetID
	// Assume gathering and analyzing performance data from various sources
	report.Summary = "Asset performance meets expected benchmarks."

	service.log.Info(fmt.Sprintf("Performance report generated for asset ID %d", assetID))
	return &report, nil
}

func (s *DefaultAssetService) ListAssetsForDecommission(thresholdAge int) ([]*models.Assets, error) {
	var assets []*models.Assets
	if err := s.DB.Where("DATEDIFF(CURRENT_DATE, purchase_date) > ?", thresholdAge*365).Find(&assets).Error; err != nil {
		return nil, fmt.Errorf("listing assets for decommission based on age threshold %d years: %w", thresholdAge, err)
	}
	return assets, nil
}

func (s *DefaultAssetService) ListOverdueMaintenances() ([]*models.AssetMaintenanceSchedule, error) {
	var schedules []*models.AssetMaintenanceSchedule
	if err := s.DB.Where("scheduled_date < ? AND status != 'Completed'", time.Now()).Find(&schedules).Error; err != nil {
		return nil, fmt.Errorf("listing overdue maintenances: %w", err)
	}
	return schedules, nil
}

// Add more methods as needed based on the application's requirements and the data model.
func (service *DefaultAssetService) RecordLifecycleEvent(event models.AssetLifecycleEvent) error {
	return service.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&event).Error; err != nil {
			service.log.Error(fmt.Sprintf("Error recording lifecycle event for asset ID %d: %v", event.AssetID, err))
			return fmt.Errorf("recording lifecycle event for asset ID %d: %w", event.AssetID, err)
		}
		service.log.Info(fmt.Sprintf("Lifecycle event recorded for asset ID %d", event.AssetID))
		return nil
	})
}

func (service *DefaultAssetService) LogAssetRepair(repair models.AssetRepairLog) error {
	return service.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&repair).Error; err != nil {
			service.log.Error(fmt.Sprintf("Error logging repair for asset ID %d: %v", repair.AssetID, err))
			return fmt.Errorf("logging repair for asset ID %d: %w", repair.AssetID, err)
		}
		service.log.Info(fmt.Sprintf("Repair logged for asset ID %d", repair.AssetID))
		return nil
	})
}

func (service *DefaultAssetService) CalculateAssetDepreciation2(assetID uint) error {
	return service.DB.Transaction(func(tx *gorm.DB) error {
		var asset models.Assets
		if err := tx.First(&asset, assetID).Error; err != nil {
			return fmt.Errorf("asset not found for ID %d: %w", assetID, err)
		}

		// Assuming straight-line depreciation for simplicity
		annualDepreciation := asset.PurchaseCost / float64(asset.UsefulLife)
		accumulatedDepreciation := annualDepreciation * float64(time.Now().Year()-asset.PurchaseDate.Year())
		newBookValue := asset.PurchaseCost - accumulatedDepreciation

		if err := tx.Model(&asset).Update("book_value", newBookValue).Error; err != nil {
			return fmt.Errorf("updating book value for asset %d: %w", assetID, err)
		}
		return nil
	})
}

func (service *DefaultAssetService) PerformComplianceAudit() ([]models.ComplianceAuditLog, error) {
	// Placeholder for actual compliance audit logic
	// This method would involve querying the database for assets not meeting compliance criteria and returning the results.
	var auditResults []models.ComplianceAuditLog
	// Example: Find all assets overdue for maintenance
	if err := service.DB.Where("maintenance_due_date < ?", time.Now()).Find(&auditResults).Error; err != nil {
		service.log.Error(fmt.Sprintf("Error performing compliance audit: %v", err))
		return nil, fmt.Errorf("performing compliance audit: %w", err)
	}
	return auditResults, nil
}
