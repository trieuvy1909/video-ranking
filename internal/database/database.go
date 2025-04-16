package database

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ConnectDatabase establishes a connection to the MySQL database
func ConnectDatabase() (*gorm.DB, error) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatalf("DATABASE_URL environment variable is not set")
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	return db, nil
}
