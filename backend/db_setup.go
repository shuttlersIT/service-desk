package main

import (
	"fmt"
	"io/ioutil"
	_ "log"
	"regexp"
	"strings"

	"github.com/shuttlersit/service-desk/backend/config"
	"github.com/shuttlersit/service-desk/backend/models"

	//"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

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
*/

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
