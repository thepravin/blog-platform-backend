package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort            string
	DBHost             string
	DBUser             string
	DBPass             string
	DBName             string
	DBPort             string
	JWTSecret          string
	CORSAllowedOrigins []string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	return &Config{
		AppPort:            getEnv("APP_PORT", "8080"), // getEnv(env_key, default_value)
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBUser:             getEnv("DB_USER", "postgres"),
		DBPass:             getEnv("DB_PASSWORD", "postgres"),
		DBName:             getEnv("DB_NAME", "blogdb"),
		DBPort:             getEnv("DB_PORT", "5432"),
		JWTSecret:          getEnv("JWT_SECRET", "kfdjsfjosifsfsf85648df"),
		CORSAllowedOrigins: strings.Split(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5173/"), ","),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
