// backend/models/users.go

package models

import (
	"time"

	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	ID                     uint                  `gorm:"primaryKey" json:"user_id"`
	FirstName              string                `json:"first_name" binding:"required"`
	LastName               string                `json:"last_name" binding:"required"`
	Email                  string                `json:"staff_email" binding:"required,email"`
	Credentials            UsersLoginCredentials `json:"user_credentials" gorm:"foreignKey:UserID"`
	Phone                  string                `json:"phoneNumber" binding:"required,e164"`
	Position               Position              `json:"position_id" gorm:"embedded"`
	Department             Department            `json:"department_id" gorm:"embedded"`
	CreatedAt              time.Time             `json:"created_at"`
	UpdatedAt              time.Time             `json:"updated_at"`
	Asset                  []AssetAssignment     `json:"asset_assignment" gorm:"foreignKey:UserID"`
	RoleBase               RoleBase              `json:"role_base" gorm:"embedded"`
	ResetPasswordRequestID uint                  `json:"reset_password_reset" gorm:"embedded"`
	DeletedAt              time.Time             `json:"deleted_at"`
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
	DeletedAt    time.Time `json:"deleted_at"`
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
	DeletedAt      time.Time `json:"deleted_at"`
}

// TableName sets the table name for the Department model.
func (Department) TableName() string {
	return "department"
}

type UserStorage interface {
	CreateUser(*Users) error
	DeleteUser(int) error
	UpdateUser(*Users) error
	GetUsers() ([]*Users, error)
	GetUserByID(int) (*Users, error)
	GetUserByNumber(int) (*Users, error)
}

type PositionStorage interface {
	CreatePosition(*Position) error
	DeletePosition(int) error
	UpdatePosition(*Position) error
	GetPosition() ([]*Position, error)
	GetPositionByID(int) (*Position, error)
	GetPositionByNumber(int) (*Position, error)
}

type DepartmentStorage interface {
	CreateDepartment(*Department) error
	DeleteDepartment(int) error
	UpdateDepartment(*Department) error
	GetDepartments() ([]*Department, error)
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
func (as *UserDBModel) CreateUser(user *Users) (*Users, error) {
	result := as.DB.Create(user)

	return as.GetUserByID(uint(result.RowsAffected))
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
func (as *UserDBModel) GetAllUsers() ([]*Users, error) {
	var users []*Users
	err := as.DB.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, err
}

// GetUserByNumber retrieves an user by their user number.
func (as *UserDBModel) GetUsersByNumber(userNumber int) (*Users, error) {
	var user Users
	err := as.DB.Where("user_id = ?", userNumber).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// backend/models/users.go

// ...

// GetUserByEmail retrieves a user by their email address.
func (as *UserDBModel) GetUserByEmail(email string) (*Users, error) {
	var user Users
	err := as.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUsersByPosition retrieves users by their position.
func (as *UserDBModel) GetUsersByPosition(positionID int) ([]*Users, error) {
	var users []*Users
	err := as.DB.Where("position_id = ?", positionID).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetUsersByDepartment retrieves users by their department.
func (as *UserDBModel) GetUsersByDepartment(departmentID int) ([]*Users, error) {
	var users []*Users
	err := as.DB.Where("department_id = ?", departmentID).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// CreatePosition creates a new position.
func (as *UserDBModel) CreatePosition(position *Position) error {
	return as.DB.Create(position).Error
}

// UpdatePosition updates the details of an existing position.
func (as *UserDBModel) UpdatePosition(position *Position) error {
	return as.DB.Save(position).Error
}

// DeletePosition deletes a position from the database.
func (as *UserDBModel) DeletePosition(id uint) error {
	if err := as.DB.Delete(&Position{}, id).Error; err != nil {
		return err
	}
	return nil
}

// GetAllPositions retrieves all positions from the database.
func (as *UserDBModel) GetAllPositions() ([]*Position, error) {
	var positions []*Position
	err := as.DB.Find(&positions).Error
	if err != nil {
		return nil, err
	}
	return positions, nil
}

// GetPositionByID retrieves a position by its ID.
func (as *UserDBModel) GetPositionByID(id uint) (*Position, error) {
	var position Position
	err := as.DB.Where("id = ?", id).First(&position).Error
	if err != nil {
		return nil, err
	}
	return &position, nil
}

// GetPositionByNumber retrieves a position by its position number.
func (as *UserDBModel) GetPositionByNumber(positionNumber int) (*Position, error) {
	var position Position
	err := as.DB.Where("position_id = ?", positionNumber).First(&position).Error
	if err != nil {
		return nil, err
	}
	return &position, nil
}

// CreateDepartment creates a new department.
func (as *UserDBModel) CreateDepartment(department *Department) error {
	return as.DB.Create(department).Error
}

// UpdateDepartment updates the details of an existing department.
func (as *UserDBModel) UpdateDepartment(department *Department) error {
	return as.DB.Save(department).Error
}

// DeleteDepartment deletes a department from the database.
func (as *UserDBModel) DeleteDepartment(id uint) error {
	if err := as.DB.Delete(&Department{}, id).Error; err != nil {
		return err
	}
	return nil
}

//GetPositions retrieves all the position from the datatbase

// GetAllDepartments retrieves all departments from the database.
func (as *UserDBModel) GetAllDepartments() ([]*Department, error) {
	var departments []*Department
	err := as.DB.Find(&departments).Error
	if err != nil {
		return nil, err
	}
	return departments, nil
}

// GetDepartmentByID retrieves a department by its ID.
func (as *UserDBModel) GetDepartmentByID(id uint) (*Department, error) {
	var department Department
	err := as.DB.Where("id = ?", id).First(&department).Error
	if err != nil {
		return nil, err
	}
	return &department, nil
}

// GetDepartmentByNumber retrieves a department by its department number.
func (as *UserDBModel) GetDepartmentByNumber(departmentNumber int) (*Department, error) {
	var department Department
	err := as.DB.Where("department_id = ?", departmentNumber).First(&department).Error
	if err != nil {
		return nil, err
	}
	return &department, nil
}

///////////////////////////////////////////////////////// ASSETS DIRECT /////////////////////////////////////////////////////////////

// CreateAssetAssignment creates a new asset assignment for a user.
func (as *UserDBModel) CreateAssetAssignment(assetAssignment *AssetAssignment) error {
	return as.DB.Create(assetAssignment).Error
}

// UpdateAssetAssignment updates the details of an existing asset assignment.
func (as *UserDBModel) UpdateAssetAssignment(assetAssignment *AssetAssignment) error {
	return as.DB.Save(assetAssignment).Error
}

// DeleteAssetAssignment deletes an asset assignment from the database.
func (as *UserDBModel) DeleteAssetAssignment(id uint) error {
	if err := as.DB.Delete(&AssetAssignment{}, id).Error; err != nil {
		return err
	}
	return nil
}

// GetAssetAssignmentsByUser retrieves asset assignments for a user by their user ID.
func (as *UserDBModel) GetAssetAssignmentsByUser(userID uint) ([]*AssetAssignment, error) {
	var assetAssignments []*AssetAssignment
	err := as.DB.Where("user_id = ?", userID).Find(&assetAssignments).Error
	if err != nil {
		return nil, err
	}
	return assetAssignments, nil
}

// GetAssetAssignmentsByAsset retrieves asset assignments for an asset by its asset ID.
func (as *UserDBModel) GetAssetAssignmentsByAsset(assetID uint) ([]*AssetAssignment, error) {
	var assetAssignments []*AssetAssignment
	err := as.DB.Where("asset_id = ?", assetID).Find(&assetAssignments).Error
	if err != nil {
		return nil, err
	}
	return assetAssignments, nil
}

// GetAssetAssignmentByID retrieves an asset assignment by its ID.
func (as *UserDBModel) GetAssetAssignmentByID(id uint) (*AssetAssignment, error) {
	var assetAssignment AssetAssignment
	err := as.DB.Where("id = ?", id).First(&assetAssignment).Error
	if err != nil {
		return nil, err
	}
	return &assetAssignment, nil
}

// GetAssetAssignmentByNumber retrieves an asset assignment by its asset assignment number.
func (as *UserDBModel) GetAssetAssignmentByNumber(assetAssignmentNumber int) (*AssetAssignment, error) {
	var assetAssignment AssetAssignment
	err := as.DB.Where("asset_assignment_id = ?", assetAssignmentNumber).First(&assetAssignment).Error
	if err != nil {
		return nil, err
	}
	return &assetAssignment, nil
}

// GetAssetByID retrieves an asset by its ID.
func (as *UserDBModel) GetAssetByID(id uint) (*Assets, error) {
	var asset Assets
	err := as.DB.Where("id = ?", id).First(&asset).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

// GetAssetByNumber retrieves an asset by its asset number.
func (as *UserDBModel) GetAssetByNumber(assetNumber int) (*Assets, error) {
	var asset Assets
	err := as.DB.Where("asset_id = ?", assetNumber).First(&asset).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

// CreateUserProfile creates a profile for a user.
func (as *UserDBModel) CreateUserProfile(user *Users) error {
	return as.DB.Create(user).Error
}

// GetUserProfileByUserID retrieves a user's profile by their ID.
func (as *UserDBModel) GetUserProfileByUserID(userID uint) (*Users, error) {
	var user Users
	err := as.DB.Where("user_id = ?", userID).First(&user).Error
	return &user, err
}

// UpdateUserProfile updates the details of a user's profile.
func (as *UserDBModel) UpdateUserProfile(user *Users) error {
	return as.DB.Save(&user).Error
}

// DeleteUserProfile deletes a user's profile.
func (as *UserDBModel) DeleteUserProfile(user *Users) error {
	return as.DB.Delete(&user).Error
}

/*
// Role based access

// CreateRole creates a new role.
func (as *UserDBModel) CreateRole(role *Role) error {
    return as.DB.Create(role).Error
}

// UpdateRole updates an existing role.
func (as *UserDBModel) UpdateRole(role *Role) error {
    return as.DB.Save(role).Error
}

// DeleteRole deletes a role by ID.
func (as *UserDBModel) DeleteRole(roleID uint) error {
    return as.DB.Delete(&Role{}, roleID).Error
}

// GetRoleByID retrieves a role by its ID.
func (as *UserDBModel) GetRoleByID(roleID uint) (*Role, error) {
    var role Role
    err := as.DB.Where("id = ?", roleID).First(&role).Error
    return &role, err
}

// GetAllRoles retrieves all roles.
func (as *UserDBModel) GetAllRoles() ([]*Role, error) {
    var roles []*Role
    err := as.DB.Find(&roles).Error
    return roles, err
}
*/
