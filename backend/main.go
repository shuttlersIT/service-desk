// backend/main.go

package main

import (
	"github.com/shuttlersit/service-desk/backend/database"
	//"github.com/shuttlersit/service-desk/backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the database
	err := database.InitDatabase()
	if err != nil {
		panic("Failed to connect to the database")
	}
	// Create a Gin router
	router := gin.Default()
	// Set up API routes
	//routes.SetupRoutes(router)
	// Start the server
	router.Run(":8080")
}
