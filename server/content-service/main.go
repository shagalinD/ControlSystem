package main

import (
	"log"

	"content-service/config"
	"content-service/database"
	"content-service/handlers"
	"content-service/middleware"

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
    
    
    commentHandler := handlers.NewCommentHandler(db, cfg.JWTSecret, cfg.AuthServiceURL, cfg.ProjectDefectServiceURL)
    attachmentHandler := handlers.NewAttachmentHandler(db, cfg.JWTSecret, cfg.AuthServiceURL, cfg.ProjectDefectServiceURL, cfg.UploadPath)
    reportHandler := handlers.NewReportHandler(db, cfg.JWTSecret, cfg.AuthServiceURL, cfg.ProjectDefectServiceURL)
    
    // Protected routes
    api := r.Group("/api")
    api.Use(middleware.JWTMiddleware(cfg.JWTSecret))
    {
        // Комментарии
        comments := api.Group("/comments")
        {
            comments.GET("/defect/:defect_id", commentHandler.GetComments)
            comments.POST("/defect/:defect_id", commentHandler.CreateComment)
            comments.PUT("/:id", commentHandler.UpdateComment)
            comments.DELETE("/:id", commentHandler.DeleteComment)
        }
        
        // Вложения
        attachments := api.Group("/attachments")
        {
            attachments.POST("/defect/:defect_id", attachmentHandler.UploadAttachment)
            attachments.GET("/defect/:defect_id", attachmentHandler.GetAttachments)
            attachments.GET("/:id/download", attachmentHandler.DownloadAttachment)
            attachments.DELETE("/:id", attachmentHandler.DeleteAttachment)
        }
        
        // Отчеты
        reports := api.Group("/reports")
        {
            reports.GET("/defects", reportHandler.GetDefectsReport)
            reports.GET("/project/:project_id", reportHandler.GetProjectReport)
            reports.GET("/defects/export", reportHandler.ExportDefectsCSV)
            reports.GET("/user-activity", reportHandler.GetUserActivityReport)
        }
    }
    
    // Health check
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status":  "healthy",
            "service": "content-service",
        })
    })
    
    log.Printf("Content Service starting on port %s", cfg.ServicePort)
    if err := r.Run(":" + cfg.ServicePort); err != nil {
        log.Fatal("Failed to start content service:", err)
    }
}