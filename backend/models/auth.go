// backend/models/auth.go

package models

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

type APIAccessLog struct {
	gorm.Model
	UserID      uint      `json:"user_id" gorm:"index"`
	Endpoint    string    `json:"endpoint" gorm:"type:varchar(255)"`
	AccessTime  time.Time `json:"access_time"`
	Status      string    `json:"status" gorm:"type:varchar(50)"` // Success, Failed
	Description string    `json:"description" gorm:"type:text"`   // Detailed access log
}

func (APIAccessLog) TableName() string {
	return "api_access_logs"
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

type AuthSystem interface {
	// Authentication and Registration
	AuthenticateUser(email, password string) (*Users, error)
	AuthenticateAgent(email, password string) (*Agents, error)
	RegisterUser(user *Users, password string) (*Users, error)
	RegisterAgent(agent *Agents, password string) (*Agents, error)

	// Credentials and Password Management
	CreateUserCredentials(credentials *UsersLoginCredentials) error
	CreateAgentCredentials(credentials *AgentLoginCredentials) error
	UpdateUserCredentials(userID uint, credentials *UsersLoginCredentials) error
	UpdateAgentCredentials(agentID uint, credentials *AgentLoginCredentials) error
	GetUserCredentialsByID(id uint) (*UsersLoginCredentials, error)
	GetAgentCredentialsByID(id uint) (*AgentLoginCredentials, error)
	DeleteUserCredentials(id uint) error
	DeleteAgentCredentials(id uint) error
	ResetUserPassword(userID uint, newPassword string) error
	ResetAgentPassword(agentID uint, newPassword string) error
	GetPasswordHistoryByUserID(userID uint) ([]*PasswordHistory, error)
	GetPasswordHistoryByAgentID(agentID uint) ([]*PasswordHistory, error)

	// Password Reset Workflow
	CreatePasswordResetRequest(request *PasswordResetRequest) error
	ValidatePasswordResetToken(token string) (*PasswordResetRequest, error)
	DeletePasswordResetRequest(requestID uint) error

	// Auxiliary Services and Data Consent Management
	CreateExternalServiceIntegration(integration *ExternalServiceIntegration) error
	UpdateExternalServiceIntegration(integrationID uint, integration *ExternalServiceIntegration) error
	DeleteExternalServiceIntegration(integrationID uint) error
	GetExternalServiceIntegrationByID(id uint) (*ExternalServiceIntegration, error)
	ListExternalServiceIntegrations() ([]*ExternalServiceIntegration, error)
	RecordDataConsent(consent *DataConsent) error
	UpdateDataConsent(consentID uint, consent *DataConsent) error
	GetDataConsentByID(consentID uint) (*DataConsent, error)
	ListDataConsentsByUserID(userID uint) ([]*DataConsent, error)

	// Encryption Key and 2FA Management
	CreateEncryptionKey(key *EncryptionKey) error
	GetEncryptionKeyByID(keyID uint) (*EncryptionKey, error)
	UpdateEncryptionKey(keyID uint, key *EncryptionKey) error
	DeleteEncryptionKey(keyID uint) error
	EnableTwoFactorAuthentication(userID uint) error
	VerifyTwoFactorCode(userID uint, code string) (bool, error)

	// OAuth2 Integration and Account Lockout Mechanism
	CreateOAuth2Integration(serviceName string, credentials *OAuth2Credentials) error
	AuthenticateWithOAuth2(serviceName string, token string) (*Users, error)
	IncrementFailedLoginAttempts(userID uint) error
	CheckAccountLockoutStatus(userID uint) (bool, error)
	ResetFailedLoginAttempts(userID uint) error

	// User Activity and Real-Time Alert
	LogUserActivity(userID uint, activityType string, details string) error
	GetUserActivityLog(userID uint) ([]*UserActivityLog, error)
	SendRealTimeAlert(userID uint, alertType string, message string) error

	// Security Measures and User Segmentation
	CheckPasswordStrength(password string) (PasswordStrength, error)
	CheckRateLimit(ipAddress string, endpoint string) (bool, error)
	SetSecurityQuestions(userID uint, questions []*SecurityQuestion) error
	VerifySecurityAnswers(userID uint, answers []*SecurityAnswer) (bool, error)
	UpdateDataSharingConsents(userID uint, consents []*DataConsent) error
	GetDataSharingConsents(userID uint) ([]*DataConsent, error)
	GetUserSegments(userID uint) ([]*UserSegment, error)
	AddUserToSegment(userID uint, segmentID uint) error

	// IP Whitelisting
	AddIPToWhitelist(userID uint, ipAddress string) error
	IsIPWhitelisted(userID uint, ipAddress string) (bool, error)
}

type AgentLoginCredentials struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	AgentID      uint       `json:"agent_id" gorm:"index;not null"`
	Username     string     `gorm:"type:varchar(255);not null" json:"username"`
	Password     string     `gorm:"-" json:"-"` // Excluded from JSON responses for security
	PasswordHash string     `gorm:"-" json:"-"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
}

// TableName sets the table name for the Agent model.
func (AgentLoginCredentials) TableName() string {
	return "agent_login_credentials"
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type UsersLoginCredentials struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	UserID       uint       `json:"user_id" gorm:"index;not null"`
	Username     string     `gorm:"type:varchar(255);not null" json:"username"`
	Password     string     `gorm:"-" json:"-"` // Excluded from JSON responses for security
	PasswordHash string     `gorm:"-" json:"-"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
}

// TableName sets the table name for the UsersLoginCredentials model.
func (UsersLoginCredentials) TableName() string {
	return "users_login_credentials"
}

// AuthModel handles database operations for Auth
type AuthDBModel struct {
	DB             *gorm.DB
	UserDBModel    *UserDBModel
	AgentDBModel   *AgentDBModel
	log            Logger
	EventPublisher *EventPublisherImpl
}

// NewUserModel creates a new instance of UserModel
func NewAuthDBModel(db *gorm.DB, userDBModel *UserDBModel, agentDBModel *AgentDBModel, log Logger, eventPublisher *EventPublisherImpl) *AuthDBModel {
	return &AuthDBModel{
		DB:             db,
		UserDBModel:    userDBModel,
		AgentDBModel:   agentDBModel,
		log:            log,
		EventPublisher: eventPublisher,
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

// TableName sets the table name for the UsersLoginCredentials model.
func (LoginInfo) TableName() string {
	return "login_info_auth"
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
func (db *RegisterModel) RegisterUser2(user *Users) error {
	return db.a.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		return nil
	})
}

// AuthenticateUser checks user credentials and returns the user if they are valid.
func (db *RegisterModel) AuthenticateUser(email, password string) (*Users, error) {
	var user Users
	// Fetch user by email
	if err := db.a.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	// Compare provided password with the hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Credentials.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	return &user, nil
}

// AuthenticateAgent checks agent credentials and returns the agent if they are valid.
func (db *AuthDBModel) AuthenticateAgent(email, password string) (*Agents, error) {
	var agent Agents
	if err := db.DB.Where("email = ?", email).First(&agent).Error; err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(agent.Credentials.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	return &agent, nil
}

// RegisterUser creates a new user account in the system with hashed password.
func (db *RegisterModel) RegisterUser(user *Users) error {
	return db.a.DB.Transaction(func(tx *gorm.DB) error {
		// Hash the user's password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Credentials.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		user.Credentials.PasswordHash = string(hashedPassword)
		user.Credentials.Password = "" // Ensure the plain password is not stored

		// Create the user with hashed password
		if err := tx.Create(user).Error; err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
		return nil
	})
}

// AgentRegistration registers a new agent with hashed password.
func (db *RegisterModel) AgentRegistration(agent *Agents) error {
	return db.c.DB.Transaction(func(tx *gorm.DB) error {
		// Hash the agent's password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(agent.Credentials.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		agent.Credentials.PasswordHash = string(hashedPassword)
		agent.Credentials.Password = "" // Ensure the plain password is not stored

		// Create the agent with hashed password
		if err := tx.Create(agent).Error; err != nil {
			return fmt.Errorf("failed to create agent: %w", err)
		}
		return nil
	})
}

// Login checks user credentials against stored credentials and returns the user on success.
func (db *RegisterModel) Login(email, password string) (*Users, error) {
	var user Users
	// Fetch user by email
	if err := db.a.DB.Where("email = ?", email).Preload("Credentials").First(&user).Error; err != nil {
		return nil, err
	}
	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Credentials.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid login credentials")
	}
	return &user, nil
}

// Login checks agent credentials against stored credentials and returns the agent on success.
func (db *RegisterModel) LoginAgent(email, password string) (*Agents, error) {
	var agent Agents
	// Fetch agent by email
	if err := db.c.DB.Where("email = ?", email).Preload("Credentials").First(&agent).Error; err != nil {
		return nil, err
	}
	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(agent.Credentials.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid login credentials")
	}
	return &agent, nil
}

// ResetUserPassword retrieves all users from the database.
func (as *AuthDBModel) ResetUserPassword(userID uint, newPassword string) error {
	return as.DB.Transaction(func(tx *gorm.DB) error {
		var userCredentials UsersLoginCredentials
		if err := tx.Where("user_id = ?", userID).First(&userCredentials).Error; err != nil {
			return err // User credentials not found
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return err // Failed to hash password
		}
		userCredentials.PasswordHash = string(hashedPassword)

		if err := tx.Save(&userCredentials).Error; err != nil {
			return err // Failed to update user credentials
		}

		// Optionally, log the password change or invalidate sessions/tokens
		return nil
	})
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

// ResetUserPassword resets the password for a user identified by userID.
func (as *AuthDBModel) ResetUserPasswordMain(userID uint, newPassword string) error {
	return as.DB.Transaction(func(tx *gorm.DB) error {
		var userCredentials UsersLoginCredentials
		// Fetch user credentials by user ID
		if err := tx.Where("user_id = ?", userID).First(&userCredentials).Error; err != nil {
			return err // User credentials not found
		}

		// Hash the new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		userCredentials.PasswordHash = string(hashedPassword)

		// Update user credentials with the new hashed password
		if err := tx.Save(&userCredentials).Error; err != nil {
			return err // Failed to update user credentials
		}

		// Optionally, log the password change or invalidate sessions/tokens
		return nil
	})
}

// ResetAgentPassword resets the password for an agent identified by agentID.
func (as *AuthDBModel) ResetAgentPasswordMain(agentID uint, newPassword string) error {
	return as.DB.Transaction(func(tx *gorm.DB) error {
		var agentCredentials AgentLoginCredentials
		// Fetch agent credentials by agent ID
		if err := tx.Where("agent_id = ?", agentID).First(&agentCredentials).Error; err != nil {
			return err // Agent credentials not found
		}

		// Hash the new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		agentCredentials.PasswordHash = string(hashedPassword)

		// Update agent credentials with the new hashed password
		if err := tx.Save(&agentCredentials).Error; err != nil {
			return err // Failed to update agent credentials
		}

		// Optionally, log the password change or invalidate sessions/tokens
		return nil
	})
}

//////////////////////////////////////////////////////////////////////////////////

// Define a model for storing password reset requests
type PasswordResetRequest struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	UserID    uint            `json:"user_id" gorm:"index;not null"`
	RequestID string          `gorm:"size:255;not null;unique" json:"request_id"`
	Token     string          `gorm:"size:255;not null;unique" json:"token"`
	ExpiresAt float64         `json:"expires_at"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (PasswordResetRequest) TableName() string {
	return "password_reset_requests"
}

// GetPasswordResetRequestByToken retrieves a password reset request record by its token.
func (as *AuthDBModel) GetPasswordResetRequestByToken(token string) (*PasswordResetRequest, error) {
	var resetRequest PasswordResetRequest
	err := as.DB.Where("token = ?", token).First(&resetRequest).Error
	return &resetRequest, err
}

func (as *AuthDBModel) CreatePasswordResetRequest(request *PasswordResetRequest) error {
	return as.DB.Create(request).Error
}

// ValidatePasswordResetToken validates a password reset token and checks if it's valid and not expired.
func (as *AuthDBModel) ValidatePasswordResetToken(token string) (*PasswordResetRequest, error) {
	var request PasswordResetRequest
	// Check if token exists and is not expired
	err := as.DB.Where("token = ? AND expires_at > ?", token, time.Now()).First(&request).Error
	if err != nil {
		return nil, err // Token not valid or expired
	}
	return &request, nil
}

func (as *AuthDBModel) ValidatePasswordResetToken2(token string) (*PasswordResetRequest, error) {
	var request PasswordResetRequest
	err := as.DB.Where("token = ? AND expires_at > ?", token, time.Now()).First(&request).Error
	if err != nil {
		return nil, err // Token not valid or expired
	}
	return &request, nil
}

func (as *AuthDBModel) DeletePasswordResetRequest(requestID uint) error {
	return as.DB.Delete(&PasswordResetRequest{}, requestID).Error
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
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `json:"user_id" gorm:"index;not null"`
	PasswordHash string    `gorm:"-" json:"-"` // Excluded from JSON for security
	DateChanged  time.Time `json:"date_changed"`
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
	return as.DB.Create(integration).Error
}

func (as *AuthDBModel) GetExternalServiceIntegrationByID(id uint) (*ExternalServiceIntegration, error) {
	var integration ExternalServiceIntegration
	err := as.DB.First(&integration, id).Error
	return &integration, err
}

func (as *AuthDBModel) UpdateExternalServiceIntegration(integrationID uint, integration *ExternalServiceIntegration) error {
	return as.DB.Transaction(func(tx *gorm.DB) error {
		var existingIntegration ExternalServiceIntegration
		if err := tx.Where("id = ?", integrationID).First(&existingIntegration).Error; err != nil {
			return err // Integration not found
		}

		// Update fields in existingIntegration from integration
		// For simplicity, assuming direct assignment is possible
		// This should be replaced with actual field updates
		existingIntegration = *integration

		return tx.Save(&existingIntegration).Error
	})
}

func (as *AuthDBModel) DeleteExternalServiceIntegration(integrationID uint) error {
	return as.DB.Delete(&ExternalServiceIntegration{}, integrationID).Error
}

func (as *AuthDBModel) ListExternalServiceIntegrations() ([]*ExternalServiceIntegration, error) {
	var integrations []*ExternalServiceIntegration
	err := as.DB.Find(&integrations).Error
	return integrations, err
}

func (db *AuthDBModel) LogAPIAccess(userID uint, endpoint, status, description string) error {
	return db.DB.Create(&APIAccessLog{
		UserID:      userID,
		Endpoint:    endpoint,
		AccessTime:  time.Now(),
		Status:      status,
		Description: description,
	}).Error
}

// //// Encryption Key
func (as *AuthDBModel) CreateEncryptionKey(key *EncryptionKey) error {
	if err := as.DB.Create(key).Error; err != nil {
		log.Printf("Failed to create encryption key: %v", err)
		return fmt.Errorf("could not create encryption key: %w", err)
	}
	return nil
}

func (as *AuthDBModel) GetEncryptionKeyByID(keyID uint) (*EncryptionKey, error) {
	var key EncryptionKey
	if err := as.DB.Where("id = ?", keyID).First(&key).Error; err != nil {
		log.Printf("Encryption key not found: %v", err)
		return nil, fmt.Errorf("encryption key not found: %w", err)
	}
	return &key, nil
}

func (as *AuthDBModel) UpdateEncryptionKey(keyID uint, key *EncryptionKey) error {
	return as.DB.Transaction(func(tx *gorm.DB) error {
		var existingKey EncryptionKey
		if err := tx.Where("id = ?", keyID).First(&existingKey).Error; err != nil {
			return fmt.Errorf("encryption key not found: %w", err)
		}

		existingKey.KeyType = key.KeyType
		existingKey.KeyData = key.KeyData

		if err := tx.Save(&existingKey).Error; err != nil {
			return fmt.Errorf("failed to update encryption key: %w", err)
		}
		return nil
	})
}

func (as *AuthDBModel) DeleteEncryptionKey(keyID uint) error {
	return as.DB.Delete(&EncryptionKey{}, keyID).Error
}

func (db *AuthDBModel) UpdateUserConsent(userID uint, consentType string, granted bool) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var consent UserConsent
		err := tx.Where("user_id = ? AND type = ?", userID, consentType).FirstOrCreate(&consent).Error
		if err != nil {
			return err // Consent record not found, or failed to create
		}

		consent.Granted = granted
		if granted {
			consent.GrantedAt = time.Now()
			consent.RevokedAt = time.Time{} // Reset revoke time
		} else {
			consent.RevokedAt = time.Now()
		}

		if err := tx.Save(&consent).Error; err != nil {
			return err // Failed to update consent
		}
		return nil
	})
}

// 2FA
type TwoFactorAuthentication struct {
	gorm.Model
	UserID    uint   `json:"user_id" gorm:"index;not null"`
	SecretKey string `json:"secret_key" gorm:"not null"`
	IsEnabled bool   `json:"is_enabled" gorm:"default:false"`
}

func (TwoFactorAuthentication) TableName() string {
	return "two_factor_authentications"
}

func (db *AuthDBModel) EnableTwoFactorAuthentication(userID uint, secretKey string) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		tfa := TwoFactorAuthentication{
			UserID:    userID,
			SecretKey: secretKey,
			IsEnabled: true,
		}
		if err := tx.Create(&tfa).Error; err != nil {
			return fmt.Errorf("failed to enable 2FA: %w", err)
		}
		return nil
	})
}

// Google OAuth2
func (db *AuthDBModel) AuthenticateWithOAuth2(serviceName, token string) (*Users, error) {
	var oauth2Creds OAuth2Credentials
	if err := db.DB.Where("service_name = ? AND token = ?", serviceName, token).First(&oauth2Creds).Error; err != nil {
		return nil, fmt.Errorf("OAuth2 credentials not found: %w", err)
	}
	user, err := db.UserDBModel.GetUserByID(oauth2Creds.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}

// User Consent
type UserConsent struct {
	gorm.Model
	UserID    uint      `json:"user_id" gorm:"index;not null"`
	Type      string    `json:"type" gorm:"type:varchar(100);not null"` // e.g., "marketing", "analytics"
	Granted   bool      `json:"granted"`
	GrantedAt time.Time `json:"granted_at"`
	RevokedAt time.Time `json:"revoked_at,omitempty"`
}

func (UserConsent) TableName() string {
	return "user_consents"
}

func (UserActivityLog) TableName() string {
	return "user_activity_log"
}

func (db *AuthDBModel) UserActivityLog(userID uint, activityType, details string) error {
	ua := &UserActivityLog{
		UserID:       userID,
		ActivityType: activityType,
		Details:      details,
		Timestamp:    time.Now(),
	}
	if err := db.DB.Create(&ua).Error; err != nil {
		// Log the error for debugging purposes
		log.Printf("Error saving user activity: %v", err)
		return fmt.Errorf("failed to log user activity: %w", err)
	}
	return nil
}

// LogAction records an action performed by an agent for auditing purposes.
func (as *AuthDBModel) LogUserActivity(userID uint, activity string, details string) error {
	actionLog := UserActivityLog{
		UserID:       userID,
		ActivityType: activity,
		Details:      details,
		Timestamp:    time.Now(),
	}

	return as.DB.Create(&actionLog).Error
}

func (as *AuthDBModel) GetUserActivityLog(userID uint) ([]*UserActivityLog, error) {
	var logs []*UserActivityLog
	if err := as.DB.Where("user_id = ?", userID).Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

func (as *AuthDBModel) SendRealTimeAlert(userID uint, alertType string, message string) error {
	// Placeholder for sending an alert (e.g., email, SMS, push notification)
	// This would involve calling the respective API of the service provider
	// For demonstration purposes only, the actual implementation is skipped
	log.Printf("Sending %s alert to user %d: %s", alertType, userID, message)
	return nil // Assuming the send operation is successful
}

type PasswordStrength struct {
	Level           string   `json:"level"`                     // Example levels: Weak, Moderate, Strong
	Recommendations []string `json:"recommendations,omitempty"` // Suggestions for improving password strength
}

func (as *AuthDBModel) CheckPasswordStrength(password string) (PasswordStrength, error) {
	// Implement logic to analyze the password's strength based on criteria such as length, diversity of characters, etc.
	// This is a simplified example:
	if len(password) < 8 {
		return PasswordStrength{Level: "Weak", Recommendations: []string{"Use at least 8 characters."}}, nil
	} else if len(password) >= 8 && len(password) <= 12 {
		return PasswordStrength{Level: "Moderate", Recommendations: []string{"Use special characters and numbers to increase strength."}}, nil
	} else {
		return PasswordStrength{Level: "Strong"}, nil
	}
}

// Data Consent

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

func (as *AuthDBModel) RecordDataConsent(consent *DataConsent) error {
	if err := as.DB.Create(consent).Error; err != nil {
		log.Printf("Failed to record data consent: %v", err)
		return fmt.Errorf("could not record data consent: %w", err)
	}
	return nil
}

func (as *AuthDBModel) UpdateDataConsent(consentID uint, consent *DataConsent) error {
	return as.DB.Transaction(func(tx *gorm.DB) error {
		var existingConsent DataConsent
		if err := tx.Where("id = ?", consentID).First(&existingConsent).Error; err != nil {
			return fmt.Errorf("data consent not found: %w", err)
		}

		existingConsent.IsGranted = consent.IsGranted
		existingConsent.GrantedAt = consent.GrantedAt
		existingConsent.RevokedAt = consent.RevokedAt

		if err := tx.Save(&existingConsent).Error; err != nil {
			return fmt.Errorf("failed to update data consent: %w", err)
		}
		return nil
	})
}

func (as *AuthDBModel) GetDataConsentByID(consentID uint) (*DataConsent, error) {
	var consent DataConsent
	err := as.DB.First(&consent, consentID).Error
	return &consent, err
}

func (as *AuthDBModel) ListDataConsentsByUserID(userID uint) ([]*DataConsent, error) {
	var consents []*DataConsent
	err := as.DB.Where("user_id = ?", userID).Find(&consents).Error
	return consents, err
}

// ////////////////////////USER Segmentation

type UserSegment struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Name        string     `json:"name" gorm:"type:varchar(255);not null"`
	Description string     `gorm:"type:varchar(255);not null;unique" json:"description"` // Optional description of the segment
	IsActive    bool       `json:"is_active" gorm:"default:true"`                        // Whether the segment is currently active
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`                     // Timestamp of segment creation
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`                     // Timestamp of last update
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index"`                    // Soft delete timestamp
	// Additional metadata fields as JSON for flexible data storage
	Metadata string `json:"metadata,omitempty" gorm:"type:text;null"` // JSON string for storing additional metadata
}

func (UserSegment) TableName() string {
	return "user_segments"
}

type UserSegmentMapping struct {
	UserID    uint `gorm:"primaryKey;autoIncrement:false" json:"user_id"`
	SegmentID uint `gorm:"primaryKey;autoIncrement:false" json:"segment_id"`
}

func (UserSegmentMapping) TableName() string {
	return "user_segment_mappings"
}

func (as *AuthDBModel) GetUserSegments(userID uint) ([]*UserSegment, error) {
	var segments []*UserSegment
	err := as.DB.Joins("JOIN user_segment_mappings on user_segments.id = user_segment_mappings.segment_id").
		Where("user_segment_mappings.user_id = ?", userID).Find(&segments).Error
	return segments, err
}

func (as *AuthDBModel) AddUserToSegment(userID uint, segmentID uint) error {
	mapping := UserSegmentMapping{UserID: userID, SegmentID: segmentID}
	return as.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&mapping).Error
}

// /////////////////// RATE LIMITS ///////////////////////////////
type RateLimit struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	IPAddress    string    `gorm:"type:varchar(255);not null;index:idx_ip_address" json:"ip_address"`
	Endpoint     string    `gorm:"type:varchar(255);not null;index:idx_endpoint" json:"endpoint"`
	RequestCount int       `gorm:"default:0" json:"request_count"`
	ResetAt      time.Time `gorm:"" json:"reset_at"`
}

func (RateLimit) TableName() string {
	return "rate_limits"
}

func (as *AuthDBModel) CheckRateLimit(ipAddress string, endpoint string) (bool, error) {
	var limit RateLimit
	now := time.Now()

	err := as.DB.Where("ip_address = ? AND endpoint = ? AND reset_at > ?", ipAddress, endpoint, now).First(&limit).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// No existing rate limit record or it's expired, create or reset
		as.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "ip_address"}, {Name: "endpoint"}},
			DoUpdates: clause.AssignmentColumns([]string{"request_count", "reset_at"}),
		}).Create(&RateLimit{
			IPAddress:    ipAddress,
			Endpoint:     endpoint,
			RequestCount: 1,
			ResetAt:      now.Add(1 * time.Hour), // Reset in 1 hour
		})
		return true, nil // Within limit
	} else if err != nil {
		return false, err // Database error
	}

	if limit.RequestCount >= 100 { // Assuming limit is 100 requests per hour
		return false, nil // Rate limit exceeded
	}

	// Increment request count
	as.DB.Model(&limit).Update("request_count", gorm.Expr("request_count + ?", 1))
	return true, nil // Within limit
}

// ////////////////////////
type SecurityQuestion struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	UserID   uint   `gorm:"index;not null" json:"user_id"`
	Question string `gorm:"type:text;not null" json:"question"`
	Answer   string `gorm:"-" json:"-"` // Store hashed answers for security
}

func (SecurityQuestion) TableName() string {
	return "security_questions"
}

type SecurityAnswer struct {
	QuestionID uint   `json:"question_id"`
	Answer     string `json:"answer"`
}

func (as *AuthDBModel) SetSecurityQuestions(userID uint, questions []*SecurityQuestion) error {
	return as.DB.Transaction(func(tx *gorm.DB) error {
		for _, question := range questions {
			hashedAnswer, err := bcrypt.GenerateFromPassword([]byte(question.Answer), bcrypt.DefaultCost)
			if err != nil {
				return err
			}
			question.Answer = string(hashedAnswer)
			question.UserID = userID
			if err := tx.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(&question).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (as *AuthDBModel) VerifySecurityAnswers(userID uint, answers []*SecurityAnswer) (bool, error) {
	for _, answer := range answers {
		var question SecurityQuestion
		if err := as.DB.Where("id = ? AND user_id = ?", answer.QuestionID, userID).First(&question).Error; err != nil {
			return false, err
		}
		if err := bcrypt.CompareHashAndPassword([]byte(question.Answer), []byte(answer.Answer)); err != nil {
			return false, nil // Answer does not match
		}
	}
	return true, nil // All answers match
}

type IPWhitelist struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	IPAddress string    `gorm:"type:varchar(255);not null" json:"ip_address"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (IPWhitelist) TableName() string {
	return "ip_whitelists"
}

func (as *AuthDBModel) AddIPToWhitelist(userID uint, ipAddress string) error {
	whitelistEntry := IPWhitelist{UserID: userID, IPAddress: ipAddress}
	return as.DB.Create(&whitelistEntry).Error
}

func (as *AuthDBModel) IsIPWhitelisted(userID uint, ipAddress string) (bool, error) {
	var count int64
	as.DB.Model(&IPWhitelist{}).Where("user_id = ? AND ip_address = ?", userID, ipAddress).Count(&count)
	return count > 0, nil
}

// GoogleUserInfo represents the information received from Google's UserInfo endpoint.
type GoogleUserInfo struct {
	ID            string `json:"id"`             // Google's identifier for the user
	Email         string `json:"email"`          // User's email address
	VerifiedEmail bool   `json:"verified_email"` // If the user's email address has been verified
	Name          string `json:"name"`           // User's full name
	GivenName     string `json:"given_name"`     // User's given name (first name)
	FamilyName    string `json:"family_name"`    // User's family name (last name)
	Picture       string `json:"picture"`        // URL of the user's profile picture
	Locale        string `json:"locale"`         // Locale of the user
}

// MyCustomClaims includes standard JWT claims and a Role field
type MyCustomClaims struct {
	jwt.StandardClaims
	Role string `json:"role"`
}

/*
type AgentLoginCredentialsStorage interface {
	Create(credentials *AgentLoginCredentials) error
	Delete(id uint) error
	Update(credentials *AgentLoginCredentials) error
	FindByID(id uint) (*AgentLoginCredentials, error)
	ResetAgentPassword(uint, string) ([]*AgentLoginCredentials, error)
	CreateAgentCredentials(agentCredentials *AgentLoginCredentials) error
	GetAgentCredentialsByID(id uint) (*AgentLoginCredentials, error)
	UpdateAgentCredentials(agentCredentials *AgentLoginCredentials) error
	DeleteAgentCredentials(id uint) error
	GetAllAgentCreds() ([]*AgentLoginCredentials, error)
	CreateAgentLoginCredentials(credentials *AgentLoginCredentials) error
	GetAgentLoginCredentialsByID(id uint) (*AgentLoginCredentials, error)
	UpdateAgentLoginCredentials(credentials *AgentLoginCredentials) error
	DeleteAgentLoginCredentials(id uint) error
}

type UserLoginCredentialsStorage interface {
	Create(credentials *UsersLoginCredentials) error
	Delete(id uint) error
	Update(credentials *UsersLoginCredentials) error
	FindByID(id uint) (*UsersLoginCredentials, error)
	AuthenticateUser(email, password string) (*Users, error)
	RegisterUser(user *Users) error
	Login(email, password string) (*Users, error)
	LoginAgent(email, password string) (*Agents, error)
	ResetUserPassword(uint, string) ([]*UsersLoginCredentials, error)
	CreatePasswordResetRequest(userID uint, requestID string, token string) error
	GetPasswordHistoryByUserID(userID uint) ([]*PasswordHistory, error)
	CreateUserCredentials(userCredentials *UsersLoginCredentials) error
	UpdateUserLoginCredentials(credentials *UsersLoginCredentials) error
	eleteUserLoginCredentials(id uint) error
	GetUserCredentialsByID(id uint) (*UsersLoginCredentials, error)
	CreateAgentUserMapping(mapping *AgentUserMapping) error
	DeleteAgentUserMapping(agentID, userID uint) error
	CreatePasswordResetToken(token *PasswordResetToken) error
	GetPasswordResetRequestByToken(token string) (*PasswordResetRequest, error)
	GetPasswordResetTokenByToken(token string) (*PasswordResetToken, error)
	DeletePasswordResetToken(token string) error
}
*/

/*!------------------------------------------------------------------------------!*/
////////////////////// Google Credentials ////////////////////////
