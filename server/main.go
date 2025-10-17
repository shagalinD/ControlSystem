package main

import (
	"kopatel_online/config"
	"kopatel_online/handlers"
	"kopatel_online/postgres"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
    // Загрузка конфигурации
    cfg := config.Load()
    
    // Подключение к БД
    db, err := postgres.NewConnection(&postgres.Config{
        Host:     cfg.DBHost,
        Port:     cfg.DBPort,
        User:     cfg.DBUser,
        Password: cfg.DBPassword,
        DBName:   cfg.DBName,
    })
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    
    // Инициализация Gin
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
    
    // Инициализация обработчиков
    authHandler := handlers.NewAuthHandler(db, cfg.JWTSecret)
    projectHandler := handlers.NewProjectHandler(db)
    defectHandler := handlers.NewDefectHandler(db)
    commentHandler := handlers.NewCommentHandler(db)
    attachmentHandler := handlers.NewAttachmentHandler(db)
    userHandler := handlers.NewUserHandler(db)
    // Группа аутентификации
    authGroup := r.Group("/auth")
    {
        authGroup.POST("/register", authHandler.Register)
        authGroup.POST("/login", authHandler.Login)
    }
    
    // Защищенные маршруты
    api := r.Group("/api")
    api.Use(authHandler.AuthMiddleware())
    {
         users := api.Group("/users")
        {
            users.GET("/engineers", userHandler.GetEngineers)
            users.GET("/managers", userHandler.GetManagers)
            users.GET("", userHandler.GetAllUsers)
        }
        // Текущий пользователь
        api.GET("/me", authHandler.GetCurrentUser)
        
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

        attachments := api.Group("/attachments")
        {
            attachments.POST("/defect/:defect_id", attachmentHandler.UploadAttachment)
            attachments.GET("/defect/:defect_id", attachmentHandler.GetAttachments)
            attachments.GET("/:id/download", attachmentHandler.DownloadAttachment)
            attachments.DELETE("/:id", attachmentHandler.DeleteAttachment)
        }
        
        // Комментарии
        comments := api.Group("/comments")
        {
            comments.GET("/defect/:defect_id", commentHandler.GetComments)
            comments.POST("/defect/:defect_id", commentHandler.CreateComment)
            comments.PUT("/:id", commentHandler.UpdateComment)
            comments.DELETE("/:id", commentHandler.DeleteComment)
        }
    }
    
    // Запуск сервера
    log.Printf("Server starting on port %s", cfg.ServerPort)
    if err := r.Run(":" + cfg.ServerPort); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}