package main

import (
	"log"

	"github.com/shuttlersit/service-desk/database"
	"github.com/shuttlersit/service-desk/models"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {

	// Initialize Database
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	// Auto Migrate Database Models (if not already migrated)
	db.AutoMigrate(
		&models.Users{},
		&models.Agents{},
		&models.Unit{},
		&models.Permission{},
		&models.Teams{},
		&models.Role{},
		&models.TeamPermission{},
		&models.RoleBase{},
		&models.RolePermission{},
		&models.AgentRole{},
		&models.UserRole{},
		&models.UserAgent{},
		&models.TeamAgent{},
		&models.AgentPermission{},
		&models.Position{},
		&models.Department{},
		&models.Ticket{},
		&models.Comment{},
		&models.TicketHistoryEntry{},
		&models.RelatedTicket{},
		&models.Tag{},
		&models.SLA{},
		&models.Priority{},
		&models.Satisfaction{},
		&models.Category{},
		&models.SubCategory{},
		&models.Status{},
		&models.Policies{},
		&models.TicketMediaAttachment{},
		&models.Session{},
		&models.UserAgentMapping{},
		&models.UserAgentAccess{},
		&models.UserAgentGroup{},
		&models.GroupMember{},
		&models.Location{},
		&models.ServiceRequestComment{},
		&models.ServiceRequestHistoryEntry{},
		&models.Incident{},
		&models.IncidentHistoryEntry{},
		&models.IncidentComment{},
		&models.GoogleCredentials{},
		&models.AgentLoginCredentials{},
		&models.UsersLoginCredentials{},
		&models.PasswordHistory{},
		&models.AgentUserMapping{},
		&models.Assets{},
		&models.AssetTag{},
		&models.AssetType{},
		&models.AssetAssignment{},
		&models.ServiceRequest{},
	)
	return db
}
