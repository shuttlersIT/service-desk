// backend/database/mysql.go

package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var status string

func ConnectDB() (*gorm.DB, error) {
	// Initialize Database
	db, err := gorm.Open(sqlite.Open("root:1T$hutt!ers@tcp(mysqlDB:3306)/itsm"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	return db, err
}

func ConnectMysql() (string, *sql.DB) {

	// Replace with your database credentials
	db, err := sql.Open("mysql", "root:1T$hutt!ers@tcp(mysqlDB:3306)/itsm")
	if err != nil {
		log.Fatal(err)
		status = "Unable to connect to mysql database"
	}
	defer db.Close()

	// Check the connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	/*
		// Check if the "tickets" table exists
		if TableExists(db, "tickets") {
			fmt.Println("The 'tickets' table exists.")
			return status, db
		} else {
			fmt.Println("The 'tickets' table does not exist.")
			status = "table doesnt exit"
			return status, nil
		}
	*/
	fmt.Println("Database is on.")
	return status, db
}

// Create a table to store tickets in the database
func CreateTicketsTable(db *sql.DB) error {
	query := `
        CREATE TABLE IF NOT EXISTS tickets (
            id INT AUTO_INCREMENT PRIMARY KEY,
			subject VARCHAR(255) NOT NULL,
			description TEXT,
			category_id INT,
			sub_category_id INT,
			priority_id INT,
			sla_id INT,
			staff_id INT,
			agent_id INT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			due_at  DATE,
			asset_id INT,
			related_ticket_id INT,
			tag   VARCHAR(50) NOT NULL,
			site  VARCHAR(50) NOT NULL,
    		status ENUM('open', 'pending', 'in_progress', 'waiting for 3rd-Party-escalated', 'waiting for 3rd-Party-vendor', 'waiting for approval', 'waiting for feedback', 'closed'),
			attachment_id  INT,
    		FOREIGN KEY (category_id) REFERENCES category(id),
    		FOREIGN KEY (sub_category_id) REFERENCES sub_category(id),
    		FOREIGN KEY (staff_id) REFERENCES staff(id),
    		FOREIGN KEY (agent_id) REFERENCES agents(id),
    		FOREIGN KEY (asset_id) REFERENCES assets(id),
    		FOREIGN KEY (related_ticket_id) REFERENCES tickets(id),
    		FOREIGN KEY (sla_id) REFERENCES sla(id),
    		FOREIGN KEY (priority_id) REFERENCES priority(id),
    		FOREIGN KEY (attachment_id) REFERENCES attachments(id)
        )`
	_, err := db.Exec(query)
	return err
}

// Check if a table exists in the database
func TableExists(db *sql.DB, tableName string) bool {
	query := "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?"
	var count int
	err := db.QueryRow(query, tableName).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	return count > 0
}
