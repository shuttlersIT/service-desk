// backend/services/user_service.go

package services

import (
	"github.com/shuttlersit/service-desk/backend/models"
	"gorm.io/gorm"
)

// AdvertisementServiceInterface provides methods for managing advertisements.
type UserServiceInterface interface {
	CreateUser(user *models.Users) (*models.Users, error)
	UpdateUser(user *models.Users) (*models.Users, error)
	GetUserByID(id uint) (*models.Users, error)
	DeleteUser(userID uint) (bool, error)
	GetAllUsers() *[]models.Users
}

// DefaultAdvertisementService is the default implementation of AdvertisementService
type DefaultUserService struct {
	DB          *gorm.DB
	UserDBModel *models.UserDBModel
	// Add any dependencies or data needed for the service
}

// NewDefaultAdvertisementService creates a new DefaultAdvertisementService.
func NewDefaultUserService(users *models.UserDBModel) *DefaultUserService {
	return &DefaultUserService{
		UserDBModel: users,
	}
}

// GetAllUsers retrieves all users.
func (ps *DefaultUserService) GetAllUsers() (*[]models.Users, error) {
	users, err := ps.UserDBModel.GetAllUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}

// CreateUser creates a new user.
func (ps *DefaultUserService) CreateUser(user *models.Users) error {
	err := ps.UserDBModel.CreateUser(user)
	if err != nil {
		return err
	}
	return nil
}

// CreateUser creates a new user.
func (ps *DefaultUserService) GetUserByID(id uint) (*models.Users, error) {
	user, err := ps.UserDBModel.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser updates an existing users.
func (ps *DefaultUserService) UpdateUser(user *models.Users) (*models.Users, error) {
	err := ps.UserDBModel.UpdateUser(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// DeleteUser deletes an user by ID.
func (ps *DefaultUserService) DeleteUser(userID uint) (bool, error) {
	status := false
	err := ps.UserDBModel.DeleteUser(userID)
	if err != nil {
		return status, err
	}
	status = true
	return status, nil
}
