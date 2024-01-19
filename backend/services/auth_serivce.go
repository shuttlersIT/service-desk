// backend/services/advertisement_service.go

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
}

// DefaultAuthService is the default implementation of AuthService
type DefaultAuthService struct {
	DB          *gorm.DB
	AuthDBModel *models.AuthDBModel
	UserDBModel *models.UserDBModel
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

	// Hash the user's password before storing it in the database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Credentials.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash password")
	}
	user.Credentials.Password = string(hashedPassword)
	erro := a.UserDBModel.CreateUser(user)
	if erro != nil {
		return nil, "", fmt.Errorf("failed to create users")
	}
	er := a.AuthDBModel.CreateUserCredentials(&user.Credentials)
	if er != nil {
		return nil, "", fmt.Errorf("failed to create users credentials")
	}

	e := a.UserDBModel.CreateUser(user)
	if e != nil {
		return nil, "", e
	}
	// Generate a JWT token for successful login
	token, err := generateJWTToken(user.ID)
	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}

// User login
func (a *DefaultAuthService) Login(login *LoginInfo) (string, error) {
	loginInfo := login
	var user models.Users
	if err := a.DB.Where("email = ?", loginInfo.Email).First(&user).Error; err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Credentials.Password), []byte(loginInfo.Password)); err != nil {
		return "", fmt.Errorf("invalid email or password")
	}
	// Generate a JWT token for successful login
	token, err := generateJWTToken(user.ID)
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
