// backend/controllers/incident_controllers.go

package controllers

import (
	"net/http"
	"strconv"
	_ "time"

	"github.com/shuttlersit/service-desk/backend/models"
	"github.com/shuttlersit/service-desk/backend/services"

	"github.com/gin-gonic/gin"
)

type IncidentController struct {
	IncidentService *services.DefaultIncidentService
}

func NewIncidentController(service *services.DefaultIncidentService) *IncidentController {
	return &IncidentController{
		IncidentService: service,
	}
}

// GetAllIncidentsHandler retrieves all incidents from database.
func (ctrl *IncidentController) GetAllIncidentsHandler(c *gin.Context) {
	incidents, err := ctrl.IncidentService.GetAllIncidents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve incidents"})
		return
	}

	c.JSON(http.StatusOK, incidents)
}

// CreateIncidentHandler creates a new incident.
func (ctrl *IncidentController) CreateIncidentHandler(c *gin.Context) {
	var incident models.Incident
	if err := c.ShouldBindJSON(&incident); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.IncidentService.CreateIncident(&incident); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, incident)
}

// UpdateIncidentHandler updates an existing incident.
func (ctrl *IncidentController) UpdateIncidentHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	incidentID := uint(id)
	var updatedIncident models.Incident

	if err := c.ShouldBindJSON(&updatedIncident); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch the incident by ID to ensure it exists.
	existingIncident, err := ctrl.IncidentService.GetIncidentByID(incidentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Incident not found"})
		return
	}

	// Update the incident details.
	existingIncident.Title = updatedIncident.Title
	existingIncident.Description = updatedIncident.Description
	existingIncident.Category = updatedIncident.Category
	existingIncident.Priority = updatedIncident.Priority
	existingIncident.Tags = updatedIncident.Tags
	existingIncident.AttachmentURL = updatedIncident.AttachmentURL
	existingIncident.HasAttachments = updatedIncident.HasAttachments
	existingIncident.Severity = updatedIncident.Severity

	if err := ctrl.IncidentService.UpdateIncident(existingIncident); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, existingIncident)
}

// GetIncidentByIDHandler retrieves an incident by its ID.
func (ctrl *IncidentController) GetIncidentByIDHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	incidentID := uint(id)

	incident, err := ctrl.IncidentService.GetIncidentByID(incidentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Incident not found"})
		return
	}

	c.JSON(http.StatusOK, incident)
}

// AssignIncidentToTeamHandler assigns an incident to a team.
func (ctrl *IncidentController) AssignIncidentToTeamHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	incidentID := uint(id)
	tID, _ := strconv.Atoi(c.Param("team_id"))
	teamID := uint(tID)

	if err := ctrl.IncidentService.AssignIncidentToTeam(incidentID, teamID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Incident assigned to team successfully"})
}

// ResolveIncidentHandler resolves an incident.
func (ctrl *IncidentController) ResolveIncidentHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	incidentID := uint(id)

	if err := ctrl.IncidentService.ResolveIncident(incidentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Incident resolved successfully"})
}

// GetIncidentComments retrieves comments to an incident
func (ctrl *IncidentController) GetIncidentCommentsHandler(c *gin.Context) {
	id := c.Param("incident_id")
	incidentID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid incident ID"})
		return
	}

	// Call the GetServiceRequestComments method from the ServiceRequestService.
	comments, err := ctrl.IncidentService.GetIncidentComments(uint(incidentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve service request comments"})
		return
	}

	c.JSON(http.StatusOK, comments)
}

// AddIncidentCommentHandler adds a comment to an incident.
func (ctrl *IncidentController) AddIncidentCommentHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	incidentID := uint(id)
	var comment models.IncidentComment

	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.IncidentService.AddIncidentComment(incidentID, comment.Comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Comment added successfully"})
}

// GetIncidentsBySeverityHandler retrieves incidents by severity.
func (ctrl *IncidentController) GetIncidentsBySeverityHandler(c *gin.Context) {
	severity := c.Param("severity")

	incidents, err := ctrl.IncidentService.GetIncidentsBySeverity(severity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, incidents)
}

// GetIncidentHistoryHandler retrieves the history of an incident.
func (ctrl *IncidentController) GetIncidentHistoryHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	incidentID := uint(id)
	history, err := ctrl.IncidentService.GetIncidentHistory(incidentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}

// CloseIncidentHandler closes an incident.
func (ctrl *IncidentController) CloseIncidentHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	incidentID := uint(id)

	if err := ctrl.IncidentService.CloseIncident(incidentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Incident closed successfully"})
}

// ReopenIncidentHandler reopens a closed incident.
func (ctrl *IncidentController) ReopenIncidentHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	incidentID := uint(id)
	if err := ctrl.IncidentService.ReopenIncident(incidentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Incident reopened successfully"})
}

// AssignIncidentToUserHandler assigns an incident to a user.
func (ctrl *IncidentController) AssignIncidentToUserHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	incidentID := uint(id)
	uID, _ := strconv.Atoi(c.Param("user_id"))
	userID := uint(uID)

	if err := ctrl.IncidentService.AssignIncidentToUser(incidentID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Incident assigned to user successfully"})
}

// GetIncidentsByUserHandler retrieves incidents by user ID.
func (ctrl *IncidentController) GetIncidentsByUserHandler(c *gin.Context) {
	uID, _ := strconv.Atoi(c.Param("user_id"))
	userID := uint(uID)

	incidents, err := ctrl.IncidentService.GetIncidentsByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, incidents)
}

// GetOpenIncidentsHandler retrieves all open incidents.
func (ctrl *IncidentController) GetOpenIncidentsHandler(c *gin.Context) {
	incidents, err := ctrl.IncidentService.GetOpenIncidents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, incidents)
}

// GetAssignedIncidentsByUserHandler retrieves incidents assigned to a user.
func (ctrl *IncidentController) GetAssignedIncidentsByUserHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("user_id"))
	userID := uint(id)

	incidents, err := ctrl.IncidentService.GetAssignedIncidentsByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, incidents)
}

// GetIncidentsByUserAndStatusHandler retrieves incidents by user ID and status.
func (ctrl *IncidentController) GetIncidentsByUserAndStatusHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("user_id"))
	userID := uint(id)
	status := c.Param("status")

	incidents, err := ctrl.IncidentService.GetIncidentsByUserAndStatus(userID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, incidents)
}

// GetIncidentAssigneeHandler retrieves the user assigned to handle an incident.
func (ctrl *IncidentController) GetIncidentAssigneeHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	incidentID := uint(id)

	user, err := ctrl.IncidentService.GetIncidentAssignee(incidentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetIncidentsByTeamAndStatusHandler retrieves incidents by team ID and status.
func (ctrl *IncidentController) GetIncidentsByTeamAndStatusHandler(c *gin.Context) {
	tid, _ := strconv.Atoi(c.Param("team_id"))
	teamID := uint(tid)
	status := c.Param("status")

	incidents, err := ctrl.IncidentService.GetIncidentsByTeamAndStatus(teamID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, incidents)
}

// GetIncidentBySubjectHandler retrieves an incident by its subject.
func (ctrl *IncidentController) GetIncidentBySubjectHandler(c *gin.Context) {
	subject := c.Param("subject")

	incident, err := ctrl.IncidentService.GetIncidentBySubject(subject)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, incident)
}

// DeleteIncidentHandler deletes an incident by its ID.
func (ctrl *IncidentController) DeleteIncidentHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	incidentID := uint(id)

	if err := ctrl.IncidentService.DeleteIncident(incidentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Incident deleted successfully"})
}

// NewIncidentHistoryEntryHandler creates a new history entry for an incident.
func (ctrl *IncidentController) NewIncidentHistoryEntryHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	incidentID := uint(id)
	status := c.Param("status")

	if err := ctrl.IncidentService.NewIncidentHistoryEntry(incidentID, status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "History entry created successfully"})
}

// UpdateIncidentStatusHandler updates the status of an incident.
func (ctrl *IncidentController) UpdateIncidentStatusHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	incidentID := uint(id)
	status := c.Param("status")

	if err := ctrl.IncidentService.UpdateIncidentStatus(incidentID, status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Incident status updated successfully"})
}

// GetIncidentStatsHandler retrieves incident statistics.
func (ctrl *IncidentController) GetIncidentStatsHandler(c *gin.Context) {
	stats, err := ctrl.IncidentService.GetIncidentStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetIncidentsByTeamHandler retrieves incidents assigned to a specific team.
func (ctrl *IncidentController) GetIncidentsByTeamHandler(c *gin.Context) {
	tid, _ := strconv.Atoi(c.Param("team_id"))
	teamID := uint(tid)

	incidents, err := ctrl.IncidentService.GetIncidentsByTeam(teamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, incidents)
}

// GetIncidentsByStatusHandler retrieves incidents by status.
func (ctrl *IncidentController) GetIncidentsByStatusHandler(c *gin.Context) {
	status := c.Param("status")

	incidents, err := ctrl.IncidentService.GetIncidentsByStatus(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, incidents)
}

// GetIncidentsWithAttachmentsHandler retrieves incidents with attachments.
func (ctrl *IncidentController) GetIncidentsWithAttachmentsHandler(c *gin.Context) {
	incidents, err := ctrl.IncidentService.GetIncidentsWithAttachments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, incidents)
}

// UnassignIncidentHandler unassigns an incident from a user.
func (ctrl *IncidentController) UnassignIncidentHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	incidentID := uint(id)

	if err := ctrl.IncidentService.UnassignIncident(incidentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Incident unassigned successfully"})
}

// GetUnassignedIncidentsHandler retrieves unassigned incidents.
func (ctrl *IncidentController) GetUnassignedIncidentsHandler(c *gin.Context) {
	incidents, err := ctrl.IncidentService.GetUnassignedIncidents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, incidents)
}

// GetAssignedIncidentsHandler retrieves incidents assigned to a specific user.
func (ctrl *IncidentController) GetAssignedIncidentsHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("user_id"))
	userID := uint(id)

	incidents, err := ctrl.IncidentService.GetAssignedIncidents(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, incidents)
}

// GetIncidentsByCategoryHandler retrieves incidents by category.
func (ctrl *IncidentController) GetIncidentsByCategoryHandler(c *gin.Context) {
	category := c.Param("category")

	incidents, err := ctrl.IncidentService.GetIncidentsByCategory(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, incidents)
}

// GetIncidentsByPriorityHandler retrieves incidents by priority.
func (ctrl *IncidentController) GetIncidentsByPriorityHandler(c *gin.Context) {
	priority := c.Param("priority")

	incidents, err := ctrl.IncidentService.GetIncidentsByPriority(priority)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, incidents)
}

// GetIncidentsByTagHandler retrieves incidents by tag.
func (ctrl *IncidentController) GetIncidentsByTagHandler(c *gin.Context) {
	tag := c.Param("tag")

	incidents, err := ctrl.IncidentService.GetIncidentsByTag(tag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, incidents)
}
