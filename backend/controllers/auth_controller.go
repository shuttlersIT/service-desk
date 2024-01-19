// backend/controllers/auth_controllers.go

package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shuttlersit/service-desk/backend/models"
	"github.com/shuttlersit/service-desk/backend/services"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthController struct {
	AuthService *services.DefaultAuthService
}

func NewAuthController(authService *services.DefaultAuthService) *AuthController {
	return &AuthController{
		AuthService: authService,
	}
}

func (ac *AuthController) Registration(ctx *gin.Context) {
	var user models.Users
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newUser, token, err := ac.AuthService.Registration(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "token": token, "loggedInUser": newUser})
}

func (ac *AuthController) Login(ctx *gin.Context) {
	var loginInfo *services.LoginInfo
	loginInfo.Email = ctx.PostForm("email")
	loginInfo.Password = ctx.PostForm("secret")
	if err := ctx.BindJSON(&loginInfo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := ac.AuthService.Login(loginInfo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (ac *AuthController) Logout(ctx *gin.Context) {
	// Implement logout logic here
}

/////////////////////////////////////////////////////////////////////////////////////////////

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "YOUR_REDIRECT_URI",
		ClientID:     "YOUR_CLIENT_ID",
		ClientSecret: "YOUR_CLIENT_SECRET",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	users []User
)

type User struct {
	ID       int
	Username string
	Email    string
	Password string
}

// User registration route
func GoogleRegister(c *gin.Context) {
	var newUser User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// For simplicity, we store passwords in plain text. In a real application, hash and salt passwords for security.
	newUser.ID = len(users) + 1
	users = append(users, newUser)

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// Google Sign-In route
func GoogleLogin(c *gin.Context) {
	url := googleOauthConfig.AuthCodeURL("", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusFound, url)
}

/*
// Google Sign-In callback route
func GoogleCallback(c *gin.Context) {
	code := c.DefaultQuery("code", "")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code parameter"})
		return
	}

	token, err := googleOauthConfig.Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	// Fetch user data from Google using the token
	client := googleOauthConfig.Client(c, token)
	response, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer response.Body.Close()

	// Parse user info from the response (you might want to validate and sanitize this data)
	var googleUser GoogleUserInfo
	if err := json.NewDecoder(response.Body).Decode(&googleUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user info"})
		return
	}

	// Check if a user with the same email exists
	var existingUser User
	for _, u := range users {
		if u.Email == googleUser.Email {
			existingUser = u
			break
		}
	}

	if existingUser.ID != 0 {
		// Handle user login or session management here
	} else {
		// Create a new user account with googleUser data
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sign-In with Google Successful", "user": googleUser})
}
*/
