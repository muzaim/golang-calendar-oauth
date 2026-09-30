package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort                  string
	AppEnv                   string
	DBHost                   string
	DBPort                   string
	DBUser                   string
	DBPassword               string
	DBName                   string
	JWTSecret                string
	JWTAccessExpirationHours int
	JWTRefreshExpirationDays int
	GoogleClientID           string
	GoogleClientSecret       string
	GoogleRedirectURL        string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("Notice: .env file not found, reading from environment variables")
	}

	accessExpHours, _ := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "24"))
	refreshExpDays, _ := strconv.Atoi(getEnv("JWT_REFRESH_EXPIRATION_DAYS", "7"))

	cfg := &Config{
		AppPort:                  getEnv("APP_PORT", "8080"),
		AppEnv:                   getEnv("APP_ENV", "development"),
		DBHost:                   getEnv("DB_HOST", "127.0.0.1"),
		DBPort:                   getEnv("DB_PORT", "3306"),
		DBUser:                   getEnv("DB_USER", "root"),
		DBPassword:               getEnv("DB_PASSWORD", ""),
		DBName:                   getEnv("DB_NAME", "calendar_service"),
		JWTSecret:                getEnv("JWT_SECRET", "supersecretjwtkeychangeinproduction"),
		JWTAccessExpirationHours: accessExpHours,
		JWTRefreshExpirationDays: refreshExpDays,
		GoogleClientID:           getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret:       getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:        getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/api/auth/google/callback"),
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}
