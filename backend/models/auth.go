// backend/models/auth.go

package models

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type APIToken struct {
	gorm.Model
	UserID      uint       `json:"user_id" gorm:"index;not null"`
	Token       string     `json:"token" gorm:"size:255;not null;unique"`
	Description string     `json:"description,omitempty" gorm:"size:255"`
	ExpiresAt   time.Time  `json:"expires_at"`
	LastUsedAt  *time.Time `json:"last_used_at,omitempty"`
}

func (APIToken) TableName() string {
	return "api_tokens"
}

type APIRequestLog struct {
	gorm.Model
	UserID       uint   `json:"user_id" gorm:"index;not null"`
	Endpoint     string `json:"endpoint" gorm:"type:varchar(255);not null"`
	Method       string `json:"method" gorm:"type:varchar(50);not null"`
	StatusCode   int    `json:"status_code" gorm:"not null"`
	RequestTime  int64  `json:"request_time" gorm:"not null"` // Request duration in milliseconds
	RequestBody  string `json:"request_body,omitempty" gorm:"type:text"`
	ResponseBody string `json:"response_body,omitempty" gorm:"type:text"`
}

func (APIRequestLog) TableName() string {
	return "api_request_logs"
}

type ExternalServiceToken struct {
	gorm.Model
	ServiceName string     `json:"service_name" gorm:"type:varchar(255);not null"`
	Token       string     `json:"token" gorm:"type:text;not null"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

func (ExternalServiceToken) TableName() string {
	return "external_service_tokens"
}

type AgentLoginCredentialsStorage interface {
	Create(credentials *AgentLoginCredentials) error
	Delete(id uint) error
	Update(credentials *AgentLoginCredentials) error
	FindByID(id uint) (*AgentLoginCredentials, error)
}

type AgentLoginCredentials struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	AgentID     uint       `json:"agent_id" gorm:"index;not null"`
	Username    string     `gorm:"type:varchar(255);not null" json:"username"`
	Password    string     `gorm:"-" json:"-"` // Excluded from JSON responses for security
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
}

// TableName sets the table name for the Agent model.
func (AgentLoginCredentials) TableName() string {
	return "agent_login_credentials"
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type UsersLoginCredentials struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
	UserID   uint   `json:"user_id"`
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
	c *AgentDBModel
}

// NewUserModel creates a new instance of UserModel
func NewRegisterModel(a *UserDBModel, b *AuthDBModel) *RegisterModel {
	return &RegisterModel{
		a: a,
		b: b,
	}
}

// RegisterUser creates a new user account in the system.
func (db *UserDBModel) RegisterUser(user *Users) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		return nil
	})
}

// AuthenticateUser checks user credentials and returns the user if they are valid.
func (db *UserDBModel) AuthenticateUser(email, password string) (*Users, error) {
	var user Users
	if err := db.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	// Assume CheckPassword is a method to verify the password; implement accordingly.
	if !CheckPassword(user.PasswordHash, password) {
		return nil, errors.New("invalid credentials")
	}
	return &user, nil
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

// Register New User
func (ab *RegisterModel) AgentRegistration(agent *Agents) (*Agents, error) {
	// Hash the user's password before storing it in the database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(agent.Credentials.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password")
	}
	agent.Credentials.Password = string(hashedPassword)
	newAgent, erro := ab.c.CreateAgent(agent)
	if erro != nil {
		return nil, fmt.Errorf("failed to create users")
	}
	er := ab.b.CreateAgentCredentials(&newAgent.Credentials)
	if er != nil {
		return nil, fmt.Errorf("failed to create users credentials")
	}
	e := ab.c.UpdateAgent(newAgent)
	if e != nil {
		return nil, fmt.Errorf("unable to update new users credentials")
	}

	return newAgent, nil
}

// AuthenticateUser checks user credentials and returns the user if they are valid.
func (db *UserDBModel) AuthenticateUser(email, password string) (*Users, error) {
	var user Users
	if err := db.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	// Assume CheckPassword is a method to verify the password; implement accordingly.
	if !CheckPassword(user.PasswordHash, password) {
		return nil, errors.New("invalid credentials")
	}
	return &user, nil
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

// Login User
func (a *AuthDBModel) LoginAgent(login *LoginInfo) (*Agents, error) {
	loginInfo := login
	var agent Agents
	if err := a.DB.Where("email = ?", loginInfo.Email).First(&agent).Error; err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(agent.Credentials.Password), []byte(loginInfo.Password)); err != nil {
		return nil, err
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
func (as *AuthDBModel) ResetAgentPassword(uint, string) ([]*AgentLoginCredentials, error) {
	var agentsCredentials []*AgentLoginCredentials
	err := as.DB.Find(&agentsCredentials).Error
	return agentsCredentials, err
}

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

//////////////////////////////////////////////////////////////////////////////////
//Password Requests

// Define a model for storing password reset requests
type PasswordResetRequest struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	AgentID   uint            `json:"agent_id" gorm:"index;not null"`
	RequestID string          `gorm:"size:255;not null;unique" json:"request_id"`
	Token     string          `gorm:"size:255;not null;unique" json:"token"`
	ExpiresAt time.Time       `json:"expires_at"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (PasswordResetRequest) TableName() string {
	return "password_reset_requests"
}

// CreatePasswordResetRequest creates a new password reset request record.
func (as *AuthDBModel) CreatePasswordResetRequest(userID uint, requestID uint, token string) error {
	resetRequest := &PasswordResetRequest{
		UserID:    userID,
		RequestID: requestID,
		Token:     token,
		// Set other fields as needed, such as expiration time
	}
	return as.DB.Create(resetRequest).Error
}

// GetPasswordResetRequestByToken retrieves a password reset request record by its token.
func (as *AuthDBModel) GetPasswordResetRequestByToken(token string) (*PasswordResetRequest, error) {
	var resetRequest PasswordResetRequest
	err := as.DB.Where("token = ?", token).First(&resetRequest).Error
	return &resetRequest, err
}

// DeletePasswordResetRequest deletes a password reset request record by its token.
func (as *AuthDBModel) DeletePasswordResetRequest(token string) error {
	return as.DB.Where("token = ?", token).Delete(&PasswordResetRequest{}).Error
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

// /////////////////////////////////////////////////////////////////////////////////////
// Password History
type PasswordHistory struct {
	ID           uint            `gorm:"primaryKey" json:"id"`
	UserID       uint            `json:"user_id" gorm:"index;not null"`
	PasswordHash string          `gorm:"-" json:"-"` // Excluded from JSON for security
	DateChanged  time.Time       `json:"date_changed"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	DeletedAt    *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (PasswordHistory) TableName() string {
	return "password_history"
}

type PasswordResetToken struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	AgentID   uint            `json:"agent_id" gorm:"index;not null"`
	Token     string          `gorm:"size:255;not null;unique" json:"token"`
	ExpiresAt time.Time       `json:"expires_at"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (PasswordResetToken) TableName() string {
	return "password_reset_tokens"
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
	ID        uint            `gorm:"primaryKey" json:"id"`
	AgentID   uint            `json:"agent_id" gorm:"index;not null"`
	UserID    uint            `json:"user_id" gorm:"index;not null"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName sets the table name for the AgentUserMapping model.
func (AgentUserMapping) TableName() string {
	return "agentUserMapping"
}

// CreateAgentUserMapping creates a mapping between an agent and a user.
func (as *AuthDBModel) CreateAgentUserMapping(mapping *AgentUserMapping) error {
	return as.DB.Create(mapping).Error
}

// DeleteAgentUserMapping deletes a mapping between an agent and a user.
func (as *AuthDBModel) DeleteAgentUserMapping(agentID, userID uint) error {
	return as.DB.Where("agent_id = ? AND user_id = ?", agentID, userID).Delete(&AgentUserMapping{}).Error
}

type DataConsent struct {
	ID          uint       `gorm:"primaryKey"`
	UserID      uint       `gorm:"index;not null" json:"user_id"`
	ConsentType string     `gorm:"type:varchar(255);not null" json:"consent_type"` // E.g., "analytics", "personalization"
	IsGranted   bool       `json:"is_granted"`
	GrantedAt   time.Time  `json:"granted_at"`
	RevokedAt   *time.Time `json:"revoked_at,omitempty"`
}

func (DataConsent) TableName() string {
	return "data_consents"
}

type EncryptionKey struct {
	ID        uint      `gorm:"primaryKey"`
	OwnerID   uint      `gorm:"index;not null" json:"owner_id"`             // Could be a user or community
	KeyType   string    `gorm:"type:varchar(100);not null" json:"key_type"` // E.g., "AES", "RSA"
	KeyData   string    `gorm:"type:text;not null" json:"key_data"`         // Encrypted key data
	CreatedAt time.Time `json:"created_at"`
}

func (EncryptionKey) TableName() string {
	return "encryption_keys"
}

type APIGateway struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"type:varchar(255);not null" json:"name"`
	Endpoint  string `gorm:"type:text;not null" json:"endpoint"`    // URL to the gateway
	AuthToken string `gorm:"type:text" json:"auth_token,omitempty"` // Optional token for accessing the gateway
	IsActive  bool   `gorm:"default:true" json:"is_active"`
}

func (APIGateway) TableName() string {
	return "api_gateways"
}

type ExternalServiceIntegration struct {
	IntegrationID     uint       `gorm:"primaryKey"`
	ServiceName       string     `gorm:"unique;not null" json:"service_name"`
	ApiKey            string     `gorm:"not null" json:"api_key"`
	IntegrationConfig string     `gorm:"type:json;not null" json:"integration_config"`
	IsActive          bool       `gorm:"default:true" json:"is_active"`
	LastSync          *time.Time `json:"last_sync,omitempty"`
}

func (ExternalServiceIntegration) TableName() string {
	return "external_service_integration"
}

func (as *AuthDBModel) CreateExternalServiceIntegration(integration *ExternalServiceIntegration) error {
	return as.DB.Create(&integration).Error
}

func (as *AuthDBModel) GetExternalServiceIntegrationByID(id uint) (*ExternalServiceIntegration, error) {
	var integration *ExternalServiceIntegration
	result := as.DB.First(&integration, id)
	return integration, result.Error
}

// Implement Update and Delete similarly.
