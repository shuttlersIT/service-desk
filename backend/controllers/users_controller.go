// backend/controllers/users_controllers.go

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

func NewUserController(userService *services.DefaultUserService) *UserController {
	return &UserController{
		UserService: userService,
	}
}

// CreateUserHandler handles the HTTP request to create a new user.
func (ctrl *UserController) CreateUserHandler(c *gin.Context) {
	var user models.Users
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err := ctrl.UserService.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// UpdateUserHandler handles the HTTP request to update an existing user.
func (ctrl *UserController) UpdateUserHandler(c *gin.Context) {
	var user models.Users
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	updatedUser, err := ctrl.UserService.UpdateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

// GetUserByIDHandler retrieves a user by their ID.
func (ctrl *UserController) GetUserByIDHandler(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("user_id"))

	user, err := ctrl.UserService.GetUserByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUserHandler deletes a user by their ID.
func (ctrl *UserController) DeleteUserHandler(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("user_id"))

	deleted, err := ctrl.UserService.DeleteUser(uint(userID))
	if err != nil || !deleted {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found or failed to delete"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// GetAllUsersHandler retrieves all users.
func (ctrl *UserController) GetAllUsersHandler(c *gin.Context) {
	users, _ := ctrl.UserService.GetAllUsers()
	c.JSON(http.StatusOK, users)
}

// CreatePositionHandler creates a new position.
func (ctrl *UserController) CreatePositionHandler(c *gin.Context) {
	var position models.Position
	if err := c.ShouldBindJSON(&position); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err := ctrl.UserService.CreatePosition(&position)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create position"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Position created successfully"})
}

// DeletePositionHandler deletes a position by its ID.
func (ctrl *UserController) DeletePositionHandler(c *gin.Context) {
	positionID, _ := strconv.Atoi(c.Param("position_id"))

	err := ctrl.UserService.DeletePosition(uint(positionID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Position not found or failed to delete"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Position deleted successfully"})
}

// UpdatePositionHandler updates an existing position.
func (ctrl *UserController) UpdatePositionHandler(c *gin.Context) {
	var position models.Position
	if err := c.ShouldBindJSON(&position); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err := ctrl.UserService.UpdatePosition(&position)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update position"})
		return
	}

	c.JSON(http.StatusOK, position)
}

// GetPositionsHandler retrieves all positions.
func (ctrl *UserController) GetPositionsHandler(c *gin.Context) {
	positions, err := ctrl.UserService.GetAllPositions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve positions"})
		return
	}

	c.JSON(http.StatusOK, positions)
}

// GetPositionByIDHandler retrieves a position by its ID.
func (ctrl *UserController) GetPositionByIDHandler(c *gin.Context) {
	positionID, _ := strconv.Atoi(c.Param("position_id"))

	position, err := ctrl.UserService.GetPositionByID(uint(positionID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Position not found"})
		return
	}

	c.JSON(http.StatusOK, position)
}

// GetPositionByNumberHandler retrieves a position by its number.
func (ctrl *UserController) GetPositionByNumberHandler(c *gin.Context) {
	positionNumber, _ := strconv.Atoi(c.Param("position_number"))

	position, err := ctrl.UserService.GetPositionByNumber(positionNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Position not found"})
		return
	}

	c.JSON(http.StatusOK, position)
}

// CreateDepartmentHandler creates a new department.
func (ctrl *UserController) CreateDepartmentHandler(c *gin.Context) {
	var department models.Department
	if err := c.ShouldBindJSON(&department); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err := ctrl.UserService.CreateDepartment(&department)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create department"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Department created successfully"})
}

// DeleteDepartmentHandler deletes a department by its ID.
func (ctrl *UserController) DeleteDepartmentHandler(c *gin.Context) {
	departmentID, _ := strconv.Atoi(c.Param("department_id"))

	err := ctrl.UserService.DeleteDepartment(uint(departmentID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Department not found or failed to delete"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Department deleted successfully"})
}

// UpdateDepartmentHandler updates an existing department.
func (ctrl *UserController) UpdateDepartmentHandler(c *gin.Context) {
	var department models.Department
	if err := c.ShouldBindJSON(&department); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err := ctrl.UserService.UpdateDepartment(&department)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update department"})
		return
	}

	c.JSON(http.StatusOK, department)
}

// GetDepartmentsHandler retrieves all departments.
func (ctrl *UserController) GetDepartmentsHandler(c *gin.Context) {
	departments, err := ctrl.UserService.GetAllDepartments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve departments"})
		return
	}

	c.JSON(http.StatusOK, departments)
}

// GetDepartmentByIDHandler retrieves a department by its ID.
func (ctrl *UserController) GetDepartmentByIDHandler(c *gin.Context) {
	departmentID, _ := strconv.Atoi(c.Param("department_id"))

	department, err := ctrl.UserService.GetDepartmentByID(uint(departmentID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Department not found"})
		return
	}

	c.JSON(http.StatusOK, department)
}

// GetDepartmentByNumberHandler retrieves a department by its number.
func (ctrl *UserController) GetDepartmentByNumberHandler(c *gin.Context) {
	departmentNumber, _ := strconv.Atoi(c.Param("department_number"))

	department, err := ctrl.UserService.GetDepartmentByNumber(departmentNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Department not found"})
		return
	}

	c.JSON(http.StatusOK, department)
}

////////////////////////////////////////////////////////////////////////////////////////
