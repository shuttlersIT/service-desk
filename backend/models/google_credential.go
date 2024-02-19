// models/google_credentials.go

package models

import (
	"context"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

type GoogleCredentials struct {
	Cid     string `json:"cid"`
	Csecret string `json:"csecret"`
}

// TableName sets the table name for the Credentials model.
func (GoogleCredentials) TableName() string {
	return "googleCredentials"
}

type Credentials struct {
	Cid     string `json:"cid"`
	Csecret string `json:"csecret"`
}

type OAuth2Credentials struct {
	gorm.Model
	UserID      uint      `json:"user_id" gorm:"index;not null"`
	ServiceName string    `json:"service_name" gorm:"type:varchar(255);not null"`
	Token       string    `json:"token" gorm:"type:text;not null"`
	ExpiresAt   time.Time `json:"expires_at"`
}

func (OAuth2Credentials) TableName() string {
	return "oauth2_credentials"
}

type GoogeAuthStorage interface {
	CreateGoogleCred(*GoogleCredentials) error
	DeleteGoogleCred(int) error
	UpdateGoogleCred(*GoogleCredentials) error
	GetGoogleCreds() ([]*GoogleCredentials, error)
	GetGoogleCredByID(int) (*GoogleCredentials, error)
	GetGoogleCredByNumber(int) (*GoogleCredentials, error)
}

var oauth2Config = &oauth2.Config{
	ClientID:     "YOUR_CLIENT_ID",
	ClientSecret: "YOUR_CLIENT_SECRET",
	RedirectURL:  "YOUR_REDIRECT_URL",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

var googleOauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8080/auth/google/callback",
	ClientID:     "YOUR_CLIENT_ID",
	ClientSecret: "YOUR_CLIENT_SECRET",
	Scopes:       []string{"email", "profile"},
	Endpoint:     google.Endpoint,
}

// UserModel handles database operations for User
type GoogleCredentialsDBModel struct {
	DB                *gorm.DB
	log               *Logger
	AuthDBModel       *AuthDBModel
	googleOauthConfig *GoogleCredentials
}

// NewUserModel creates a new instance of UserModel
func NewGoogleCredentialsDBModel(db *gorm.DB, AuthDBModel *AuthDBModel, log *Logger) *GoogleCredentialsDBModel {
	return &GoogleCredentialsDBModel{
		DB:          db,
		AuthDBModel: AuthDBModel,
		log:         log,
	}
}

func googleAuthHandler(c *gin.Context) {
	url := googleOauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func googleAuthCallbackHandler(c *gin.Context) {
	code := c.Query("code")
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to exchange token"})
		return
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read user info"})
		return
	}

	// Here you should parse the contents and create or update the user in your DB
	// For example, you could JSON unmarshal contents and extract user info

	c.String(http.StatusOK, string(contents)) // For demonstration purposes
}

var jwtKey = []byte("your_secret_key")

func generateJWT(email string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	return tokenString, err
}

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

/*
func (as *AuthDBModel) CreateOAuth2Integration(serviceName string, credentials *OAuth2Credentials) error {
	// Assuming OAuth2Credentials includes fields like ClientID, ClientSecret
	return as.DB.Create(&OAuth2Integration{
		ServiceName: serviceName,
		Credentials: *credentials,
	}).Error
}

func (as *AuthDBModel) AuthenticateWithOAuth2(serviceName string, token string) (*Users, error) {
	// Logic to authenticate a user with an external OAuth2 service
	// This might involve token validation with the external service
	// Example code skipped for brevity
	return &Users{}, nil
}

func (as *AuthDBModel) IncrementFailedLoginAttempts(userID uint) error {
	// Increment failed login attempts for user and check for lockout
	// Example code skipped for brevity
	return nil
}

func (as *AuthDBModel) CheckAccountLockoutStatus(userID uint) (bool, error) {
	// Check if user account is locked due to too many failed login attempts
	// Example code skipped for brevity
	return false, nil
}

func (as *AuthDBModel) ResetFailedLoginAttempts(userID uint) error {
	// Reset the failed login attempts counter for a user
	// Example code skipped for brevity
	return nil
}
*/
