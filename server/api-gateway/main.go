package main

import (
	"log"

	"api-gateway/config"
	"api-gateway/handlers"
	"api-gateway/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
    cfg := config.Load()
    
    r := gin.Default()
    
     // CORS middleware
    r.Use(func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
        c.Header("Access-Control-Allow-Credentials", "true")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        
        c.Next()
    })

    // Rate limiting configuration
    rateLimitConfig := middleware.RateLimitConfig{
        GlobalRate:   cfg.RateLimitGlobal,
        AuthRate:     cfg.RateLimitAuth,
        APIRate:      cfg.RateLimitAPI,
        UploadRate:   cfg.RateLimitUpload,
        UseRedis:     cfg.UseRedisRateLimit,
        RedisURL:     cfg.RedisURL,
        TrustProxy:   cfg.TrustProxy,
    }
    
    // Whitelist middleware (если есть whitelist IP)
    if len(cfg.WhitelistIPs) > 0 {
        r.Use(middleware.WhitelistMiddleware(cfg.WhitelistIPs))
    }

     // Rate limiting middleware
    r.Use(middleware.RateLimitMiddleware(rateLimitConfig))
    

    // JWT middleware для защищенных маршрутов
    r.Use(middleware.JWTMiddleware(cfg.JWTSecret))
    
    proxyHandler := handlers.NewProxyHandler(
        cfg.AuthServiceURL,
        cfg.ProjectDefectServiceURL,
        cfg.ContentServiceURL,
    )
    
    // Маршрутизация запросов
    // Аутентификация
    auth := r.Group("/auth")
    {
        auth.POST("/register", proxyHandler.AuthProxy())
        auth.POST("/login", proxyHandler.AuthProxy())
    }
    
    // API маршруты
    api := r.Group("/api")
    {
        // Пользователи
        api.GET("/me", proxyHandler.AuthProxy())
        
        // Проекты и дефекты
        projects := api.Group("/projects")
        {
            projects.GET("", proxyHandler.ProjectDefectProxy())
            projects.GET("/:id", proxyHandler.ProjectDefectProxy())
            projects.POST("", proxyHandler.ProjectDefectProxy())
            projects.PUT("/:id", proxyHandler.ProjectDefectProxy())
            projects.DELETE("/:id", proxyHandler.ProjectDefectProxy())
        }
        
        defects := api.Group("/defects")
        {
            defects.GET("", proxyHandler.ProjectDefectProxy())
            defects.GET("/my", proxyHandler.ProjectDefectProxy())
            defects.GET("/:id", proxyHandler.ProjectDefectProxy())
            defects.POST("", proxyHandler.ProjectDefectProxy())
            defects.PUT("/:id", proxyHandler.ProjectDefectProxy())
            defects.PATCH("/:id/status", proxyHandler.ProjectDefectProxy())
            defects.DELETE("/:id", proxyHandler.ProjectDefectProxy())
        }
        
        // Комментарии
        comments := api.Group("/comments")
        {
            comments.GET("/defect/:defect_id", proxyHandler.ContentProxy())
            comments.POST("/defect/:defect_id", proxyHandler.ContentProxy())
            comments.PUT("/:id", proxyHandler.ContentProxy())
            comments.DELETE("/:id", proxyHandler.ContentProxy())
        }
        
        // Вложения
        attachments := api.Group("/attachments")
        {
            attachments.POST("/defect/:defect_id", proxyHandler.ContentProxy())
            attachments.GET("/defect/:defect_id", proxyHandler.ContentProxy())
            attachments.GET("/:id/download", proxyHandler.ContentProxy())
            attachments.DELETE("/:id", proxyHandler.ContentProxy())
        }
        
        // Отчеты
        reports := api.Group("/reports")
        {
            reports.GET("/defects", proxyHandler.ContentProxy())
            reports.GET("/project/:project_id", proxyHandler.ContentProxy())
            reports.GET("/defects/export", proxyHandler.ContentProxy())
            reports.GET("/user-activity", proxyHandler.ContentProxy())
        }

        users := api.Group("/users")
        {
            users.GET("/engineers", proxyHandler.AuthProxy())
            users.GET("/managers", proxyHandler.AuthProxy())
            users.GET("", proxyHandler.AuthProxy())
            users.GET("/:id", proxyHandler.AuthProxy())
            users.PUT("/:id", proxyHandler.AuthProxy())
        }
    }
    
    // Health check gateway
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status":  "healthy",
            "service": "api-gateway",
        })
    })
    
    // Health check агрегированный
    r.GET("/health/all", func(c *gin.Context) {
        // TODO: Добавить проверку здоровья всех сервисов
        c.JSON(200, gin.H{
            "status": "healthy",
            "services": gin.H{
                "api-gateway": "healthy",
                "auth-service": "unknown",
                "project-defect-service": "unknown", 
                "content-service": "unknown",
            },
        })
    })
    
    log.Printf("API Gateway starting on port %s", cfg.GatewayPort)
    log.Printf("Rate limiting enabled: Global=%s, Auth=%s, API=%s, Upload=%s", 
        cfg.RateLimitGlobal, cfg.RateLimitAuth, cfg.RateLimitAPI, cfg.RateLimitUpload)
        
    if err := r.Run(":" + cfg.GatewayPort); err != nil {
        log.Fatal("Failed to start API Gateway:", err)
    }
    
}