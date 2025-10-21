package db

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"event-management/config"
	"event-management/models"
)

var DB *gorm.DB

// ConnectDatabase connects to the database
func ConnectDatabase(cfg *config.Config) {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connection successful.")

	// Migrate the schema
	DB.AutoMigrate(&models.User{}, &models.Event{}, &models.Guest{})
	log.Println("Database migration successful.")
}
