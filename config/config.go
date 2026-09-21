package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
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
	BaseURL      string
}

var App *Config

func Load() *Config {
	loadEnvFile()

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
		BaseURL:      getEnv("BASE_URL", "http://localhost:8080"),
	}

	return App
}

func loadEnvFile() {
	// Try multiple locations for .env file
	envPaths := []string{
		".env",                          // Current working directory
		"../.env",                       // Parent directory
		"/etc/prangibar/.env",          // System config directory
	}

	// Try to get executable directory (works on Linux/Windows)
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		envPaths = append(envPaths, filepath.Join(exeDir, ".env"))
	}

	// Try to get source directory (for development)
	if _, filename, _, ok := runtime.Caller(0); ok {
		srcDir := filepath.Dir(filename)
		envPaths = append(envPaths, filepath.Join(srcDir, ".env"))
	}

	for _, path := range envPaths {
		if _, err := os.Stat(path); err == nil {
			if err := godotenv.Load(path); err == nil {
				log.Printf("Loaded .env from: %s", path)
				return
			}
		}
	}

	log.Println("No .env file found, using environment variables")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
