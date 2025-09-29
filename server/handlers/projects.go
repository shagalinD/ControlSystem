package handlers

import (
	"kopatel_online/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetProject - получение проекта по ID
type ProjectHandler struct {
    Handler
}

func NewProjectHandler(db *gorm.DB) *ProjectHandler {
    return &ProjectHandler{Handler: Handler{DB: db}}
}

// GetProjects - получение списка проектов
func (h *ProjectHandler) GetProjects(c *gin.Context) {
    var projects []models.Project
    
    query := h.DB.Preload("Manager")
    
    // Пагинация
    page := 1
    pageSize := 20
    
    if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&projects).Error; err != nil {
        h.internalError(c, "Failed to fetch projects")
        return
    }
    
    h.success(c, gin.H{
        "projects": projects,
        "page":     page,
        "pageSize": pageSize,
    }, "Projects retrieved successfully")
}

// GetProject - получение проекта по ID
func (h *ProjectHandler) GetProject(c *gin.Context) {
    projectID := c.Param("id")
    
    var project models.Project
    if err := h.DB.Preload("Manager").First(&project, projectID).Error; err != nil {
        h.notFound(c, "Project not found")
        return
    }
    
    h.success(c, gin.H{
        "project": project,
    }, "Project retrieved successfully")
}

// CreateProject - создание нового проекта
func (h *ProjectHandler) CreateProject(c *gin.Context) {
    var req models.ProjectCreateRequest
    if !h.validateRequest(c, &req) {
        return
    }
    
    // Получаем текущего пользователя из контекста (теперь метод доступен)
    user, err := h.GetUserFromContext(c)
    if err != nil {
        h.unauthorized(c, "User not authenticated")
        return
    }
    
    // Логируем кто создает проект (для отладки)
    println("Project created by user:", user.Email)
    
    // Проверяем, что менеджер существует
    var manager models.User
    if err := h.DB.First(&manager, req.ManagerID).Error; err != nil {
        h.badRequest(c, "Manager not found")
        return
    }
    
    project := models.Project{
        Name:        req.Name,
        Description: req.Description,
        ManagerID:   req.ManagerID,
    }
    
    if err := h.DB.Create(&project).Error; err != nil {
        h.internalError(c, "Failed to create project")
        return
    }
    
    // Загружаем связанные данные
    h.DB.Preload("Manager").First(&project, project.ID)
    
    h.success(c, gin.H{
        "project": project,
    }, "Project created successfully")
}

// UpdateProject - обновление проекта
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
    projectID := c.Param("id")
    
    var project models.Project
    if err := h.DB.First(&project, projectID).Error; err != nil {
        h.notFound(c, "Project not found")
        return
    }
    
    var req models.ProjectCreateRequest
    if !h.validateRequest(c, &req) {
        return
    }
    
    // Обновляем поля
    project.Name = req.Name
    project.Description = req.Description
    project.ManagerID = req.ManagerID
    
    if err := h.DB.Save(&project).Error; err != nil {
        h.internalError(c, "Failed to update project")
        return
    }
    
    h.DB.Preload("Manager").First(&project, project.ID)
    
    h.success(c, gin.H{
        "project": project,
    }, "Project updated successfully")
}

// DeleteProject - удаление проекта
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
    projectID := c.Param("id")
    
    var project models.Project
    if err := h.DB.First(&project, projectID).Error; err != nil {
        h.notFound(c, "Project not found")
        return
    }
    
    // Проверяем, нет ли связанных дефектов
    var defectsCount int64
    h.DB.Model(&models.Defect{}).Where("project_id = ?", projectID).Count(&defectsCount)
    if defectsCount > 0 {
        h.badRequest(c, "Cannot delete project with existing defects")
        return
    }
    
    if err := h.DB.Delete(&project).Error; err != nil {
        h.internalError(c, "Failed to delete project")
        return
    }
    
    h.success(c, nil, "Project deleted successfully")
}