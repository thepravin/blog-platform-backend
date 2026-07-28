package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort            string
	DBWriteHost        string
	DBReadHost         string
	DBPort             string
	DBReadPort         string
	DBUser             string
	DBPass             string
	DBName             string
	JWTSecret          string
	CORSAllowedOrigins []string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	return &Config{
		AppPort:            getEnv("APP_PORT", "8080"),
		DBWriteHost:        getEnv("DB_WRITE_HOST", "localhost"),
		DBReadHost:         getEnv("DB_READ_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "5432"),
		DBReadPort:         getEnv("DB_READ_PORT", "5433"),
		DBUser:             getEnv("DB_USER", "postgres"),
		DBPass:             getEnv("DB_PASSWORD", "postgres"),
		DBName:             getEnv("DB_NAME", "blog_db"),
		JWTSecret:          getEnv("JWT_SECRET", "kfdjsfjosifsfsf85648df"),
		CORSAllowedOrigins: strings.Split(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5173"), ","),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
