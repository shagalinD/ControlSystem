package handlers

import (
	"kopatel_online/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CommentHandler struct {
    Handler
}

func NewCommentHandler(db *gorm.DB) *CommentHandler {
    return &CommentHandler{
        Handler: *NewHandler(db),
    }
}
// GetComments - получение комментариев для дефекта
func (h *CommentHandler) GetComments(c *gin.Context) {
    defectID := c.Param("defect_id")
    
    var comments []models.Comment
    
    query := h.DB.
        Preload("Author").
        Where("defect_id = ?", defectID).
        Order("created_at DESC")
    
    // Пагинация
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))
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

// CreateComment - создание комментария
func (h *CommentHandler) CreateComment(c *gin.Context) {
    defectID := c.Param("defect_id")
    
    var req models.CommentCreateRequest
    if !h.validateRequest(c, &req) {
        return
    }
    
    // Получаем текущего пользователя
    user, err := h.GetUserFromContext(c)
    if err != nil {
        h.unauthorized(c, "User not authenticated")
        return
    }
    
    // Проверяем существование дефекта
    var defect models.Defect
    if err := h.DB.First(&defect, defectID).Error; err != nil {
        h.notFound(c, "Defect not found")
        return
    }
    
    comment := models.Comment{
        Text:     req.Text,
        DefectID: req.DefectID,
        AuthorID: user.ID,
    }
    
    if err := h.DB.Create(&comment).Error; err != nil {
        h.internalError(c, "Failed to create comment")
        return
    }
    
    // Загружаем автора для ответа
    h.DB.Preload("Author").First(&comment, comment.ID)
    
    h.success(c, gin.H{
        "comment": comment,
    }, "Comment created successfully")
}

// UpdateComment - обновление комментария
func (h *CommentHandler) UpdateComment(c *gin.Context) {
    commentID := c.Param("id")
    
    var comment models.Comment
    if err := h.DB.First(&comment, commentID).Error; err != nil {
        h.notFound(c, "Comment not found")
        return
    }
    
    // Получаем текущего пользователя
    user, err := h.GetUserFromContext(c)
    if err != nil {
        h.unauthorized(c, "User not authenticated")
        return
    }
    
    // Проверяем, что пользователь является автором комментария
    if comment.AuthorID != user.ID {
        h.error(c, http.StatusForbidden, "You can only edit your own comments")
        return
    }
    
    var req struct {
        Text string `json:"text" validate:"required"`
    }
    
    if !h.validateRequest(c, &req) {
        return
    }
    
    comment.Text = req.Text
    
    if err := h.DB.Save(&comment).Error; err != nil {
        h.internalError(c, "Failed to update comment")
        return
    }
    
    h.DB.Preload("Author").First(&comment, comment.ID)
    
    h.success(c, gin.H{
        "comment": comment,
    }, "Comment updated successfully")
}

// DeleteComment - удаление комментария
func (h *CommentHandler) DeleteComment(c *gin.Context) {
    commentID := c.Param("id")
    
    var comment models.Comment
    if err := h.DB.First(&comment, commentID).Error; err != nil {
        h.notFound(c, "Comment not found")
        return
    }
    
    // Получаем текущего пользователя
    user, err := h.GetUserFromContext(c)
    if err != nil {
        h.unauthorized(c, "User not authenticated")
        return
    }
    
    // Проверяем, что пользователь является автором комментария
    if comment.AuthorID != user.ID {
        h.error(c, http.StatusForbidden, "You can only delete your own comments")
        return
    }
    
    if err := h.DB.Delete(&comment).Error; err != nil {
        h.internalError(c, "Failed to delete comment")
        return
    }
    
    h.success(c, nil, "Comment deleted successfully")
}