package controllers

import (
	"github.com/shuttlersit/service-desk/backend/models"
	"github.com/shuttlersit/service-desk/backend/services"

	"github.com/gin-gonic/gin"
)

type IncidentController struct {
	Service services.IncidentService
}

func (controller *IncidentController) CreateIncident(c *gin.Context) {
	// Parse request and create a new incident report
	var incident models.Incident
	if err := c.ShouldBindJSON(&incident); err != nil {
		// Handle validation errors
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := controller.Service.CreateIncident(&incident); err != nil {
		// Handle database errors
		c.JSON(500, gin.H{"error": "Failed to create incident report"})
		return
	}

	c.JSON(200, incident)
}

// Implement other methods like UpdateIncident, GetIncidentByID, GetIncidentsBySeverity, and AssignIncidentToTeam.
