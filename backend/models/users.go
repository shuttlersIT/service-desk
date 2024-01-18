// backend/models/users.go

package models

import (
	"time"

	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	ID          uint                  `gorm:"primaryKey" json:"user_id"`
	FirstName   string                `json:"first_name" binding:"required"`
	LastName    string                `json:"last_name" binding:"required"`
	Email       string                `json:"staff_email" binding:"required,email"`
	Credentials UsersLoginCredentials `json:"user_credentials" gorm:"foreignKey:UserID"`
	Phone       string                `json:"phoneNumber" binding:"required,e164"`
	Position    Position              `json:"position_id" gorm:"embedded"`
	Department  Department            `json:"department_id" gorm:"embedded"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
	Asset       []AssetAssignment     `json:"asset_assignment" gorm:"foreignKey:UserID"`
}

// TableName sets the table name for the Users model.
func (Users) TableName() string {
	return "users"
}

type Position struct {
	gorm.Model
	PositionID   int       `gorm:"primaryKey" json:"position_id"`
	PositionName string    `json:"position_name"`
	CadreName    string    `json:"cadre_name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName sets the table name for the Position model.
func (Position) TableName() string {
	return "position"
}

type Department struct {
	gorm.Model
	DepartmentID   int       `gorm:"primaryKey" json:"department_id"`
	DepartmentName string    `json:"department_name"`
	Emoji          string    `json:"emoji"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// TableName sets the table name for the Department model.
func (Department) TableName() string {
	return "department"
}

type UserStorage interface {
	CreateUser(*Users) error
	DeleteUser(int) error
	UpdateUser(*Users) error
	GetUsers() (*[]Users, error)
	GetUserByID(int) (*Users, error)
	GetUserByNumber(int) (*Users, error)
}

type PositionStorage interface {
	CreatePosition(*Position) error
	DeletePosition(int) error
	UpdatePosition(*Position) error
	GetPosition() (*[]Position, error)
	GetPositionByID(int) (*Position, error)
	GetPositionByNumber(int) (*Position, error)
}

type DepartmentStorage interface {
	CreateDepartment(*Department) error
	DeleteDepartment(int) error
	UpdateDepartment(*Department) error
	GetDepartments() (*[]Department, error)
	GetDepartmentByID(int) (*Department, error)
	GetDepartmentByNumber(int) (*Department, error)
}

// UserModel handles database operations for User
type UserDBModel struct {
	DB *gorm.DB
}

// NewUserModel creates a new instance of UserModel
func NewUserDBModel(db *gorm.DB) *UserDBModel {
	return &UserDBModel{
		DB: db,
	}
}

// CreateUser creates a new user.
func (as *UserDBModel) CreateUser(user *Users) error {
	return as.DB.Create(user).Error
}

// GetUserByID retrieves a user by its ID.
func (as *UserDBModel) GetUserByID(id uint) (*Users, error) {
	var user Users
	err := as.DB.Where("id = ?", id).First(&user).Error
	return &user, err
}

// UpdateUser updates the details of an existing user.
func (as *UserDBModel) UpdateUser(user *Users) error {
	if err := as.DB.Save(user).Error; err != nil {
		return err
	}
	return nil
}

// DeleteUser deletes a user from the database.
func (as *UserDBModel) DeleteUser(id uint) error {
	if err := as.DB.Delete(&Users{}, id).Error; err != nil {
		return err
	}
	return nil
}

// GetAllUsers retrieves all users from the database.
func (as *UserDBModel) GetAllUsers() (*[]Users, error) {
	var users []Users
	err := as.DB.Find(&users).Error
	return &users, err
}
