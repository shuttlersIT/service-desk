// backend/main.go

package main

import (
	"fmt"
	"io/ioutil"
	_ "log"
	"os"
	"regexp"
	"strings"
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

// events := models.NewEventsDBModel(db, log)

func main() {

	// Initialize Gin Engine
	r := gin.Default()
	fmt.Printf("server on")
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))
	gin.SetMode(gin.ReleaseMode)
	fmt.Printf("gin mode")

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

	db := database.SetupDatabase(configuration, log)

	/*
		// Initialize and migrate the database
		d, err := InitializeDatabase(configuration, log)
		if err != nil {
			fmt.Printf("Failed to initialize database: %v\n", err)
			return
		}
		defer CloseDB()

		fmt.Println("Database initialized and migrated successfully!")
	*/

	/* // Legacy DB Init
	// Auto Migrate Database Models (if not already migrated)
	db, err := InitDB(configuration, log)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error initializing database: %v", err))
	}

	// Auto Migrate Database Models (if not already migrated)
	err = AutoMigrateModels(db)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error migrating database models: %v", err))
	}
	fmt.Printf("gorm migration successful")
	*/

	/*
		// Check if database tables have been setup
		fmt.Printf(configuration.DBSetupStatus)
		if configuration.DBSetupStatus == "false" {
			//err := setup(db, log)
			//if err != nil {
			log.Fatal(fmt.Sprintf("Error setting up database tables: %v", err))
			//}

			//configuration.DBSetupStatus = "true"
			er := config.DBStatusUpdate(configuration)
			if er != nil {
				log.Info("Error updating DB setup status: %v", er)
			}
			log.Info(fmt.Sprintf("DB SETUP STATUS: %v", configuration.DBSetupStatus))
		}
	*/
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

	return db, nil

}

// ////////////////////////////////////////////////////////////////
// executeSQLScript executes a SQL script in a transaction
func executeSQLScript(db *gorm.DB, sqlScript string, logger models.Logger) error {
	// Splitting the script into individual statements
	statements := strings.Split(sqlScript, ";\n")

	// Starting a transaction
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}
	defer commitOrRollback(tx, logger)

	// Executing each statement
	count := 0
	for i, statement := range statements {
		trimmedStatement := strings.TrimSpace(statement)
		if trimmedStatement != "" {
			if err := tx.Exec(trimmedStatement).Error; err != nil {
				// Log the error and rollback the transaction
				logger.LogWithStackTrace(models.ErrorLevel, err, "Error executing statement [%d]: %s", i, trimmedStatement)
				return fmt.Errorf("failed to execute statement [%d] '%s': %w", i, trimmedStatement, err)
			}
			count = count + 1
		}
	}
	logger.Info("successfully executed [%d] statements", count)
	return nil
}

// commitOrRollback commits or rolls back a transaction based on error presence
func commitOrRollback(tx *gorm.DB, logger models.Logger) {
	// Recover from panics and commit or rollback the transaction
	if r := recover(); r != nil {
		tx.Rollback()
		logger.LogWithStackTrace(models.ErrorLevel, fmt.Errorf("%v", r), "Panic occurred during transaction")
		panic(r) // Re-throw panic after rollback
	} else if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		logger.Error("Error committing transaction:", err)
	}
}

// readAndSortScript reads a MySQL script from a file and returns statements sorted topologically
func readAndSortScript(filePath string) (string, error) {
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Sort statements topologically
	sortedScript, err := topologicalSort(string(fileContent))
	if err != nil {
		return "", err
	}

	return sortedScript, nil
}

// topologicalSort performs a topological sort on a SQL script
func topologicalSort(sqlScript string) (string, error) {
	// Parse SQL script to identify dependencies
	dependencies, err := parseDependencies(sqlScript)
	if err != nil {
		return "", fmt.Errorf("failed to parse dependencies: %w", err)
	}

	// Create a map to represent the graph of dependencies
	graph := make(map[string][]string)
	for _, dep := range dependencies {
		graph[dep.Table] = append(graph[dep.Table], dep.Dependency)
	}

	// Initialize visited and stack
	visited := make(map[string]bool)
	stack := make([]string, 0)

	// Perform DFS to create topological order
	for table := range graph {
		if !visited[table] {
			err := dfs(table, graph, visited, &stack)
			if err != nil {
				return "", fmt.Errorf("failed to perform DFS: %w", err)
			}
		}
	}

	// Reverse the stack to get the correct order
	reverse(stack)

	// Reorder the statements based on the topological order
	var reorderedStatements []string
	for _, table := range stack {
		reorderedStatements = append(reorderedStatements, getStatementsForTable(table, sqlScript))
	}

	return strings.Join(reorderedStatements, ";\n"), nil
}

// parseDependencies parses table dependencies from SQL script
func parseDependencies(sqlScript string) ([]config.Dependency, error) {
	var dependencies []config.Dependency

	// Regular expression to match FOREIGN KEY constraints
	re := regexp.MustCompile(`FOREIGN KEY \(\w+\) REFERENCES (\w+)`)

	// Find all matches in the script
	matches := re.FindAllStringSubmatch(sqlScript, -1)

	// Extract dependencies from matches
	for _, match := range matches {
		if len(match) == 2 {
			dependency := config.Dependency{
				Table:      match[1],
				Dependency: match[1],
			}
			dependencies = append(dependencies, dependency)
		}
	}

	return dependencies, nil
}

// dfs performs depth-first search for topological sort
func dfs(table string, graph map[string][]string, visited map[string]bool, stack *[]string) error {
	visited[table] = true

	for _, neighbor := range graph[table] {
		if !visited[neighbor] {
			err := dfs(neighbor, graph, visited, stack)
			if err != nil {
				return err
			}
		}
	}

	*stack = append(*stack, table)
	return nil
}

// reverse reverses the order of elements in a string slice
func reverse(stack []string) {
	for i, j := 0, len(stack)-1; i < j; i, j = i+1, j-1 {
		stack[i], stack[j] = stack[j], stack[i]
	}
}

// getStatementsForTable retrieves SQL statements related to a table
func getStatementsForTable(table string, sqlScript string) string {
	// Regular expression to find statements related to a table
	re := regexp.MustCompile(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s \(.*?;\n\)`, table))

	// Find the match in the script
	match := re.FindString(sqlScript)

	return match
}

// setup function for executing the MySQL deployment script
func setup(db *gorm.DB, logger models.Logger) error {
	sortedScript, err := readAndSortScript("service-desk.mysql")
	if err != nil {
		logger.LogWithStackTrace(models.ErrorLevel, err, "Error reading and sorting SQL script")
		return err
	}

	// Start a transaction for executing the script
	tx := db.Begin()
	if tx.Error != nil {
		logger.LogWithStackTrace(models.ErrorLevel, tx.Error, "Error starting transaction")
		return tx.Error
	}
	defer commitOrRollback(tx, logger)

	// Execute the sorted SQL script within the transaction
	err = executeSQLScript(tx, sortedScript, logger)
	if err != nil {
		logger.LogWithStackTrace(models.ErrorLevel, err, "Error executing SQL script")
		return err
	}

	logger.Info("MySQL deployment script executed successfully.")
	return nil
}

/*
// readConfig reads configuration from environment variables
func readConfig() (*database.Config, error) {
	config := &database.Config{
		DBUsername: os.Getenv("docker"),
		DBPassword: os.Getenv("itrootpassword"),
		DBHost:     os.Getenv("db"),
		DBPort:     os.Getenv("3306"),
		DBName:     os.Getenv("itsm"),
	}

	// Check if required config fields are set
	if config.DBUsername == "" || config.DBPassword == "" || config.DBHost == "" || config.DBPort == "" || config.DBName == "" {
		return nil, fmt.Errorf("missing required configuration values")
	}

	return config, nil
}

// createGormConnection creates a Gorm database connection
func createGormConnection(config *database.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config.DBUsername, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}


// New function to check if database setup is required
func isDBSetupRequired(db *gorm.DB, config *config.Config, logger models.Logger) (bool, error) {
	// Check if the database tables have been set up
	if config.DBSetupStatus == "false" {
		return true, nil
	}

	// If the tables exist, check if they match the expected structure
	if !database.TablesMatchStructure(db, logger) {
		return true, nil
	}

	return false, nil
}

// New function to execute the setup process
func executeSetup(db *gorm.DB, config *config.Config, logger models.Logger) error {
	sortedScript, err := readAndSortScript("service-desk.mysql")
	if err != nil {
		logger.LogWithStackTrace(models.ErrorLevel, err, "Error reading and sorting SQL script")
		return err
	}

	// Start a transaction for executing the script
	tx := db.Begin()
	if tx.Error != nil {
		logger.LogWithStackTrace(models.ErrorLevel, tx.Error, "Error starting transaction")
		return tx.Error
	}
	defer commitOrRollback(tx, logger)

	// Execute the sorted SQL script within the transaction
	err = executeSQLScript(tx, sortedScript, logger)
	if err != nil {
		logger.LogWithStackTrace(models.ErrorLevel, err, "Error executing SQL script")
		return err
	}

	logger.Info("MySQL deployment script executed successfully.")

	// Update the database setup status
	config.DBSetupStatus = "true"
	//err = config.DBStatusUpdate(config)
	//if err != nil {
	//	logger.Info("Error updating DB setup status: %v", err)
	//}

	logger.Info(fmt.Sprintf("DB SETUP STATUS: %v", config.DBSetupStatus))
	return nil
}

// ...

// Modified InitDB function
func InitDB(config *config.Config, log models.Logger) (*gorm.DB, error) {
	// Initialize Database
	db, err := database.InitializeMySQLConnection(config, log)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error initializing database: %v", err))
		return nil, err
	}
	fmt.Printf("db connected")

	// Check if database setup is required
	setupRequired, err := isDBSetupRequired(db, config, log)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error checking database setup requirement: %v", err))
		return nil, err
	}

	if setupRequired {
		// Execute the setup process including the deployment script
		err := executeSetup(db, config, log)
		if err != nil {
			log.Fatal(fmt.Sprintf("Error executing setup: %v", err))
			return nil, err
		}
	}

	// Auto Migrate Database Models (if not already migrated)
	if err := AutoMigrateModels(db); err != nil {
		log.Fatal(fmt.Sprintf("Error migrating database models: %v", err))
		return nil, err
	}
	fmt.Printf("gorm migration successful")

	return db, nil
}
*/

// InitializeDatabase initializes the MySQL database with models and migrates them.
func InitializeDatabase(config *config.Config, log models.Logger) (*gorm.DB, error) {
	db, err := database.InitializeDB(config, log)
	if err != nil {
		return nil, err
	}
	fmt.Println("DB INITIALIZED")

	// Migrate models to the database
	if err := database.AutoMigrateModels(db, log); err != nil {
		return nil, err
	}
	fmt.Println("MODELS MIGRATED")

	return db, nil
}

// CloseDB closes the MySQL database connection.
func CloseDB() {
	if db != nil {
		dbSQL, err := db.DB()
		if err == nil {
			dbSQL.Close()
		}
	}
}
