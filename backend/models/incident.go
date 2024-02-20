// backend/models/incident_db_model.go

package models

import (
	"time"

	"gorm.io/gorm"
)

// Incident represents an incident report.
type Incident struct {
	gorm.Model                 // This includes ID, CreatedAt, UpdatedAt, and DeletedAt
	AssignedTo     *uint       `gorm:"index" json:"assigned_to,omitempty"`
	UserID         uint        `json:"user_id" gorm:"not null;index"`
	ReportedBy     uint        `gorm:"index" json:"reported_by"`
	Reporter       Users       `gorm:"foreignKey:ReportedBy" json:"-"`
	Title          string      `json:"title" gorm:"size:255;not null"`
	Description    string      `json:"description" gorm:"type:text;not null"`
	CategoryID     *uint       `gorm:"index" json:"category_id,omitempty"`
	Category       Category    `gorm:"foreignKey:CategoryID" json:"-"`
	SubCategoryID  uint        `gorm:"index" json:"sub_category_id,omitempty"`
	SubCategory    SubCategory `gorm:"foreignKey:SubCategoryID" json:"-"`
	Priority       string      `json:"priority" gorm:"size:50;not null"`
	Tags           []Tag       `json:"tags" gorm:"type:text[]"` // Use pq.StringArray for PostgreSQL; adjust for MySQL if necessary
	AttachmentURL  string      `json:"attachment_url" gorm:"size:255"`
	HasAttachments bool        `json:"has_attachments"`
	Severity       string      `gorm:"type:enum('Low', 'Medium', 'High', 'Critical');not null" json:"severity"`
	Status         string      `gorm:"type:enum('Open', 'Investigating', 'Resolved', 'Closed');not null" json:"status"`
	ResolvedAt     *time.Time  `json:"resolved_at"`
	ClosedAt       *time.Time  `json:"closed_at,omitempty"`
	TicketID       *uint       `json:"ticket_id" gorm:"foreignKey:TicketID"`
}

func (Incident) TableName() string {
	return "incidents"
}

// IncidentHistoryEntry represents a historical entry related to an incident.
type IncidentHistoryEntry struct {
	gorm.Model                  // Includes ID, CreatedAt, UpdatedAt, and DeletedAt automatically
	IncidentID  uint            `json:"incident_id" gorm:"not null;index"`
	Description string          `json:"description" gorm:"type:text;not null"`
	Status      string          `json:"status" gorm:"size:100;not null"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (IncidentHistoryEntry) TableName() string {
	return "incident_history_entry"
}

// IncidentComment represents a comment made on an incident.
type IncidentComment struct {
	gorm.Model                 // Includes ID, CreatedAt, UpdatedAt, and DeletedAt automatically
	IncidentID uint            `json:"incident_id" gorm:"not null;index"`
	Comment    string          `json:"comment" gorm:"type:text;not null"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
	DeletedAt  *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (IncidentComment) TableName() string {
	return "incident_comment"
}

// IncidentStorage defines the methods for managing incidents.
type IncidentStorage interface {
	GetAllIncidents() ([]*Incident, error)
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
	DB             *gorm.DB
	log            Logger
	EventPublisher *EventPublisherImpl
}

// NewIncidentDBModel creates a new instance of IncidentDBModel.
func NewIncidentDBModel(db *gorm.DB, log Logger, eventPublisher *EventPublisherImpl) *IncidentDBModel {
	return &IncidentDBModel{
		DB:             db,
		log:            log,
		EventPublisher: eventPublisher,
	}
}

// GetAllIncidents retrieves all service requests from the database.
func (idm *IncidentDBModel) GetAllIncidents() ([]*Incident, error) {
	var incidents []*Incident
	err := idm.DB.Find(&incidents).Error
	return incidents, err
}

// CreateIncident creates a new incident report.
func (idm *IncidentDBModel) CreateIncident(incident *Incident) error {
	return idm.DB.Create(incident).Error
}

// ReportIncident records a new incident in the database.
func (db *IncidentDBModel) ReportIncident(incident *Incident) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(incident).Error; err != nil {
			return err
		}
		return nil
	})
}

// UpdateIncidentStatus updates the status of an existing incident.
func (db *IncidentDBModel) UpdateIncidentStatus(incidentID uint, status string) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var incident Incident
		if err := tx.First(&incident, incidentID).Error; err != nil {
			return err
		}
		incident.Status = status
		if err := tx.Save(&incident).Error; err != nil {
			return err
		}
		return nil
	})
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
func (im *IncidentDBModel) UpdateIncidentCategory(incidentID uint, newCategory string) error {
	// Assuming GetIncidentByID fetches the incident directly
	incident, err := im.GetIncidentByID(incidentID)
	if err != nil {
		return err // Handle error if incident is not found
	}

	var c Category
	c.Name = newCategory

	// Update the category directly
	incident.Category = c

	// Assuming UpdateIncident updates the incident based on the struct's current state
	if err := im.UpdateIncident(incident); err != nil {
		return err // Handle error if the update fails
	}

	return nil // Return nil on success
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

// UpdateIncidentTags updates the tags of an incident by creating new tags if they do not exist
// and associating them with the given incident.
func (im *IncidentDBModel) UpdateIncidentTags(incidentID uint, tagNames []string) error {
	return im.DB.Transaction(func(tx *gorm.DB) error {
		// Fetch the incident by ID with its current tags loaded
		var incident Incident
		if err := tx.Preload("Tags").First(&incident, incidentID).Error; err != nil {
			return err // Incident not found or DB error
		}

		// Process each tagName to ensure the tag exists or create it
		var tagsToUpdate []*Tag
		for _, tagName := range tagNames {
			var tag Tag
			// Find an existing tag or create a new one if it doesn't exist
			if err := tx.FirstOrCreate(&tag, Tag{Name: tagName}).Error; err != nil {
				return err // Error handling tag
			}
			tagsToUpdate = append(tagsToUpdate, &tag)
		}

		// Associate the incident with the new set of tags
		// This replaces the incident's current tags with the new ones
		if err := tx.Model(&incident).Association("Tags").Replace(tagsToUpdate); err != nil {
			return err // Error updating incident tags
		}

		return nil // Success
	})
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
