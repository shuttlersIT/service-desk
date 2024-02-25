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
	CreatePosition(*models.Position) error
	DeletePosition(int) error
	UpdatePosition(*models.Position) error
	GetAllPositions() ([]*models.Position, error)
	GetPositionByID(int) (*models.Position, error)
	GetPositionByNumber(int) (*models.Position, error)
	CreateDepartment(*models.Department) error
	DeleteDepartment(int) error
	UpdateDepartment(*models.Department) error
	GetAllDepartments() ([]*models.Department, error)
	GetDepartmentByID(int) (*models.Department, error)
	GetDepartmentByNumber(int) (*models.Department, error)
}

// DefaultUserService is the default implementation of UserService
type DefaultUserService struct {
	DB             *gorm.DB
	UserDBModel    *models.UserDBModel
	log            models.Logger
	EventPublisher models.EventPublisherImpl
	// Add any dependencies or data needed for the service
}

// NewDefaultUserService creates a new DefaultUserService.
func NewDefaultUserService(users *models.UserDBModel, log models.Logger, eventPublisher models.EventPublisherImpl) *DefaultUserService {
	return &DefaultUserService{
		UserDBModel:    users,
		EventPublisher: eventPublisher,
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

// CreatePosition creates a new position.
func (us *DefaultUserService) CreatePosition(position *models.Position) error {
	return us.UserDBModel.CreatePosition(position)
}

// DeletePosition deletes a position by its ID.
func (us *DefaultUserService) DeletePosition(positionID uint) error {
	return us.UserDBModel.DeletePosition(positionID)
}

// UpdatePosition updates an existing position.
func (us *DefaultUserService) UpdatePosition(position *models.Position) error {
	return us.UserDBModel.UpdatePosition(position)
}

// GetPosition retrieves all positions.
func (us *DefaultUserService) GetAllPositions() ([]*models.Position, error) {
	return us.UserDBModel.GetAllPositions()
}

// GetPositionByID retrieves a position by its ID.
func (us *DefaultUserService) GetPositionByID(positionID uint) (*models.Position, error) {
	return us.UserDBModel.GetPositionByID(positionID)
}

// GetPositionByNumber retrieves a position by its number.
func (us *DefaultUserService) GetPositionByNumber(positionNumber int) (*models.Position, error) {
	return us.UserDBModel.GetPositionByNumber(positionNumber)
}

// CreateDepartment creates a new department.
func (us *DefaultUserService) CreateDepartment(department *models.Department) error {
	return us.UserDBModel.CreateDepartment(department)
}

// DeleteDepartment deletes a department by its ID.
func (us *DefaultUserService) DeleteDepartment(departmentID uint) error {
	return us.UserDBModel.DeleteDepartment(departmentID)
}

// UpdateDepartment updates an existing department.
func (us *DefaultUserService) UpdateDepartment(department *models.Department) error {
	return us.UserDBModel.UpdateDepartment(department)
}

// GetDepartments retrieves all departments.
func (us *DefaultUserService) GetAllDepartments() ([]*models.Department, error) {
	return us.UserDBModel.GetAllDepartments()
}

// GetDepartmentByID retrieves a department by its ID.
func (us *DefaultUserService) GetDepartmentByID(departmentID uint) (*models.Department, error) {
	return us.UserDBModel.GetDepartmentByID(departmentID)
}

// GetDepartmentByNumber retrieves a department by its number.
func (us *DefaultUserService) GetDepartmentByNumber(departmentNumber int) (*models.Department, error) {
	return us.UserDBModel.GetDepartmentByNumber(departmentNumber)
}
