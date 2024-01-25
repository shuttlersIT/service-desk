// backend/models/incident_db_model.go

package models

import (
	"time"

	"gorm.io/gorm"
)

// Incident represents an incident report.
type Incident struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         uint      `json:"user_id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Category       string    `json:"category"`
	Priority       string    `json:"priority"`
	Tags           []string  `gorm:"type:text[]" json:"tags"`
	AttachmentURL  string    `json:"attachment_url"`
	HasAttachments bool      `json:"has_attachments"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Severity       string    `json:"severity"`
}

// IncidentHistoryEntry represents a historical entry related to an incident.
type IncidentHistoryEntry struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	IncidentID  uint      `json:"incident_id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// IncidentComment represents a comment made on an incident.
type IncidentComment struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	IncidentID uint      `json:"incident_id"`
	Comment    string    `json:"comment"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// IncidentStorage defines the methods for managing incidents.
type IncidentStorage interface {
	CreateIncident(incident *Incident) error
	UpdateIncident(incident *Incident) error
	GetIncidentByID(incidentID uint) (*Incident, error)
	GetUserIncidents(userID uint) ([]*Incident, error)
	CloseIncident(incidentID uint) error
	GetIncidentHistory(incidentID uint) ([]*IncidentHistoryEntry, error)
	AddCommentToIncident(incidentID uint, comment string) error
	GetIncidentComments(incidentID uint) ([]*IncidentComment, error)
	NewIncidentHistoryEntry(incidentID uint, status string) error
	UpdateIncidentStatus(incidentID uint, status string) error
	GetOpenIncidents() ([]*Incident, error)                  // New method to retrieve all open incidents
	GetIncidentStats() (map[string]int, error)               // New method to retrieve incident statistics
	GetIncidentsByStatus(status string) ([]*Incident, error) // New method to retrieve incidents by status
	DeleteIncident(incidentID uint) error
	GetIncidentsByUserAndStatus(userID uint, status string) ([]*Incident, error)
	GetIncidentsByCategory(category string) ([]*Incident, error)
	GetIncidentsByPriority(priority string) ([]*Incident, error)
	GetIncidentsByTag(tag string) ([]*Incident, error)
	AssignIncidentToUser(incidentID uint, userID uint) error
	UnassignIncident(incidentID uint) error
	GetUnassignedIncidents() ([]*Incident, error)
	GetAssignedIncidents(userID uint) ([]*Incident, error)
	GetIncidentsWithAttachments() ([]*Incident, error)
	// Alternative Methods
	GetAssignedIncidents2(userID uint) ([]*Incident, error)
	GetIncidentsByTeamAndStatus2(teamID uint, status string) ([]*Incident, error)
	GetIncidentBySubject2(subject string) (*Incident, error)
	DeleteIncident2(incidentID uint) error
}

// IncidentDBModel handles database operations for incidents.
type IncidentDBModel struct {
	DB *gorm.DB
}

// NewIncidentDBModel creates a new instance of IncidentDBModel.
func NewIncidentDBModel(db *gorm.DB) *IncidentDBModel {
	return &IncidentDBModel{
		DB: db,
	}
}

// CreateIncident creates a new incident report.
func (idm *IncidentDBModel) CreateIncident(incident *Incident) error {
	return idm.DB.Create(incident).Error
}

// UpdateIncident updates an existing incident report.
func (idm *IncidentDBModel) UpdateIncident(incident *Incident) error {
	if err := idm.DB.Save(incident).Error; err != nil {
		return err
	}
	return nil
}

// GetIncidentByID retrieves an incident report by its ID.
func (idm *IncidentDBModel) GetIncidentByID(incidentID uint) (*Incident, error) {
	var incident Incident
	err := idm.DB.Where("id = ?", incidentID).First(&incident).Error
	return &incident, err
}

// GetUserIncidents retrieves all incidents reported by a user.
func (idm *IncidentDBModel) GetUserIncidents(userID uint) ([]*Incident, error) {
	var incidents []*Incident
	err := idm.DB.Where("user_id = ?", userID).Find(&incidents).Error
	return incidents, err
}

// CloseIncident closes an incident report by updating its status.
func (idm *IncidentDBModel) CloseIncident(incidentID uint) error {
	// Assuming that "Closed" is a valid status value
	return idm.DB.Model(&Incident{}).Where("id = ?", incidentID).Update("Status", "Closed").Error
}

// AddCommentToIncident adds a comment to an incident report.
func (idm *IncidentDBModel) AddCommentToIncident(incidentID uint, comment string) error {
	incidentComment := &IncidentComment{
		IncidentID: incidentID,
		Comment:    comment,
		CreatedAt:  time.Now(),
	}
	return idm.DB.Create(incidentComment).Error
}

// NewIncidentHistoryEntry creates a new history entry for an incident.
func (idm *IncidentDBModel) NewIncidentHistoryEntry(incidentID uint, status string) error {
	entry := &IncidentHistoryEntry{
		IncidentID: incidentID,
		Status:     status,
		UpdatedAt:  time.Now(),
	}
	return idm.DB.Create(entry).Error
}

// UpdateIncidentStatus updates the status of an incident.
func (idm *IncidentDBModel) UpdateIncidentStatus(incidentID uint, status string) error {
	return idm.DB.Model(&Incident{}).Where("id = ?", incidentID).Update("Status", status).Error
}

// GetOpenIncidents retrieves all open incidents.
func (idm *IncidentDBModel) GetOpenIncidents() ([]*Incident, error) {
	var openIncidents []*Incident
	err := idm.DB.Where("status = ?", "Open").Find(&openIncidents).Error
	return openIncidents, err
}

// GetIncidentStats retrieves incident statistics.
func (idm *IncidentDBModel) GetIncidentStats() (map[string]int, error) {
	// Count incidents per status
	var stats []struct {
		Status string
		Count  int
	}

	if err := idm.DB.Table("incidents").Select("status, COUNT(*) as count").Group("status").Scan(&stats).Error; err != nil {
		return nil, err
	}

	// Convert the result into a map
	incidentStats := make(map[string]int)
	for _, stat := range stats {
		incidentStats[stat.Status] = stat.Count
	}

	return incidentStats, nil
}

// GetIncidentsByStatus retrieves incidents with a specific status.
func (idm *IncidentDBModel) GetIncidentsByStatus(status string) ([]*Incident, error) {
	var incidents []*Incident
	err := idm.DB.Where("status = ?", status).Find(&incidents).Error
	return incidents, err
}

// DeleteIncident deletes an incident from the database.
func (idm *IncidentDBModel) DeleteIncident(incidentID uint) error {
	if err := idm.DB.Delete(&Incident{}, incidentID).Error; err != nil {
		return err
	}
	return nil
}

// GetIncidentsByUserAndStatus retrieves incidents for a specific user and status.
func (idm *IncidentDBModel) GetIncidentsByUserAndStatus(userID uint, status string) ([]*Incident, error) {
	var incidents []*Incident
	err := idm.DB.Where("user_id = ? AND status = ?", userID, status).Find(&incidents).Error
	return incidents, err
}

// GetIncidentsByCategory retrieves incidents with a specific category.
func (idm *IncidentDBModel) GetIncidentsByCategory(category string) ([]*Incident, error) {
	var incidents []*Incident
	err := idm.DB.Where("category = ?", category).Find(&incidents).Error
	return incidents, err
}

// GetIncidentsByPriority retrieves incidents with a specific priority level.
func (idm *IncidentDBModel) GetIncidentsByPriority(priority string) ([]*Incident, error) {
	var incidents []*Incident
	err := idm.DB.Where("priority = ?", priority).Find(&incidents).Error
	return incidents, err
}

// GetIncidentsByTag retrieves incidents with a specific tag.
func (im *IncidentDBModel) GetIncidentsByTag(tag string) ([]*Incident, error) {
	var incidents []*Incident
	err := im.DB.Where("tags LIKE ?", "%"+tag+"%").Find(&incidents).Error
	if err != nil {
		return nil, err
	}
	return incidents, nil
}

// AssignIncidentToUser assigns an incident to a specific user.
func (im *IncidentDBModel) AssignIncidentToUser(incidentID uint, userID uint) error {
	incident, err := im.GetIncidentByID(incidentID)
	if err != nil {
		return err
	}
	incident.UserID = userID
	if err := im.UpdateIncident(incident); err != nil {
		return err
	}
	return nil
}

// UnassignIncident unassigns an incident from a user.
func (im *IncidentDBModel) UnassignIncident(incidentID uint) error {
	incident, err := im.GetIncidentByID(incidentID)
	if err != nil {
		return err
	}
	incident.UserID = 0
	if err := im.UpdateIncident(incident); err != nil {
		return err
	}
	return nil
}

// GetUnassignedIncidents retrieves all unassigned incidents.
func (im *IncidentDBModel) GetUnassignedIncidents() ([]*Incident, error) {
	var incidents []*Incident
	err := im.DB.Where("user_id IS NULL").Find(&incidents).Error
	if err != nil {
		return nil, err
	}
	return incidents, nil
}

// GetAssignedIncidents retrieves all incidents assigned to a specific user.
func (im *IncidentDBModel) GetAssignedIncidents(userID uint) ([]*Incident, error) {
	var incidents []*Incident
	err := im.DB.Where("user_id = ?", userID).Find(&incidents).Error
	if err != nil {
		return nil, err
	}
	return incidents, nil
}

// GetIncidentsWithAttachments retrieves all incidents that have attachments.
func (im *IncidentDBModel) GetIncidentsWithAttachments() ([]*Incident, error) {
	var incidents []*Incident
	err := im.DB.Where("has_attachments = true").Find(&incidents).Error
	if err != nil {
		return nil, err
	}
	return incidents, nil
}

// UpdateIncidentCategory updates the category of an incident.
func (im *IncidentDBModel) UpdateIncidentCategory(incidentID uint, category string) error {
	incident, err := im.GetIncidentByID(incidentID)
	if err != nil {
		return err
	}
	incident.Category = category
	if err := im.UpdateIncident(incident); err != nil {
		return err
	}
	return nil
}

// UpdateIncidentPriority updates the priority of an incident.
func (im *IncidentDBModel) UpdateIncidentPriority(incidentID uint, priority string) error {
	incident, err := im.GetIncidentByID(incidentID)
	if err != nil {
		return err
	}
	incident.Priority = priority
	if err := im.UpdateIncident(incident); err != nil {
		return err
	}
	return nil
}

// UpdateIncidentTags updates the tags of an incident.
func (im *IncidentDBModel) UpdateIncidentTags(incidentID uint, tags []string) error {
	incident, err := im.GetIncidentByID(incidentID)
	if err != nil {
		return err
	}
	incident.Tags = tags
	if err := im.UpdateIncident(incident); err != nil {
		return err
	}
	return nil
}

// GetIncidentsByUserAndCategory retrieves incidents with a specific category assigned to a user.
func (im *IncidentDBModel) GetIncidentsByUserAndCategory(userID uint, category string) ([]*Incident, error) {
	var incidents []*Incident
	err := im.DB.Where("user_id = ? AND category = ?", userID, category).Find(&incidents).Error
	if err != nil {
		return nil, err
	}
	return incidents, nil
}

// GetIncidentsByUserAndPriority retrieves incidents with a specific priority assigned to a user.
func (im *IncidentDBModel) GetIncidentsByUserAndPriority(userID uint, priority string) ([]*Incident, error) {
	var incidents []*Incident
	err := im.DB.Where("user_id = ? AND priority = ?", userID, priority).Find(&incidents).Error
	if err != nil {
		return nil, err
	}
	return incidents, nil
}

// GetIncidentsByUserAndTag retrieves incidents with a specific tag assigned to a user.
func (im *IncidentDBModel) GetIncidentsByUserAndTag(userID uint, tag string) ([]*Incident, error) {
	var incidents []*Incident
	err := im.DB.Where("user_id = ? AND tags LIKE ?", userID, "%"+tag+"%").Find(&incidents).Error
	if err != nil {
		return nil, err
	}
	return incidents, nil
}

// AddAttachmentToIncident adds an attachment to an incident.
func (im *IncidentDBModel) AddAttachmentToIncident(incidentID uint, attachmentURL string) error {
	incident, err := im.GetIncidentByID(incidentID)
	if err != nil {
		return err
	}
	incident.AttachmentURL = attachmentURL
	incident.HasAttachments = true
	if err := im.UpdateIncident(incident); err != nil {
		return err
	}
	return nil
}

// RemoveAttachmentFromIncident removes an attachment from an incident.
func (im *IncidentDBModel) RemoveAttachmentFromIncident(incidentID uint) error {
	incident, err := im.GetIncidentByID(incidentID)
	if err != nil {
		return err
	}
	incident.AttachmentURL = ""
	incident.HasAttachments = false
	if err := im.UpdateIncident(incident); err != nil {
		return err
	}
	return nil
}

// GetIncidentsWithHighPriority retrieves incidents with high priority.
func (im *IncidentDBModel) GetIncidentsWithHighPriority() ([]*Incident, error) {
	var incidents []*Incident
	err := im.DB.Where("priority = ?", "High").Find(&incidents).Error
	if err != nil {
		return nil, err
	}
	return incidents, nil
}

// GetIncidentsWithMediumPriority retrieves incidents with medium priority.
func (im *IncidentDBModel) GetIncidentsWithMediumPriority() ([]*Incident, error) {
	var incidents []*Incident
	err := im.DB.Where("priority = ?", "Medium").Find(&incidents).Error
	if err != nil {
		return nil, err
	}
	return incidents, nil
}

// GetIncidentsWithLowPriority retrieves incidents with low priority.
func (im *IncidentDBModel) GetIncidentsWithLowPriority() ([]*Incident, error) {
	var incidents []*Incident
	err := im.DB.Where("priority = ?", "Low").Find(&incidents).Error
	if err != nil {
		return nil, err
	}
	return incidents, nil
}

// GetIncidentsWithCriticalPriority retrieves incidents with critical priority.
func (im *IncidentDBModel) GetIncidentsWithCriticalPriority() ([]*Incident, error) {
	var incidents []*Incident
	err := im.DB.Where("priority = ?", "Critical").Find(&incidents).Error
	if err != nil {
		return nil, err
	}
	return incidents, nil
}

// CreateIncidentHistoryEntry creates a new history entry for an incident.
func (im *IncidentDBModel) CreateIncidentHistoryEntry(incidentID uint, description string) error {
	entry := &IncidentHistoryEntry{
		IncidentID:  incidentID,
		Description: description,
		CreatedAt:   time.Now(),
	}
	if err := im.DB.Create(entry).Error; err != nil {
		return err
	}
	return nil
}

// GetIncidentsWithoutAttachments retrieves incidents that do not have attachments.
func (im *IncidentDBModel) GetIncidentsWithoutAttachments() ([]*Incident, error) {
	var incidents []*Incident
	err := im.DB.Where("has_attachments = ?", false).Find(&incidents).Error
	if err != nil {
		return nil, err
	}
	return incidents, nil
}

func (s *IncidentDBModel) ResolveIncident(incidentID uint) error {
	// Implement logic to mark an incident as resolved in the database.
	if err := s.DB.Model(&Incident{}).Where("id = ?", incidentID).Update("status", "Resolved").Error; err != nil {
		return err
	}
	return nil
}

func (s *IncidentDBModel) AddIncidentComment(incidentID uint, comment string) error {
	// Implement logic to add a comment to an incident in the database.
	commentEntry := &IncidentComment{
		IncidentID: incidentID,
		Comment:    comment,
	}
	if err := s.DB.Create(commentEntry).Error; err != nil {
		return err
	}
	return nil
}

func (s *IncidentDBModel) GetIncidentComments(incidentID uint) ([]*IncidentComment, error) {
	// Implement logic to retrieve comments for an incident from the database.
	var comments []*IncidentComment
	err := s.DB.Where("incident_id = ?", incidentID).Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (s *IncidentDBModel) GetIncidentHistory(incidentID uint) ([]*IncidentHistoryEntry, error) {
	// Implement logic to retrieve the history of an incident from the database.
	var history []*IncidentHistoryEntry
	err := s.DB.Where("incident_id = ?", incidentID).Find(&history).Error
	if err != nil {
		return nil, err
	}
	return history, nil
}

func (s *IncidentDBModel) GetIncidentsBySeverity(severity string) ([]*Incident, error) {
	// Implement logic to retrieve incident reports by severity from the database.
	var incidents []*Incident
	err := s.DB.Where("severity = ?", severity).Find(&incidents).Error
	if err != nil {
		return nil, err
	}
	return incidents, nil
}

func (s *IncidentDBModel) AssignIncidentToTeam(incidentID uint, teamID uint) error {
	// Implement logic to assign an incident report to a specific team in the database.
	// You can update the "TeamID" field of the incident report.
	if err := s.DB.Model(&Incident{}).Where("id = ?", incidentID).Update("team_id", teamID).Error; err != nil {
		return err
	}
	return nil
}

func (s *IncidentDBModel) ReopenIncident(incidentID uint) error {
	// Implement logic to reopen a closed incident in the database.
	if err := s.DB.Model(&Incident{}).Where("id = ?", incidentID).Update("status", "Open").Error; err != nil {
		return err
	}
	return nil
}

func (s *IncidentDBModel) GetIncidentsByUser(userID uint) ([]*Incident, error) {
	// Implement logic to retrieve incident reports by user ID from the database.
	var incidents []*Incident
	err := s.DB.Where("user_id = ?", userID).Find(&incidents).Error
	if err != nil {
		return nil, err
	}
	return incidents, nil
}

func (s *IncidentDBModel) GetAssignedIncidentsByUser(userID uint) ([]*Incident, error) {
	// Implement logic to retrieve incident reports assigned to a specific user from the database.
	var incidents []*Incident
	err := s.DB.Where("assigned_user_id = ?", userID).Find(&incidents).Error
	if err != nil {
		return nil, err
	}
	return incidents, nil
}

func (s *IncidentDBModel) AddIncidentHistoryEntry(incidentID uint, entry *IncidentHistoryEntry) error {
	// Implement logic to add a history entry to an incident in the database.
	entry.IncidentID = incidentID
	if err := s.DB.Create(entry).Error; err != nil {
		return err
	}
	return nil
}

func (s *IncidentDBModel) GetIncidentAssignedTeam(incidentID uint) (*Teams, error) {
	// Implement logic to retrieve the team assigned to handle an incident from the database.
	var team Teams
	err := s.DB.Model(&Incident{}).Where("id = ?", incidentID).Association("Team").Find(&team)
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (s *IncidentDBModel) GetIncidentsByTeam(teamID uint) ([]*Incident, error) {
	// Implement logic to retrieve incident reports assigned to a specific team from the database.
	var incidents []*Incident
	err := s.DB.Where("team_id = ?", teamID).Find(&incidents).Error
	if err != nil {
		return nil, err
	}
	return incidents, nil
}

func (s *IncidentDBModel) GetIncidentAssignee(incidentID uint) (*Users, error) {
	// Implement logic to retrieve the user assigned to handle a specific incident from the database.
	var user Users
	err := s.DB.Model(&Incident{}).Where("id = ?", incidentID).Association("AssignedUser").Find(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *IncidentDBModel) GetIncidentsByTeamAndStatus(teamID uint, status string) ([]*Incident, error) {
	// Implement logic to retrieve incident reports assigned to a specific team with a given status from the database.
	var incidents []*Incident
	err := s.DB.Where("team_id = ? AND status = ?", teamID, status).Find(&incidents).Error
	if err != nil {
		return nil, err
	}
	return incidents, nil
}

func (s *IncidentDBModel) GetIncidentBySubject(subject string) (*Incident, error) {
	// Implement logic to retrieve an incident report by its subject from the database.
	var incident Incident
	err := s.DB.Where("subject = ?", subject).First(&incident).Error
	if err != nil {
		return nil, err
	}
	return &incident, nil
}

// GetAssignedIncidents retrieves incident reports assigned to a specific user.
func (idm *IncidentDBModel) GetAssignedIncidents2(userID uint) ([]*Incident, error) {
	var incidents []*Incident
	err := idm.DB.Where("user_id = ?", userID).Find(&incidents).Error
	return incidents, err
}

// GetIncidentsByTeamAndStatus retrieves incident reports assigned to a specific team with a given status.
func (idm *IncidentDBModel) GetIncidentsByTeamAndStatus2(teamID uint, status string) ([]*Incident, error) {
	var incidents []*Incident
	err := idm.DB.Where("team_id = ? AND status = ?", teamID, status).Find(&incidents).Error
	return incidents, err
}

// GetIncidentBySubject retrieves an incident report by its subject.
func (idm *IncidentDBModel) GetIncidentBySubject2(subject string) (*Incident, error) {
	var incident Incident
	err := idm.DB.Where("subject = ?", subject).First(&incident).Error
	return &incident, err
}

// DeleteIncident deletes an incident report from the database.
func (idm *IncidentDBModel) DeleteIncident2(incidentID uint) error {
	err := idm.DB.Delete(&Incident{}, incidentID).Error
	return err
}
