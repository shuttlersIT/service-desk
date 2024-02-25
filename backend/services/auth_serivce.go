package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/shuttlersit/service-desk/backend/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthServiceInterface defines the contract for authentication services.
type AuthServiceInterface interface {
	RegisterUser(user *models.Users, password string) (*models.Users, string, error)
	Login1(loginInfo *LoginInfo) (string, *models.Users, error)
	ResetUserPassword(userID uint, newPassword string) error
	Login2(credentials *models.UsersLoginCredentials) (*models.UserSession, error)
	Logout(sessionID string) error
	ValidateSession(sessionID string) (bool, error)
}

type AuthService struct {
	DB             *gorm.DB
	AuthDBModel    *models.AuthDBModel
	log            models.Logger
	EventPublisher models.EventPublisherImpl
}

// NewAuthService creates a new instance of AuthService with the given DB connection.
func NewAuthService(db *gorm.DB, authDBModel *models.AuthDBModel, log models.Logger, eventPublisher models.EventPublisherImpl) *AuthService {
	return &AuthService{
		DB:             db,
		AuthDBModel:    authDBModel,
		log:            log,
		EventPublisher: eventPublisher,
	}
}

// RegisterUser handles new user registration.
func (service *AuthService) RegisterUser(user *models.Users, password string) (*models.Users, string, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}
	user.Credentials.PasswordHash = string(hashedPassword)

	// Save user to DB
	if err := service.DB.Create(user).Error; err != nil {
		return nil, "", err
	}

	// Generate JWT token for the new user
	token, err := generateJWTToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// Login authenticates a user and returns a JWT token.
func (service *AuthService) Login(loginInfo *LoginInfo) (string, *models.Users, error) {
	var user models.Users
	if err := service.DB.Where("email = ?", loginInfo.Email).First(&user).Error; err != nil {
		return "", nil, err
	}

	// Compare the provided password with the hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Credentials.PasswordHash), []byte(loginInfo.Password)); err != nil {
		return "", nil, errors.New("incorrect password")
	}

	// Generate JWT token for authenticated user
	token, err := generateJWTToken(user.ID)
	if err != nil {
		return "", nil, err
	}

	return token, &user, nil
}

// ResetUserPassword updates a user's password.
func (service *AuthService) ResetUserPassword(userID uint, newPassword string) error {
	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update user's password in the DB
	if err := service.DB.Model(&models.Users{}).Where("id = ?", userID).Update("password_hash", string(hashedPassword)).Error; err != nil {
		return err
	}

	return nil
}

// ResetAgentPassword updates a agent's password.
func (service *AuthService) ResetAgentPassword(agentID uint, newPassword string) error {
	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update agent's password in the DB
	if err := service.DB.Model(&models.Agents{}).Where("id = ?", agentID).Update("password_hash", string(hashedPassword)).Error; err != nil {
		return err
	}

	return nil
}

// RegisterUser handles new user registration.
func (service *AuthService) RegisterAgent(agent *models.Agents, password string) (*models.Agents, string, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}
	agent.Credentials.PasswordHash = string(hashedPassword)

	// Save user to DB
	if err := service.DB.Create(agent).Error; err != nil {
		return nil, "", err
	}

	// Generate JWT token for the new user
	token, err := generateJWTToken(agent.ID)
	if err != nil {
		return nil, "", err
	}

	return agent, token, nil
}

// generateJWTToken creates a JWT token for user identification.
func generateJWTToken(userID uint) (string, error) {
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		Issuer:    "AuthService",
		Subject:   fmt.Sprintf("%d", userID),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte("YourSecretKeyHere")) // Use a secure method to store and retrieve your secret key

	return signedToken, err
}
