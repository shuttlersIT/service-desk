package models

import (
	"time"

	"gorm.io/gorm"
)

type ServiceRequest struct {
	gorm.Model
	Title         string     `gorm:"size:255;not null" json:"title"`                       // Title of the service request
	Description   string     `gorm:"type:text" json:"description"`                         // Detailed description of the request
	UserID        uint       `gorm:"index;not null" json:"user_id"`                        // User who made the request
	Status        string     `gorm:"size:100;not null" json:"status"`                      // Status of the request, e.g., "open", "closed"
	CategoryID    uint       `gorm:"index;not null" json:"category_id"`                    // Category of the service request
	SubCategoryID uint       `gorm:"index" json:"subcategory_id,omitempty"`                // Optional sub-category
	LocationID    uint       `gorm:"index;not null" json:"location_id"`                    // Location of the service request
	CompletedAt   *time.Time `json:"completed_at,omitempty"`                               // Time when the service request was completed
	AssignedTo    *uint      `gorm:"index;type:int unsigned" json:"assignee_id,omitempty"` // Optional assignee of the service request
	ServiceType   string     `gorm:"size:100" json:"service_type"`                         // Type of the service requested
	Priority      string     `gorm:"size:50;not null" json:"priority"`                     // Priority of the service request
	TicketID      *uint      `gorm:"index" json:"ticket_id,omitempty"`
	// Removed embedded Location struct to normalize data structure and reference by ID instead
}

type Location struct {
	gorm.Model
	LocationName string `gorm:"size:255;not null" json:"location_name"` // Name of the location
}

func (Location) TableName() string {
	return "locations"
}

type ServiceRequestComment struct {
	gorm.Model       `gorm:"primaryKey" json:"id"`
	ServiceRequestID uint   `gorm:"index;not null" json:"service_request_id"` // Foreign key linking to the service request
	Comment          string `gorm:"type:text;not null" json:"comment"`        // Text of the comment
}

func (ServiceRequestComment) TableName() string {
	return "service_request_comments"
}

type ServiceCatalogItem struct {
	gorm.Model
	Name        string `json:"name" gorm:"type:varchar(255);not null"`
	Description string `json:"description" gorm:"type:text"`
	ServiceTime int    `json:"service_time"` // Expected time to deliver the service in minutes
	Category    string `json:"category" gorm:"type:varchar(100);not null"`
}

func (ServiceCatalogItem) TableName() string {
	return "service_catalog_item"
}

type ProblemManagementRecord struct {
	gorm.Model
	Title        string     `json:"title" gorm:"type:varchar(255);not null"`
	Description  string     `json:"description" gorm:"type:text"`
	IdentifiedAt time.Time  `json:"identified_at"`
	ResolvedAt   *time.Time `json:"resolved_at,omitempty"`
	RootCause    string     `json:"root_cause" gorm:"type:text"`
	Resolution   string     `json:"resolution" gorm:"type:text"`
	Status       string     `json:"status" gorm:"type:varchar(100);not null"` // E.g., "Open", "Under Investigation", "Resolved"
}

func (ProblemManagementRecord) TableName() string {
	return "problem_management_record"
}

type ServiceRequestHistoryEntry struct {
	gorm.Model       `gorm:"primaryKey" json:"id"`
	ServiceRequestID uint   `gorm:"index;not null" json:"service_request_id"` // Link to the service request
	Status           string `gorm:"size:100;not null" json:"status"`          // New status after the change
}

func (ServiceRequestHistoryEntry) TableName() string {
	return "service_request_history_entries"
}

type ChangeRequest struct {
	gorm.Model
	RequesterID uint       `json:"requester_id" gorm:"index;not null"`
	Title       string     `json:"title" gorm:"type:varchar(255);not null"`
	Description string     `json:"description" gorm:"type:text"`
	RequestedAt time.Time  `json:"requested_at"`
	Status      string     `json:"status" gorm:"type:varchar(100);not null"` // E.g., "New", "Approved", "Implemented", "Rejected"
	Approval    string     `json:"approval" gorm:"type:varchar(100)"`        // E.g., "Pending", "Approved", "Rejected"
	ImplementBy *time.Time `json:"implement_by,omitempty"`
}

type ServiceRequestStorage interface {
	CreateServiceRequest(request *ServiceRequest) error
	UpdateServiceRequest(request *ServiceRequest) error
	GetAllServiceRequests() ([]*ServiceRequest, error)
	GetServiceRequestByID(requestID uint) (*ServiceRequest, error)
	GetUserServiceRequests(userID uint) ([]*ServiceRequest, error)
	CloseServiceRequest(requestID uint) error
	ReopenServiceRequest(requestID uint) error
	GetServiceRequestHistory(requestID uint) ([]*ServiceRequestHistoryEntry, error)
	AddCommentToServiceRequest(requestID uint, comment string) error
	GetServiceRequestComments(requestID uint) ([]*ServiceRequestComment, error)
	GetOpenServiceRequests() ([]*ServiceRequest, error)
	GetClosedServiceRequests() ([]*ServiceRequest, error)

	// Additional methods
	GetUserServiceRequestCountByCategory(userID uint) (map[uint]int, error)
	GetUserServiceRequestCountBySubCategory(userID uint) (map[uint]int, error)
	GetServiceRequestsByCategoryAndUser(categoryID, userID uint) ([]*ServiceRequest, error)
	GetServiceRequestsBySubCategoryAndUser(subCategoryID, userID uint) ([]*ServiceRequest, error)
	GetServiceRequestCountByCategoryAndUser(userID uint) (map[uint]int, error)
	GetServiceRequestCountBySubCategoryAndUser(userID uint) (map[uint]int, error)
	GetServiceRequestCountByCategoryAndStatus(status string) (map[uint]int, error)
	GetServiceRequestCountBySubCategoryAndStatus(status string) (map[uint]int, error)
	GetUserAssignedServiceRequests(userID uint) ([]*ServiceRequest, error)
	AssignServiceRequestToUser(requestID, userID uint) error
	GetServiceRequestsByPriority(priority string) ([]*ServiceRequest, error)
	GetServiceRequestCountByPriority() (map[string]int, error)
	GetServiceRequestsByUserAndPriority(userID, priority string) ([]*ServiceRequest, error)
	GetServiceRequestCountByUserAndPriority(userID, priority string) (map[string]int, error)
	GetServiceRequestCountByCategoryAndPriority(categoryID, priority string) (map[uint]int, error)
	GetServiceRequestsByLocation(locationID uint) ([]*ServiceRequest, error)
	GetServiceRequestCountByLocation() (map[uint]int, error)
	GetServiceRequestsByUserAndLocation(userID, locationID uint) ([]*ServiceRequest, error)
	GetServiceRequestCountByUserAndLocation(userID uint) (map[uint]int, error)
	GetServiceRequestsByCategoryAndLocation(categoryID, locationID uint) ([]*ServiceRequest, error)
	GetServiceRequestCountByCategoryAndLocation(categoryID uint) (map[uint]int, error)
	GetServiceRequestsByPriorityAndLocation(priority string, locationID uint) ([]*ServiceRequest, error)
	GetServiceRequestCountByPriorityAndLocation(locationID uint) (map[string]int, error)

	// DeleteServiceRequest deletes a service request from the database by its ID.
	DeleteServiceRequest(id uint) error

	// GetServiceRequestsByUser retrieves all service requests for a specific user.
	GetServiceRequestsByUser(userID uint) ([]*ServiceRequest, error)

	// GetServiceRequestsByCategory retrieves all service requests for a specific category.
	GetServiceRequestsByCategory(categoryID uint) ([]*ServiceRequest, error)

	// GetServiceRequestsBySubCategory retrieves all service requests for a specific sub-category.
	GetServiceRequestsBySubCategory(subCategoryID uint) ([]*ServiceRequest, error)

	// GetServiceRequestsByStatus retrieves all service requests with a specific status.
	GetServiceRequestsByStatus(status string) ([]*ServiceRequest, error)

	// GetServiceRequestCountByStatus returns the count of service requests for a given status.
	GetServiceRequestCountByStatus(status string) (int, error)
}

// ServiceRequestDBModel is the database model for service requests.
type ServiceRequestDBModel struct {
	DB             *gorm.DB
	log            Logger
	EventPublisher *EventPublisherImpl
}

// NewServiceRequestDBModel creates a new ServiceRequestDBModel with the provided GORM database instance.
func NewServiceRequestDBModel(db *gorm.DB, log Logger, eventPublisher *EventPublisherImpl) *ServiceRequestDBModel {
	return &ServiceRequestDBModel{
		DB:             db,
		log:            log,
		EventPublisher: eventPublisher,
	}
}

// CreateServiceRequest initializes a new service request.
func (db *ServiceRequestDBModel) CreateServiceRequest2(request *ServiceRequest) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(request).Error; err != nil {
			return err
		}
		return nil
	})
}

// GetAllServiceRequests retrieves all service requests from the database.
func (srm *ServiceRequestDBModel) GetAllServiceRequests() ([]*ServiceRequest, error) {
	var serviceRequests []*ServiceRequest
	err := srm.DB.Find(&serviceRequests).Error
	return serviceRequests, err
}

// GetServiceRequestCountBySubCategoryAndStatus retrieves the count of service requests grouped by sub-category and status.
func (srm *ServiceRequestDBModel) GetServiceRequestCountBySubCategoryAndStatus(status string) (map[uint]int, error) {
	var counts []struct {
		SubCategoryID uint
		Count         int
	}

	err := srm.DB.Model(&ServiceRequest{}).
		Joins("JOIN service_request_subcategories ON service_requests.subcategory_id = service_request_subcategories.id").
		Where("service_requests.status = ?", status).
		Group("service_request_subcategories.id").
		Select("service_request_subcategories.id as subcategory_id, COUNT(service_requests.id) as count").
		Scan(&counts).Error

	if err != nil {
		return nil, err
	}

	result := make(map[uint]int)
	for _, count := range counts {
		result[count.SubCategoryID] = count.Count
	}

	return result, nil
}

// GetServiceRequestCountByCategoryAndPriority returns the count of service requests for a given category and grouped by priority.
func (srm *ServiceRequestDBModel) GetServiceRequestCountByCategoryAndPriority(categoryID uint) (map[string]int, error) {
	var counts []struct {
		Priority string
		Count    int
	}

	err := srm.DB.Model(&ServiceRequest{}).
		Where("category_id = ?", categoryID).
		Group("priority").
		Select("priority, COUNT(id) as count").
		Scan(&counts).Error

	if err != nil {
		return nil, err
	}

	result := make(map[string]int)
	for _, count := range counts {
		result[count.Priority] = count.Count
	}

	return result, nil
}

// GetUserAssignedServiceRequests returns service requests assigned to a user.
func (srm *ServiceRequestDBModel) GetUserAssignedServiceRequests(userID uint) ([]*ServiceRequest, error) {
	var requests []*ServiceRequest

	err := srm.DB.Where("assigned_to = ?", userID).Find(&requests).Error
	if err != nil {
		return nil, err
	}

	return requests, nil
}

// AssignServiceRequestToUser assigns a service request to a user.
func (srm *ServiceRequestDBModel) AssignServiceRequestToUser(requestID, userID uint) error {
	return srm.DB.Model(&ServiceRequest{}).
		Where("id = ?", requestID).
		Update("assigned_to", userID).
		Error
}

// GetServiceRequestsByPriority returns service requests filtered by priority.
func (srm *ServiceRequestDBModel) GetServiceRequestsByPriority(priority string) ([]*ServiceRequest, error) {
	var requests []*ServiceRequest

	err := srm.DB.Where("priority = ?", priority).Find(&requests).Error
	if err != nil {
		return nil, err
	}

	return requests, nil
}

// GetServiceRequestCountByPriority returns the count of service requests grouped by priority.
func (srm *ServiceRequestDBModel) GetServiceRequestCountByPriority() (map[string]int, error) {
	var counts []struct {
		Priority string
		Count    int
	}

	err := srm.DB.Model(&ServiceRequest{}).
		Group("priority").
		Select("priority, COUNT(id) as count").
		Scan(&counts).Error

	if err != nil {
		return nil, err
	}

	result := make(map[string]int)
	for _, count := range counts {
		result[count.Priority] = count.Count
	}

	return result, nil
}

// GetServiceRequestsByUserAndPriority returns service requests for a user filtered by priority.
func (srm *ServiceRequestDBModel) GetServiceRequestsByUserAndPriority(userID uint, priority string) ([]*ServiceRequest, error) {
	var requests []*ServiceRequest

	err := srm.DB.Where("user_id = ? AND priority = ?", userID, priority).Find(&requests).Error
	if err != nil {
		return nil, err
	}

	return requests, nil
}

// GetServiceRequestCountByUserAndPriority returns the count of service requests for a user grouped by priority.
func (srm *ServiceRequestDBModel) GetServiceRequestCountByUserAndPriority(userID uint) (map[string]int, error) {
	var counts []struct {
		Priority string
		Count    int
	}

	err := srm.DB.Model(&ServiceRequest{}).
		Where("user_id = ?", userID).
		Group("priority").
		Select("priority, COUNT(id) as count").
		Scan(&counts).Error

	if err != nil {
		return nil, err
	}

	result := make(map[string]int)
	for _, count := range counts {
		result[count.Priority] = count.Count
	}

	return result, nil
}

// GetServiceRequestsByCategoryAndPriority returns service requests filtered by category and priority.
func (srm *ServiceRequestDBModel) GetServiceRequestsByCategoryAndPriority(categoryID uint, priority string) ([]*ServiceRequest, error) {
	var requests []*ServiceRequest

	err := srm.DB.Where("category_id = ? AND priority = ?", categoryID, priority).Find(&requests).Error
	if err != nil {
		return nil, err
	}

	return requests, nil
}

// GetServiceRequestsByLocation returns service requests filtered by location.
func (srm *ServiceRequestDBModel) GetServiceRequestsByLocation(locationID uint) ([]*ServiceRequest, error) {
	var requests []*ServiceRequest

	err := srm.DB.Where("location_id = ?", locationID).Find(&requests).Error
	if err != nil {
		return nil, err
	}

	return requests, nil
}

// GetServiceRequestCountByLocation returns the count of service requests grouped by location.
func (srm *ServiceRequestDBModel) GetServiceRequestCountByLocation() (map[uint]int, error) {
	var counts []struct {
		LocationID uint
		Count      int
	}

	err := srm.DB.Model(&ServiceRequest{}).
		Group("location_id").
		Select("location_id, COUNT(id) as count").
		Scan(&counts).Error

	if err != nil {
		return nil, err
	}

	result := make(map[uint]int)
	for _, count := range counts {
		result[count.LocationID] = count.Count
	}

	return result, nil
}

func (s *ServiceRequestDBModel) ReopenServiceRequest(requestID uint) error {
	// Implement logic to reopen a closed service request by updating its status to "Open" in the database.
	if err := s.DB.Model(&ServiceRequest{}).Where("id = ?", requestID).Update("status", "Open").Error; err != nil {
		return err
	}
	return nil
}

func (s *ServiceRequestDBModel) GetOpenServiceRequests() ([]*ServiceRequest, error) {
	// Implement logic to retrieve all open service requests from the database.
	var requests []*ServiceRequest
	err := s.DB.Where("status = ?", "Open").Find(&requests).Error
	if err != nil {
		return nil, err
	}
	return requests, nil
}

func (s *ServiceRequestDBModel) GetClosedServiceRequests() ([]*ServiceRequest, error) {
	// Implement logic to retrieve all closed service requests from the database.
	var requests []*ServiceRequest
	err := s.DB.Where("status = ?", "Closed").Find(&requests).Error
	if err != nil {
		return nil, err
	}
	return requests, nil
}

func (s *ServiceRequestDBModel) DeleteServiceRequest(requestID uint) error {
	// Implement logic to delete a service request by its ID from the database.
	if err := s.DB.Delete(&ServiceRequest{}, requestID).Error; err != nil {
		return err
	}
	return nil
}

func (s *ServiceRequestDBModel) UpdateServiceRequestStatus2(requestID uint, status string) error {
	// Implement logic to update the status of a service request by its ID in the database.
	if err := s.DB.Model(&ServiceRequest{}).Where("id = ?", requestID).Update("status", status).Error; err != nil {
		return err
	}
	return nil
}

// GetServiceRequestsByCategory retrieves service requests by category ID.
func (srm *ServiceRequestDBModel) GetServiceRequestsByCategory(categoryID uint) ([]*ServiceRequest, error) {
	var requests []*ServiceRequest
	err := srm.DB.Where("category_id = ?", categoryID).Find(&requests).Error
	if err != nil {
		return nil, err
	}
	return requests, nil
}

// GetServiceRequestsBySubCategory retrieves service requests by sub-category ID.
func (srm *ServiceRequestDBModel) GetServiceRequestsBySubCategory(subCategoryID uint) ([]*ServiceRequest, error) {
	var requests []*ServiceRequest
	err := srm.DB.Where("sub_category_id = ?", subCategoryID).Find(&requests).Error
	if err != nil {
		return nil, err
	}
	return requests, nil
}

// GetUserClosedServiceRequests retrieves closed service requests associated with a user.
func (srm *ServiceRequestDBModel) GetUserClosedServiceRequests(userID uint) ([]*ServiceRequest, error) {
	var requests []*ServiceRequest
	err := srm.DB.Where("user_id = ? AND status = ?", userID, "Closed").Find(&requests).Error
	if err != nil {
		return nil, err
	}
	return requests, nil
}

// GetUserOpenServiceRequests retrieves open service requests associated with a user.
func (srm *ServiceRequestDBModel) GetUserOpenServiceRequests(userID uint) ([]*ServiceRequest, error) {
	var requests []*ServiceRequest
	err := srm.DB.Where("user_id = ? AND status = ?", userID, "Open").Find(&requests).Error
	if err != nil {
		return nil, err
	}
	return requests, nil
}

// GetServiceRequestsByStatus retrieves service requests by status.
func (srm *ServiceRequestDBModel) GetServiceRequestsByStatus(status string) ([]*ServiceRequest, error) {
	var requests []*ServiceRequest
	err := srm.DB.Where("status = ?", status).Find(&requests).Error
	if err != nil {
		return nil, err
	}
	return requests, nil
}

// GetServiceRequestCountByStatus returns the count of service requests grouped by status.
func (srm *ServiceRequestDBModel) GetServiceRequestCountByStatus() (map[string]int, error) {
	var result []struct {
		Status string
		Count  int
	}
	err := srm.DB.Model(&ServiceRequest{}).Select("status, count(*) as count").Group("status").Scan(&result).Error
	if err != nil {
		return nil, err
	}

	countMap := make(map[string]int)
	for _, r := range result {
		countMap[r.Status] = r.Count
	}
	return countMap, nil
}

// GetServiceRequestsByCategoryAndStatus retrieves service requests by category ID and status.
func (srm *ServiceRequestDBModel) GetServiceRequestsByCategoryAndStatus(categoryID uint, status string) ([]*ServiceRequest, error) {
	var requests []*ServiceRequest
	err := srm.DB.Where("category_id = ? AND status = ?", categoryID, status).Find(&requests).Error
	if err != nil {
		return nil, err
	}
	return requests, nil
}

// GetServiceRequestsBySubCategoryAndStatus retrieves service requests by sub-category ID and status.
func (srm *ServiceRequestDBModel) GetServiceRequestsBySubCategoryAndStatus(subCategoryID uint, status string) ([]*ServiceRequest, error) {
	var requests []*ServiceRequest
	err := srm.DB.Where("sub_category_id = ? AND status = ?", subCategoryID, status).Find(&requests).Error
	if err != nil {
		return nil, err
	}
	return requests, nil
}

// GetServiceRequestCountByCategory returns the count of service requests grouped by category.
func (srm *ServiceRequestDBModel) GetServiceRequestCountByCategory() (map[uint]int, error) {
	var result []struct {
		CategoryID uint
		Count      int
	}
	err := srm.DB.Model(&ServiceRequest{}).Select("category_id, count(*) as count").Group("category_id").Scan(&result).Error
	if err != nil {
		return nil, err
	}

	countMap := make(map[uint]int)
	for _, r := range result {
		countMap[r.CategoryID] = r.Count
	}
	return countMap, nil
}

// GetServiceRequestCountBySubCategory returns the count of service requests grouped by sub-category.
func (srm *ServiceRequestDBModel) GetServiceRequestCountBySubCategory() (map[uint]int, error) {
	var result []struct {
		SubCategoryID uint
		Count         int
	}
	err := srm.DB.Model(&ServiceRequest{}).Select("sub_category_id, count(*) as count").Group("sub_category_id").Scan(&result).Error
	if err != nil {
		return nil, err
	}

	countMap := make(map[uint]int)
	for _, r := range result {
		countMap[r.SubCategoryID] = r.Count
	}
	return countMap, nil
}

// GetServiceRequestsByUserAndStatus retrieves service requests by user ID and status.
func (srm *ServiceRequestDBModel) GetServiceRequestsByUserAndStatus(userID uint, status string) ([]*ServiceRequest, error) {
	var requests []*ServiceRequest
	err := srm.DB.Where("user_id = ? AND status = ?", userID, status).Find(&requests).Error
	if err != nil {
		return nil, err
	}
	return requests, nil
}

// GetUserServiceRequestCount returns the count of service requests associated with a user.
func (srm *ServiceRequestDBModel) GetUserServiceRequestCount(userID uint) (uint, error) {
	var count int64
	err := srm.DB.Model(&ServiceRequest{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return uint(count), nil
}

// GetUserServiceRequestCountByStatus returns the count of service requests associated with a user grouped by status.
func (srm *ServiceRequestDBModel) GetUserServiceRequestCountByStatus(userID uint) (map[string]int, error) {
	var result []struct {
		Status string
		Count  int
	}
	err := srm.DB.Model(&ServiceRequest{}).Select("status, count(*) as count").Where("user_id = ?", userID).Group("status").Scan(&result).Error
	if err != nil {
		return nil, err
	}

	countMap := make(map[string]int)
	for _, r := range result {
		countMap[r.Status] = r.Count
	}
	return countMap, nil
}

// GetUserServiceRequestCountByCategory returns the count of service requests for a user grouped by category.
func (srm *ServiceRequestDBModel) GetUserServiceRequestCountByCategory(userID uint) (map[uint]int, error) {
	var counts []struct {
		CategoryID uint
		Count      int
	}

	err := srm.DB.Model(&ServiceRequest{}).
		Joins("JOIN service_request_categories ON service_requests.category_id = service_request_categories.id").
		Where("service_requests.user_id = ?", userID).
		Group("service_request_categories.id").
		Select("service_request_categories.id as category_id, COUNT(service_requests.id) as count").
		Scan(&counts).Error

	if err != nil {
		return nil, err
	}

	result := make(map[uint]int)
	for _, count := range counts {
		result[count.CategoryID] = count.Count
	}

	return result, nil
}

// GetUserServiceRequestCountBySubCategory returns the count of service requests for a user grouped by sub-category.
func (srm *ServiceRequestDBModel) GetUserServiceRequestCountBySubCategory(userID uint) (map[uint]int, error) {
	var counts []struct {
		SubCategoryID uint
		Count         int
	}

	err := srm.DB.Model(&ServiceRequest{}).
		Joins("JOIN service_request_subcategories ON service_requests.subcategory_id = service_request_subcategories.id").
		Where("service_requests.user_id = ?", userID).
		Group("service_request_subcategories.id").
		Select("service_request_subcategories.id as subcategory_id, COUNT(service_requests.id) as count").
		Scan(&counts).Error

	if err != nil {
		return nil, err
	}

	result := make(map[uint]int)
	for _, count := range counts {
		result[count.SubCategoryID] = count.Count
	}

	return result, nil
}

// GetServiceRequestsByCategoryAndUser returns service requests for a user filtered by category.
func (srm *ServiceRequestDBModel) GetServiceRequestsByCategoryAndUser(categoryID, userID uint) ([]*ServiceRequest, error) {
	var requests []*ServiceRequest

	err := srm.DB.Where("user_id = ? AND category_id = ?", userID, categoryID).Find(&requests).Error
	if err != nil {
		return nil, err
	}

	return requests, nil
}

// GetServiceRequestsBySubCategoryAndUser returns service requests for a user filtered by sub-category.
func (srm *ServiceRequestDBModel) GetServiceRequestsBySubCategoryAndUser(subCategoryID, userID uint) ([]*ServiceRequest, error) {
	var requests []*ServiceRequest

	err := srm.DB.Where("user_id = ? AND subcategory_id = ?", userID, subCategoryID).Find(&requests).Error
	if err != nil {
		return nil, err
	}

	return requests, nil
}

// GetServiceRequestCountByCategoryAndUser returns the count of service requests for a user grouped by category.
func (srm *ServiceRequestDBModel) GetServiceRequestCountByCategoryAndUser(userID uint) (map[uint]int, error) {
	var counts []struct {
		CategoryID uint
		Count      int
	}

	err := srm.DB.Model(&ServiceRequest{}).
		Joins("JOIN service_request_categories ON service_requests.category_id = service_request_categories.id").
		Where("service_requests.user_id = ?", userID).
		Group("service_request_categories.id").
		Select("service_request_categories.id as category_id, COUNT(service_requests.id) as count").
		Scan(&counts).Error

	if err != nil {
		return nil, err
	}

	result := make(map[uint]int)
	for _, count := range counts {
		result[count.CategoryID] = count.Count
	}

	return result, nil
}

// GetServiceRequestCountBySubCategoryAndUser returns the count of service requests for a user grouped by sub-category.
func (srm *ServiceRequestDBModel) GetServiceRequestCountBySubCategoryAndUser(userID uint) (map[uint]int, error) {
	var counts []struct {
		SubCategoryID uint
		Count         int
	}

	err := srm.DB.Model(&ServiceRequest{}).
		Joins("JOIN service_request_subcategories ON service_requests.subcategory_id = service_request_subcategories.id").
		Where("service_requests.user_id = ?", userID).
		Group("service_request_subcategories.id").
		Select("service_request_subcategories.id as subcategory_id, COUNT(service_requests.id) as count").
		Scan(&counts).Error

	if err != nil {
		return nil, err
	}

	result := make(map[uint]int)
	for _, count := range counts {
		result[count.SubCategoryID] = count.Count
	}

	return result, nil
}

// GetServiceRequestCountByCategoryAndStatus returns the count of service requests for a given status and grouped by category.
func (srm *ServiceRequestDBModel) GetServiceRequestCountByCategoryAndStatus(status string) (map[uint]int, error) {
	var counts []struct {
		CategoryID uint
		Count      int
	}

	err := srm.DB.Model(&ServiceRequest{}).
		Joins("JOIN service_request_categories ON service_requests.category_id = service_request_categories.id").
		Where("service_requests.status = ?", status).
		Group("service_request_categories.id").
		Select("service_request_categories.id as category_id, COUNT(service_requests.id) as count").
		Scan(&counts).Error

	if err != nil {
		return nil, err
	}

	result := make(map[uint]int)
	for _, count := range counts {
		result[count.CategoryID] = count.Count
	}

	return result, nil
}

// GetServiceRequestCountBySubCategoryAndStatus returns the count of service requests for a given status and grouped by sub-category.
func (srm *ServiceRequestDBModel) GetServiceRequestCountBySubCategoryAndStatus2(status string) (map[uint]int, error) {
	var counts []struct {
		SubCategoryID uint
		Count         int
	}

	err := srm.DB.Model(&ServiceRequest{}).
		Joins("JOIN service_request_subcategories ON service_requests.subcategory_id = service_request_subcategories.id").
		Where("service_requests.status = ?", status).
		Group("service_request_subcategories.id").
		Select("service_request_subcategories.id as subcategory_id, COUNT(service_requests.id) as count").
		Scan(&counts).Error

	if err != nil {
		return nil, err
	}

	result := make(map[uint]int)
	for _, count := range counts {
		result[count.SubCategoryID] = count.Count
	}

	return result, nil
}

// GetUserAssignedServiceRequests returns service requests assigned to a user.
func (srm *ServiceRequestDBModel) GetUserAssignedServiceRequests2(userID uint) ([]*ServiceRequest, error) {
	var requests []*ServiceRequest

	err := srm.DB.Where("assigned_to = ?", userID).Find(&requests).Error
	if err != nil {
		return nil, err
	}

	return requests, nil
}

// AssignServiceRequestToUser assigns a service request to a user.
func (srm *ServiceRequestDBModel) AssignServiceRequestToUser2(requestID, userID uint) error {
	return srm.DB.Model(&ServiceRequest{}).
		Where("id = ?", requestID).
		Update("assigned_to", userID).
		Error
}

// GetServiceRequestsByPriority returns service requests filtered by priority.
func (srm *ServiceRequestDBModel) GetServiceRequestsByPriority2(priority string) ([]*ServiceRequest, error) {
	var requests []*ServiceRequest

	err := srm.DB.Where("priority = ?", priority).Find(&requests).Error
	if err != nil {
		return nil, err
	}

	return requests, nil
}

// GetServiceRequestCountByPriority returns the count of service requests grouped by priority.
func (srm *ServiceRequestDBModel) GetServiceRequestCountByPriority2() (map[string]int, error) {
	var counts []struct {
		Priority string
		Count    int
	}

	err := srm.DB.Model(&ServiceRequest{}).
		Group("priority").
		Select("priority, COUNT(id) as count").
		Scan(&counts).Error

	if err != nil {
		return nil, err
	}

	result := make(map[string]int)
	for _, count := range counts {
		result[count.Priority] = count.Count
	}

	return result, nil
}

// GetServiceRequestsByUserAndPriority returns service requests for a user filtered by priority.
func (srm *ServiceRequestDBModel) GetServiceRequestsByUserAndPriority2(userID uint, priority string) ([]*ServiceRequest, error) {
	var requests []*ServiceRequest

	err := srm.DB.Where("user_id = ? AND priority = ?", userID, priority).Find(&requests).Error
	if err != nil {
		return nil, err
	}

	return requests, nil
}

// GetServiceRequestCountByUserAndPriority returns the count of service requests for a user grouped by priority.
func (srm *ServiceRequestDBModel) GetServiceRequestCountByUserAndPriority2(userID uint) (map[string]int, error) {
	var counts []struct {
		Priority string
		Count    int
	}

	err := srm.DB.Model(&ServiceRequest{}).
		Where("user_id = ?", userID).
		Group("priority").
		Select("priority, COUNT(id) as count").
		Scan(&counts).Error

	if err != nil {
		return nil, err
	}

	result := make(map[string]int)
	for _, count := range counts {
		result[count.Priority] = count.Count
	}

	return result, nil
}

// GetServiceRequestCountByCategoryAndPriority returns the count of service requests grouped by category and priority.
func (srm *ServiceRequestDBModel) GetServiceRequestCountByCategoryAndPriority2(categoryID uint) (map[string]int, error) {
	var counts []struct {
		CategoryID uint
		Priority   string
		Count      int
	}

	err := srm.DB.Model(&ServiceRequest{}).
		Where("category_id = ?", categoryID).
		Group("priority").
		Select("category_id, priority, COUNT(id) as count").
		Scan(&counts).Error

	if err != nil {
		return nil, err
	}

	result := make(map[string]int)
	for _, count := range counts {
		result[count.Priority] = count.Count
	}

	return result, nil
}

// GetServiceRequestsByUserAndLocation returns service requests for a user filtered by location.
func (srm *ServiceRequestDBModel) GetServiceRequestsByUserAndLocation(userID, locationID uint) ([]*ServiceRequest, error) {
	var requests []*ServiceRequest

	err := srm.DB.Where("user_id = ? AND location_id = ?", userID, locationID).Find(&requests).Error
	if err != nil {
		return nil, err
	}

	return requests, nil
}

// GetServiceRequestCountByUserAndLocation returns the count of service requests for a user grouped by location.
func (srm *ServiceRequestDBModel) GetServiceRequestCountByUserAndLocation(userID uint) (map[uint]int, error) {
	var counts []struct {
		LocationID uint
		Count      int
	}

	err := srm.DB.Model(&ServiceRequest{}).
		Where("user_id = ?", userID).
		Group("location_id").
		Select("location_id, COUNT(id) as count").
		Scan(&counts).Error

	if err != nil {
		return nil, err
	}

	result := make(map[uint]int)
	for _, count := range counts {
		result[count.LocationID] = count.Count
	}

	return result, nil
}

// GetServiceRequestsByCategoryAndLocation returns service requests filtered by category and location.
func (srm *ServiceRequestDBModel) GetServiceRequestsByCategoryAndLocation(categoryID, locationID uint) ([]*ServiceRequest, error) {
	var requests []*ServiceRequest

	err := srm.DB.Where("category_id = ? AND location_id = ?", categoryID, locationID).Find(&requests).Error
	if err != nil {
		return nil, err
	}

	return requests, nil
}

// GetServiceRequestCountByCategoryAndLocation returns the count of service requests by category and location.
func (srm *ServiceRequestDBModel) GetServiceRequestCountByCategoryAndLocation(categoryID uint) (map[uint]int, error) {
	var counts []struct {
		LocationID uint
		Count      int
	}

	query := srm.DB.Model(&ServiceRequest{}).
		Select("location_id, COUNT(*) as count").
		Where("category_id = ?", categoryID).
		Group("location_id").
		Scan(&counts)

	if query.Error != nil {
		return nil, query.Error
	}

	locationCount := make(map[uint]int)
	for _, count := range counts {
		locationCount[count.LocationID] = count.Count
	}

	return locationCount, nil
}

// GetServiceRequestsByPriorityAndLocation retrieves service requests by priority and location.
func (srm *ServiceRequestDBModel) GetServiceRequestsByPriorityAndLocation(priority string, locationID uint) ([]*ServiceRequest, error) {
	var requests []*ServiceRequest
	err := srm.DB.Where("priority = ? AND location_id = ?", priority, locationID).Find(&requests).Error
	if err != nil {
		return nil, err
	}
	return requests, nil
}

// GetServiceRequestCountByPriorityAndLocation returns the count of service requests for a specific priority and location.
func (srm *ServiceRequestDBModel) GetServiceRequestCountByPriorityAndLocation(priority string, locationID uint) (map[string]int, error) {
	var result []struct {
		Priority   string
		LocationID uint
		Count      int
	}

	err := srm.DB.Table("service_requests").
		Select("priority, location_id, COUNT(*) as count").
		Where("priority = ? AND location_id = ?", priority, locationID).
		Group("priority, location_id").
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	countMap := make(map[string]int)
	for _, r := range result {
		countMap[r.Priority] = r.Count
	}

	return countMap, nil
}

// GetServiceRequestHistory retrieves the history entries of a service request by its ID.
func (srm *ServiceRequestDBModel) GetServiceRequestHistory(requestID uint) ([]*ServiceRequestHistoryEntry, error) {
	var history []*ServiceRequestHistoryEntry
	err := srm.DB.Where("service_request_id = ?", requestID).Find(&history).Error
	if err != nil {
		return nil, err
	}
	return history, nil
}

// AddCommentToServiceRequest adds a comment to a service request.
func (srm *ServiceRequestDBModel) AddCommentToServiceRequest(requestID uint, comment string) error {
	commentEntry := &ServiceRequestComment{
		ServiceRequestID: requestID,
		Comment:          comment,
	}
	if err := srm.DB.Create(commentEntry).Error; err != nil {
		return err
	}
	return nil
}

// GetServiceRequestComments retrieves comments for a service request by its ID.
func (srm *ServiceRequestDBModel) GetServiceRequestComments(requestID uint) ([]*ServiceRequestComment, error) {
	var comments []*ServiceRequestComment
	err := srm.DB.Where("service_request_id = ?", requestID).Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// CreateServiceRequest creates a new service request in the database.
func (srm *ServiceRequestDBModel) CreateServiceRequest(request *ServiceRequest) error {
	if err := srm.DB.Create(request).Error; err != nil {
		return err
	}
	return nil
}

// UpdateServiceRequest updates an existing service request in the database.
func (srm *ServiceRequestDBModel) UpdateServiceRequest(request *ServiceRequest) error {
	if err := srm.DB.Save(request).Error; err != nil {
		return err
	}
	return nil
}

// GetServiceRequestByID retrieves a service request by its ID.
func (srm *ServiceRequestDBModel) GetServiceRequestByID(requestID uint) (*ServiceRequest, error) {
	var request ServiceRequest
	err := srm.DB.Where("id = ?", requestID).First(&request).Error
	if err != nil {
		return nil, err
	}
	return &request, nil
}

// GetUserServiceRequests retrieves all service requests associated with a user from the database.
func (srm *ServiceRequestDBModel) GetUserServiceRequests(userID uint) ([]*ServiceRequest, error) {
	var requests []*ServiceRequest
	err := srm.DB.Where("user_id = ?", userID).Find(&requests).Error
	if err != nil {
		return nil, err
	}
	return requests, nil
}

// CloseServiceRequest closes a service request by setting its status to "Closed" in the database.
func (srm *ServiceRequestDBModel) CloseServiceRequest(requestID uint) error {
	return srm.DB.Model(&ServiceRequest{}).Where("id = ?", requestID).Update("status", "Closed").Error
}

// SubmitServiceRequest captures a new service request from a user or system.
func (db *ServiceRequestDBModel) SubmitServiceRequest(request *ServiceRequest) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(request).Error; err != nil {
			return err
		}
		return nil
	})
}

// UpdateServiceRequestStatus changes the status of a service request, e.g., to "Completed".
func (db *ServiceRequestDBModel) UpdateServiceRequestStatus(requestID uint, status string) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&ServiceRequest{}).Where("id = ?", requestID).Update("status", status).Error; err != nil {
			return err
		}
		return nil
	})
}
