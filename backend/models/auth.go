// backend/models/auth.go

package models

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type APIToken struct {
	gorm.Model
	UserID      uint       `gorm:"not null;index" json:"user_id"`
	User        Users      `gorm:"foreignKey:UserID" json:"-"`
	Token       string     `gorm:"size:255;not null;unique" json:"token"`
	ExpiresAt   time.Time  `gorm:"type:datetime" json:"expires_at"`
	Description string     `gorm:"size:255" json:"description"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
}

type AgentLoginCredentialsStorage interface {
	Create(credentials *AgentLoginCredentials) error
	Delete(id uint) error
	Update(credentials *AgentLoginCredentials) error
	FindByID(id uint) (*AgentLoginCredentials, error)
}

type AgentLoginCredentials struct {
	gorm.Model
	ID          uint       `gorm:"primaryKey" json:"_"`
	Username    string     `json:"username"`
	Password    string     `json:"password"`
	AgentID     uint       `json:"agent_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   time.Time  `json:"deleted_at"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
}

// TableName sets the table name for the Agent model.
func (AgentLoginCredentials) TableName() string {
	return "agent_login_credentials"
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
	return "users_login_credentials"
}

type UserLoginCredentialsStorage interface {
	Create(credentials *UsersLoginCredentials) error
	Delete(id uint) error
	Update(credentials *UsersLoginCredentials) error
	FindByID(id uint) (*UsersLoginCredentials, error)
}

// AuthModel handles database operations for Auth
type AuthDBModel struct {
	DB           *gorm.DB
	UserDBModel  *UserDBModel
	AgentDBModel *AgentDBModel
}

// NewUserModel creates a new instance of UserModel
func NewAuthDBModel(db *gorm.DB, userDBModel *UserDBModel, agentDBModel *AgentDBModel) *AuthDBModel {
	return &AuthDBModel{
		DB:           db,
		UserDBModel:  userDBModel,
		AgentDBModel: agentDBModel,
	}
}

/////////////////////////////////////////////// USERS //////////////////////////////////////////////////////////

// CreateUserCredentials creates new user credentials.
func (as *AuthDBModel) CreateUserCredentials(userCredentials *UsersLoginCredentials) error {
	return as.DB.Create(userCredentials).Error
}

// GetUserCredentialsByID retrieves user credentials by their ID.
func (as *AuthDBModel) GetUserCredentialsByID(id uint) (*UsersLoginCredentials, error) {
	var userCredentials UsersLoginCredentials
	err := as.DB.Where("id = ?", id).First(&userCredentials).Error
	return &userCredentials, err
}

// UpdateUserCredentials updates user credentials.
func (as *AuthDBModel) UpdateUserCredentials(userCredentials *UsersLoginCredentials) error {
	return as.DB.Save(userCredentials).Error
}

// DeleteUserCredentials deletes user credentials by their ID.
func (as *AuthDBModel) DeleteUserCredentials(id uint) error {
	return as.DB.Delete(&UsersLoginCredentials{}, id).Error
}

type LoginInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterModel struct {
	a *UserDBModel
	b *AuthDBModel
	c *AgentDBModel
}

// NewUserModel creates a new instance of UserModel
func NewRegisterModel(a *UserDBModel, b *AuthDBModel) *RegisterModel {
	return &RegisterModel{
		a: a,
		b: b,
	}
}

// Registration registers a new user.
func (ab *RegisterModel) Registration(user *Users) (*Users, error) {
	// Hash the user's password before storing it in the database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.HashedPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password")
	}
	user.HashedPassword = string(hashedPassword)

	newUser, erro := ab.a.CreateUser(user)
	if erro != nil {
		return nil, fmt.Errorf("failed to create user")
	}

	er := ab.b.CreateUserCredentials(&UsersLoginCredentials{
		Username: user.Email, // You can customize the username field
		Password: user.HashedPassword,
		UserID:   newUser.ID,
	})
	if er != nil {
		return nil, fmt.Errorf("failed to create user credentials")
	}

	e := ab.a.UpdateUser(newUser)
	if e != nil {
		return nil, fmt.Errorf("unable to update user credentials")
	}

	return newUser, nil
}

// AgentRegistration registers a new agent.
func (ab *RegisterModel) AgentRegistration(agent *Agents) (*Agents, error) {
	// Hash the agent's password before storing it in the database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(agent.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password")
	}
	agent.PasswordHash = string(hashedPassword)

	newAgent, erro := ab.c.CreateAgent(agent)
	if erro != nil {
		return nil, fmt.Errorf("failed to create agent")
	}

	er := ab.b.CreateAgentCredentials(&AgentLoginCredentials{
		Username: agent.Email, // You can customize the username field
		Password: agent.PasswordHash,
		AgentID:  newAgent.ID,
	})
	if er != nil {
		return nil, fmt.Errorf("failed to create agent credentials")
	}

	e := ab.c.UpdateAgent(newAgent)
	if e != nil {
		return nil, fmt.Errorf("unable to update agent credentials")
	}

	return newAgent, nil
}

// LoginUser authenticates a user based on email and password.
func (a *AuthDBModel) LoginUser(login *LoginInfo) (*Users, error) {
	loginInfo := login
	var user Users

	if err := a.DB.Where("email = ?", loginInfo.Email).First(&user).Error; err != nil {
		return nil, err // User not found
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(loginInfo.Password)); err != nil {
		return nil, err // Password does not match
	}

	return &user, nil
}

// LoginAgent authenticates an agent based on email and password.
func (a *AuthDBModel) LoginAgent(login *LoginInfo) (*Agents, error) {
	loginInfo := login
	var agent Agents

	if err := a.DB.Where("email = ?", loginInfo.Email).First(&agent).Error; err != nil {
		return nil, err // Agent not found
	}

	if err := bcrypt.CompareHashAndPassword([]byte(agent.PasswordHash), []byte(loginInfo.Password)); err != nil {
		return nil, err // Password does not match
	}

	return &agent, nil
}

// ResetUserPassword retrieves all users from the database.
func (as *AuthDBModel) ResetUserPassword(uint, string) ([]*UsersLoginCredentials, error) {
	var usersCredentials []*UsersLoginCredentials
	err := as.DB.Find(&usersCredentials).Error
	return usersCredentials, err
}

/////////////////////////////////////////////// AGENTS //////////////////////////////////////////////////////////

// ResetAgentPassword retrieves all agents from the database.
func (as *AuthDBModel) ResetAgentPassword() ([]*AgentLoginCredentials, error) {
	var agentsCredentials []*AgentLoginCredentials
	err := as.DB.Find(&agentsCredentials).Error
	return agentsCredentials, err
}

// CreateAgent creates new agent credentials.
func (as *AuthDBModel) CreateAgentCredentials(agentCredentials *AgentLoginCredentials) error {
	return as.DB.Create(agentCredentials).Error
}

// GetAgentCredentialsByID retrieves agent credentials by their ID.
func (as *AuthDBModel) GetAgentCredentialsByID(id uint) (*AgentLoginCredentials, error) {
	var agentCredentials AgentLoginCredentials
	err := as.DB.Where("id = ?", id).First(&agentCredentials).Error
	return &agentCredentials, err
}

// UpdateAgentCredentials updates agent credentials.
func (as *AuthDBModel) UpdateAgentCredentials(agentCredentials *AgentLoginCredentials) error {
	return as.DB.Save(agentCredentials).Error
}

// DeleteAgentCredentials deletes agent credentials by their ID.
func (as *AuthDBModel) DeleteAgentCredentials(id uint) error {
	return as.DB.Delete(&AgentLoginCredentials{}, id).Error
}

// GetAllUsers retrieves all users from the database.
func (as *AuthDBModel) GetAllAgentCreds() ([]*AgentLoginCredentials, error) {
	var agentCredentials []*AgentLoginCredentials
	err := as.DB.Find(&agentCredentials).Error
	return agentCredentials, err
}

//////////////////////////////////////////////////////////////////////////////////
//Password Requests

// Define a model for storing password reset requests

type PasswordResetRequest struct {
	gorm.Model
	UserID    uint      `json:"user_id"`
	RequestID uint      `json:"request_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// CreatePasswordResetRequest creates a new password reset request record.
func (as *AuthDBModel) CreatePasswordResetRequest(request *PasswordResetRequest) error {
	return as.DB.Create(request).Error
}

// GetPasswordResetRequestByToken retrieves a password reset request record by its token.
func (as *AuthDBModel) GetPasswordResetRequestByToken(token string) (*PasswordResetRequest, error) {
	var request PasswordResetRequest
	err := as.DB.Where("token = ?", token).First(&request).Error
	return &request, err
}

// DeletePasswordResetRequest deletes a password reset request record by its token.
func (as *AuthDBModel) DeletePasswordResetRequest(token string) error {
	return as.DB.Where("token = ?", token).Delete(&PasswordResetRequest{}).Error
}

func (PasswordResetRequest) TableName() string {
	return "password_reset_requests"
}

// CreateAgentLoginCredentials creates new agent login credentials.
func (as *AuthDBModel) CreateAgentLoginCredentials(credentials *AgentLoginCredentials) error {
	return as.DB.Create(credentials).Error
}

// GetAgentLoginCredentialsByID retrieves agent login credentials by their ID.
func (as *AuthDBModel) GetAgentLoginCredentialsByID(id uint) (*AgentLoginCredentials, error) {
	var credentials AgentLoginCredentials
	err := as.DB.Where("id = ?", id).First(&credentials).Error
	return &credentials, err
}

// UpdateAgentLoginCredentials updates agent login credentials.
func (as *AuthDBModel) UpdateAgentLoginCredentials(credentials *AgentLoginCredentials) error {
	return as.DB.Save(credentials).Error
}

// DeleteAgentLoginCredentials deletes agent login credentials by their ID.
func (as *AuthDBModel) DeleteAgentLoginCredentials(id uint) error {
	return as.DB.Delete(&AgentLoginCredentials{}, id).Error
}

// CreateUserLoginCredentials creates new user login credentials.
func (as *AuthDBModel) CreateUserLoginCredentials(credentials *UsersLoginCredentials) error {
	return as.DB.Create(credentials).Error
}

// GetUserLoginCredentialsByID retrieves user login credentials by their ID.
func (as *AuthDBModel) GetUserLoginCredentialsByID(id uint) (*UsersLoginCredentials, error) {
	var credentials UsersLoginCredentials
	err := as.DB.Where("id = ?", id).First(&credentials).Error
	return &credentials, err
}

// UpdateUserLoginCredentials updates user login credentials.
func (as *AuthDBModel) UpdateUserLoginCredentials(credentials *UsersLoginCredentials) error {
	return as.DB.Save(credentials).Error
}

// DeleteUserLoginCredentials deletes user login credentials by their ID.
func (as *AuthDBModel) DeleteUserLoginCredentials(id uint) error {
	return as.DB.Delete(&UsersLoginCredentials{}, id).Error
}

// /////////////////////////////////////////////////////////////////////////////////////
// Password History
type PasswordHistory struct {
	gorm.Model
	UserID      uint      `json:"user_id"`
	Password    string    `json:"-"`
	DateChanged time.Time `json:"date_changed"`
}

// CreatePasswordHistory creates a new password history entry.
func (as *AuthDBModel) CreatePasswordHistory(history *PasswordHistory) error {
	return as.DB.Create(history).Error
}

// GetPasswordHistoryByUserID retrieves password history for a user by their ID.
func (as *AuthDBModel) GetPasswordHistoryByUserID(userID uint) ([]*PasswordHistory, error) {
	var history []*PasswordHistory
	err := as.DB.Where("user_id = ?", userID).Find(&history).Error
	return history, err
}

// //////////////////////////////////////////////////////////////////////////////////////////
// Agent User Mapping
type AgentUserMapping struct {
	gorm.Model
	AgentID uint `json:"agent_id"`
	UserID  uint `json:"user_id"`
}

// TableName sets the table name for the AgentUserMapping model.
func (AgentUserMapping) TableName() string {
	return "agentUserMapping"
}

type PasswordResetToken struct {
	gorm.Model
	Token     string    `json:"token"`
	UserID    uint      `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (PasswordResetToken) TableName() string {
	return "password_reset_tokens"
}

// CreatePasswordResetToken creates a new password reset token.
func (as *AuthDBModel) CreatePasswordResetToken(token *PasswordResetToken) error {
	return as.DB.Create(token).Error
}

// GetPasswordResetTokenByToken retrieves a password reset token by its token string.
func (as *AuthDBModel) GetPasswordResetTokenByToken(token string) (*PasswordResetToken, error) {
	var resetToken PasswordResetToken
	err := as.DB.Where("token = ?", token).First(&resetToken).Error
	return &resetToken, err
}

// DeletePasswordResetToken deletes a password reset token by its token string.
func (as *AuthDBModel) DeletePasswordResetToken(token string) error {
	return as.DB.Where("token = ?", token).Delete(&PasswordResetToken{}).Error
}

// CreateAgentUserMapping creates a mapping between an agent and a user.
func (as *AuthDBModel) CreateAgentUserMapping(mapping *AgentUserMapping) error {
	return as.DB.Create(mapping).Error
}

// DeleteAgentUserMapping deletes a mapping between an agent and a user by agent and user IDs.
func (as *AuthDBModel) DeleteAgentUserMapping(agentID, userID uint) error {
	return as.DB.Where("agent_id = ? AND user_id = ?", agentID, userID).Delete(&AgentUserMapping{}).Error
}

// GetAllUserCredentials retrieves all user login credentials from the database.
func (as *AuthDBModel) GetAllUserCredentials() ([]*UsersLoginCredentials, error) {
	var userCredentials []*UsersLoginCredentials
	err := as.DB.Find(&userCredentials).Error
	return userCredentials, err
}

// GetAllAgentCredentials retrieves all agent login credentials from the database.
func (as *AuthDBModel) GetAllAgentCredentials() ([]*AgentLoginCredentials, error) {
	var agentCredentials []*AgentLoginCredentials
	err := as.DB.Find(&agentCredentials).Error
	return agentCredentials, err
}

func (am *AuthDBModel) VerifyTwoFactorAuth(agentID uint, token string) (bool, error) {
	// Implementation details...
}

// In auth.go
type OAuthProvider interface {
	ExchangeToken(code string) (token string, err error)
	GetUserInfo(token string) (user *OAuthUserInfo, err error)
}

type OAuthUserInfo struct {
	Email string
	Name  string
}

func NewOAuthProvider(providerName string) OAuthProvider {
	// Factory method to return a specific OAuth provider implementation
	// based on the provider name, e.g., Google, Facebook, etc.
}

func (am *AuthDBModel) OAuthLogin(providerName string, code string) (*Agents, error) {
	provider := NewOAuthProvider(providerName)
	token, err := provider.ExchangeToken(code)
	if err != nil {
		return nil, err
	}

	userInfo, err := provider.GetUserInfo(token)
	if err != nil {
		return nil, err
	}

	// Check if user already exists, create a new one if not, and return the user
}
