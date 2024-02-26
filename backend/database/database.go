// backend/database/database.go

package database

import (
	"errors"
	"fmt"
	_ "log"

	"github.com/shuttlersit/service-desk/backend/config"
	"github.com/shuttlersit/service-desk/backend/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// createGormConnection creates a Gorm database connection
func CreateGormDsn(config *config.Config) (string, error) {
	//dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config.DBUsername, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config.DBUsername, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
	fmt.Printf(("%v : success\n"), dsn)
	return dsn, nil
}

var db *gorm.DB

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
	defer sqlDB.Close()

	err = sqlDB.Ping()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	table := "users"
	fmt.Println("Connected to the MySQL database")
	log.Info(fmt.Sprintf("Connecting to MySQL with DSN: %s\n", dsn))
	if !tableExists(db, table, log) {
		log.Error(fmt.Sprintf("Table %s not found in the database", table))
		return nil, errors.New("database schema incomplete")
	}

	return db, nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return db
}

// Check if a table exists in the database
func tableExists(db *gorm.DB, tableName string, log models.Logger) bool {
	query := "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?"
	var count int
	err := db.Exec(query, tableName).Scan(&count)
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
