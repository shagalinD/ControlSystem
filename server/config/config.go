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
    ServerPort string
    Env        string // "development" | "production"
}

func Load() *Config {
    config := &Config{
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnv("DB_PORT", "5432"),
        DBUser:     getEnv("DB_USER", "postgres"),
        DBPassword: getEnv("DB_PASSWORD", ""),
        DBName:     getEnv("DB_NAME", "defect_management"),
        JWTSecret:  getEnv("JWT_SECRET", "development-secret-key"),
        ServerPort: getEnv("SERVER_PORT", "8080"),
        Env:        getEnv("ENV", "development"),
    }

		if err := config.validate(); err != nil {
        log.Fatal("Config validation failed:", err)
    }
    
    return config
}

func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}

func (c *Config) validate() error {
    if c.DBPassword == "" && c.Env == "production" {
        return fmt.Errorf("DB_PASSWORD is required in production")
    }
    if c.JWTSecret == "development-secret-key" && c.Env == "production" {
        return fmt.Errorf("JWT_SECRET must be changed in production")
    }
    return nil
}