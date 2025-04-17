package database

import (
	"log"
	"os"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"errors"
)

// ConnectDatabase establishes a connection to the MySQL database
func ConnectDatabase() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, errors.New("DATABASE_URL environment variable is not set")
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	return db, nil
}
