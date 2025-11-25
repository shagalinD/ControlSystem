package handlers

import (
	"content-service/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CommentHandler struct {
    Handler
}

func NewCommentHandler(db *gorm.DB, jwtSecret, authServiceURL, projectDefectServiceURL string) *CommentHandler {
    return &CommentHandler{
        Handler: *NewHandler(db, jwtSecret, authServiceURL, projectDefectServiceURL),
    }
}

func (h *CommentHandler) GetComments(c *gin.Context) {
    defectID := c.Param("defect_id")
    
    var comments []models.Comment
    
    query := h.DB.
        Where("defect_id = ?", defectID).
        Order("created_at DESC")
    
    // Пагинация
    page, pageSize := h.getPaginationParams(c)
    offset := (page - 1) * pageSize
    
    var total int64
    query.Model(&models.Comment{}).Count(&total)
    
    if err := query.Offset(offset).Limit(pageSize).Find(&comments).Error; err != nil {
        h.internalError(c, "Failed to fetch comments")
        return
    }
    
    h.success(c, gin.H{
        "comments": comments,
        "pagination": gin.H{
            "page":       page,
            "page_size":  pageSize,
            "total":      total,
            "total_pages": (int(total) + pageSize - 1) / pageSize,
        },
    }, "Comments retrieved successfully")
}

func (h *CommentHandler) CreateComment(c *gin.Context) {
    defectIDStr := c.Param("defect_id")
    
 	// Конвертируем defect_id из string в uint
    defectID, err := strconv.ParseUint(defectIDStr, 10, 32)
    if err != nil {
        h.badRequest(c, "Invalid defect ID")
        return
    }

    var req models.CommentCreateRequest
    if !h.validateRequest(c, &req) {
        return
    }
    
    // Проверяем, что defect_id из пути совпадает с defect_id в теле запроса
    if req.DefectID != uint(defectID) {
        h.badRequest(c, "Defect ID in path and body must match")
        return
    }
    
    userID, _, err := h.GetUserFromContext(c)
    if err != nil {
        h.unauthorized(c, "User not authenticated")
        return
    }
    
    // TODO: Проверить существование дефекта через project-defect-service
    
    comment := models.Comment{
        Text:     req.Text,
        DefectID: req.DefectID,
        AuthorID: userID,
    }
    
    if err := h.DB.Create(&comment).Error; err != nil {
        h.internalError(c, "Failed to create comment")
        return
    }
    
    h.success(c, gin.H{
        "comment": comment,
    }, "Comment created successfully")
}

func (h *CommentHandler) UpdateComment(c *gin.Context) {
    commentID := c.Param("id")
    
    var comment models.Comment
    if err := h.DB.First(&comment, commentID).Error; err != nil {
        h.notFound(c, "Comment not found")
        return
    }
    
    userID, _, err := h.GetUserFromContext(c)
    if err != nil {
        h.unauthorized(c, "User not authenticated")
        return
    }
    
    // Проверяем, что пользователь является автором комментария
    if comment.AuthorID != userID {
        h.error(c, http.StatusForbidden, "You can only edit your own comments")
        return
    }
    
    var req struct {
        Text string `json:"text" binding:"required"`
    }
    
    if !h.validateRequest(c, &req) {
        return
    }
    
    comment.Text = req.Text
    
    if err := h.DB.Save(&comment).Error; err != nil {
        h.internalError(c, "Failed to update comment")
        return
    }
    
    h.success(c, gin.H{
        "comment": comment,
    }, "Comment updated successfully")
}

func (h *CommentHandler) DeleteComment(c *gin.Context) {
    commentID := c.Param("id")
    
    var comment models.Comment
    if err := h.DB.First(&comment, commentID).Error; err != nil {
        h.notFound(c, "Comment not found")
        return
    }
    
    userID, userRole, err := h.GetUserFromContext(c)
    if err != nil {
        h.unauthorized(c, "User not authenticated")
        return
    }
    
    // Проверяем, что пользователь является автором комментария или менеджером
    if comment.AuthorID != userID && userRole != "manager" {
        h.error(c, http.StatusForbidden, "You can only delete your own comments")
        return
    }
    
    if err := h.DB.Delete(&comment).Error; err != nil {
        h.internalError(c, "Failed to delete comment")
        return
    }
    
    h.success(c, nil, "Comment deleted successfully")
}