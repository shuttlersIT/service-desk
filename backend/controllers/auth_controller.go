// backend/controllers/auth_controllers.go

package controllers

import (
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shuttlersit/service-desk/backend/models"
	"github.com/shuttlersit/service-desk/backend/services"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthController struct {
	AuthService *services.DefaultAuthService
	AuthDBModel *models.AuthDBModel
}

func NewAuthController(authService *services.DefaultAuthService) *AuthController {
	return &AuthController{
		AuthService: authService,
	}
}

// Register handles user registration
func (ac *AuthController) Register(c *gin.Context) {
	var user models.Users
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUser, token, err := ac.AuthService.Registration(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": newUser, "token": token})
}

// Login handles user login
func (ac *AuthController) Login(c *gin.Context) {
	var loginInfo models.LoginInfo
	if err := c.ShouldBindJSON(&loginInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := ac.AuthService.Login(&loginInfo)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Login handles user login
func (ac *AuthController) LoginAgent(c *gin.Context) {
	var loginInfo models.LoginInfo
	if err := c.ShouldBindJSON(&loginInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := ac.AuthService.LoginAgent(&loginInfo)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Register handles user registration
func (ac *AuthController) RegisterAgent(c *gin.Context) {
	var agent models.Agents
	if err := c.ShouldBindJSON(&agent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newAgent, token, err := ac.AuthService.AgentRegistration(&agent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"agent": newAgent, "token": token})
}

// ResetPassword handles resetting the user's password
func (ac *AuthController) ResetPassword(c *gin.Context) {
	// Parse user ID from request parameters
	userIDParam := c.Param("user_id")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var newPassword struct {
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&newPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = ac.AuthService.ResetUserPassword(uint(userID), newPassword.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

// ResetPassword handles resetting the user's password
func (ac *AuthController) ResetAgentPassword(c *gin.Context) {
	// Parse user ID from request parameters
	agentIDParam := c.Param("agent_id")
	agentID, err := strconv.Atoi(agentIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid agent ID"})
		return
	}

	var newPassword struct {
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&newPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = ac.AuthService.ResetAgentPassword(uint(agentID), newPassword.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

// ResetUserPassword resets the user's password.
func (ac *AuthController) ResetUserPassword(ctx *gin.Context) {
	userID, _ := strconv.Atoi(ctx.Param("id"))
	var newPassword struct {
		NewPassword string `json:"new_password"`
	}
	if err := ctx.BindJSON(&newPassword); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := ac.AuthService.ResetUserPassword(uint(userID), newPassword.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset password"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

// ResetUserPassword2 resets the user's password using an alternative implementation.
func (ac *AuthController) ResetUserPassword2(ctx *gin.Context) {
	userID, _ := strconv.Atoi(ctx.Param("id"))
	var newPassword struct {
		NewPassword string `json:"new_password"`
	}
	if err := ctx.BindJSON(&newPassword); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := ac.AuthService.ResetUserPassword(uint(userID), newPassword.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset password"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Password reset successfully (Alternative method)"})
}

func (ac *AuthController) Logout(ctx *gin.Context) {
	// Implement logout logic here
}

// RequestPasswordResetToken sends a password reset token to the user's email.
func (ac *AuthController) RequestPasswordResetToken(ctx *gin.Context) {
	var userEmail struct {
		Email string `json:"email"`
	}
	if err := ctx.BindJSON(&userEmail); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Validate the user's email and send a password reset token via email.
	// Implement this logic using your preferred email service.
	// You can generate a unique token and send it to the user's email.
	// Don't forget to include a link that allows the user to reset their password.
	// Once the email is sent, respond with a success message.
	ctx.JSON(http.StatusOK, gin.H{"message": "Password reset token sent successfully"})
}

// ResetPasswordWithToken handles resetting the user's password using a token
func (ac *AuthController) ResetPasswordWithToken2(c *gin.Context) {
	token := c.Param("token")
	var newPassword struct {
		Email       string `json:"email"`
		NewPassword string `json:"password"`
	}
	if err := c.ShouldBindJSON(&newPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ac.AuthService.ResetPasswordWithToken(token, newPassword.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

// ResetPasswordWithToken resets the user's password using a valid token.
func (ac *AuthController) ResetPasswordWithToken(ctx *gin.Context) {
	token := ctx.Param("token")
	var newPassword struct {
		NewPassword string `json:"new_password"`
	}
	if err := ctx.BindJSON(&newPassword); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Verify the validity of the password reset token.
	// If the token is valid, update the user's password.
	// Otherwise, respond with an error message.
	err := ac.AuthService.ResetPasswordWithToken(token, newPassword.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

// ChangePassword allows authenticated users to change their password.
func (ac *AuthController) ChangePassword(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(uint) // Get the user's ID from the token or session.
	var newPassword struct {
		NewPassword string `json:"new_password"`
	}
	if err := ctx.BindJSON(&newPassword); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Update the user's password with the new password.
	err := ac.AuthService.ResetUserPassword(userID, newPassword.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to change password"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
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

// ResetPasswordWithToken handles resetting the user's password using a token
func (ac *AuthController) ResetPasswordWithToken3(c *gin.Context) {
	token := c.Param("token")
	var newPassword struct {
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&newPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ac.AuthService.ResetPasswordWithToken(token, newPassword.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}
