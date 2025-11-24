package handlers

import (
	"fmt"
	"net/http"
	"project-defect-service/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DefectHandler struct {
    Handler
}

func NewDefectHandler(db *gorm.DB, jwtSecret, authServiceURL string) *DefectHandler {
    return &DefectHandler{
        Handler: *NewHandler(db, jwtSecret, authServiceURL),
    }
}

func (h *DefectHandler) GetDefects(c *gin.Context) {
    var defects []models.Defect
    
    query := h.DB
    
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
    page, pageSize := h.getPaginationParams(c)
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

func (h *DefectHandler) GetDefect(c *gin.Context) {
    defectID := c.Param("id")
    
    var defect models.Defect
    if err := h.DB.
        Preload("History").
        First(&defect, defectID).Error; err != nil {
        h.notFound(c, "Defect not found")
        return
    }
    
    h.success(c, gin.H{
        "defect": defect,
    }, "Defect retrieved successfully")
}

func (h *DefectHandler) CreateDefect(c *gin.Context) {
    var req models.DefectCreateRequest
    if !h.validateRequest(c, &req) {
        return
    }
    
    userID, _, err := h.GetUserFromContext(c)
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
    
    defect := models.Defect{
        Title:       req.Title,
        Description: req.Description,
        Status:      models.StatusNew,
        Priority:    req.Priority,
        Deadline:    req.Deadline,
        ProjectID:   req.ProjectID,
        AuthorID:    userID,
        AssigneeID:  req.AssigneeID,
    }
    
    if err := h.DB.Create(&defect).Error; err != nil {
        h.internalError(c, "Failed to create defect")
        return
    }
    
    h.success(c, gin.H{
        "defect": defect,
    }, "Defect created successfully")
}

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
    
    userID, _, err := h.GetUserFromContext(c)
    if err != nil {
        h.unauthorized(c, "User not authenticated")
        return
    }
    
    // Логируем изменения
    if req.Title != nil && *req.Title != defect.Title {
        h.logDefectChange(defect.ID, userID, "title", defect.Title, *req.Title)
        defect.Title = *req.Title
    }
    if req.Description != nil && *req.Description != defect.Description {
        h.logDefectChange(defect.ID, userID, "description", defect.Description, *req.Description)
        defect.Description = *req.Description
    }
    if req.Status != nil && *req.Status != defect.Status {
        h.logDefectChange(defect.ID, userID, "status", string(defect.Status), string(*req.Status))
        defect.Status = *req.Status
    }
    if req.Priority != nil && *req.Priority != defect.Priority {
        h.logDefectChange(defect.ID, userID, "priority", string(defect.Priority), string(*req.Priority))
        defect.Priority = *req.Priority
    }
    if req.Deadline != nil && *req.Deadline != *defect.Deadline {
        oldDeadline := "none"
        if defect.Deadline != nil && !defect.Deadline.IsZero() {
            oldDeadline = defect.Deadline.Format("2006-01-02")
        }
        newDeadline := "none"
        if !req.Deadline.IsZero() {
            newDeadline = req.Deadline.Format("2006-01-02")
        }
        h.logDefectChange(defect.ID, userID, "deadline", oldDeadline, newDeadline)
        defect.Deadline = req.Deadline
    }
    if req.AssigneeID != nil && *req.AssigneeID != *defect.AssigneeID {
        oldAssignee := "none"
        if defect.AssigneeID != nil {
            oldAssignee = fmt.Sprintf("%d", *defect.AssigneeID)
        }
        newAssignee := "none"
        if *req.AssigneeID != 0 {
            newAssignee = fmt.Sprintf("%d", *req.AssigneeID)
        }
        h.logDefectChange(defect.ID, userID, "assignee", oldAssignee, newAssignee)
        defect.AssigneeID = req.AssigneeID
    }
    
    if err := h.DB.Save(&defect).Error; err != nil {
        h.internalError(c, "Failed to update defect")
        return
    }
    
    h.DB.Preload("History").First(&defect, defect.ID)
    
    h.success(c, gin.H{
        "defect": defect,
    }, "Defect updated successfully")
}

func (h *DefectHandler) UpdateDefectStatus(c *gin.Context) {
    defectID := c.Param("id")
    
    var defect models.Defect
    if err := h.DB.First(&defect, defectID).Error; err != nil {
        h.notFound(c, "Defect not found")
        return
    }
    
    var req struct {
        Status models.DefectStatus `json:"status" binding:"required,oneof=new in_progress on_review closed cancelled"`
    }
    
    if !h.validateRequest(c, &req) {
        return
    }
    
    userID, _, err := h.GetUserFromContext(c)
    if err != nil {
        h.unauthorized(c, "User not authenticated")
        return
    }
    
    // Логируем изменение статуса
    h.logDefectChange(defect.ID, userID, "status", string(defect.Status), string(req.Status))
    defect.Status = req.Status
    
    if err := h.DB.Save(&defect).Error; err != nil {
        h.internalError(c, "Failed to update defect status")
        return
    }
    
    h.success(c, gin.H{
        "defect": defect,
    }, "Defect status updated successfully")
}

func (h *DefectHandler) GetMyDefects(c *gin.Context) {
    userID, _, err := h.GetUserFromContext(c)
    if err != nil {
        h.unauthorized(c, "User not authenticated")
        return
    }
    
    var defects []models.Defect
    
    query := h.DB.
        Where("author_id = ? OR assignee_id = ?", userID, userID)
    
    // Фильтрация
    if status := c.Query("status"); status != "" {
        query = query.Where("status = ?", status)
    }
    
    // Пагинация
    page, pageSize := h.getPaginationParams(c)
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

func (h *DefectHandler) DeleteDefect(c *gin.Context) {
    defectID := c.Param("id")
    
    var defect models.Defect
    if err := h.DB.First(&defect, defectID).Error; err != nil {
        h.notFound(c, "Defect not found")
        return
    }
    
    userID, userRole, err := h.GetUserFromContext(c)
    if err != nil {
        h.unauthorized(c, "User not authenticated")
        return
    }
    
    // Проверяем права - только автор или менеджер могут удалять
    if defect.AuthorID != userID && userRole != "manager" {
        h.error(c, http.StatusForbidden, "You can only delete your own defects")
        return
    }
    
    // Удаляем связанную историю
    if err := h.DB.Where("defect_id = ?", defectID).Delete(&models.DefectHistory{}).Error; err != nil {
        h.internalError(c, "Failed to delete defect history")
        return
    }
    
    if err := h.DB.Delete(&defect).Error; err != nil {
        h.internalError(c, "Failed to delete defect")
        return
    }
    
    h.success(c, nil, "Defect deleted successfully")
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