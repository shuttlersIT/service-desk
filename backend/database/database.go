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
	return db
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
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return fmt.Errorf("failed to connect to database: %v", err)
		}
		DB = db
		// maximum idle connections and maximum open connections to avoid resource leaks.
		sqlDB, err := DB.DB() //
		if err != nil {
			return fmt.Errorf("failed to get DB instance: %v", err)
		}

		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)

		return db, nil
	}

	// Use exponential backoff for database connection
	backoffCfg := backoff.NewExponentialBackOff()
	backoffCfg.MaxElapsedTime = 10 * time.Second
	if err := backoff.Retry(operation, backoffCfg); err != nil {
		log.Error("Failed to establish a connection to the database:", err)
		fmt.Errorf("exceeded max retry attempts: %v", err)
		return nil, err
	}

	return db, nil
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
func SetupDatabase(config *config.Config, log models.Logger) {
	// Initialize the database connection
	db, err := InitializeDB(config, log)
	if err != nil {
		log.Fatal("Error initializing database:", err)
		return
	}
	defer CloseDB()

	// Start a transaction
	tx := db.Begin()
	if tx.Error != nil {
		log.Fatal("Error starting transaction:", tx.Error)
		return
	}

	err = AutoMigrateModels(tx)
	if err != nil {
		tx.Rollback()
		log.Fatal("Error migrating models:", err)
		return
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

// Temp Additional Migration Steps:
// additionalMigrationSteps includes any additional steps beyond AutoMigrate
func additionalMigrationSteps(db *gorm.DB, log models.Logger) error {
	// Add any additional migration steps here
	log.Info("Performing additional migration steps...")

	// Example 1: Adding an index on the 'email' column in the 'users' table
	if err := addIndex(db, "users", "idx_email", "email"); err != nil {
		return err
	}

	// Example 2: Adding a foreign key constraint between 'user_id' in 'comments' and 'id' in 'users'
	if err := addForeignKeyConstraint(db, "comments", "fk_comments_user_id", "user_id", "users(id)"); err != nil {
		return err
	}

	// Example 3: Adding a check constraint to ensure 'priority' in 'tickets' is within a specific range
	if err := addCheckConstraint(db, "tickets", "chk_priority_range", "priority >= 1 AND priority <= 5"); err != nil {
		return err
	}

	// Example 4: Creating a new table 'audit_logs' with a timestamp column
	if err := createTable(db, "audit_logs", "id SERIAL PRIMARY KEY, created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP"); err != nil {
		return err
	}

	// Example 5: Renaming the column 'old_name' to 'new_name' in the 'users' table
	if err := renameColumn(db, "users", "old_name", "new_name"); err != nil {
		return err
	}

	// Example 6: Dropping an existing column 'obsolete_column' from the 'tickets' table
	if err := dropColumn(db, "tickets", "obsolete_column"); err != nil {
		return err
	}

	// Example 7: Modifying the data type of 'amount' column in 'transactions' table
	if err := modifyColumnType(db, "transactions", "amount", "DECIMAL(10,2)"); err != nil {
		return err
	}

	// Example 8: Adding a unique constraint on a combination of columns in the 'orders' table
	if err := addUniqueConstraint(db, "orders", "uq_order_customer", "order_id, customer_id"); err != nil {
		return err
	}

	// Example 9: Creating a composite index on multiple columns in the 'sales' table
	if err := addIndex(db, "sales", "idx_sales_product_customer", "product_id, customer_id"); err != nil {
		return err
	}

	// Example 10: Adding a trigger to automatically update 'last_modified' timestamp in 'articles'
	if err := addTrigger(db, "articles", "trg_articles_last_modified", "BEFORE UPDATE", "SET NEW.last_modified = CURRENT_TIMESTAMP"); err != nil {
		return err
	}

	// Example 11: Creating a view that combines data from 'users' and 'orders' tables
	if err := createView(db, "customer_orders", "SELECT u.*, o.order_id FROM users u JOIN orders o ON u.id = o.customer_id"); err != nil {
		return err
	}

	// Example 12: Adding a default value to the 'status' column in the 'tasks' table
	if err := setDefaultValue(db, "tasks", "status", "pending"); err != nil {
		return err
	}

	// Example 13: Creating a stored procedure for calculating ticket response time
	if err := createStoredProcedure(db, "calculate_ticket_response_time", "CREATE OR REPLACE FUNCTION calculate_ticket_response_time() RETURNS TRIGGER AS $$ BEGIN NEW.response_time := NEW.updated_at - NEW.created_at; RETURN NEW; END; $$ LANGUAGE PLPGSQL"); err != nil {
		return err
	}

	// Example 14: Adding a unique index on 'username' in 'users' table
	if err := addUniqueIndex(db, "users", "idx_users_username", "username"); err != nil {
		return err
	}

	// Example 15: Adding a spatial index on the 'location' column in 'places' table
	if err := addSpatialIndex(db, "places", "idx_places_location", "location"); err != nil {
		return err
	}

	// Example 16: Modifying the length of 'description' column in 'products' table
	if err := modifyColumnLength(db, "products", "description", 500); err != nil {
		return err
	}

	// Example 17: Adding a unique constraint on 'employee_id' in 'employees' table
	if err := addUniqueConstraint(db, "employees", "uq_employees_employee_id", "employee_id"); err != nil {
		return err
	}

	// Example 18: Adding a foreign key constraint with ON DELETE CASCADE
	if err := addForeignKeyConstraint(db, "orders", "fk_orders_customer_id", "customer_id", "customers(id) ON DELETE CASCADE"); err != nil {
		return err
	}

	// Example 19: Renaming the 'old_table' to 'new_table'
	if err := renameTable(db, "old_table", "new_table"); err != nil {
		return err
	}

	// Example 20: Dropping an existing index 'idx_old_index' from the 'some_table' table
	if err := dropIndex(db, "some_table", "idx_old_index"); err != nil {
		return err
	}

	log.Info("Additional migration steps completed successfully.")
	return nil
}

// addIndex adds an index on the specified column in the given table
func addIndex(db *gorm.DB, tableName, indexName, columnName string) error {
	query := fmt.Sprintf("CREATE INDEX %s ON %s (%s)", indexName, tableName, columnName)
	if err := db.Exec(query).Error; err != nil {
		log.Error("Error adding index:", err)
		return err
	}
	log.Info(fmt.Sprintf("Index '%s' added on column '%s' in table '%s'", indexName, columnName, tableName))
	return nil
}

// addForeignKey adds a foreign key constraint on the specified column in the given table
func addForeignKey(db *gorm.DB, tableName, columnName, reference, onDelete, onUpdate string) error {
	query := fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT fk_%s FOREIGN KEY (%s) REFERENCES %s ON DELETE %s ON UPDATE %s",
		tableName, columnName, columnName, reference, onDelete, onUpdate)
	if err := db.Exec(query).Error; err != nil {
		log.Error("Error adding foreign key constraint:", err)
		return err
	}
	log.Info(fmt.Sprintf("Foreign key constraint added on column '%s' in table '%s' referencing '%s'", columnName, tableName, reference))
	return nil
}

// addUniqueConstraint adds a unique constraint on the specified column in the given table
func addUniqueConstraint(db *gorm.DB, tableName, constraintName, columnName string) error {
	query := fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s UNIQUE (%s)", tableName, constraintName, columnName)
	if err := db.Exec(query).Error; err != nil {
		log.Error("Error adding unique constraint:", err)
		return err
	}
	log.Info(fmt.Sprintf("Unique constraint '%s' added on column '%s' in table '%s'", constraintName, columnName, tableName))
	return nil
}

// setDefaultValue sets a default value for the specified column in the given table
func setDefaultValue(db *gorm.DB, tableName, columnName, defaultValue string) error {
	query := fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s SET DEFAULT '%s'", tableName, columnName, defaultValue)
	if err := db.Exec(query).Error; err != nil {
		log.Error("Error setting default value:", err)
		return err
	}
	log.Info(fmt.Sprintf("Default value '%s' set for column '%s' in table '%s'", defaultValue, columnName, tableName))
	return nil
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
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
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
func InitializeDatabase(dbURL string) (*gorm.DB, error) {
	db, err := InitializeDB(dbURL)
	if err != nil {
		return nil, err
	}

	// Migrate models to the database
	if err := MigrateModels(); err != nil {
		return nil, err
	}

	return db, nil
}
