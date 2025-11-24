package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
    GatewayPort         string
    AuthServiceURL      string
    ProjectDefectServiceURL string
    ContentServiceURL   string
    JWTSecret           string
}

func Load() *Config {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using system environment variables")
    }
    
    return &Config{
        GatewayPort:         getEnv("GATEWAY_PORT", "8080"),
        AuthServiceURL:      getEnv("AUTH_SERVICE_URL", "http://auth-service:8081"),
        ProjectDefectServiceURL: getEnv("PROJECT_DEFECT_SERVICE_URL", "http://project-defect-service:8082"),
        ContentServiceURL:   getEnv("CONTENT_SERVICE_URL", "http://content-service:8083"),
        JWTSecret:           getEnv("JWT_SECRET", "development-secret-key"),
    }
}

func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}