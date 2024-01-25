package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/shuttlersit/service-desk/backend/models"
	"github.com/shuttlersit/service-desk/backend/services"

	"github.com/gin-gonic/gin"
)

type ServiceRequestController struct {
	ServiceRequest *services.DefaultServiceRequestService
}

func NewServiceRequestController(service *services.DefaultServiceRequestService) *ServiceRequestController {
	return &ServiceRequestController{
		ServiceRequest: service,
	}
}

// CreateServiceRequestHandler handles the creation of a new service request.
func (ctrl *ServiceRequestController) CreateServiceRequestHandler(c *gin.Context) {
	var request models.ServiceRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.ServiceRequest.CreateServiceRequest(&request); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create service request"})
		return
	}

	c.JSON(http.StatusCreated, request)
}

// GetServiceRequestByIDHandler retrieves a service request by its ID.
func (ctrl *ServiceRequestController) GetServiceRequestByIDHandler(c *gin.Context) {
	id := c.Param("id")
	requestID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service request ID"})
		return
	}

	serviceRequest, err := ctrl.ServiceRequest.GetServiceRequestByID(uint(requestID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service request not found"})
		return
	}

	c.JSON(http.StatusOK, serviceRequest)
}

// UpdateServiceRequestHandler updates an existing service request.
func (ctrl *ServiceRequestController) UpdateServiceRequestHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	request, err := ctrl.ServiceRequest.GetServiceRequestByID(uint(id))
	request.UpdatedAt = time.Now()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service request ID"})
		return
	}

	var updatedRequest models.ServiceRequest
	if err := c.ShouldBindJSON(&updatedRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the UpdateServiceRequest method from the ServiceRequestService.
	if err := ctrl.ServiceRequest.UpdateServiceRequest(&updatedRequest); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update service request"})
		return
	}

	c.JSON(http.StatusOK, updatedRequest)
}

// CloseServiceRequestHandler closes a service request by its ID.
func (ctrl *ServiceRequestController) CloseServiceRequestHandler(c *gin.Context) {
	id := c.Param("id")
	requestID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service request ID"})
		return
	}

	// Call the CloseServiceRequest method from the ServiceRequestService.
	if err := ctrl.ServiceRequest.CloseServiceRequest(uint(requestID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to close service request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service request closed successfully"})
}

// ReopenServiceRequestHandler reopens a closed service request by its ID.
func (ctrl *ServiceRequestController) ReopenServiceRequestHandler(c *gin.Context) {
	id := c.Param("id")
	requestID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service request ID"})
		return
	}

	// Call the ReopenServiceRequest method from the ServiceRequestService.
	if err := ctrl.ServiceRequest.ReopenServiceRequest(uint(requestID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reopen service request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service request reopened successfully"})
}

// GetServiceRequestByIDHandler retrieves all service request from the database.
func (ctrl *ServiceRequestController) GetAllServiceRequestsHandler(c *gin.Context) {
	serviceRequest, err := ctrl.ServiceRequest.GetAllServiceRequests()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve service requests"})
		return
	}

	c.JSON(http.StatusOK, serviceRequest)
}

// GetServiceRequestHistoryHandler retrieves the history entries of a service request by its ID.
func (ctrl *ServiceRequestController) GetServiceRequestHistoryHandler(c *gin.Context) {
	id := c.Param("id")
	requestID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service request ID"})
		return
	}

	// Call the GetServiceRequestHistory method from the ServiceRequestService.
	historyEntries, err := ctrl.ServiceRequest.GetServiceRequestHistory(uint(requestID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve service request history"})
		return
	}

	c.JSON(http.StatusOK, historyEntries)
}

// AddCommentToServiceRequestHandler adds a comment to a service request by its ID.
func (ctrl *ServiceRequestController) AddCommentToServiceRequestHandler(c *gin.Context) {
	id := c.Param("id")
	requestID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service request ID"})
		return
	}

	var commentRequest struct {
		Comment string `json:"comment"`
	}

	if err := c.ShouldBindJSON(&commentRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the AddCommentToServiceRequest method from the ServiceRequestService.
	if err := ctrl.ServiceRequest.AddCommentToServiceRequest(uint(requestID), commentRequest.Comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add comment to service request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment added successfully"})
}

// GetServiceRequestCommentsHandler retrieves comments of a service request by its ID.
func (ctrl *ServiceRequestController) GetServiceRequestCommentsHandler(c *gin.Context) {
	id := c.Param("service_request_id")
	requestID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service request ID"})
		return
	}

	// Call the GetServiceRequestComments method from the ServiceRequestService.
	comments, err := ctrl.ServiceRequest.GetServiceRequestComments(uint(requestID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve service request comments"})
		return
	}

	c.JSON(http.StatusOK, comments)
}

// GetServiceRequestsByCategoryHandler retrieves service requests by category.
func (ctrl *ServiceRequestController) GetServiceRequestsByCategoryHandler(c *gin.Context) {
	categoryID := c.Param("categoryID")
	categoryIDUint, err := strconv.ParseUint(categoryID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	// Call the GetServiceRequestsByCategory method from the ServiceRequestService.
	serviceRequests, err := ctrl.ServiceRequest.GetServiceRequestsByCategory(uint(categoryIDUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve service requests by category"})
		return
	}

	c.JSON(http.StatusOK, serviceRequests)
}

// GetServiceRequestsBySubcategoryHandler retrieves service requests by subcategory ID.
func (ctrl *ServiceRequestController) GetServiceRequestsBySubcategoryHandler(c *gin.Context) {
	subcategoryID := c.Param("subcategoryID")
	subcategoryIDUint, err := strconv.ParseUint(subcategoryID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subcategory ID"})
		return
	}

	// Call the GetServiceRequestsBySubCategory method from the ServiceRequestService.
	serviceRequests, err := ctrl.ServiceRequest.GetServiceRequestsBySubCategory(uint(subcategoryIDUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve service requests by subcategory"})
		return
	}

	c.JSON(http.StatusOK, serviceRequests)
}

// GetServiceRequestsByStatusHandler retrieves service requests by status.
func (ctrl *ServiceRequestController) GetServiceRequestsByStatusHandler(c *gin.Context) {
	status := c.Param("status")

	// Call the GetServiceRequestsByStatus method from the ServiceRequestService.
	serviceRequests, err := ctrl.ServiceRequest.GetServiceRequestsByStatus(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve service requests by status"})
		return
	}

	c.JSON(http.StatusOK, serviceRequests)
}

// GetServiceRequestsByPriorityHandler retrieves service requests by priority.
func (ctrl *ServiceRequestController) GetServiceRequestsByPriorityHandler(c *gin.Context) {
	priority := c.Param("priority")

	// Call the GetServiceRequestsByPriority method from the ServiceRequestService.
	serviceRequests, err := ctrl.ServiceRequest.GetServiceRequestsByPriority(priority)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve service requests by priority"})
		return
	}

	c.JSON(http.StatusOK, serviceRequests)
}

// GetServiceRequestsByUserHandler retrieves service requests by user ID.
func (ctrl *ServiceRequestController) GetServiceRequestsByUserHandler(c *gin.Context) {
	userID := c.Param("userID")
	userIDUint, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Call the GetServiceRequestsByUser method from the ServiceRequestService.
	serviceRequests, err := ctrl.ServiceRequest.GetUserServiceRequests(uint(userIDUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve service requests by user"})
		return
	}

	c.JSON(http.StatusOK, serviceRequests)
}

// GetServiceRequestsByCategoryAndStatusHandler retrieves service requests by both category and status.
func (ctrl *ServiceRequestController) GetServiceRequestsByCategoryAndStatusHandler(c *gin.Context) {
	categoryID := c.Param("categoryID")
	categoryIDUint, err := strconv.ParseUint(categoryID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	status := c.Param("status")

	// Call the GetServiceRequestsByCategoryAndStatus method from the ServiceRequestService.
	serviceRequests, err := ctrl.ServiceRequest.GetServiceRequestsByCategoryAndStatus(uint(categoryIDUint), status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve service requests by category and status"})
		return
	}

	c.JSON(http.StatusOK, serviceRequests)
}

// GetServiceRequestsBySubCategoryAndStatusHandler retrieves service requests by both subcategory and status.
func (ctrl *ServiceRequestController) GetServiceRequestsBySubCategoryAndStatusHandler(c *gin.Context) {
	subcategoryID := c.Param("subcategoryID")
	subcategoryIDUint, err := strconv.ParseUint(subcategoryID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subcategory ID"})
		return
	}

	status := c.Param("status")

	// Call the GetServiceRequestsBySubCategoryAndStatus method from the ServiceRequestService.
	serviceRequests, err := ctrl.ServiceRequest.GetServiceRequestsBySubCategoryAndStatus(uint(subcategoryIDUint), status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve service requests by subcategory and status"})
		return
	}

	c.JSON(http.StatusOK, serviceRequests)
}

// GetServiceRequestsBySubCategoryHandler retrieves service requests by sub-category.
func (ctrl *ServiceRequestController) GetServiceRequestsBySubCategoryHandler(c *gin.Context) {
	subCategoryID := c.Param("subCategoryID")
	subCategoryIDUint, err := strconv.ParseUint(subCategoryID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sub-category ID"})
		return
	}

	// Call the GetServiceRequestsBySubCategory method from the ServiceRequestService.
	serviceRequests, err := ctrl.ServiceRequest.GetServiceRequestsBySubCategory(uint(subCategoryIDUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve service requests by sub-category"})
		return
	}

	c.JSON(http.StatusOK, serviceRequests)
}

// GetServiceRequestsByCategoryAndSubCategoryHandler retrieves service requests for a specific category and sub-category.
func (ctrl *ServiceRequestController) GetServiceRequestsByCategoryAndSubCategoryHandler(c *gin.Context) {
	categoryID := c.Param("categoryID")
	categoryIDUint, err := strconv.ParseUint(categoryID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	subCategoryID := c.Param("subCategoryID")
	subCategoryIDUint, err := strconv.ParseUint(subCategoryID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sub-category ID"})
		return
	}

	// Call the GetServiceRequestsByCategoryAndSubCategory method from the ServiceRequestService.
	serviceRequests, err := ctrl.ServiceRequest.GetServiceRequestsByCategoryAndSubCategory(uint(categoryIDUint), uint(subCategoryIDUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve service requests for the category and sub-category"})
		return
	}

	c.JSON(http.StatusOK, serviceRequests)
}

// GetServiceRequestsByStatusAndLocationHandler retrieves service requests for a specific status and location.
func (ctrl *ServiceRequestController) GetServiceRequestsByStatusAndLocationHandler(c *gin.Context) {
	status := c.Param("status")
	locationID := c.Param("locationID")
	locationIDUint, err := strconv.ParseUint(locationID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid location ID"})
		return
	}

	// Call the GetServiceRequestsByStatusAndLocation method from the ServiceRequestService.
	serviceRequests, err := ctrl.ServiceRequest.GetServiceRequestsByStatusAndLocation(status, uint(locationIDUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve service requests for the status and location"})
		return
	}

	c.JSON(http.StatusOK, serviceRequests)
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

	existingRequest, err := c.ServiceRequest.GetServiceRequestByID(uint(requestID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Service request not found"})
		return
	}

	existingRequest.Title = request.Title
	existingRequest.Description = request.Description
	existingRequest.Status = request.Status
	existingRequest.CategoryID = request.CategoryID
	existingRequest.SubCategoryID = request.SubCategoryID

	if err := c.ServiceRequest.UpdateServiceRequest(existingRequest); err != nil {
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

	requests, err := c.ServiceRequest.GetUserServiceRequests(uint(userID))
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

	if err := c.ServiceRequest.CloseServiceRequest(uint(requestID)); err != nil {
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

	history, err := c.ServiceRequest.GetServiceRequestHistory(uint(requestID))
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

	if err := c.ServiceRequest.AddCommentToServiceRequest(uint(requestID), comment.Comment); err != nil {
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

	comments, err := c.ServiceRequest.GetServiceRequestComments(uint(requestID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, comments)
}

// GetOpenServiceRequestsHandler retrieves all open service requests.
func (ctrl *ServiceRequestController) GetOpenServiceRequestsHandler(c *gin.Context) {
	// Call the GetOpenServiceRequests method from the ServiceRequestService.
	openServiceRequests, err := ctrl.ServiceRequest.GetOpenServiceRequests()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve open service requests"})
		return
	}

	c.JSON(http.StatusOK, openServiceRequests)
}

// GetClosedServiceRequestsHandler retrieves all closed service requests.
func (ctrl *ServiceRequestController) GetClosedServiceRequestsHandler2(c *gin.Context) {
	// Call the GetClosedServiceRequests method from the ServiceRequestService.
	closedServiceRequests, err := ctrl.ServiceRequest.GetClosedServiceRequests()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve closed service requests"})
		return
	}

	c.JSON(http.StatusOK, closedServiceRequests)
}

// DeleteServiceRequestHandler deletes a service request by its ID.
func (ctrl *ServiceRequestController) DeleteServiceRequestHandler(c *gin.Context) {
	id := c.Param("id")
	requestID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service request ID"})
		return
	}

	// Call the DeleteServiceRequest method from the ServiceRequestService.
	if err := ctrl.ServiceRequest.DeleteServiceRequest(uint(requestID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete service request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service request deleted successfully"})
}

// UpdateServiceRequestStatusHandler updates the status of a service request by its ID.
func (ctrl *ServiceRequestController) UpdateServiceRequestStatusHandler(c *gin.Context) {
	id := c.Param("id")
	requestID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service request ID"})
		return
	}

	var statusRequest struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&statusRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the UpdateServiceRequestStatus method from the ServiceRequestService.
	if err := ctrl.ServiceRequest.UpdateServiceRequestStatus(uint(requestID), statusRequest.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update service request status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service request status updated successfully"})
}

// GetServiceRequestCountByStatusHandler retrieves the count of service requests by their statuses.
func (ctrl *ServiceRequestController) GetServiceRequestCountByStatusHandler(c *gin.Context) {
	status := c.Param("status")
	// Call the GetServiceRequestCountByStatus method from the ServiceRequestService.
	counts, err := ctrl.ServiceRequest.GetServiceRequestCountByStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve service request counts by status"})
		return
	}

	for i, s := range counts {
		if i == status {

			c.JSON(http.StatusOK, s)
			return
		}
	}

	c.JSON(http.StatusOK, 0)
}

// GetServiceRequestCountByStatusHandler retrieves the count of service requests by their status.
func (ctrl *ServiceRequestController) GetServiceRequestCountByStatusHandler2(c *gin.Context) {
	// Call the GetServiceRequestCountByStatus method from the ServiceRequestService.
	counts, err := ctrl.ServiceRequest.GetServiceRequestCountByStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve service request counts by status"})
		return
	}

	c.JSON(http.StatusOK, counts)
}

// GetServiceRequestCountByCategoryHandler retrieves the count of service requests by category.
func (ctrl *ServiceRequestController) GetServiceRequestCountByCategoryHandler(c *gin.Context) {
	categoryID := c.Param("categoryID")
	categoryIDint, err := strconv.Atoi(categoryID)
	categoryIDUint := uint(categoryIDint)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	// Call the GetServiceRequestCountByCategory method from the ServiceRequestService.
	counts, err := ctrl.ServiceRequest.GetServiceRequestCountByCategory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve service request counts by category"})
		return
	}

	for i, s := range counts {
		if i == categoryIDUint {

			c.JSON(http.StatusOK, s)
			return
		}
	}

	c.JSON(http.StatusOK, 0)
}

// GetServiceRequestCountBySubCategoryHandler retrieves the count of service requests by sub-category.
func (ctrl *ServiceRequestController) GetServiceRequestCountBySubCategoryHandler(c *gin.Context) {
	subCategoryID := c.Param("subCategoryID")
	subCategoryIDint, err := strconv.Atoi(subCategoryID)
	subCategoryIDUint := uint(subCategoryIDint)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sub-category ID"})
		return
	}

	// Call the GetServiceRequestCountBySubCategory method from the ServiceRequestService.
	counts, err := ctrl.ServiceRequest.GetServiceRequestCountBySubCategory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve service request counts by sub-category"})
		return
	}

	for i, s := range counts {
		if i == subCategoryIDUint {

			c.JSON(http.StatusOK, s)
			return
		}
	}

	c.JSON(http.StatusOK, 0)
}

// GetUserClosedServiceRequestsHandler retrieves closed service requests for a specific user.
func (ctrl *ServiceRequestController) GetUserClosedServiceRequestsHandler(c *gin.Context) {
	userID := c.Param("userID")
	userIDUint, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Call the GetUserClosedServiceRequests method from the ServiceRequestService.
	serviceRequests, err := ctrl.ServiceRequest.GetUserClosedServiceRequests(uint(userIDUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve closed service requests for the user"})
		return
	}

	c.JSON(http.StatusOK, serviceRequests)
}

// GetUserOpenServiceRequestsHandler retrieves open service requests for a specific user.
func (ctrl *ServiceRequestController) GetUserOpenServiceRequestsHandler(c *gin.Context) {
	userID := c.Param("userID")
	userIDUint, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Call the GetUserOpenServiceRequests method from the ServiceRequestService.
	serviceRequests, err := ctrl.ServiceRequest.GetUserOpenServiceRequests(uint(userIDUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve open service requests for the user"})
		return
	}

	c.JSON(http.StatusOK, serviceRequests)
}

// GetServiceRequestsByUserAndStatusHandler retrieves service requests for a specific user with a given status.
func (ctrl *ServiceRequestController) GetServiceRequestsByUserAndStatusHandler(c *gin.Context) {
	userID := c.Param("userID")
	userIDUint, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	status := c.Param("status")

	// Call the GetServiceRequestsByUserAndStatus method from the ServiceRequestService.
	serviceRequests, err := ctrl.ServiceRequest.GetServiceRequestsByUserAndStatus(uint(userIDUint), status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve service requests for the user with the given status"})
		return
	}

	c.JSON(http.StatusOK, serviceRequests)
}

// GetUserServiceRequestCountHandler retrieves the count of service requests for a specific user.
func (ctrl *ServiceRequestController) GetUserServiceRequestCountHandler(c *gin.Context) {
	userID := c.Param("userID")
	userIDUint, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Call the GetUserServiceRequestCount method from the ServiceRequestService.
	count, err := ctrl.ServiceRequest.GetUserServiceRequestCount(uint(userIDUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve the count of service requests for the user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}

// GetUserServiceRequestCountByStatusHandler retrieves the count of service requests for a specific user based on status.
func (ctrl *ServiceRequestController) GetUserServiceRequestCountByStatusHandler(c *gin.Context) {
	userID := c.Param("userID")
	userIDUint, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Call the GetUserServiceRequestCountByStatus method from the ServiceRequestService.
	countByStatus, err := ctrl.ServiceRequest.GetUserServiceRequestCountByStatus(uint(userIDUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve the count of service requests by status for the user"})
		return
	}

	c.JSON(http.StatusOK, countByStatus)
}

// GetServiceRequestsByLocationHandler retrieves service requests by location.
func (ctrl *ServiceRequestController) GetServiceRequestsByLocationHandler(c *gin.Context) {
	locationID := c.Param("locationID")
	locationIDUint, err := strconv.ParseUint(locationID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid location ID"})
		return
	}

	// Call the GetServiceRequestsByLocation method from the ServiceRequestService.
	serviceRequests, err := ctrl.ServiceRequest.GetServiceRequestsByLocation(uint(locationIDUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve service requests by location"})
		return
	}

	c.JSON(http.StatusOK, serviceRequests)
}

// GetServiceRequestCountByLocationHandler retrieves the count of service requests for each location.
func (ctrl *ServiceRequestController) GetServiceRequestCountByLocationHandler(c *gin.Context) {
	// Call the GetServiceRequestCountByLocation method from the ServiceRequestService.
	countByLocation, err := ctrl.ServiceRequest.GetServiceRequestCountByLocation()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve the count of service requests by location"})
		return
	}

	c.JSON(http.StatusOK, countByLocation)
}

// GetServiceRequestsByUserAndLocationHandler retrieves service requests for a specific user and location.
func (ctrl *ServiceRequestController) GetServiceRequestsByUserAndLocationHandler(c *gin.Context) {
	userID := c.Param("userID")
	userIDUint, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	locationID := c.Param("locationID")
	locationIDUint, err := strconv.ParseUint(locationID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid location ID"})
		return
	}

	// Call the GetServiceRequestsByUserAndLocation method from the ServiceRequestService.
	serviceRequests, err := ctrl.ServiceRequest.GetServiceRequestsByUserAndLocation(uint(userIDUint), uint(locationIDUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve service requests for the user and location"})
		return
	}

	c.JSON(http.StatusOK, serviceRequests)
}

// GetServiceRequestCountByUserAndLocationHandler retrieves the count of service requests for each user and location.
func (ctrl *ServiceRequestController) GetServiceRequestCountByUserAndLocationHandler(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("user_id"))
	// Call the GetServiceRequestCountByUserAndLocation method from the ServiceRequestService.
	countByUserAndLocation, err := ctrl.ServiceRequest.GetServiceRequestCountByUserAndLocation(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve the count of service requests by user and location"})
		return
	}

	c.JSON(http.StatusOK, countByUserAndLocation)
}

// GetServiceRequestsByCategoryAndLocationHandler retrieves service requests for a specific category and location.
func (ctrl *ServiceRequestController) GetServiceRequestsByCategoryAndLocationHandler(c *gin.Context) {
	categoryID := c.Param("categoryID")
	categoryIDUint, err := strconv.ParseUint(categoryID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	locationID := c.Param("locationID")
	locationIDUint, err := strconv.ParseUint(locationID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid location ID"})
		return
	}

	// Call the GetServiceRequestsByCategoryAndLocation method from the ServiceRequestService.
	serviceRequests, err := ctrl.ServiceRequest.GetServiceRequestsByCategoryAndLocation(uint(categoryIDUint), uint(locationIDUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve service requests for the category and location"})
		return
	}

	c.JSON(http.StatusOK, serviceRequests)
}

// GetServiceRequestCountByCategoryAndLocationHandler retrieves the count of service requests for each category and location.
func (ctrl *ServiceRequestController) GetServiceRequestCountByCategoryAndLocationHandler(c *gin.Context) {
	categoryID, _ := strconv.Atoi(c.Param("category_id"))

	// Call the GetServiceRequestCountByCategoryAndLocation method from the ServiceRequestService.
	countByCategoryAndLocation, err := ctrl.ServiceRequest.GetServiceRequestCountByCategoryAndLocation(uint(categoryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve the count of service requests by category and location"})
		return
	}

	c.JSON(http.StatusOK, countByCategoryAndLocation)
}

// GetServiceRequestsByPriorityAndLocationHandler retrieves service requests for a specific priority and location.
func (ctrl *ServiceRequestController) GetServiceRequestsByPriorityAndLocationHandler(c *gin.Context) {
	priority := c.Param("priority")
	locationID := c.Param("locationID")
	locationIDUint, err := strconv.ParseUint(locationID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid location ID"})
		return
	}

	// Call the GetServiceRequestsByPriorityAndLocation method from the ServiceRequestService.
	serviceRequests, err := ctrl.ServiceRequest.GetServiceRequestsByPriorityAndLocation(priority, uint(locationIDUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve service requests for the priority and location"})
		return
	}

	c.JSON(http.StatusOK, serviceRequests)
}

// GetServiceRequestCountByPriorityAndLocationHandler retrieves the count of service requests for each priority and location.
func (ctrl *ServiceRequestController) GetServiceRequestCountByPriorityAndLocationHandler(c *gin.Context) {
	priority := c.Param("priority")
	locationID, _ := strconv.Atoi(c.Param("location"))

	// Call the GetServiceRequestCountByPriorityAndLocation method from the ServiceRequestService.
	countByPriorityAndLocation, err := ctrl.ServiceRequest.GetServiceRequestCountByPriorityAndLocation(priority, uint(locationID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve the count of service requests by priority and location"})
		return
	}

	c.JSON(http.StatusOK, countByPriorityAndLocation)
}

// GetServiceRequestCountByCategoryHandler retrieves the count of service requests by category.
func (ctrl *ServiceRequestController) GetServiceRequestCountByCategoryHandler2(c *gin.Context) {
	// Call the GetServiceRequestCountByCategory method from the ServiceRequestService.
	counts, err := ctrl.ServiceRequest.GetServiceRequestCountByCategory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve service request counts by category"})
		return
	}

	c.JSON(http.StatusOK, counts)
}

// GetServiceRequestCountBySubCategoryHandler retrieves the count of service requests by sub-category.
func (ctrl *ServiceRequestController) GetServiceRequestCountBySubCategoryHandler2(c *gin.Context) {
	// Call the GetServiceRequestCountBySubCategory method from the ServiceRequestService.
	counts, err := ctrl.ServiceRequest.GetServiceRequestCountBySubCategory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve service request counts by sub-category"})
		return
	}

	c.JSON(http.StatusOK, counts)
}

// GetClosedServiceRequestsHandler retrieves all closed service requests.
func (ctrl *ServiceRequestController) GetClosedServiceRequestsHandler(c *gin.Context) {
	// Call the GetClosedServiceRequests method from the ServiceRequestService.
	closedRequests, err := ctrl.ServiceRequest.GetClosedServiceRequests()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve closed service requests"})
		return
	}

	c.JSON(http.StatusOK, closedRequests)
}

// DeleteServiceRequestHandler deletes a service request by its ID.
func (ctrl *ServiceRequestController) DeleteServiceRequestHandler2(c *gin.Context) {
	id := c.Param("id")
	requestID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service request ID"})
		return
	}

	// Call the DeleteServiceRequest method from the ServiceRequestService.
	if err := ctrl.ServiceRequest.DeleteServiceRequest(uint(requestID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete service request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service request deleted successfully"})
}

// UpdateServiceRequestStatusHandler updates the status of a service request by its ID.
func (ctrl *ServiceRequestController) UpdateServiceRequestStatusHandler2(c *gin.Context) {
	id := c.Param("id")
	requestID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service request ID"})
		return
	}

	var statusUpdate struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the UpdateServiceRequestStatus method from the ServiceRequestService.
	if err := ctrl.ServiceRequest.UpdateServiceRequestStatus(uint(requestID), statusUpdate.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update service request status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service request status updated successfully"})
}

// GetServiceRequestsByPriorityAndLocationHandler retrieves service requests by priority and location.
func (ctrl *ServiceRequestController) GetServiceRequestsByPriorityAndLocationHandler2(c *gin.Context) {
	priority := c.Param("priority")
	locationID := c.Param("locationID")
	locationIDUint, err := strconv.ParseUint(locationID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid location ID"})
		return
	}

	// Call the GetServiceRequestsByPriorityAndLocation method from the ServiceRequestService.
	serviceRequests, err := ctrl.ServiceRequest.GetServiceRequestsByPriorityAndLocation(priority, uint(locationIDUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve service requests by priority and location"})
		return
	}

	c.JSON(http.StatusOK, serviceRequests)
}

// GetServiceRequestCountByPriorityAndLocationHandler retrieves the count of service requests by priority and location.
func (ctrl *ServiceRequestController) GetServiceRequestCountByPriorityAndLocationHandler2(c *gin.Context) {
	locationID := c.Param("locationID")
	priorityID := c.Param("priorityID")
	locationIDUint, err := strconv.ParseUint(locationID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid location ID"})
		return
	}

	// Call the GetServiceRequestCountByPriorityAndLocation method from the ServiceRequestService.
	counts, err := ctrl.ServiceRequest.GetServiceRequestCountByPriorityAndLocation(priorityID, uint(locationIDUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve service request counts by priority and location"})
		return
	}

	c.JSON(http.StatusOK, counts)
}

// GetServiceRequestsByCategoryAndSubCategoryHandler retrieves service requests for a specific category and sub-category.
func (ctrl *ServiceRequestController) GetServiceRequestsByCategoryAndSubCategoryHandler2(c *gin.Context) {
	categoryID := c.Param("categoryID")
	categoryIDUint, err := strconv.ParseUint(categoryID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	subCategoryID := c.Param("subCategoryID")
	subCategoryIDUint, err := strconv.ParseUint(subCategoryID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sub-category ID"})
		return
	}

	// Call the GetServiceRequestsByCategoryAndSubCategory method from the ServiceRequestService.
	serviceRequests, err := ctrl.ServiceRequest.GetServiceRequestsByCategoryAndSubCategory(uint(categoryIDUint), uint(subCategoryIDUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve service requests for the category and sub-category"})
		return
	}

	c.JSON(http.StatusOK, serviceRequests)
}

// GetServiceRequestsByStatusAndLocationHandler retrieves service requests for a specific status and location.
func (ctrl *ServiceRequestController) GetServiceRequestsByStatusAndLocationHandler2(c *gin.Context) {
	status := c.Param("status")
	locationID := c.Param("locationID")
	locationIDUint, err := strconv.ParseUint(locationID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid location ID"})
		return
	}

	// Call the GetServiceRequestsByStatusAndLocation method from the ServiceRequestService.
	serviceRequests, err := ctrl.ServiceRequest.GetServiceRequestsByStatusAndLocation(status, uint(locationIDUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve service requests for the status and location"})
		return
	}

	c.JSON(http.StatusOK, serviceRequests)
}
