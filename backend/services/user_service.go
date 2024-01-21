// backend/services/user_service.go

package services

import (
	"github.com/shuttlersit/service-desk/backend/models"
	"gorm.io/gorm"
)

// UserServiceInterface provides methods for managing users.
type UserServiceInterface interface {
	CreateUser(user *models.Users) (*models.Users, error)
	UpdateUser(user *models.Users) (*models.Users, error)
	GetUserByID(id uint) (*models.Users, error)
	DeleteUser(userID uint) (bool, error)
	GetAllUsers() []*models.Users
}

// DefaultUserService is the default implementation of UserService
type DefaultUserService struct {
	DB          *gorm.DB
	UserDBModel *models.UserDBModel
	// Add any dependencies or data needed for the service
}

// NewDefaultUserService creates a new DefaultUserService.
func NewDefaultUserService(users *models.UserDBModel) *DefaultUserService {
	return &DefaultUserService{
		UserDBModel: users,
	}
}

// GetAllUsers retrieves all users.
func (ps *DefaultUserService) GetAllUsers() ([]*models.Users, error) {
	users, err := ps.UserDBModel.GetAllUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}

// CreateUser creates a new user.
func (ps *DefaultUserService) CreateUser(user *models.Users) error {
	_, err := ps.UserDBModel.CreateUser(user)
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

func (us *DefaultUserService) UpdateUserProfile(userID uint, updatedUser *models.Users) error {
	// Retrieve the user by userID
	user, err := us.UserDBModel.GetUserByID(userID)
	if err != nil {
		return err
	}

	// Update user profile details
	user.FirstName = updatedUser.FirstName
	user.LastName = updatedUser.LastName
	user.Email = updatedUser.Email
	user.Phone = updatedUser.Phone

	// Save the updated user profile
	err = us.UserDBModel.UpdateUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (us *DefaultUserService) UpdateUserProfile2(userID uint, updatedUser *models.Users) error {
	// Update the user's profile details
	err := us.UserDBModel.UpdateUser(updatedUser)
	if err != nil {
		return err
	}

	return nil
}

func (us *DefaultUserService) DeleteUser2(userID uint) error {
	// Delete a user by user ID
	err := us.UserDBModel.DeleteUser(userID)
	if err != nil {
		return err
	}

	return nil
}
