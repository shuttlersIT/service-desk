package database

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() error {
	db, err := gorm.Open(sqlite.Open("service-desk.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	DB = db
	return nil
}
