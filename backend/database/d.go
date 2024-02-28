// backend/database/database.go

package database

/*
import (
	"fmt"
	_ "log"

	"github.com/shuttlersit/service-desk/backend/models"
	"gorm.io/gorm"
)

var DB *gorm.DB

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
*/
