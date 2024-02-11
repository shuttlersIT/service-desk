// backend/services/auth_service.go

package services

import (
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

// AuthServiceInterface provides methods for managing auth.
type AuthServiceInterface interface {
	Registration(user *models.Users) (*models.Users, string, error)
	Login(login *LoginInfo) (string, error)
	ResetUserPassword(userID uint, newPassword string) error
	ResetUserPassword2(userID uint, newPassword string) error
}

// DefaultAuthService is the default implementation of AuthService
type DefaultAuthService struct {
	DB            *gorm.DB
	AuthDBModel   *models.AuthDBModel
	UserDBModel   *models.UserDBModel
	RegisterModel *models.RegisterModel
	// Add any dependencies or data needed for the service
}

// NewDefaultAuthService creates a new DefaultAuthService.
func NewDefaultAuthService(authDBModel *models.AuthDBModel) *DefaultAuthService {
	return &DefaultAuthService{
		AuthDBModel: authDBModel,
	}
}

// User registration
func (a *DefaultAuthService) Registration(user *models.Users) (*models.Users, string, error) {

	u, e := a.RegisterModel.Registration(user)
	if e != nil {
		return nil, "", e
	}
	// Generate a JWT token for successful login
	token, err := generateJWTToken(u.ID)
	if err != nil {
		return nil, "", err
	}
	return u, token, nil
}

// User registration
func (a *DefaultAuthService) AgentRegistration(agent *models.Agents) (*models.Agents, string, error) {

	u, e := a.RegisterModel.AgentRegistration(agent)
	if e != nil {
		return nil, "", e
	}
	// Generate a JWT token for successful login
	token, err := generateJWTToken(u.ID)
	if err != nil {
		return nil, "", err
	}
	return u, token, nil
}

// User login
func (a *DefaultAuthService) Login(loginInfo *models.LoginInfo) (string, error) {
	user, err := a.AuthDBModel.Login(loginInfo)
	if err != nil {
		return "", err
	}
	// Generate a JWT token for successful login
	token, err := generateJWTToken(user.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

// User login
func (a *DefaultAuthService) LoginAgent(loginInfo *models.LoginInfo) (string, error) {
	agent, err := a.AuthDBModel.LoginAgent(loginInfo)
	if err != nil {
		return "", err
	}
	// Generate a JWT token for successful login
	token, err := generateJWTToken(agent.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

// Define a function to generate JWT token
func generateJWTToken(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	})
	return token.SignedString([]byte("your-secret-key"))
}

// backend/services/user_service.go

func (us *DefaultAuthService) ResetUserPassword(userID uint, newPassword string) error {
	// Reset the user's password
	_, err := us.AuthDBModel.ResetUserPassword(userID, newPassword)
	if err != nil {
		return err
	}

	return nil
}

func (us *DefaultAuthService) ResetAgentPassword(agentID uint, newPassword string) error {
	// Reset the agent's password
	_, err := us.AuthDBModel.ResetAgentPassword(agentID, newPassword)
	if err != nil {
		return err
	}

	return nil
}

func (us *DefaultUserService) ResetUserPassword2(userID uint, newPassword string) error {
	// Retrieve the user by userID
	user, err := us.UserDBModel.GetUserByID(userID)
	if err != nil {
		return err
	}

	// Update the user's password with the new password
	user.Credentials.Password = newPassword

	// Save the updated user
	err = us.UserDBModel.UpdateUser(user)
	if err != nil {
		return err
	}

	return nil
}

// ResetPasswordWithToken resets the user's password using a valid token.
func (a *DefaultAuthService) ResetPasswordWithToken(token string, newPassword string) error {
	// Token validation logic (customize this based on your token handling)
	claims, err := validatePasswordResetToken(token)
	if err != nil {
		return err
	}

	// Retrieve the user by their email (you should implement this method)
	user, err := a.UserDBModel.GetUserByEmail(claims["email"].(string))
	if err != nil {
		return err
	}

	// Check if the token has a valid reset password request ID claim
	requestIDClaim, ok := claims["request_id"].(float64)
	if !ok {
		return fmt.Errorf("token missing or invalid request_id claim")
	}
	requestID := uint(requestIDClaim)

	// Verify that the user has a pending reset password request with the same request ID
	if user.ResetPasswordRequestID != requestID {
		return fmt.Errorf("invalid or expired token")
	}

	// Update the user's password with the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Credentials.Password = string(hashedPassword)

	// Clear the reset password request ID
	user.ResetPasswordRequestID = 0

	// Save the updated user
	if err := a.UserDBModel.UpdateUser(user); err != nil {
		return err
	}

	return nil
}

// Helper function to validate the password reset token.
func validatePasswordResetToken(token string) (jwt.MapClaims, error) {
	// Parse and validate the JWT token
	claims := jwt.MapClaims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token signing method")
		}
		// Replace "your-secret-key" with your actual secret key used for token signing
		return []byte("your-secret-key"), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	if !tkn.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	// Ensure that the token contains the required claims (customize as needed)
	if claims["email"] == nil || claims["request_id"] == nil {
		return nil, fmt.Errorf("token missing required claims")
	}

	// Add more claim validations as needed

	return claims, nil
}

func (a *DefaultAuthService) SyncDataFromExternalService(integrationID uint) error {
	integration, err := a.AuthDBModel.GetExternalServiceIntegrationByID(integrationID)
	if err != nil {
		return err
	}
	// Use integration.ApiKey to sync data...
	return nil
}

func (a *DefaultAuthService) CreateExternalServiceIntegration(req *models.ExternalServiceIntegration) error {
	integration, err := a.AuthDBModel.GetExternalServiceIntegrationByID(integrationID)
	if err != nil {
		return err
	}
	// Use integration.ApiKey to sync data...
	return nil
}
