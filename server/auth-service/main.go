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
    
    authHandler := handlers.NewAuthHandler(db, cfg.JWTSecret)
    userHandler := handlers.NewUserHandler(db, cfg.JWTSecret)
    
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
        // Текущий пользователь
        api.GET("/me", authHandler.GetCurrentUser)
        
        // Пользователи
        users := api.Group("/users")
        {
            users.GET("/engineers", userHandler.GetEngineers)
            users.GET("/managers", userHandler.GetManagers)
            users.GET("", userHandler.GetAllUsers)
            users.GET("/:id", userHandler.GetUserByID)
            users.PUT("/:id", userHandler.UpdateUserData)
        }
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