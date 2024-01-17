package controllers

import (

	// "github.com/gin-gonic/gin"
	"github.com/shuttlersit/service-desk/backend/models"
)

type UserController struct {
	User *models.UserDBModel
}

func NewUserController() *UserController {
	return &UserController{}
}

// Implement controller methods like GetUsers, CreateUsers, GetUser, UpdateUser, DeleteUser
