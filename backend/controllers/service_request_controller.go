package controllers

import (
	"github.com/shuttlersit/service-desk/backend/models"
	"github.com/shuttlersit/service-desk/backend/services"

	"github.com/gin-gonic/gin"
)

type ServiceRequestController struct {
	Service services.ServiceRequestService
}

func (controller *ServiceRequestController) CreateServiceRequest(c *gin.Context) {
	// Parse request and create a new service request
	var request models.ServiceRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		// Handle validation errors
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := controller.Service.CreateServiceRequest(&request); err != nil {
		// Handle database errors
		c.JSON(500, gin.H{"error": "Failed to create service request"})
		return
	}

	c.JSON(200, request)
}

// Implement other methods like UpdateServiceRequest, GetServiceRequestByID, and GetUserServiceRequests.
