package database

import (
	"blog_platform/config"
	"blog_platform/internal/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	// 1. Dynamically build the DSN (Data Source Name) for the Primary (Write) DB
	dsnWrite := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBWriteHost, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBPort)

	// 2. Connect to the PRIMARY (Write) database
	db, err := gorm.Open(postgres.Open(dsnWrite), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("Failed to connect database %w", err)
	}

	// 3. Dynamically build the DSN for the Replica (Read) DB
	dsnRead := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBReadHost, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBReadPort)

	// 4. Configure the DB Resolver plugin for the REPLICA (Read)
	resolverConfig := dbresolver.Register(dbresolver.Config{
		Replicas: []gorm.Dialector{postgres.Open(dsnRead)},
		Sources:  []gorm.Dialector{postgres.Open(dsnWrite)},
		Policy:   dbresolver.RandomPolicy{},
	})

	// 3. Attach the plugin
	err = db.Use(resolverConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to use dbresolver plugin %w", err)
	}

	// Enable the uuid-ossp extension for uuid_generate_v4()
	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)

	// Setup explicitly the many2many join table for soft deletes
	err = db.SetupJoinTable(&models.Post{}, "Tags", &models.PostTag{})
	if err != nil {
		log.Fatalf("failed to setup join table: %v", err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Comment{},
		&models.Tag{},
		&models.Reaction{},
		&models.PostView{},
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
