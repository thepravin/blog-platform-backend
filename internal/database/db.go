package database

import (
	"blog_platform/config"
	"blog_platform/internal/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", cfg.DBHost, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("Failed to connect database %w", err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Comment{},
		&models.Tag{},
		&models.Reaction{},
	)

	if err != nil {
		log.Fatalf("auto migration failed: %v", err)
	}

	log.Print("Connected to PostgreSQL and auto migration successfully!")

	return db, nil
}

func CloseDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("⚠️ Failed to close DB: %v", err)
	}

	sqlDB.Close()
	log.Printf("Darabase connection closed")
}
