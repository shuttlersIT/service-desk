// backend/main.go

package main

import (
	"log"

	"github.com/shuttlersit/service-desk/backend/controllers"
	"github.com/shuttlersit/service-desk/backend/models"
	"github.com/shuttlersit/service-desk/backend/routes"
	"github.com/shuttlersit/service-desk/backend/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Gin Engine
	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	// Initialize Database
	/*db, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	// Auto Migrate Database Models (if not already migrated)
	db.AutoMigrate(
		&models.ServiceRequest{},
		&models.ServiceRequestHistoryEntry{},
		&models.ServiceRequestComment{},
		&models.Incident{},
		&models.Agents{},
		&models.Users{},
		&models.Ticket{},
		&models.Assets{},
		&models.GoogleCredentials{},
		&models.Session{},
	)*/

	db := InitDB()

	// Public routes
	r.GET("/auth/google", controllers.GoogleLogin)
	r.GET("/auth/google/callback", controllers.GoogleAuthCallback)

	agentDBModel := models.NewAgentDBModel(db)
	userDBModel := models.NewUserDBModel(db)
	ticketDBModel := models.NewTicketDBModel(db)
	assetAssignmentDBModel := models.NewAssetAssignmentDBModel(db)
	assetDBModel := models.NewAssetDBModel(db, assetAssignmentDBModel)
	authDBModel := models.NewAuthDBModel(db, userDBModel, agentDBModel)
	serviceRequestDBModel := models.NewServiceRequestDBModel(db)
	incidentDBModel := models.NewIncidentDBModel(db)

	// Initialize Services
	userService := services.NewDefaultUserService(userDBModel)
	agentService := services.NewDefaultAgentService(agentDBModel)
	ticketService := services.NewDefaultTicketingService(ticketDBModel)
	assetService := services.NewDefaultAssetService(assetDBModel, assetAssignmentDBModel)
	authService := services.NewDefaultAuthService(authDBModel)
	serviceRequestService := services.NewDefaultServiceRequestService(serviceRequestDBModel)
	incidentService := services.NewDefaultIncidentService(incidentDBModel)

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
