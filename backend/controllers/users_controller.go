package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shuttlersit/service-desk/backend/models"
	"github.com/shuttlersit/service-desk/backend/services"
)

type UserController struct {
	UserService *services.DefaultUserService
}

func NewUserDBController(userService *services.DefaultUserService) *UserController {
	return &UserController{
		UserService: userService,
	}
}

// CreateUser handles the HTTP request to create a new user.
func (pc *UserController) CreateUser(ctx *gin.Context) {
	var newUser models.Users
	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err := pc.UserService.CreateUser(&newUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

// GetUserByID handles the HTTP request to retrieve a user by ID.
func (pc *UserController) GetUserByID(ctx *gin.Context) {
	userID, _ := strconv.Atoi(ctx.Param("id"))
	user, err := pc.UserService.GetUserByID(uint(userID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

// UpdateUser handles PUT /users/:id route.
func (pc *UserController) UpdateUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var ad models.Users
	if err := ctx.ShouldBindJSON(&ad); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ad.ID = uint(id)

	updatedAd, err := pc.UserService.UpdateUser(&ad)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedAd)
}

// DeleteUser handles DELETE /users/:id route.
func (pc *UserController) DeleteUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	status, err := pc.UserService.DeleteUser(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, status)
}

// GetAgentByID handles the HTTP request to retrieve a agents by ID.
func (pc *UserController) GetAllUsers(ctx *gin.Context) {
	users, err := pc.UserService.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "users not found"})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

//c.Params.ByName("id")
//ctx.Param("id")

// Implement controller methods like GetUsers, CreateUsers, GetUser, UpdateUser, DeleteUser
