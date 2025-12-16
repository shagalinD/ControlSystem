package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
    GatewayPort         string
    AuthServiceURL      string
    ProjectDefectServiceURL string
    ContentServiceURL   string
    JWTSecret           string
    
    // Rate limiting configuration
    RateLimitGlobal     string
    RateLimitAuth       string
    RateLimitAPI        string
    RateLimitUpload     string
    UseRedisRateLimit   bool
    RedisURL            string
    TrustProxy          bool
    WhitelistIPs        []string
}

func Load() *Config {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using system environment variables")
    }
    
    config := &Config{
        GatewayPort:         getEnv("GATEWAY_PORT", "8080"),
        AuthServiceURL:      getEnv("AUTH_SERVICE_URL", "http://auth-service:8081"),
        ProjectDefectServiceURL: getEnv("PROJECT_DEFECT_SERVICE_URL", "http://project-defect-service:8082"),
        ContentServiceURL:   getEnv("CONTENT_SERVICE_URL", "http://content-service:8083"),
        JWTSecret:           getEnv("JWT_SECRET", "development-secret-key"),
        
        // Rate limiting
        RateLimitGlobal:     getEnv("RATE_LIMIT_GLOBAL", "1000-H"),
        RateLimitAuth:       getEnv("RATE_LIMIT_AUTH", "10-M"),
        RateLimitAPI:        getEnv("RATE_LIMIT_API", "100-M"),
        RateLimitUpload:     getEnv("RATE_LIMIT_UPLOAD", "20-H"),
        UseRedisRateLimit:   getEnv("USE_REDIS_RATELIMIT", "false") == "true",
        RedisURL:            getEnv("REDIS_URL", "redis://redis:6379"),
        TrustProxy:          getEnv("TRUST_PROXY", "false") == "true",
        WhitelistIPs:        parseWhitelist(getEnv("WHITELIST_IPS", "")),
    }
    
    return config
}

func parseWhitelist(whitelistStr string) []string {
    if whitelistStr == "" {
        return []string{}
    }
    
    var ips []string
    for _, ip := range strings.Split(whitelistStr, ",") {
        trimmed := strings.TrimSpace(ip)
        if trimmed != "" {
            ips = append(ips, trimmed)
        }
    }
    return ips
}

func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}