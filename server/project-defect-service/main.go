package main

import (
	"log"

	"project-defect-service/config"
	"project-defect-service/database"
	"project-defect-service/handlers"
	"project-defect-service/middleware"

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
    
    projectHandler := handlers.NewProjectHandler(db, cfg.JWTSecret, cfg.AuthServiceURL)
    defectHandler := handlers.NewDefectHandler(db, cfg.JWTSecret, cfg.AuthServiceURL)
    
    // Protected routes
    api := r.Group("/api")
    api.Use(middleware.JWTMiddleware(cfg.JWTSecret))
    {
        // Проекты
        projects := api.Group("/projects")
        {
            projects.GET("", projectHandler.GetProjects)
            projects.GET("/:id", projectHandler.GetProject)
            projects.POST("", projectHandler.CreateProject)
            projects.PUT("/:id", projectHandler.UpdateProject)
            projects.DELETE("/:id", projectHandler.DeleteProject)
        }
        
        // Дефекты
        defects := api.Group("/defects")
        {
            defects.GET("", defectHandler.GetDefects)
            defects.GET("/my", defectHandler.GetMyDefects)
            defects.GET("/:id", defectHandler.GetDefect)
            defects.POST("", defectHandler.CreateDefect)
            defects.PUT("/:id", defectHandler.UpdateDefect)
            defects.PATCH("/:id/status", defectHandler.UpdateDefectStatus)
            defects.DELETE("/:id", defectHandler.DeleteDefect)
        }
    }
    
    // Health check
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status":  "healthy",
            "service": "project-defect-service",
        })
    })
    
    log.Printf("Project-Defect Service starting on port %s", cfg.ServicePort)
    if err := r.Run(":" + cfg.ServicePort); err != nil {
        log.Fatal("Failed to start project-defect service:", err)
    }
}