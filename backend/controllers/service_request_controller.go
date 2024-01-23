package controllers

import (
	"net/http"
	"strconv"

	"github.com/shuttlersit/service-desk/backend/models"
	"github.com/shuttlersit/service-desk/backend/services"

	"github.com/gin-gonic/gin"
)

type ServiceRequestController struct {
	Service *services.DefaultServiceRequestService
}

func NewServiceRequestController(service *services.DefaultServiceRequestService) *ServiceRequestController {
	return &ServiceRequestController{
		Service: service,
	}
}

func (c *ServiceRequestController) CreateServiceRequest(ctx *gin.Context) {
	var request models.ServiceRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.Service.CreateServiceRequest(&request); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, request)
}

func (c *ServiceRequestController) UpdateServiceRequest(ctx *gin.Context) {
	requestID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	var request models.ServiceRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingRequest, err := c.Service.GetServiceRequestByID(uint(requestID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Service request not found"})
		return
	}

	existingRequest.Title = request.Title
	existingRequest.Description = request.Description
	existingRequest.Status = request.Status
	existingRequest.CategoryID = request.CategoryID
	existingRequest.SubCategoryID = request.SubCategoryID

	if err := c.Service.UpdateServiceRequest(existingRequest); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, existingRequest)
}

func (c *ServiceRequestController) GetUserRequests(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("user_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	requests, err := c.Service.GetUserServiceRequests(uint(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, requests)
}

func (c *ServiceRequestController) CloseServiceRequest(ctx *gin.Context) {
	requestID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	if err := c.Service.CloseServiceRequest(uint(requestID)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Service request closed successfully"})
}

func (c *ServiceRequestController) GetServiceRequestHistory(ctx *gin.Context) {
	requestID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	history, err := c.Service.GetServiceRequestHistory(uint(requestID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, history)
}

func (c *ServiceRequestController) AddCommentToServiceRequest(ctx *gin.Context) {
	requestID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	var comment models.ServiceRequestComment
	if err := ctx.ShouldBindJSON(&comment); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.Service.AddCommentToServiceRequest(uint(requestID), comment.Comment); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, comment)
}

func (c *ServiceRequestController) GetServiceRequestComments(ctx *gin.Context) {
	requestID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	comments, err := c.Service.GetServiceRequestComments(uint(requestID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, comments)
}
