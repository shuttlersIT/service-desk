// backend/models/users.go

package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Users2 struct {
	gorm.Model
	ID                     uint                  `gorm:"primaryKey" json:"id"`
	FirstName              string                `gorm:"size:255;not null" json:"first_name"`
	LastName               string                `gorm:"size:255;not null" json:"last_name"`
	Email                  string                `gorm:"size:255;not null;unique" json:"email"`
	Phone                  *string               `gorm:"size:20;null" json:"phone,omitempty"`
	PositionID             uint                  `gorm:"index" json:"position_id,omitempty"`
	DepartmentID           uint                  `gorm:"index" json:"department_id,omitempty"`
	IsActive               bool                  `gorm:"default:true" json:"is_active"`
	Roles                  []Role                `gorm:"many2many:user_roles;" json:"roles,omitempty"`
	Projects               []Project             `gorm:"many2many:project_members;" json:"projects,omitempty"`
	Position               Position              `gorm:"foreignKey:PositionID" json:"-"`
	Department             Department            `gorm:"foreignKey:DepartmentID" json:"-"`
	ProfilePic             string                `gorm:"size:255" json:"profile_pic,omitempty"`
	Credentials            UsersLoginCredentials `gorm:"embedded" json:"username,omitempty"` // Excluded from JSON responses
	ResetPasswordRequestID *uint                 `json:"reset_password_request_id,omitempty"`
	Processed              bool                  `json:"processed,omitempty"`
	LastLoginAt            *time.Time            `json:"last_login_at,omitempty"`
}

type Users struct {
	gorm.Model
	ID           uint                  `gorm:"primaryKey" json:"id"`
	FirstName    string                `gorm:"size:255;not null" json:"first_name" binding:"required"`
	LastName     string                `gorm:"size:255;not null" json:"last_name" binding:"required"`
	Email        string                `gorm:"size:255;not null;unique" json:"email" binding:"required,email"`
	Phone        *string               `gorm:"size:20;null" json:"phone,omitempty" binding:"omitempty,e164"`
	PositionID   uint                  `json:"position_id,omitempty" gorm:"type:int unsigned"`
	DepartmentID uint                  `json:"department_id,omitempty" gorm:"type:int unsigned"`
	IsActive     bool                  `gorm:"default:true" json:"is_active"`
	Roles        []Role                `gorm:"many2many:user_roles;" json:"roles,omitempty"`
	ProfilePic   string                `gorm:"size:255" json:"profile_pic,omitempty"`
	Credentials  UsersLoginCredentials `gorm:"embedded" json:"username,omitempty"` // Excluded from JSON responses
	LastLoginAt  *time.Time            `gorm:"type:datetime" json:"last_login_at,omitempty"`
	DeletedAt    *gorm.DeletedAt       `gorm:"index" json:"deleted_at,omitempty"`
}

func (Users) TableName() string {
	return "users"
}

type UserProfile struct {
	UserID          uint      `gorm:"primaryKey;autoIncrement:false" json:"user_id"`
	Bio             string    `gorm:"type:text" json:"bio,omitempty"`
	AvatarURL       string    `gorm:"type:text" json:"avatar_url,omitempty"`
	Preferences     string    `gorm:"type:text" json:"preferences,omitempty"`      // Stored as JSON
	PrivacySettings string    `gorm:"type:text" json:"privacy_settings,omitempty"` // Stored as JSON
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (UserProfile) TableName() string {
	return "user_profiles"
}

type Department struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	Name        string          `gorm:"size:255;not null;unique" json:"name"`
	Description string          `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Department) TableName() string {
	return "departments"
}

type Position struct {
	ID           uint            `gorm:"primaryKey" json:"id"`
	Name         string          `gorm:"size:255;not null;unique" json:"name"`
	Description  string          `gorm:"type:text" json:"description,omitempty"`
	DepartmentID uint            `json:"department_id" gorm:"index;not null"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	DeletedAt    *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Position) TableName() string {
	return "positions"
}

type UserRole struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	UserID    uint            `gorm:"not null index;type:int unsigned" json:"user_id"`
	RoleID    uint            `gorm:"not null index;type:int unsigned" json:"role_id"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (UserRole) TableName() string {
	return "user_roles"
}

type ProjectAssignment struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	UserID    uint            `json:"user_id" gorm:"index;not null"`
	ProjectID uint            `json:"project_id" gorm:"index;not null"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName sets the table name for the Project model.
func (ProjectAssignment) TableName() string {
	return "project_assignment"
}

type Activity struct {
	ID           uint            `gorm:"primaryKey" json:"id"`
	UserID       uint            `json:"user_id" gorm:"index;not null"`
	Description  string          `gorm:"type:text" json:"description"`
	ActivityType string          `gorm:"type:varchar(100);not null" json:"activity_type"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	DeletedAt    *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Activity) TableName() string {
	return "activities"
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
	DB             *gorm.DB
	log            Logger
	EventPublisher EventPublisherImpl
}

// NewUserModel creates a new instance of UserModel
func NewUserDBModel(db *gorm.DB, log Logger, eventPublisher EventPublisherImpl) *UserDBModel {
	return &UserDBModel{
		DB:             db,
		log:            log,
		EventPublisher: eventPublisher,
	}
}

// CreateUser adds a new user to the database with transactional integrity.
func (db *UserDBModel) CreateUser(user *Users) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		return nil
	})
}

// UpdateUserDetails updates specified fields of a user's information.
func (db *UserDBModel) UpdateUserDetails(userID uint, updates map[string]interface{}) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Users{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
			return err
		}
		return nil
	})
}

// GetUserByID retrieves a user by its ID.
func (as *UserDBModel) GetUserByID(id uint) (*Users, error) {
	var user Users
	err := as.DB.Where("id = ?", id).First(&user).Error
	return &user, err
}

// UpdateUser updates the details of an existing user with transactional integrity.
func (db *UserDBModel) UpdateUser(user *Users) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(user).Error; err != nil {
			return err
		}
		return nil
	})
}

// DeleteUser deletes a user from the database with transactional integrity.
func (db *UserDBModel) DeleteUser(id uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&Users{}, id).Error; err != nil {
			return err
		}
		return nil
	})
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

// CreatePosition creates a new position with transactional integrity.
func (db *UserDBModel) CreatePosition(position *Position) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(position).Error; err != nil {
			return err
		}
		return nil
	})
}

// UpdatePosition updates the details of an existing position with transactional integrity.
func (db *UserDBModel) UpdatePosition(position *Position) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(position).Error; err != nil {
			return err
		}
		return nil
	})
}

// DeletePosition deletes a position from the database with transactional integrity.
func (db *UserDBModel) DeletePosition(id uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&Position{}, id).Error; err != nil {
			return err
		}
		return nil
	})
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

// CreateDepartment creates a new department with transactional integrity.
func (db *UserDBModel) CreateDepartment(department *Department) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(department).Error; err != nil {
			return err
		}
		return nil
	})
}

// UpdateDepartment updates the details of an existing department with transactional integrity.
func (db *UserDBModel) UpdateDepartment(department *Department) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(department).Error; err != nil {
			return err
		}
		return nil
	})
}

// DeleteDepartment deletes a department from the database with transactional integrity.
func (db *UserDBModel) DeleteDepartment(id uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&Department{}, id).Error; err != nil {
			return err
		}
		return nil
	})
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

// CreateUserProfile creates a profile for a user with transactional integrity.
func (db *UserDBModel) CreateUserProfile(userProfile *Users) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(userProfile).Error; err != nil {
			return err
		}
		return nil
	})
}

// GetUserProfileByUserID retrieves a user's profile by their ID without needing transactions as it's a read operation.
func (db *UserDBModel) GetUserProfileByUserID(userID uint) (*Users, error) {
	var userProfile Users
	err := db.DB.Where("user_id = ?", userID).First(&userProfile).Error
	return &userProfile, err
}

// UpdateUserProfile updates the details of a user's profile with transactional integrity.
func (db *UserDBModel) UpdateUserProfile(userProfile *Users) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(userProfile).Error; err != nil {
			return err
		}
		return nil
	})
}

// DeleteUserProfile deletes a user's profile with transactional integrity.
func (db *UserDBModel) DeleteUserProfile(userID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&Users{}, userID).Error; err != nil {
			return err
		}
		return nil
	})
}

// AssignRoleToUser assigns a specific role to a user with transactional integrity.
func (db *UserDBModel) AssignRoleToUser(userID, roleID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		userRole := UserRole{UserID: userID, RoleID: roleID}
		if err := tx.Create(&userRole).Error; err != nil {
			return err
		}
		return nil
	})
}

// RevokeRoleFromUser removes a role assignment from a user with transactional integrity.
func (db *UserDBModel) RevokeRoleFromUser(userID, roleID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&UserRole{}).Error; err != nil {
			return err
		}
		return nil
	})
}

// AssignRolesToUsers assigns multiple roles to multiple users in a batch operation with transactional integrity.
func (db *UserDBModel) AssignRolesToUsers(userIDs, roleIDs []uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		for _, userID := range userIDs {
			for _, roleID := range roleIDs {
				userRole := UserRole{UserID: userID, RoleID: roleID}
				if err := tx.Create(&userRole).Error; err != nil {
					// If the role is already assigned to the user, consider handling or ignoring the error based on your application logic.
					return err
				}
			}
		}
		return nil
	})
}

// ProcessUsersAndRoles performs a complex operation involving users and roles.
// ProcessUsersAndRoles performs a complex operation involving users and roles.
/*func (db *UserDBModel) ProcessUsersAndRoles() error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var users []Users
		// Example: Fetch users that meet certain criteria.
		if err := tx.Where("active = ?", true).Find(&users).Error; err != nil {
			// Explicitly rollback the transaction on error (optional with GORM as it handles rollback on error automatically).
			tx.Rollback() // This is optional as GORM automatically rolls back on error.
			return fmt.Errorf("failed to fetch active users: %w", err)
		}

		// Process each user (pseudo-code).
		for _, user := range users {
			// Example processing logic.
			user.Processed = true
			if err := tx.Save(&user).Error; err != nil {
				tx.Rollback() // This is optional.
				return fmt.Errorf("failed to save processed user [%d]: %w", user.ID, err)
			}
		}

		// Explicitly commit the transaction (optional with GORM as it automatically commits if no errors occurred).
		// This can be omitted because if the function returns nil, GORM commits the transaction.
		return nil // If everything went well, GORM commits the transaction.
	})
}*/

// UpdateMultipleUserProfiles updates profiles for multiple users, handling errors individually.
func (db *UserDBModel) UpdateMultipleUserProfiles(userProfiles []Users) ([]error, error) {
	var errors []error

	for _, userProfile := range userProfiles {
		err := db.DB.Transaction(func(tx *gorm.DB) error {
			return tx.Save(&userProfile).Error
		})
		if err != nil {
			errors = append(errors, err)
			// Optionally, log the error or handle it as per your application's requirements.
		}
	}

	if len(errors) > 0 {
		return errors, fmt.Errorf("some user profile updates failed")
	}

	return nil, nil
}

// AssignUserToProject assigns a user to a project and updates related records accordingly.
func (db *UserDBModel) AssignUserToProject(userID, projectID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		// Step 1: Create a new project assignment record.
		assignment := ProjectAssignment{UserID: userID, ProjectID: projectID}
		if err := tx.Create(&assignment).Error; err != nil {
			return err
		}

		// Step 2: Update related records, such as project member counts or user assignment counts.
		if err := tx.Model(&Project{}).Where("id = ?", projectID).UpdateColumn("member_count", gorm.Expr("member_count + ?", 1)).Error; err != nil {
			return err
		}

		// Step 3: Optionally, perform other updates or checks, such as verifying user availability or project capacity.
		return nil
	})
}

// GetUsersWithActiveProjects retrieves users who have active projects, demonstrating a join operation.
func (db *UserDBModel) GetUsersWithActiveProjects() ([]Users, error) {
	var users []Users
	err := db.DB.Joins("JOIN project_assignments ON project_assignments.user_id = users.id").
		Joins("JOIN projects ON projects.id = project_assignments.project_id AND projects.status = ?", "active").
		Distinct("users.*").
		Find(&users).Error
	return users, err
}

// ActivateUsers activates a list of users by their IDs with transactional integrity.
func (db *UserDBModel) ActivateUsers(userIDs []uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Users{}).Where("id IN ?", userIDs).Update("active", true).Error; err != nil {
			return err
		}
		return nil
	})
}

// DeactivateUsers deactivates a list of users by their IDs with transactional integrity.
func (db *UserDBModel) DeactivateUsers(userIDs []uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Users{}).Where("id IN ?", userIDs).Update("active", false).Error; err != nil {
			return err
		}
		return nil
	})
}

/*
// PublishUserActivity utilizes an event-driven approach to handle user activities, aiming for eventual consistency.
func (db *UserDBModel) PublishUserActivity(userID uint, activity UserActivityLog) error {
	// Record the activity in the local database first.
	err := db.DB.Create(&activity).Error
	if err != nil {
		return err
	}

	// Publish an event to a message queue for further processing, e.g., updating a search index or analytics.
	err = messageQueue.Publish("user-activity", activity)
	if err != nil {
		// Handle message queue errors. Depending on the use case, you might choose to log the error and proceed,
		// ensuring the database operation isn't rolled back due to issues outside the core transaction.
		return err
	}

	return nil
}

// MigrateUserData migrates user data from one schema to another with transactional integrity.
func (db *UserDBModel) MigrateUserData() error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var users []Users
		if err := tx.Find(&users).Error; err != nil {
			return err
		}

		for _, oldUser := range users {
			newUser := transformUser(oldUser) // Assume transformUser converts an old user record to a new schema.
			if err := tx.Create(&newUser).Error; err != nil {
				// Handle error, potentially logging failed conversions without halting the entire migration.
				return err
			}
		}
		return nil
	})
}

// PromoteUser checks user criteria and promotes the user if they meet certain conditions.
func (db *UserDBModel) PromoteUser(userID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var user Users
		// Fetch user with their current status and metrics.
		if err := tx.Where("id = ?", userID).First(&user).Error; err != nil {
			return err
		}

		// Check if user meets promotion criteria.
		if user.MeetsPromotionCriteria() { // Assuming MeetsPromotionCriteria is a method that checks user's eligibility for promotion.
			// Apply promotion logic.
			if err := tx.Model(&user).Update("status", "promoted").Error; err != nil {
				return err
			}
			// Additional logic for promotion, e.g., notifying the user, updating related records, etc.
		}

		return nil
	})
}
*/

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
