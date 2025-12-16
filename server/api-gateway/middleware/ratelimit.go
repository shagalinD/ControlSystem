package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	ginlimiter "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"github.com/ulule/limiter/v3/drivers/store/redis"

	goredis "github.com/redis/go-redis/v9" // Для Redis клиента
)

// RateLimitConfig конфигурация rate limiting
type RateLimitConfig struct {
    // Общий лимит для всех endpoints
    GlobalRate   string
    // Лимит для аутентификации
    AuthRate     string
    // Лимит для API endpoints
    APIRate      string
    // Лимит для загрузки файлов
    UploadRate   string
    // Использовать Redis для распределенного лимита
    UseRedis     bool
    RedisURL     string
    // TrustProxy доверять прокси headers
    TrustProxy   bool
}

// DefaultConfig дефолтная конфигурация
func DefaultConfig() RateLimitConfig {
    return RateLimitConfig{
        GlobalRate:   "1000-H",  // 1000 запросов в час
        AuthRate:     "10-M",    // 10 запросов в минуту
        APIRate:      "100-M",   // 100 запросов в минуту
        UploadRate:   "20-H",    // 20 загрузок в час
        UseRedis:     false,
        RedisURL:     "redis://localhost:6379",
        TrustProxy:   false,
    }
}

// GetClientIP извлекает реальный IP клиента
func GetClientIP(c *gin.Context, trustProxy bool) string {
    if trustProxy {
        // Пробуем получить IP из заголовков прокси
        if ip := c.GetHeader("X-Real-IP"); ip != "" {
            return ip
        }
        if ip := c.GetHeader("X-Forwarded-For"); ip != "" {
            // Берем первый IP из цепочки
            parts := strings.Split(ip, ",")
            if len(parts) > 0 {
                return strings.TrimSpace(parts[0])
            }
        }
    }
    
    // Возвращаем IP из запроса
    ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
    if err != nil {
        return c.Request.RemoteAddr
    }
    return ip
}

// parseRedisURL парсит URL Redis для go-redis/v9
func parseRedisURL(url string) (*goredis.Options, error) {
    // Используем встроенный парсер go-redis
    opts, err := goredis.ParseURL(url)
    if err != nil {
        return nil, err
    }
    return opts, nil
}

// RateLimitMiddleware создает middleware для rate limiting
func RateLimitMiddleware(config RateLimitConfig) gin.HandlerFunc {
    // Создаем store для rate limiter
    var store limiter.Store
    
    if config.UseRedis {
        // Используем Redis для распределенного rate limiting
        opts, err := parseRedisURL(config.RedisURL)
        if err != nil {
            // Fallback на memory store
            fmt.Printf("Redis connection error: %v, falling back to memory store\n", err)
            store = memory.NewStore()
        } else {
            // Создаем Redis клиент
            client := goredis.NewClient(opts)
            
            // Проверяем подключение
            ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
            defer cancel()
            
            if _, err := client.Ping(ctx).Result(); err != nil {
                fmt.Printf("Redis ping failed: %v, falling back to memory store\n", err)
                store = memory.NewStore()
            } else {
                // Создаем Redis store для limiter
                redisStore, err := redis.NewStoreWithOptions(client, limiter.StoreOptions{
                    Prefix:   "api_gateway_ratelimit",
                })
                if err != nil {
                    fmt.Printf("Redis store creation error: %v, falling back to memory store\n", err)
                    store = memory.NewStore()
                } else {
                    store = redisStore
                    fmt.Println("Using Redis store for rate limiting")
                }
            }
        }
    } else {
        // In-memory store (для standalone)
        store = memory.NewStore()
    }
    
    // Создаем rate limiter
    rate, err := limiter.NewRateFromFormatted(config.GlobalRate)
    if err != nil {
        // Fallback на разумные лимиты
        rate = limiter.Rate{
            Limit:  1000,
            Period: time.Hour,
        }
    }
    
    // Базовый middleware
    baseMiddleware := ginlimiter.NewMiddleware(
        limiter.New(store, rate),
        ginlimiter.WithKeyGetter(func(c *gin.Context) string {
            // Ключ на основе IP + UserID если есть
            ip := GetClientIP(c, config.TrustProxy)
            
            // Добавляем user_id если пользователь аутентифицирован
            if userID, exists := c.Get("user_id"); exists {
                return fmt.Sprintf("%s:%v", ip, userID)
            }
            
            return ip
        }),
        ginlimiter.WithErrorHandler(func(c *gin.Context, err error) {
            c.JSON(http.StatusInternalServerError, gin.H{
                "success": false,
                "error":   "Internal server error",
            })
        }),
        ginlimiter.WithLimitReachedHandler(func(c *gin.Context) {
            c.JSON(http.StatusTooManyRequests, gin.H{
                "success": false,
                "error":   "Too many requests. Please try again later.",
                "retry_after": c.GetHeader("X-RateLimit-Reset"),
            })
        }),
    )
    
    return func(c *gin.Context) {
        // Проверяем если нужно пропустить rate limiting
        if skip, exists := c.Get("skip_ratelimit"); exists && skip.(bool) {
            c.Next()
            return
        }
        
        // Определяем какая конфигурация лимита использовать
        path := c.Request.URL.Path
        
        var rateStr string
        
        switch {
        case strings.HasPrefix(path, "/auth/"):
            rateStr = config.AuthRate
        case strings.HasPrefix(path, "/api/attachments") && c.Request.Method == "POST":
            rateStr = config.UploadRate
        case strings.HasPrefix(path, "/api/"):
            rateStr = config.APIRate
        default:
            rateStr = config.GlobalRate
        }
        
        // Если лимит отличается от глобального, создаем отдельный лимитер
        if rateStr != config.GlobalRate {
            customRate, err := limiter.NewRateFromFormatted(rateStr)
            if err == nil {
                customLimiter := limiter.New(store, customRate)
                customMiddleware := ginlimiter.NewMiddleware(
                    customLimiter,
                    ginlimiter.WithKeyGetter(func(c *gin.Context) string {
                        ip := GetClientIP(c, config.TrustProxy)
                        // Добавляем суффикс для отдельного бакета
                        suffix := ""
                        
                        switch {
                        case strings.HasPrefix(path, "/auth/"):
                            suffix = ":auth"
                        case strings.HasPrefix(path, "/api/attachments") && c.Request.Method == "POST":
                            suffix = ":upload"
                        case strings.HasPrefix(path, "/api/"):
                            suffix = ":api"
                        }
                        
                        if userID, exists := c.Get("user_id"); exists {
                            return fmt.Sprintf("%s:%v%s", ip, userID, suffix)
                        }
                        return ip + suffix
                    }),
                )
                customMiddleware(c)
                return
            }
        }
        
        // Используем глобальный middleware
        baseMiddleware(c)
    }
}

// WhitelistMiddleware middleware для исключения IP из rate limiting
func WhitelistMiddleware(whitelist []string) gin.HandlerFunc {
    return func(c *gin.Context) {
        clientIP := GetClientIP(c, false)
        
        // Проверяем если IP в whitelist
        for _, ip := range whitelist {
            if ip == clientIP {
                // Пропускаем rate limiting для этого IP
                c.Set("skip_ratelimit", true)
                break
            }
        }
        
        c.Next()
    }
}