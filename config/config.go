package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort      int
	AppName      string
	DBHost       string
	DBPort       int
	DBDatabase   string
	DBUser       string
	DBPassword   string
	JWTSecret    string
	JWTExpiresIn time.Duration
	FrontendURL  string
	CORSOrigin   string
}

var App *Config

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	port, _ := strconv.Atoi(getEnv("APP_PORT", "3000"))
	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "3306"))
	jwtExp, _ := strconv.Atoi(getEnv("JWT_EXPIRES_IN_HOURS", "168"))

	App = &Config{
		AppPort:      port,
		AppName:      getEnv("APP_NAME", "Backend pendataan perusahaan SE2026-L.UB"),
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:       dbPort,
		DBDatabase:   getEnv("DB_DATABASE", "prangibar"),
		DBUser:       getEnv("DB_USER", "root"),
		DBPassword:   getEnv("DB_PASSWORD", ""),
		JWTSecret:    getEnv("JWT_SECRET", "default-secret-change-me"),
		JWTExpiresIn: time.Duration(jwtExp) * time.Hour,
		FrontendURL:  getEnv("FRONTEND_URL", "http://localhost:5173"),
		CORSOrigin:   getEnv("CORS_ORIGIN", "http://localhost:5173"),
	}

	return App
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
