package services

import (
	"github.com/shuttlersit/service-desk/backend/models"
)

type IncidentService interface {
	CreateIncident(incident *models.Incident) error
	UpdateIncident(incident *models.Incident) error
	GetIncidentByID(incidentID uint) (*models.Incident, error)
	GetIncidentsBySeverity(severity string) ([]*models.Incident, error)
	AssignIncidentToTeam(incidentID uint, teamID uint) error
}

type DefaultIncidentService struct {
	// Implement methods defined in IncidentService interface
}

// Implement the CreateIncident, UpdateIncident, GetIncidentByID,
// GetIncidentsBySeverity, and AssignIncidentToTeam methods.
