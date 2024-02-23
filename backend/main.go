// backend/main.go

package main

import (
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

	gin.SetMode(gin.ReleaseMode)

	// Auto Migrate Database Models (if not already migrated)
	db := InitDB()

	log := models.NewLogger()
	//events := models.NewEventsDBModel(db, log)
	eventPublisher := models.NewEventPublisher()

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
	assetService := services.NewDefaultAssetService(assetDBModel, assetAssignmentDBModel, *log, eventPublisher)
	authService := services.NewAuthService(db, authDBModel, *log, eventPublisher)
	serviceRequestService := services.NewDefaultServiceRequestService(db, serviceRequestDBModel, *log, eventPublisher)
	incidentService := services.NewDefaultIncidentService(db, incidentDBModel, *log, eventPublisher)

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
