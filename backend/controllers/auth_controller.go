package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shuttlersit/service-desk/backend/models"
	"github.com/shuttlersit/service-desk/backend/services"
)

type AuthController struct {
	AuthService *services.DefaultAuthService
}

func NewAuthController(authService *services.DefaultAuthService) *AuthController {
	return &AuthController{
		AuthService: authService,
	}
}

// User registration
func (a *AuthController) Registration(c *gin.Context) {
	var user models.Users
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newUser, token, err := a.AuthService.Registration(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "token": token, "loggedInUser": newUser})
}

// User login
func (a *AuthController) Login(c *gin.Context) {
	var loginInfo *services.LoginInfo
	loginInfo.Email = c.PostForm("email")
	loginInfo.Password = c.PostForm("secret")
	if err := c.BindJSON(&loginInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := a.AuthService.Login(loginInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
