package handlers

import (
	"net/http"
	"project-defect-service/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectHandler struct {
    Handler
}

func NewProjectHandler(db *gorm.DB, jwtSecret, authServiceURL string) *ProjectHandler {
    return &ProjectHandler{
        Handler: *NewHandler(db, jwtSecret, authServiceURL),
    }
}

func (h *ProjectHandler) GetProjects(c *gin.Context) {
    var projects []models.Project
    
    query := h.DB
    page, pageSize := h.getPaginationParams(c)
    offset := (page - 1) * pageSize
    
    var total int64
    query.Model(&models.Project{}).Count(&total)
    
    if err := query.Offset(offset).Limit(pageSize).Find(&projects).Error; err != nil {
        h.internalError(c, "Failed to fetch projects")
        return
    }
    
    h.success(c, gin.H{
        "projects": projects,
        "pagination": gin.H{
            "page":       page,
            "page_size":  pageSize,
            "total":      total,
            "total_pages": (int(total) + pageSize - 1) / pageSize,
        },
    }, "Projects retrieved successfully")
}

func (h *ProjectHandler) GetProject(c *gin.Context) {
    projectID := c.Param("id")
    
    var project models.Project
    if err := h.DB.First(&project, projectID).Error; err != nil {
        h.notFound(c, "Project not found")
        return
    }
    
    h.success(c, gin.H{
        "project": project,
    }, "Project retrieved successfully")
}

func (h *ProjectHandler) CreateProject(c *gin.Context) {
    var req models.ProjectCreateRequest
    if !h.validateRequest(c, &req) {
        return
    }
    
    userID, userRole, err := h.GetUserFromContext(c)
    if err != nil {
        h.unauthorized(c, "User not authenticated")
        return
    }
    
    // Проверяем права - только менеджеры могут создавать проекты
    if userRole != "manager" {
        h.error(c, http.StatusForbidden, "Only managers can create projects")
        return
    }
    
    project := models.Project{
        Name:        req.Name,
        Description: req.Description,
        ManagerID:   userID, // Менеджер - текущий пользователь
    }
    
    if err := h.DB.Create(&project).Error; err != nil {
        h.internalError(c, "Failed to create project")
        return
    }
    
    h.success(c, gin.H{
        "project": project,
    }, "Project created successfully")
}

func (h *ProjectHandler) UpdateProject(c *gin.Context) {
    projectID := c.Param("id")
    
    var project models.Project
    if err := h.DB.First(&project, projectID).Error; err != nil {
        h.notFound(c, "Project not found")
        return
    }
    
    userID, userRole, err := h.GetUserFromContext(c)
    if err != nil {
        h.unauthorized(c, "User not authenticated")
        return
    }
    
    // Проверяем права - только автор проекта (менеджер) может редактировать
    if project.ManagerID != userID && userRole != "manager" {
        h.error(c, http.StatusForbidden, "You can only edit your own projects")
        return
    }
    
    var req models.ProjectCreateRequest
    if !h.validateRequest(c, &req) {
        return
    }
    
    project.Name = req.Name
    project.Description = req.Description
    
    if err := h.DB.Save(&project).Error; err != nil {
        h.internalError(c, "Failed to update project")
        return
    }
    
    h.success(c, gin.H{
        "project": project,
    }, "Project updated successfully")
}

func (h *ProjectHandler) DeleteProject(c *gin.Context) {
    projectID := c.Param("id")
    
    var project models.Project
    if err := h.DB.First(&project, projectID).Error; err != nil {
        h.notFound(c, "Project not found")
        return
    }
    
    userID, userRole, err := h.GetUserFromContext(c)
    if err != nil {
        h.unauthorized(c, "User not authenticated")
        return
    }
    
    // Проверяем права - только автор проекта (менеджер) может удалять
    if project.ManagerID != userID && userRole != "manager" {
        h.error(c, http.StatusForbidden, "You can only delete your own projects")
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