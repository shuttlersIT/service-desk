// backend/main.go

package main

import (
	"fmt"
	_ "log"
	"os"
	"time"

	"github.com/shuttlersit/service-desk/backend/config"
	"github.com/shuttlersit/service-desk/backend/controllers"
	"github.com/shuttlersit/service-desk/backend/database"
	"github.com/shuttlersit/service-desk/backend/models"
	"github.com/shuttlersit/service-desk/backend/routes"
	"github.com/shuttlersit/service-desk/backend/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//events := models.NewEventsDBModel(db, log)

func main() {
	logConfig := models.LoggerConfig{
		MaxSizeMB:    10,
		MaxBackups:   3,
		RotatePeriod: 24 * time.Hour,
		LogLevel:     models.DebugLevel, // Adjust the log level based on your requirements
		LogOutput:    os.Stdout,
		LogFormat:    models.TextFormat, // Adjust the log format based on your requirements
	}

	log := models.NewLoggerWithConfig(logConfig)
	eventPublisher := models.NewEventPublisher()

	configuration, err := config.ReadConfig()
	if err != nil {
		log.Fatal(fmt.Sprintf("Error reading configuration: %v", err))
	}

	// Initialize Gin Engine
	r := gin.Default()
	fmt.Printf("server on")
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))
	gin.SetMode(gin.ReleaseMode)
	fmt.Printf("gin mode")

	// Auto Migrate Database Models (if not already migrated)
	db, err := InitDB(configuration, log)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error initializing database: %v", err))
	}

	/*
		if err := AutoMigrateModels(db); err != nil {
			log.Fatal(fmt.Sprintf("Error migrating database models: %v", err))
		}
		fmt.Printf("db and gorm migration done")
	*/

	// Check if database tables have been setup
	fmt.Printf(configuration.DBSetupStatus)
	if configuration.DBSetupStatus == "false" {
		err := setup(db, log)
		if err != nil {
			log.Fatal(fmt.Sprintf("Error setting up database tables: %v", err))
		}

		configuration.DBSetupStatus = "true"
		er := config.DBStatusUpdate(configuration)
		if er != nil {
			log.Info("Error updating DB setup status: %v", er)
		}
		log.Info(fmt.Sprintf("DB SETUP STATUS: %v", configuration.DBSetupStatus))
	}

	fmt.Printf("db service running")

	// Public routes
	r.GET("/auth/google", controllers.GoogleLogin)
	//r.GET("/auth/google/callback", controllers.GoogleAuthCallback)

	agentDBModel := models.NewAgentDBModel(db, log, eventPublisher)
	userDBModel := models.NewUserDBModel(db, log, eventPublisher)
	ticketDBModel := models.NewTicketDBModel(db, log, eventPublisher)
	ticketCommentDBModel := models.NewTicketCommentDBModel(db, log, eventPublisher)
	ticketHistoryDBModel := models.NewTicketHistoryEntryDBModel(db, log, eventPublisher)
	assetAssignmentDBModel := models.NewAssetAssignmentDBModel(db, log)
	assetDBModel := models.NewAssetDBModel(db, assetAssignmentDBModel, log, eventPublisher)
	authDBModel := models.NewAuthDBModel(db, userDBModel, agentDBModel, log, eventPublisher)
	//googleAuth := models.NewGoogleCredentialsDBModel(db, authDBModel, log)
	serviceRequestDBModel := models.NewServiceRequestDBModel(db, log, eventPublisher)
	incidentDBModel := models.NewIncidentDBModel(db, log, eventPublisher)

	// Initialize Services
	userService := services.NewDefaultUserService(userDBModel, log, eventPublisher)
	agentService := services.NewDefaultAgentService(db, eventPublisher, log, agentDBModel)
	ticketService := services.NewDefaultTicketingService(db, ticketDBModel, ticketCommentDBModel, ticketHistoryDBModel, userDBModel, agentDBModel, eventPublisher, log)
	assetService := services.NewDefaultAssetService(assetDBModel, assetAssignmentDBModel, log, eventPublisher)
	authService := services.NewAuthService(db, authDBModel, log, eventPublisher)
	serviceRequestService := services.NewDefaultServiceRequestService(db, serviceRequestDBModel, log, eventPublisher)
	incidentService := services.NewDefaultIncidentService(db, incidentDBModel, log, eventPublisher)

	// Initialize Controllers
	agentsController := controllers.NewAgentController(agentService)
	usersController := controllers.NewUserController(userService)
	ticketsController := controllers.NewTicketController(ticketService)
	assetsController := controllers.NewAssetController(assetService)
	serviceRequestController := controllers.NewServiceRequestController(serviceRequestService)
	incidentController := controllers.NewIncidentController(incidentService)
	authController := controllers.NewAuthController(authService)

	// Setup Routes
	routes.SetupAgentRoutes(r, agentsController)
	routes.SetupUserRoutes(r, usersController)
	routes.SetupTicketRoutes(r, ticketsController)
	routes.SetupAssetsRoutes(r, assetsController)
	routes.SetupServiceRequestRoutes(r, serviceRequestController)
	routes.SetupIncidentRoutes(r, incidentController)
	routes.SetupAuthRoutes(r, authController)
	routes.SetupOpenRoutes(r, authController)

	// Run the application
	if err := r.Run(":7788"); err != nil {
		log.Fatal("Error running the server: ", err)
	}
}

func InitDB(config *config.Config, log models.Logger) (*gorm.DB, error) {
	// Initialize Database
	db, err := database.InitializeMySQLConnection(config, log)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error initializing database: %v", err))
		return nil, err
	}
	fmt.Printf("db connected")

	// Auto Migrate Database Models (if not already migrated)
	if err := AutoMigrateModels(db); err != nil {
		log.Fatal(fmt.Sprintf("Error migrating database models: %v", err))
		return nil, err
	}
	fmt.Printf("gorm migration successful")

	return db, nil

}

// Auto Migrate Database Models (if not already migrated)
func AutoMigrateModels(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.Users{},
		&models.UserRole{},
		&models.Position{},
		&models.Department{},
		&models.UserProfile{},
		&models.ProjectAssignment{},
		&models.Activity{},
		/*&models.Agents{},
		&models.Unit{},
		&models.Permission{},
		&models.Teams{},
		&models.Role{},
		&models.TeamPermission{},
		&models.RoleBase{},
		&models.RolePermission{},
		&models.AgentRole{},
		&models.UserAgent{},
		&models.TeamAgent{},
		&models.AgentPermission{},
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
		&models.ServiceRequest{}, */
	)
	return err
}
