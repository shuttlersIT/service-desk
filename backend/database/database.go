// backend/database/database.go

package database

import (
	"fmt"
	_ "log"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/shuttlersit/service-desk/backend/config"
	"github.com/shuttlersit/service-desk/backend/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// createGormConnection creates a Gorm database connection
func CreateGormDsn2(config *config.Config) (string, error) {
	//dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config.DBUsername, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config.DBUsername, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
	fmt.Printf(("%v : success\n"), dsn)
	return dsn, nil
}

// Change the signature of CreateGormDsn to return an error
func CreateGormDsn(config *config.Config) (string, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config.DBUsername, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
	return dsn, nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {

	return DB
}

// InitializeDB initializes the MySQL database connection.
func InitializeDB(config *config.Config, log models.Logger) (*gorm.DB, error) {
	dsn, err := CreateGormDsn(config)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	// Establish a connection to the database
	// var err error

	operation := func() error {
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return fmt.Errorf("failed to connect to database: %v", err)
		}
		DB = db
		// maximum idle connections and maximum open connections to avoid resource leaks.
		sqlDB, err := DB.DB()
		if err != nil {
			return fmt.Errorf("failed to get DB instance: %v", err)
		}

		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)

		DB = db

		return nil
	}

	// Use exponential backoff for database connection
	backoffCfg := backoff.NewExponentialBackOff()
	backoffCfg.MaxElapsedTime = 10 * time.Second
	if err := backoff.Retry(operation, backoffCfg); err != nil {
		log.Error("Failed to establish a connection to the database:", err)
		return nil, fmt.Errorf("exceeded max retry attempts: %v", err)
	}

	return DB, nil
}

// AutoMigrateModels performs auto-migration of all models within a transaction
func AutoMigrateModels(db *gorm.DB, log models.Logger) error {
	log.Info("Starting database auto-migration...")

	// Begin a transaction
	tx := db.Begin()
	if tx.Error != nil {
		log.Error("Error starting transaction:", tx.Error)
		return tx.Error
	}

	// Migrate all models
	if err := migrateAllModels(tx, log); err != nil {
		tx.Rollback()
		log.Error("Error migrating models:", err)
		return err
	}

	/*
		// Additional migration steps if needed
		if err := additionalMigrationSteps(tx, log); err != nil {
		    tx.Rollback()
		    log.Error("Error in additional migration steps:", err)
		    return err
		}
	*/

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		log.Error("Error committing transaction:", err)
		return err
	}

	log.Info("Database auto-migration completed successfully.")
	return nil
}

// migrateAllModels performs the actual migration of all models
func migrateAllModels(db *gorm.DB, log models.Logger) error {
	log.Info("Migrating all models...")

	err := db.AutoMigrate(
		&models.Users{},
		&models.UserRole{},
		&models.Position{},
		&models.Department{},
		&models.UserProfile{},
		&models.ProjectAssignment{},
		&models.Activity{},
		&models.Agents{},
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
		&models.Vendor{},
	)

	if err != nil {
		log.Error("Error auto-migrating models:", err)
		return err
	}

	log.Info("All models migrated successfully.")
	return nil
}

// CloseDB closes the database connection
func CloseDB() {
	sqlDB, _ := DB.DB()
	sqlDB.Close()
}

// SetupDatabase initializes and migrates the database
func SetupDatabase(config *config.Config, log models.Logger) *gorm.DB {
	// Initialize the database connection
	db, err := InitializeDB(config, log)
	if err != nil {
		log.Fatal("Error initializing database:", err)
		return nil
	}
	defer CloseDB()

	// Start a transaction
	tx := db.Begin()
	if tx.Error != nil {
		log.Fatal("Error starting transaction:", tx.Error)
		return nil
	}

	err = AutoMigrateModels(tx, log)
	if err != nil {
		tx.Rollback()
		log.Fatal("Error migrating models:", err)
		return nil
	}

	tx.Commit()
	fmt.Println("Database setup completed successfully.")
}

// Check if a table exists in the database
func tableExists(db *gorm.DB, tableName string, log models.Logger) bool {
	query := fmt.Sprintf("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = %v", tableName)
	fmt.Println("Query:", query) // Add this line to print the query
	var count int
	err := db.Exec(query).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	return count > 0
}

/*
// InitDatabase initializes the database connection
func InitDatabase() error {
	// Connect to the database (you can modify the DSN)
	dsn := "root:1T$hutt!ers@tcp(db:3307)/itsm"
	var err error
	db, err = sql.Open("mysql", dsn) // Change the driver and DSN as needed
	if err != nil {
		log.Fatal(err)
		return err
	}

	// Check if the database connection is successful
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Println("Connected to the database")

	// Check if the "users" table exists
	if tableExists(db, "users") {
		fmt.Println("The 'users' table exists.")
	} else {
		fmt.Println("The 'users' table does not exist.")
	}

	return nil
}

// GetDB returns the database instance
func GetDB() *sql.DB {
	return db
}

// CloseDatabase closes the database connection
func CloseDatabase() {
	if db != nil {
		db.Close()
		fmt.Println("Database connection closed")
	}
}
*/

func InitializeMySQLConnection(config *config.Config, log models.Logger) (*gorm.DB, error) {
	// Connect to the MySQL database
	//dsn := "docker:itrootpassword@tcp(db:3306)/itsm" // MySQL connection details
	//var err error
	dsn, err := CreateGormDsn(config)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Info(fmt.Sprintf(("%v : success\n"), dsn))
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// Check if the connection is successful
	sqlDB, err := db.DB()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	//defer sqlDB.Close()

	err = sqlDB.Ping()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	//table := "users"
	fmt.Println("Connected to the MySQL database")
	log.Info(fmt.Sprintf("Connecting to MySQL with DSN: %s\n", dsn))
	/*
		if !tableExists(db, table, log) {
			log.Error(fmt.Sprintf("Table %s not found in the database", table))
			return nil, errors.New("database schema incomplete")
		}
	*/
	return db, nil
}

// InitializeDatabase initializes the MySQL database with models and migrates them.
func InitializeDatabase(config *config.Config, log models.Logger) (*gorm.DB, error) {
	db, err := InitializeDB(config, log)
	if err != nil {
		return nil, err
	}

	// Migrate models to the database
	if err := AutoMigrateModels(db, log); err != nil {
		return nil, err
	}

	return db, nil
}
