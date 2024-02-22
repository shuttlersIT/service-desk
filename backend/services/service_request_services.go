// Service Requests service_request_service.go

package services

import (
	"fmt"

	"github.com/shuttlersit/service-desk/models"
	"gorm.io/gorm"
)

type ServiceRequestService interface {
	CreateServiceRequest(request *models.ServiceRequest) error
	UpdateServiceRequest(request *models.ServiceRequest) error
	GetServiceRequestByID(requestID uint) (*models.ServiceRequest, error)
	GetUserServiceRequests(userID uint) ([]*models.ServiceRequest, error)
	CloseServiceRequest(requestID uint) error
	ReopenServiceRequest(requestID uint) error
	GetServiceRequestHistory(requestID uint) ([]*models.ServiceRequestHistoryEntry, error)
	AddCommentToServiceRequest(requestID uint, comment string) error
	GetServiceRequestComments(requestID uint) ([]*models.ServiceRequestComment, error)
	GetOpenServiceRequests() ([]*models.ServiceRequest, error)
	GetClosedServiceRequests() ([]*models.ServiceRequest, error)
	GetAllServiceRequests() ([]*models.ServiceRequest, error)
	DeleteServiceRequest(requestID uint) error
	UpdateServiceRequestStatus(requestID uint, status string) error
	GetServiceRequestsByCategory(categoryID uint) ([]*models.ServiceRequest, error)
	GetServiceRequestsBySubCategory(subCategoryID uint) ([]*models.ServiceRequest, error)
	GetUserClosedServiceRequests(userID uint) ([]*models.ServiceRequest, error)
	GetUserOpenServiceRequests(userID uint) ([]*models.ServiceRequest, error)
	GetServiceRequestsByStatus(status string) ([]*models.ServiceRequest, error)
	GetServiceRequestCountByStatus() (map[string]int, error)
	GetServiceRequestsByCategoryAndStatus(categoryID uint, status string) ([]*models.ServiceRequest, error)
	GetServiceRequestsBySubCategoryAndStatus(subCategoryID uint, status string) ([]*models.ServiceRequest, error)
	GetServiceRequestCountByCategory() (map[uint]int, error)
	GetServiceRequestCountBySubCategory() (map[uint]int, error)
	GetServiceRequestsByUserAndStatus(userID uint, status string) ([]*models.ServiceRequest, error)
	GetUserServiceRequestCount(userID uint) (int, error)
	GetUserServiceRequestCountByStatus(userID uint) (map[string]int, error)
	GetServiceRequestsByLocation(locationID uint) ([]*models.ServiceRequest, error)
	GetServiceRequestCountByLocation() (map[uint]int, error)
	GetServiceRequestsByUserAndLocation(userID, locationID uint) ([]*models.ServiceRequest, error)
	GetServiceRequestCountByUserAndLocation(userID uint) (map[uint]int, error)
	GetServiceRequestsByCategoryAndLocation(categoryID, locationID uint) ([]*models.ServiceRequest, error)
	GetServiceRequestCountByCategoryAndLocation(categoryID uint) (map[uint]int, error)
	GetServiceRequestsByPriorityAndLocation(priority string, locationID uint) ([]*models.ServiceRequest, error)
	GetServiceRequestCountByPriorityAndLocation(locationID uint) (map[string]int, error)
	GetServiceRequestsByCategoryAndSubCategory(categoryID uint, subCategoryID uint) ([]*models.ServiceRequest, error)
	GetServiceRequestsByStatusAndLocation(status string, locationID uint) ([]*models.ServiceRequest, error)
}

type DefaultServiceRequestService struct {
	DB                    *gorm.DB
	ServiceRequestDBModel *models.ServiceRequestDBModel
	log                   models.PrintLogger
	EventPublisher        *models.EventPublisherImpl
}

func NewDefaultServiceRequestService(db *gorm.DB, serviceRequestDBModel *models.ServiceRequestDBModel, log models.PrintLogger, eventPublisher *models.EventPublisherImpl) *DefaultServiceRequestService {
	return &DefaultServiceRequestService{
		DB:                    db,
		ServiceRequestDBModel: serviceRequestDBModel,
		log:                   log,
		EventPublisher:        eventPublisher,
	}
}

func (s *DefaultServiceRequestService) CreateServiceRequest(request *models.ServiceRequest) error {
	// You can directly call the corresponding method from the DB model.
	return s.ServiceRequestDBModel.CreateServiceRequest(request)
}

func (s *DefaultServiceRequestService) GetServiceRequestByID(requestID uint) (*models.ServiceRequest, error) {
	// You can directly call the corresponding method from the DB model.
	return s.ServiceRequestDBModel.GetServiceRequestByID(requestID)
}

func (s *DefaultServiceRequestService) UpdateServiceRequest(request *models.ServiceRequest) error {
	// You can directly call the corresponding method from the DB model.
	return s.ServiceRequestDBModel.UpdateServiceRequest(request)
}

func (s *DefaultServiceRequestService) GetUserServiceRequests(userID uint) ([]*models.ServiceRequest, error) {
	// You can directly call the corresponding method from the DB model.
	return s.ServiceRequestDBModel.GetUserServiceRequests(userID)
}

func (s *DefaultServiceRequestService) CloseServiceRequest(requestID uint) error {
	// You can directly call the corresponding method from the DB model.
	return s.ServiceRequestDBModel.CloseServiceRequest(requestID)
}

func (s *DefaultServiceRequestService) ReopenServiceRequest(requestID uint) error {
	// You can directly call the corresponding method from the DB model.
	return s.ServiceRequestDBModel.ReopenServiceRequest(requestID)
}

func (s *DefaultServiceRequestService) GetAllServiceRequests() ([]*models.ServiceRequest, error) {
	// You can directly call the corresponding method from the DB model.
	return s.ServiceRequestDBModel.GetAllServiceRequests()
}

func (s *DefaultServiceRequestService) GetServiceRequestHistory(requestID uint) ([]*models.ServiceRequestHistoryEntry, error) {
	// You can directly call the corresponding method from the DB model.
	return s.ServiceRequestDBModel.GetServiceRequestHistory(requestID)
}

func (s *DefaultServiceRequestService) AddCommentToServiceRequest(requestID uint, comment string) error {
	// You can directly call the corresponding method from the DB model.
	return s.ServiceRequestDBModel.AddCommentToServiceRequest(requestID, comment)
}
func (s *DefaultServiceRequestService) GetServiceRequestComments(requestID uint) ([]*models.ServiceRequestComment, error) {
	// You can directly call the corresponding method from the DB model.
	return s.ServiceRequestDBModel.GetServiceRequestComments(requestID)
}

func (s *DefaultServiceRequestService) GetOpenServiceRequests() ([]*models.ServiceRequest, error) {
	return s.ServiceRequestDBModel.GetOpenServiceRequests()
}

func (s *DefaultServiceRequestService) GetClosedServiceRequests() ([]*models.ServiceRequest, error) {
	return s.ServiceRequestDBModel.GetClosedServiceRequests()
}

func (s *DefaultServiceRequestService) DeleteServiceRequest(requestID uint) error {
	return s.ServiceRequestDBModel.DeleteServiceRequest(requestID)
}

func (s *DefaultServiceRequestService) UpdateServiceRequestStatus(requestID uint, status string) error {
	return s.ServiceRequestDBModel.UpdateServiceRequestStatus(requestID, status)
}

func (s *DefaultServiceRequestService) GetServiceRequestsByCategory(categoryID uint) ([]*models.ServiceRequest, error) {
	return s.ServiceRequestDBModel.GetServiceRequestsByCategory(categoryID)
}

func (s *DefaultServiceRequestService) GetServiceRequestsBySubCategory(subCategoryID uint) ([]*models.ServiceRequest, error) {
	return s.ServiceRequestDBModel.GetServiceRequestsBySubCategory(subCategoryID)
}

func (s *DefaultServiceRequestService) GetUserClosedServiceRequests(userID uint) ([]*models.ServiceRequest, error) {
	return s.ServiceRequestDBModel.GetUserClosedServiceRequests(userID)
}

func (s *DefaultServiceRequestService) GetUserOpenServiceRequests(userID uint) ([]*models.ServiceRequest, error) {
	return s.ServiceRequestDBModel.GetUserOpenServiceRequests(userID)
}

func (s *DefaultServiceRequestService) GetServiceRequestsByStatus(status string) ([]*models.ServiceRequest, error) {
	return s.ServiceRequestDBModel.GetServiceRequestsByStatus(status)
}

func (s *DefaultServiceRequestService) GetServiceRequestCountByStatus() (map[string]int, error) {
	return s.ServiceRequestDBModel.GetServiceRequestCountByStatus()
}

func (s *DefaultServiceRequestService) GetServiceRequestsByCategoryAndStatus(categoryID uint, status string) ([]*models.ServiceRequest, error) {
	return s.ServiceRequestDBModel.GetServiceRequestsByCategoryAndStatus(categoryID, status)
}

func (s *DefaultServiceRequestService) GetServiceRequestsBySubCategoryAndStatus(subCategoryID uint, status string) ([]*models.ServiceRequest, error) {
	return s.ServiceRequestDBModel.GetServiceRequestsBySubCategoryAndStatus(subCategoryID, status)
}

func (s *DefaultServiceRequestService) GetServiceRequestCountByCategory() (map[uint]int, error) {
	return s.ServiceRequestDBModel.GetServiceRequestCountByCategory()
}

func (s *DefaultServiceRequestService) GetServiceRequestCountBySubCategory() (map[uint]int, error) {
	return s.ServiceRequestDBModel.GetServiceRequestCountBySubCategory()
}

func (s *DefaultServiceRequestService) GetServiceRequestsByUserAndStatus(userID uint, status string) ([]*models.ServiceRequest, error) {
	return s.ServiceRequestDBModel.GetServiceRequestsByUserAndStatus(userID, status)
}

func (s *DefaultServiceRequestService) GetUserServiceRequestCount(userID uint) (uint, error) {
	return s.ServiceRequestDBModel.GetUserServiceRequestCount(userID)
}

func (s *DefaultServiceRequestService) GetUserServiceRequestCountByStatus(userID uint) (map[string]int, error) {
	return s.ServiceRequestDBModel.GetUserServiceRequestCountByStatus(userID)
}

func (s *DefaultServiceRequestService) GetUserServiceRequestCountByCategory(userID uint) (map[uint]int, error) {
	return s.ServiceRequestDBModel.GetUserServiceRequestCountByCategory(userID)
}

func (s *DefaultServiceRequestService) GetUserServiceRequestCountBySubCategory(userID uint) (map[uint]int, error) {
	return s.ServiceRequestDBModel.GetUserServiceRequestCountBySubCategory(userID)
}

func (s *DefaultServiceRequestService) GetServiceRequestsByCategoryAndUser(categoryID, userID uint) ([]*models.ServiceRequest, error) {
	return s.ServiceRequestDBModel.GetServiceRequestsByCategoryAndUser(categoryID, userID)
}

func (s *DefaultServiceRequestService) GetServiceRequestsBySubCategoryAndUser(subCategoryID, userID uint) ([]*models.ServiceRequest, error) {
	return s.ServiceRequestDBModel.GetServiceRequestsBySubCategoryAndUser(subCategoryID, userID)
}

func (s *DefaultServiceRequestService) GetServiceRequestCountByCategoryAndUser(userID uint) (map[uint]int, error) {
	return s.ServiceRequestDBModel.GetServiceRequestCountByCategoryAndUser(userID)
}

func (s *DefaultServiceRequestService) GetServiceRequestCountBySubCategoryAndUser(userID uint) (map[uint]int, error) {
	return s.ServiceRequestDBModel.GetServiceRequestCountBySubCategoryAndUser(userID)
}

func (s *DefaultServiceRequestService) GetServiceRequestCountByCategoryAndStatus(status string) (map[uint]int, error) {
	return s.ServiceRequestDBModel.GetServiceRequestCountByCategoryAndStatus(status)
}

func (s *DefaultServiceRequestService) GetServiceRequestCountBySubCategoryAndStatus(status string) (map[uint]int, error) {
	// You'll need to implement this method using your database model (ServiceRequestDBModel).
	return s.ServiceRequestDBModel.GetServiceRequestCountBySubCategoryAndStatus(status)
}

func (s *DefaultServiceRequestService) GetUserAssignedServiceRequests(userID uint) ([]*models.ServiceRequest, error) {
	// You'll need to implement this method using your database model (ServiceRequestDBModel).
	return s.ServiceRequestDBModel.GetUserAssignedServiceRequests(userID)
}

func (s *DefaultServiceRequestService) AssignServiceRequestToUser(requestID, userID uint) error {
	// You'll need to implement this method using your database model (ServiceRequestDBModel).
	return s.ServiceRequestDBModel.AssignServiceRequestToUser(requestID, userID)
}

func (s *DefaultServiceRequestService) GetServiceRequestsByPriority(priority string) ([]*models.ServiceRequest, error) {
	// You'll need to implement this method using your database model (ServiceRequestDBModel).
	return s.ServiceRequestDBModel.GetServiceRequestsByPriority(priority)
}

func (s *DefaultServiceRequestService) GetServiceRequestCountByPriority() (map[string]int, error) {
	// You'll need to implement this method using your database model (ServiceRequestDBModel).
	return s.ServiceRequestDBModel.GetServiceRequestCountByPriority()
}

func (s *DefaultServiceRequestService) GetServiceRequestsByUserAndPriority(userID uint, priority string) ([]*models.ServiceRequest, error) {
	// You'll need to implement this method using your database model (ServiceRequestDBModel).
	return s.ServiceRequestDBModel.GetServiceRequestsByUserAndPriority(userID, priority)
}

func (s *DefaultServiceRequestService) GetServiceRequestCountByUserAndPriority(userID uint) (map[string]int, error) {
	// You'll need to implement this method using your database model (ServiceRequestDBModel).
	return s.ServiceRequestDBModel.GetServiceRequestCountByUserAndPriority(userID)
}

func (s *DefaultServiceRequestService) GetServiceRequestCountByCategoryAndPriority(categoryID uint) (map[string]int, error) {
	// You'll need to implement this method using your database model (ServiceRequestDBModel).
	return s.ServiceRequestDBModel.GetServiceRequestCountByCategoryAndPriority(categoryID)
}

func (s *DefaultServiceRequestService) GetServiceRequestsByLocation(locationID uint) ([]*models.ServiceRequest, error) {
	// You'll need to implement this method using your database model (ServiceRequestDBModel).
	return s.ServiceRequestDBModel.GetServiceRequestsByLocation(locationID)
}

func (s *DefaultServiceRequestService) GetServiceRequestCountByLocation() (map[uint]int, error) {
	// You'll need to implement this method using your database model (ServiceRequestDBModel).
	return s.ServiceRequestDBModel.GetServiceRequestCountByLocation()
}

func (s *DefaultServiceRequestService) GetServiceRequestsByUserAndLocation(userID, locationID uint) ([]*models.ServiceRequest, error) {
	// You'll need to implement this method using your database model (ServiceRequestDBModel).
	return s.ServiceRequestDBModel.GetServiceRequestsByUserAndLocation(userID, locationID)
}

func (s *DefaultServiceRequestService) GetServiceRequestCountByUserAndLocation(userID uint) (map[uint]int, error) {
	// You'll need to implement this method using your database model (ServiceRequestDBModel).
	return s.ServiceRequestDBModel.GetServiceRequestCountByUserAndLocation(userID)
}

func (s *DefaultServiceRequestService) GetServiceRequestsByCategoryAndLocation(categoryID, locationID uint) ([]*models.ServiceRequest, error) {
	// You'll need to implement this method using your database model (ServiceRequestDBModel).
	return s.ServiceRequestDBModel.GetServiceRequestsByCategoryAndLocation(categoryID, locationID)
}

func (s *DefaultServiceRequestService) GetServiceRequestCountByCategoryAndLocation(categoryID uint) (map[uint]int, error) {
	// You'll need to implement this method using your database model (ServiceRequestDBModel).
	return s.ServiceRequestDBModel.GetServiceRequestCountByCategoryAndLocation(categoryID)
}

func (s *DefaultServiceRequestService) GetServiceRequestsByPriorityAndLocation(priority string, locationID uint) ([]*models.ServiceRequest, error) {
	// You'll need to implement this method using your database model (ServiceRequestDBModel).
	return s.ServiceRequestDBModel.GetServiceRequestsByPriorityAndLocation(priority, locationID)
}

func (s *DefaultServiceRequestService) GetServiceRequestCountByPriorityAndLocation(priority string, locationID uint) (map[string]int, error) {
	// You'll need to implement this method using your database model (ServiceRequestDBModel).
	return s.ServiceRequestDBModel.GetServiceRequestCountByPriorityAndLocation(priority, locationID)
}

func (s *DefaultServiceRequestService) GetServiceRequestsByCategoryAndSubCategory(categoryID uint, subCategoryID uint) ([]*models.ServiceRequest, error) {

	requests, err := s.ServiceRequestDBModel.GetServiceRequestsByCategory(categoryID)
	if err != nil {
		return nil, err
	}

	var matchingRequest []*models.ServiceRequest
	for _, s := range requests {
		if s.SubCategoryID == subCategoryID {
			matchingRequest = append(matchingRequest, s)
		}
	}
	if len(matchingRequest) > 0 {
		return matchingRequest, nil
	}

	return nil, fmt.Errorf("no matching service request")
}

func (s *DefaultServiceRequestService) GetServiceRequestsByStatusAndLocation(status string, locationID uint) ([]*models.ServiceRequest, error) {
	requests, err := s.ServiceRequestDBModel.GetServiceRequestsByStatus(status)
	if err != nil {
		return nil, err
	}

	var matchingRequest []*models.ServiceRequest
	for _, s := range requests {
		if s.LocationID == locationID {
			matchingRequest = append(matchingRequest, s)
		}
	}
	if len(matchingRequest) > 0 {
		return matchingRequest, nil
	}

	return nil, fmt.Errorf("no matching service request")
}
