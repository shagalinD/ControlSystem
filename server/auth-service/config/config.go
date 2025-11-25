package config

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
    DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    DBName     string
    JWTSecret  string
    ServicePort string
    Env        string
}

func Load() *Config {
    config := &Config{
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnv("DB_PORT", "5432"),
        DBUser:     getEnv("DB_USER", "postgres"),
        DBPassword: getEnv("DB_PASSWORD", ""),
        DBName:     getEnv("DB_NAME", "auth_db"),
        JWTSecret:  getEnv("JWT_SECRET", "development-secret-key"),
        ServicePort: getEnv("AUTH_SERVICE_PORT", "8081"),
        Env:        getEnv("ENV", "development"),
    }
    
    if err := config.validate(); err != nil {
        log.Fatal("Config validation failed:", err)
    }
    
    return config
}

func (c *Config) validate() error {
    if c.DBUser == "" {
        return fmt.Errorf("DB_USER is required")
    }
    if c.DBName == "" {
        return fmt.Errorf("DB_NAME is required")
    }
    return nil
}

func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}