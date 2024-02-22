package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/shuttlersit/service-desk/backend/middleware"
	"github.com/shuttlersit/service-desk/backend/models"
	"github.com/shuttlersit/service-desk/backend/services"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Initialize Google OAuth configuration with environmental variables
var googleOAuthConfig = &oauth2.Config{
	RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	Endpoint:     google.Endpoint,
}

type GoogleAuthController struct {
	GoogleAuthService *services.GoogleAuthService
}

func NewAGoogleAuthController(googleAuthService *services.GoogleAuthService) *GoogleAuthController {
	return &GoogleAuthController{
		GoogleAuthService: googleAuthService,
	}
}

// GoogleLogin initiates the OAuth flow by redirecting the user to Google's OAuth server.
func GoogleLogin(c *gin.Context) {
	state := middleware.GenerateStateOauthCookie(c)
	url := googleOAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleAuthCallback handles the callback from Google after user has authenticated.
func (s *GoogleAuthController) GoogleAuthCallback(c *gin.Context) {
	// Validate state parameter matches the CSRF token in the cookie
	if !middleware.ValidateStateOauthCookie(c) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "State validation failed"})
		return
	}

	code := c.Query("code")
	token, err := googleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange token: " + err.Error()})
		return
	}

	userInfo, err := fetchGoogleUserInfo(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Process the userInfo to create or update the user in your system
	user, err := s.ProcessUserOAuth(&userInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process user info"})
		return
	}

	// Generate JWT token for the user
	jwtToken, err := middleware.GenerateJWT(user.Email, "User")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT"})
		return
	}

	// Set JWT token in session or return it in the response
	session := sessions.Default(c)
	session.Set("jwt_token", jwtToken)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User authenticated successfully",
		"token":   jwtToken,
		"user":    user,
	})
}

// processUserOAuth processes the OAuth user data, either registering a new user or updating an existing one.
func (s *GoogleAuthController) ProcessUserOAuth(userInfo *models.GoogleUserInfo) (*models.Users, error) {
	fname, lastname := getFirstAndLastName(userInfo.Name)
	//var user models.Users
	user, err := s.GoogleAuthService.AuthDBModel.UserDBModel.GetUserByEmail(userInfo.Email)
	if err != nil {
		// New user registration
		user = &models.Users{
			Email:     userInfo.Email,
			FirstName: fname,
			LastName:  lastname,
			// Populate additional fields as necessary
		}
		if err := s.GoogleAuthService.AuthDBModel.UserDBModel.CreateUser(user); err != nil {
			return nil, err
		}
	} else {
		// Existing user - update details
		user.FirstName = fname
		user.LastName = lastname
		if err := s.GoogleAuthService.AuthDBModel.UserDBModel.UpdateUser(user); err != nil {
			return nil, err
		}
	}

	return user, nil
}

// fetchGoogleUserInfo uses the OAuth token to fetch user information from Google.
func fetchGoogleUserInfo(token *oauth2.Token) (models.GoogleUserInfo, error) {
	client := googleOAuthConfig.Client(context.Background(), token)
	response, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return models.GoogleUserInfo{}, fmt.Errorf("failed to get user info: %v", err)
	}
	defer response.Body.Close()

	userInfo, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return models.GoogleUserInfo{}, fmt.Errorf("failed to read user info response: %v", err)
	}

	var user models.GoogleUserInfo
	if err := json.Unmarshal(userInfo, &user); err != nil {
		return models.GoogleUserInfo{}, fmt.Errorf("failed to unmarshal user info: %v", err)
	}

	return user, nil
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
