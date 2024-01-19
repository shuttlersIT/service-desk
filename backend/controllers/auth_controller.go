// backend/controllers/auth_controllers.go

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
