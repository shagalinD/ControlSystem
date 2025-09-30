package handlers

import (
	"io"
	"kopatel_online/models"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AttachmentHandler struct {
    Handler
    UploadPath string
}

func NewAttachmentHandler(db *gorm.DB) *AttachmentHandler {
    // Создаем папку для загрузок, если её нет
    uploadPath := "./uploads"
    os.MkdirAll(uploadPath, 0755)
    
    return &AttachmentHandler{
        Handler:    Handler{DB: db},
        UploadPath: uploadPath,
    }
}

// UploadAttachment - загрузка файла для дефекта
func (h *AttachmentHandler) UploadAttachment(c *gin.Context) {
    defectID := c.Param("defect_id")
    
    // Проверяем существование дефекта
    var defect models.Defect
    if err := h.DB.First(&defect, defectID).Error; err != nil {
        h.notFound(c, "Defect not found")
        return
    }
    
    // Получаем файл из запроса
    file, header, err := c.Request.FormFile("file")
    if err != nil {
        h.badRequest(c, "File is required")
        return
    }
    defer file.Close()
    
    // Получаем текущего пользователя
    user, err := h.GetUserFromContext(c)
    if err != nil {
        h.unauthorized(c, "User not authenticated")
        return
    }
    
    // Генерируем уникальное имя файла
    ext := filepath.Ext(header.Filename)
    filename := strconv.FormatInt(time.Now().UnixNano(), 10) + ext
    filepath := filepath.Join(h.UploadPath, filename)
    
    // Сохраняем файл
    out, err := os.Create(filepath)
    if err != nil {
        h.internalError(c, "Failed to save file")
        return
    }
    defer out.Close()
    
    _, err = io.Copy(out, file)
    if err != nil {
        h.internalError(c, "Failed to save file")
        return
    }
    
    // Создаем запись в БД
    attachment := models.Attachment{
        DefectID:   defect.ID,
        Filename:   header.Filename,
        Filepath:   filename,
        FileSize:   header.Size,
        MimeType:   header.Header.Get("Content-Type"),
        UploadedBy: user.ID,
    }
    
    if err := h.DB.Create(&attachment).Error; err != nil {
        // Удаляем файл если не удалось сохранить в БД
        os.Remove(filepath)
        h.internalError(c, "Failed to save attachment info")
        return
    }
    
    h.success(c, gin.H{
        "attachment": attachment,
    }, "File uploaded successfully")
}

// GetAttachments - получение списка вложений для дефекта
func (h *AttachmentHandler) GetAttachments(c *gin.Context) {
    defectID := c.Param("defect_id")
    
    var attachments []models.Attachment
    if err := h.DB.Preload("Uploader").
        Where("defect_id = ?", defectID).
        Find(&attachments).Error; err != nil {
        h.internalError(c, "Failed to fetch attachments")
        return
    }
    
    h.success(c, gin.H{
        "attachments": attachments,
    }, "Attachments retrieved successfully")
}

// DownloadAttachment - скачивание файла
func (h *AttachmentHandler) DownloadAttachment(c *gin.Context) {
    attachmentID := c.Param("id")
    
    var attachment models.Attachment
    if err := h.DB.First(&attachment, attachmentID).Error; err != nil {
        h.notFound(c, "Attachment not found")
        return
    }
    
    filePath := filepath.Join(h.UploadPath, attachment.Filepath)
    
    // Проверяем существование файла
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        h.notFound(c, "File not found")
        return
    }
    
    // Устанавливаем заголовки для скачивания
    c.Header("Content-Disposition", "attachment; filename="+attachment.Filename)
    c.Header("Content-Type", attachment.MimeType)
    c.File(filePath)
}

// DeleteAttachment - удаление вложения
func (h *AttachmentHandler) DeleteAttachment(c *gin.Context) {
    attachmentID := c.Param("id")
    
    var attachment models.Attachment
    if err := h.DB.First(&attachment, attachmentID).Error; err != nil {
        h.notFound(c, "Attachment not found")
        return
    }
    
    // Получаем текущего пользователя
    user, err := h.GetUserFromContext(c)
    if err != nil {
        h.unauthorized(c, "User not authenticated")
        return
    }
    
    // Проверяем права (только автор загрузки или менеджер)
    if attachment.UploadedBy != user.ID && user.Role.RoleName != "manager" {
        h.error(c, http.StatusForbidden, "You can only delete your own attachments")
        return
    }
    
    filePath := filepath.Join(h.UploadPath, attachment.Filepath)
    
    // Удаляем файл
    if err := os.Remove(filePath); err != nil {
        h.internalError(c, "Failed to delete file")
        return
    }
    
    // Удаляем запись из БД
    if err := h.DB.Delete(&attachment).Error; err != nil {
        h.internalError(c, "Failed to delete attachment record")
        return
    }
    
    h.success(c, nil, "Attachment deleted successfully")
}