package main

import (
	"log"

	"auth-service/config"
	"auth-service/database"
	"auth-service/handlers"
	"auth-service/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
    cfg := config.Load()
    
    db, err := database.NewConnection(&database.Config{
        Host:     cfg.DBHost,
        Port:     cfg.DBPort,
        User:     cfg.DBUser,
        Password: cfg.DBPassword,
        DBName:   cfg.DBName,
    })
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    
    r := gin.Default()
    
    // CORS middleware
    r.Use(func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        
        c.Next()
    })
    
    authHandler := handlers.NewAuthHandler(db, cfg.JWTSecret)
    
    // Public routes
    authGroup := r.Group("/auth")
    {
        authGroup.POST("/register", authHandler.Register)
        authGroup.POST("/login", authHandler.Login)
    }
    
    // Protected routes
    api := r.Group("/api")
    api.Use(middleware.JWTMiddleware(cfg.JWTSecret))
    {
        api.GET("/me", authHandler.GetCurrentUser)
    }
    
    // Health check
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status":  "healthy",
            "service": "auth-service",
        })
    })
    
    log.Printf("Auth Service starting on port %s", cfg.ServicePort)
    if err := r.Run(":" + cfg.ServicePort); err != nil {
        log.Fatal("Failed to start auth service:", err)
    }
}