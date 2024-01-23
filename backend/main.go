// backend/main.go

package main

import (
	"log"

	"github.com/shuttlersit/service-desk/backend/controllers"
	"github.com/shuttlersit/service-desk/backend/database"
	"github.com/shuttlersit/service-desk/backend/models"
	"github.com/shuttlersit/service-desk/backend/routes"
	"github.com/shuttlersit/service-desk/backend/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Gin Engine
	r := gin.Default()

	// Initialize Database
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	// Auto Migrate Database Models (if not already migrated)
	db.AutoMigrate(
		&models.ServiceRequest{},
		&models.ServiceRequestHistoryEntry{},
		&models.ServiceRequestComment{},
	)
	serviceRequestDBModel := models.NewServiceRequestDBModel(db)

	// Initialize Services
	serviceRequestService := services.NewDefaultServiceRequestService(db, serviceRequestDBModel)

	// Initialize Controllers
	serviceRequestController := controllers.NewServiceRequestController(serviceRequestService)

	// Setup Routes
	routes.SetupServiceRequestRoutes(r, serviceRequestController)

	// Run the application
	r.Run(":6195")
}
