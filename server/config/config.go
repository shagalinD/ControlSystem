package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
    DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    DBName     string
    JWTSecret  string
    ServerPort string
    Env        string
}

func Load() *Config {
    // Загружаем .env файл
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using system environment variables")
    }
    
    config := &Config{
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnv("DB_PORT", "5432"),
        DBUser:     getEnv("DB_USER", "postgres"),
        DBPassword: getEnv("DB_PASSWORD", ""),
        DBName:     getEnv("DB_NAME", "kopatel_online"),
        JWTSecret:  getEnv("JWT_SECRET", "development-secret-key"),
        ServerPort: getEnv("SERVER_PORT", "8080"),
        Env:        getEnv("ENV", "development"),
    }
    
    // Логируем конфигурацию (без пароля)
    log.Printf("Database config: host=%s, port=%s, user=%s, dbname=%s", 
        config.DBHost, config.DBPort, config.DBUser, config.DBName)
    
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