package services

import (
	"github.com/shuttlersit/service-desk/backend/models"
)

type ServiceRequestService interface {
	CreateServiceRequest(request *models.ServiceRequest) error
	UpdateServiceRequest(request *models.ServiceRequest) error
	GetServiceRequestByID(requestID uint) (*models.ServiceRequest, error)
	GetUserServiceRequests(userID uint) ([]*models.ServiceRequest, error)
}

type DefaultServiceRequestService struct {
	// Implement methods defined in ServiceRequestService interface
}

// Implement the CreateServiceRequest, UpdateServiceRequest, GetServiceRequestByID,
// and GetUserServiceRequests methods.
