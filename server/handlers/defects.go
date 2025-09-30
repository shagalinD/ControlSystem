package handlers

import (
	"fmt"
	"kopatel_online/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DefectHandler struct {
    Handler
}

func NewDefectHandler(db *gorm.DB) *DefectHandler {
    return &DefectHandler{Handler: Handler{DB: db}}
}

// GetDefects - получение списка дефектов с фильтрацией
func (h *DefectHandler) GetDefects(c *gin.Context) {
    var defects []models.Defect
    
    query := h.DB.Preload("Project").Preload("Author").Preload("Assignee")
    
    // Фильтрация по проекту
    if projectID := c.Query("project_id"); projectID != "" {
        query = query.Where("project_id = ?", projectID)
    }
    
    // Фильтрация по статусу
    if status := c.Query("status"); status != "" {
        query = query.Where("status = ?", status)
    }
    
    // Фильтрация по приоритету
    if priority := c.Query("priority"); priority != "" {
        query = query.Where("priority = ?", priority)
    }
    
    // Фильтрация по исполнителю
    if assigneeID := c.Query("assignee_id"); assigneeID != "" {
        query = query.Where("assignee_id = ?", assigneeID)
    }
    
    // Сортировка
    sortBy := c.DefaultQuery("sort_by", "created_at")
    order := c.DefaultQuery("order", "desc")
    query = query.Order(sortBy + " " + order)
    
    // Пагинация
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
    offset := (page - 1) * pageSize
    
    var total int64
    query.Model(&models.Defect{}).Count(&total)
    
    if err := query.Offset(offset).Limit(pageSize).Find(&defects).Error; err != nil {
        h.internalError(c, "Failed to fetch defects")
        return
    }
    
    h.success(c, gin.H{
        "defects": defects,
        "pagination": gin.H{
            "page":       page,
            "page_size":  pageSize,
            "total":      total,
            "total_pages": (int(total) + pageSize - 1) / pageSize,
        },
    }, "Defects retrieved successfully")
}

// GetDefect - получение дефекта по ID
func (h *DefectHandler) GetDefect(c *gin.Context) {
    defectID := c.Param("id")
    
    var defect models.Defect
    if err := h.DB.
        Preload("Project").
        Preload("Author").
        Preload("Assignee").
        Preload("Comments").
        Preload("Comments.Author").
        First(&defect, defectID).Error; err != nil {
        h.notFound(c, "Defect not found")
        return
    }
    
    h.success(c, gin.H{
        "defect": defect,
    }, "Defect retrieved successfully")
}

// CreateDefect - создание нового дефекта
func (h *DefectHandler) CreateDefect(c *gin.Context) {
    var req models.DefectCreateRequest
    if !h.validateRequest(c, &req) {
        return
    }
    
    // Получаем текущего пользователя
    user, err := h.GetUserFromContext(c)
    if err != nil {
        h.unauthorized(c, "User not authenticated")
        return
    }
    
    // Проверяем существование проекта
    var project models.Project
    if err := h.DB.First(&project, req.ProjectID).Error; err != nil {
        h.badRequest(c, "Project not found")
        return
    }
    
    // Проверяем исполнителя, если указан
    if req.AssigneeID != nil {
        var assignee models.User
        if err := h.DB.First(&assignee, *req.AssigneeID).Error; err != nil {
            h.badRequest(c, "Assignee not found")
            return
        }
    }
    
    defect := models.Defect{
        Title:       req.Title,
        Description: req.Description,
        Status:      models.StatusNew,
        Priority:    req.Priority,
        Deadline:    req.Deadline,
        ProjectID:   req.ProjectID,
        AuthorID:    user.ID,
        AssigneeID:  req.AssigneeID,
    }
    
    if err := h.DB.Create(&defect).Error; err != nil {
        h.internalError(c, "Failed to create defect")
        return
    }
    
    // Загружаем связанные данные
    h.DB.
        Preload("Project").
        Preload("Author").
        Preload("Assignee").
        First(&defect, defect.ID)
    
    h.success(c, gin.H{
        "defect": defect,
    }, "Defect created successfully")
}

// UpdateDefectStatus - обновление только статуса дефекта
func (h *DefectHandler) UpdateDefectStatus(c *gin.Context) {
    defectID := c.Param("id")
    
    var defect models.Defect
    if err := h.DB.First(&defect, defectID).Error; err != nil {
        h.notFound(c, "Defect not found")
        return
    }
    
    var req struct {
        Status models.DefectStatus `json:"status" validate:"required,defect_status"`
    }
    
    if !h.validateRequest(c, &req) {
        return
    }
    
    defect.Status = req.Status
    
    if err := h.DB.Save(&defect).Error; err != nil {
        h.internalError(c, "Failed to update defect status")
        return
    }
    
    h.success(c, gin.H{
        "defect": defect,
    }, "Defect status updated successfully")
}

// DeleteDefect - удаление дефекта
func (h *DefectHandler) DeleteDefect(c *gin.Context) {
    defectID := c.Param("id")
    
    var defect models.Defect
    if err := h.DB.First(&defect, defectID).Error; err != nil {
        h.notFound(c, "Defect not found")
        return
    }
    
    // Удаляем связанные комментарии
    if err := h.DB.Where("defect_id = ?", defectID).Delete(&models.Comment{}).Error; err != nil {
        h.internalError(c, "Failed to delete related comments")
        return
    }
    
    if err := h.DB.Delete(&defect).Error; err != nil {
        h.internalError(c, "Failed to delete defect")
        return
    }
    
    h.success(c, nil, "Defect deleted successfully")
}

// GetMyDefects - получение дефектов текущего пользователя
func (h *DefectHandler) GetMyDefects(c *gin.Context) {
    user, err := h.GetUserFromContext(c)
    if err != nil {
        h.unauthorized(c, "User not authenticated")
        return
    }
    
    var defects []models.Defect
    
    query := h.DB.
        Preload("Project").
        Preload("Author").
        Preload("Assignee").
        Where("author_id = ? OR assignee_id = ?", user.ID, user.ID)
    
    // Фильтрация
    if status := c.Query("status"); status != "" {
        query = query.Where("status = ?", status)
    }
    
    // Пагинация
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
    offset := (page - 1) * pageSize
    
    var total int64
    query.Model(&models.Defect{}).Count(&total)
    
    if err := query.Offset(offset).Limit(pageSize).Find(&defects).Error; err != nil {
        h.internalError(c, "Failed to fetch defects")
        return
    }
    
    h.success(c, gin.H{
        "defects": defects,
        "pagination": gin.H{
            "page":       page,
            "page_size":  pageSize,
            "total":      total,
            "total_pages": (int(total) + pageSize - 1) / pageSize,
        },
    }, "My defects retrieved successfully")
}

func (h *DefectHandler) logDefectChange(defectID uint, userID uint, field string, oldValue, newValue string) {
    history := models.DefectHistory{
        DefectID:  defectID,
        Field:     field,
        OldValue:  oldValue,
        NewValue:  newValue,
        ChangedBy: userID,
    }
    h.DB.Create(&history)
}

// Обновим метод UpdateDefect
func (h *DefectHandler) UpdateDefect(c *gin.Context) {
    defectID := c.Param("id")
    
    var defect models.Defect
    if err := h.DB.First(&defect, defectID).Error; err != nil {
        h.notFound(c, "Defect not found")
        return
    }
    
    var req models.DefectUpdateRequest
    if !h.validateRequest(c, &req) {
        return
    }
    
    user, err := h.GetUserFromContext(c)
    if err != nil {
        h.unauthorized(c, "User not authenticated")
        return
    }
    
    // Логируем изменения
    if req.Title != nil && *req.Title != defect.Title {
        h.logDefectChange(defect.ID, user.ID, "title", defect.Title, *req.Title)
        defect.Title = *req.Title
    }
    if req.Description != nil && *req.Description != defect.Description {
        h.logDefectChange(defect.ID, user.ID, "description", defect.Description, *req.Description)
        defect.Description = *req.Description
    }
    if req.Status != nil && *req.Status != defect.Status {
        h.logDefectChange(defect.ID, user.ID, "status", string(defect.Status), string(*req.Status))
        defect.Status = *req.Status
    }
    if req.Priority != nil && *req.Priority != defect.Priority {
        h.logDefectChange(defect.ID, user.ID, "priority", string(defect.Priority), string(*req.Priority))
        defect.Priority = *req.Priority
    }
    if req.AssigneeID != nil {
        oldAssignee := "none"
        if defect.AssigneeID != nil {
            oldAssignee = fmt.Sprintf("%d", *defect.AssigneeID)
        }
        newAssignee := "none"
        if *req.AssigneeID != 0 {
            newAssignee = fmt.Sprintf("%d", *req.AssigneeID)
        }
        h.logDefectChange(defect.ID, user.ID, "assignee", oldAssignee, newAssignee)
        defect.AssigneeID = req.AssigneeID
    }
    
    if err := h.DB.Save(&defect).Error; err != nil {
        h.internalError(c, "Failed to update defect")
        return
    }
    
    h.DB.
        Preload("Project").
        Preload("Author").
        Preload("Assignee").
        Preload("History").
        Preload("History.Changer").
        First(&defect, defect.ID)
    
    h.success(c, gin.H{
        "defect": defect,
    }, "Defect updated successfully")
}