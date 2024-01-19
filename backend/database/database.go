// backend/database/database.go

package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// InitDatabase initializes the database connection
func InitDatabase() error {
	// Connect to the database (you can modify the DSN)
	dsn := "username:password@tcp(your-mysql-server:3306)/database-name"
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

// Check if a table exists in the database
func tableExists(db *sql.DB, tableName string) bool {
	query := "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?"
	var count int
	err := db.QueryRow(query, tableName).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	return count > 0
}
