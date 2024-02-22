// backend/controllers/auth_controllers.go

package controllers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shuttlersit/service-desk/backend/config"
	"github.com/shuttlersit/service-desk/backend/middleware"
	"github.com/shuttlersit/service-desk/backend/models"
	"github.com/shuttlersit/service-desk/backend/services"
	"golang.org/x/oauth2"
)

// Initialize configuration globally or pass it to your controller struct
var cfg = config.LoadConfig()

type AuthController struct {
	AuthService *services.AuthService
	AuthDBModel *models.AuthDBModel
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		AuthService: authService,
	}
}

func GoogleLoginAlt(c *gin.Context) {
	state := "pseudo-random" // Generate or use a library to create a secure state
	url := cfg.GoogleOAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallbackAlt(c *gin.Context) {
	// Error handling omitted for brevity
	code := c.Query("code")
	token, _ := cfg.GoogleOAuthConfig.Exchange(context.Background(), code)
	userInfo, _ := fetchGoogleUserInfo(token)

	// Process userInfo (create or find user, generate JWT)
	jwtToken, _ := middleware.GenerateJWT(userInfo.Email, "User")

	// Redirect or send JWT as needed
	c.JSON(http.StatusOK, gin.H{"token": jwtToken})
}

// Register handles user registration
func (ac *AuthController) Register(c *gin.Context) {
	var user models.Users
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUser, token, err := ac.AuthService.RegisterUser(&user, c.Param("password"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": newUser, "token": token})
}

// Login handles user login
func (ac *AuthController) Login(c *gin.Context) {
	var loginInfo services.LoginInfo
	if err := c.ShouldBindJSON(&loginInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, _, err := ac.AuthService.Login(&loginInfo)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Login handles user login
func (ac *AuthController) LoginAgent(c *gin.Context) {
	var loginInfo services.LoginInfo
	if err := c.ShouldBindJSON(&loginInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, _, err := ac.AuthService.Login(&loginInfo)
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

	newAgent, token, err := ac.AuthService.RegisterAgent(&agent, c.Param("password"))
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

/*
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
*/

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

type User struct {
	ID       int
	Username string
	Email    string
	Password string
}

// User registration route
/*func GoogleRegister(c *gin.Context) {
	var newUser User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// For simplicity, we store passwords in plain text. In a real application, hash and salt passwords for security.
	newUser.ID = len(users) + 1
	users = append(users, newUser)

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}*/

/*
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
*/

/*

func (ac *AuthController) CreateExternalServiceIntegrationHandler(c *gin.Context) {
	var req models.ExternalServiceIntegration
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ac.AuthService.CreateExternalServiceIntegration(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, req)
}
*/
