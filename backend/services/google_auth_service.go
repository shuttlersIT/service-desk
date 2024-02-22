package services

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/shuttlersit/service-desk/models"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

// Assuming the structure of the User model includes necessary fields like Email, PasswordHash, etc.
// Assuming GoogleOAuthConfig is properly configured with the correct RedirectURL, ClientID, ClientSecret, and Scopes.

// RandToken generates a random @l length token.
func RandToken(l int) (string, error) {
	b := make([]byte, l)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func getLoginURL(state string) string {
	return conf.AuthCodeURL(state)
}

func init() {
	cid := "946670882701-dcidm9tcfdpcikpbjj8rfsb6uci22o4s.apps.googleusercontent.com"
	cs := "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"

	conf = &oauth2.Config{
		ClientID:     cid,
		ClientSecret: cs,
		RedirectURL:  "https://intelligence.shuttlers.africa/rest/oauth2-credential/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
		},
		Endpoint: google.Endpoint,
	}
}

type GoogleAuthServiceInterface interface {
	AuthenticateUser(code string) (*models.Users, error)
	RefreshAccessToken(refreshToken string) (string, error)
	RevokeAccessToken(accessToken string) error
	ProcessUserOAuth(userInfo *models.GoogleUserInfo) (*models.Users, error)
}

var c string
var conf *oauth2.Config

// GoogleAuthService handles Google OAuth authentication and user registration.
type GoogleAuthService struct {
	DB                *gorm.DB
	AuthDBModel       *models.AuthDBModel
	GoogleOAuthConfig *models.GoogleCredentialsDBModel
}

// NewGoogleAuthService creates a new instance of GoogleAuthService.
func NewGoogleAuthService(db *gorm.DB, AuthDBModel *models.AuthDBModel, googleOAuthConfig *models.GoogleCredentialsDBModel) *GoogleAuthService {
	return &GoogleAuthService{
		DB:                db,
		AuthDBModel:       AuthDBModel,
		GoogleOAuthConfig: googleOAuthConfig,
	}
}

// Registration handles new user registration.
func (s *AuthService) GoogleRegistration(user *models.Users, password string) (string, *models.Users, error) {
	// Check if user already exists
	var existingUser models.Users
	err := s.DB.Where("email = ?", user.Email).First(&existingUser).Error
	if err == nil {
		return "", nil, fmt.Errorf("user already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil, err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil, err
	}
	user.Credentials.PasswordHash = string(hashedPassword)

	// Save user to DB
	if err := s.DB.Create(user).Error; err != nil {
		return "", nil, err
	}

	// Generate JWT token
	token, err := s.generateJWTToken(user.ID)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

// Login authenticates a user and returns a JWT token.
func (s *AuthService) GoogleLogin(email, password string) (string, *models.Users, error) {
	var user models.Users
	if err := s.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return "", nil, err
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Credentials.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("incorrect password")
	}

	// Generate JWT token
	token, err := s.generateJWTToken(user.ID)
	if err != nil {
		return "", nil, err
	}

	// Set the token in the user's session or response (you can adapt this part to your application)
	//user.Token = token

	return token, &user, nil
}

// generateJWTToken generates a JWT token for authenticated sessions.
func (s *AuthService) generateJWTToken(userID uint) (string, error) {
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		Issuer:    "your_project",
		Subject:   fmt.Sprintf("%d", userID),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte("YOUR_SECRET_KEY")) // Ensure this is securely managed

	return signedToken, err
}

// fetchGoogleUserInfo fetches user info from Google using the provided OAuth token.
func fetchGoogleUserInfo(client *http.Client) (*models.GoogleUserInfo, error) {
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo models.GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

// processUserOAuth processes the OAuth user data, either registering a new user or updating an existing one.
func (s *GoogleAuthService) ProcessUserOAuth(userInfo *models.GoogleUserInfo) (*models.Users, error) {
	fname, lastname := getFirstAndLastName(userInfo.Name)
	var user models.Users
	if err := s.DB.Where("email = ?", userInfo.Email).First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {

		// New user registration
		user = models.Users{
			Email:     userInfo.Email,
			FirstName: fname,
			LastName:  lastname,
			// Populate additional fields as necessary
		}
		if err := s.DB.Create(&user).Error; err != nil {
			return nil, err
		}
	} else {
		// Existing user - update details
		user.FirstName = fname
		user.LastName = lastname
		if err := s.DB.Save(&user).Error; err != nil {
			return nil, err
		}
	}

	return &user, nil
}

/*
// AuthenticateViaGoogle processes the OAuth callback from Google.
func (s *GoogleAuthService) AuthenticateViaGoogle(c *gin.Context, code string) (*models.Users, error) {
	token, err := s.GoogleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}

	client := s.GoogleOAuthConfig.Client(context.Background(), token)
	userInfo, err := fetchGoogleUserInfo(client)
	if err != nil {
		return nil, err
	}

	user, err := s.ProcessUserOAuth(userInfo)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// AuthenticateWithGoogle processes the OAuth callback from Google.
func (s *GoogleAuthService) AuthenticateViaGoogle2(c *gin.Context, code string) (*models.Users, error) {
	token, err := googleOauthConfig.googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}

	client := GoogleOAuthConfig.Client(context.Background(), token)
	userInfo, err := fetchGoogleUserInfo(client)
	if err != nil {
		return nil, err
	}

	user, err := s.processUserOAuth2(userInfo)
	if err != nil {
		return nil, err
	}

	return user, nil
}
*/

// fetchGoogleUserInfo fetches user info from Google using the provided OAuth token.
func fetchGoogleUser2(client *http.Client) (*models.GoogleUserInfo, error) {
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo models.GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

// processOAuthUser processes the OAuth user data, either registering a new user or updating an existing one.
func (s *GoogleAuthService) processUserOAuth2(userInfo *models.GoogleUserInfo) (*models.Users, error) {
	fname, lastname := getFirstAndLastName(userInfo.Name)
	var user models.Users
	if err := s.DB.Where("email = ?", userInfo.Email).First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {

		// New user registration
		user = models.Users{
			Email:     userInfo.Email,
			FirstName: fname,
			LastName:  lastname,
			// Populate additional fields as necessary
		}
		if err := s.DB.Create(&user).Error; err != nil {
			return nil, err
		}
	} else {
		// Existing user - update details
		user.FirstName = fname
		user.LastName = lastname
		if err := s.DB.Save(&user).Error; err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func getFirstAndLastName(fullName string) (string, string) {
	// Split the full name into first name and last name
	parts := strings.Fields(fullName)
	var firstName, lastName string

	if len(parts) > 0 {
		firstName = parts[0]
	}

	if len(parts) > 1 {
		lastName = parts[len(parts)-1]
	}

	return firstName, lastName
}
