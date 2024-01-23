// backend/services/advertisement_service.go

package services

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/shuttlersit/service-desk/backend/models"
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
