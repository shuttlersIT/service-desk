// backend/models/auth.go

package models

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AgentLoginCredentialsStorage interface {
	CreateAgentLoginCredentials(*AgentLoginCredentials) error
	DeleteAgentLoginCredentials(int) error
	UpdateAgentLoginCredentials(*AgentLoginCredentials) error
	GetAgentLoginCredentials() ([]*AgentLoginCredentials, error)
	GetAgentLoginCredentialsByID(int) (*AgentLoginCredentials, error)
	GetAgentLoginCredentialsByNumber(int) (*AgentLoginCredentials, error)
}

type AgentLoginCredentials struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey" json:"_"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	AgentID   uint      `json:"agent_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName sets the table name for the Agent model.
func (AgentLoginCredentials) TableName() string {
	return "agentLoginDetails"
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type UsersLoginCredentials struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey" json:"_"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName sets the table name for the UsersLoginCredentials model.
func (UsersLoginCredentials) TableName() string {
	return "usersLoginCredentials"
}

type UserLoginCredentialsStorage interface {
	CreateUserLoginCredentials(*UsersLoginCredentials) error
	DeleteUserLoginCredentials(int) error
	UpdateUserLoginCredentials(*UsersLoginCredentials) error
	GetUsersLoginCredentials() ([]*UsersLoginCredentials, error)
	GetUserUserLoginCredentialsByID(int) (*UsersLoginCredentials, error)
	GetUserLoginCredentialsByNumber(int) (*UsersLoginCredentials, error)
}

// AuthModel handles database operations for Auth
type AuthDBModel struct {
	DB *gorm.DB
}

// NewUserModel creates a new instance of UserModel
func NewAuthDBModel(db *gorm.DB) *AuthDBModel {
	return &AuthDBModel{
		DB: db,
	}
}

/////////////////////////////////////////////// USERS //////////////////////////////////////////////////////////

// CreateUser creates a new user.
func (as *AuthDBModel) CreateUserCredentials(userCredentials *UsersLoginCredentials) error {
	return as.DB.Create(userCredentials).Error
}

// GetUserByID retrieves a user by its ID.
func (as *AuthDBModel) GetUserCredentialsByID(id uint) (*UsersLoginCredentials, error) {
	var userCredentials UsersLoginCredentials
	err := as.DB.Where("id = ?", id).First(&userCredentials).Error
	return &userCredentials, err
}

// UpdateUser updates the details of an existing user.
func (as *AuthDBModel) UpdateUserCredentials(userCredentials *UsersLoginCredentials) error {
	return as.DB.Save(userCredentials).Error
}

// DeleteUser deletes a user from the database.
func (as *AuthDBModel) DeleteUserCredentials(id uint) error {
	return as.DB.Delete(&UsersLoginCredentials{}, id).Error
}

// GetAllUsers retrieves all users from the database.
func (as *AuthDBModel) GetAllUserCreds() ([]*UsersLoginCredentials, error) {
	var usersCredentials []*UsersLoginCredentials
	err := as.DB.Find(&usersCredentials).Error
	return usersCredentials, err
}

type LoginInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterModel struct {
	a *UserDBModel
	b *AuthDBModel
}

// NewUserModel creates a new instance of UserModel
func NewRegisterModel(a *UserDBModel, b *AuthDBModel) *RegisterModel {
	return &RegisterModel{
		a: a,
		b: b,
	}
}

// Register New User
func (ab *RegisterModel) Registration(user *Users) (*Users, error) {
	// Hash the user's password before storing it in the database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Credentials.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password")
	}
	user.Credentials.Password = string(hashedPassword)
	newUser, erro := ab.a.CreateUser(user)
	if erro != nil {
		return nil, fmt.Errorf("failed to create users")
	}
	er := ab.b.CreateUserCredentials(&newUser.Credentials)
	if er != nil {
		return nil, fmt.Errorf("failed to create users credentials")
	}
	e := ab.a.UpdateUser(newUser)
	if e != nil {
		return nil, fmt.Errorf("unable to update new users credentials")
	}

	return newUser, nil
}

// Login User
func (a *AuthDBModel) Login(login *LoginInfo) (*Users, error) {
	loginInfo := login
	var user Users
	if err := a.DB.Where("email = ?", loginInfo.Email).First(&user).Error; err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Credentials.Password), []byte(loginInfo.Password)); err != nil {
		return nil, err
	}
	return &user, nil
}

// ResetUserPassword retrieves all users from the database.
func (as *AuthDBModel) ResetUserPassword(uint, string) ([]*UsersLoginCredentials, error) {
	var usersCredentials []*UsersLoginCredentials
	err := as.DB.Find(&usersCredentials).Error
	return usersCredentials, err
}

/////////////////////////////////////////////// AGENTS //////////////////////////////////////////////////////////

// CreateUser creates a new user.
func (as *AuthDBModel) CreateAgentCredentials(agentCredentials *AgentLoginCredentials) error {
	return as.DB.Create(agentCredentials).Error
}

// GetUserByID retrieves a user by its ID.
func (as *AuthDBModel) GetAgentCredentialsByID(id uint) (*AgentLoginCredentials, error) {
	var agentCredentials AgentLoginCredentials
	err := as.DB.Where("id = ?", id).First(&agentCredentials).Error
	return &agentCredentials, err
}

// UpdateUser updates the details of an existing user.
func (as *AuthDBModel) UpdateAgentCredentials(agentCredentials *AgentLoginCredentials) error {
	return as.DB.Save(agentCredentials).Error
}

// DeleteUser deletes a user from the database.
func (as *AuthDBModel) DeleteAgentCredentials(id uint) error {
	return as.DB.Delete(&AgentLoginCredentials{}, id).Error
}

// GetAllUsers retrieves all users from the database.
func (as *AuthDBModel) GetAllAgentCreds() ([]*AgentLoginCredentials, error) {
	var agentCredentials []*AgentLoginCredentials
	err := as.DB.Find(&agentCredentials).Error
	return agentCredentials, err
}
