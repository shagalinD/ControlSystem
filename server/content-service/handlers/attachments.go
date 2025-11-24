package handlers

import (
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"content-service/models"
	"content-service/storage"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AttachmentHandler struct {
    Handler
    UploadPath string
    FileStorage *storage.FileStorage
}

func NewAttachmentHandler(db *gorm.DB, jwtSecret, authServiceURL, projectDefectServiceURL, uploadPath string) *AttachmentHandler {
    fileStorage := storage.NewFileStorage(uploadPath)
    return &AttachmentHandler{
        Handler:    *NewHandler(db, jwtSecret, authServiceURL, projectDefectServiceURL),
        UploadPath: uploadPath,
        FileStorage: fileStorage,
    }
}

func (h *AttachmentHandler) UploadAttachment(c *gin.Context) {
    defectIDStr := c.Param("defect_id")
    
    defectID, err := strconv.ParseUint(defectIDStr, 10, 32)
    if err != nil {
        h.badRequest(c, "Invalid defect ID")
        return
    }
    
    // Получаем файл из запроса
    file, header, err := c.Request.FormFile("file")
    if err != nil {
        h.badRequest(c, "File is required")
        return
    }
    defer file.Close()
    
    userID, _, err := h.GetUserFromContext(c)
    if err != nil {
        h.unauthorized(c, "User not authenticated")
        return
    }
    
    // Генерируем уникальное имя файла
    ext := filepath.Ext(header.Filename)
    filename := strconv.FormatInt(time.Now().UnixNano(), 10) + ext
    
    // Сохраняем файл
    filePath, err := h.FileStorage.SaveFile(filename, file)
    if err != nil {
        h.internalError(c, "Failed to save file")
        return
    }
    
    // Создаем запись в БД
    attachment := models.Attachment{
        DefectID:   uint(defectID), 
        Filename:   header.Filename,
        Filepath:   filename,
        FileSize:   header.Size,
        MimeType:   header.Header.Get("Content-Type"),
        UploadedBy: userID,
    }
    
    if err := h.DB.Create(&attachment).Error; err != nil {
        // Удаляем файл если не удалось сохранить в БД
        h.FileStorage.DeleteFile(filePath)
        h.internalError(c, "Failed to save attachment info")
        return
    }
    
    h.success(c, gin.H{
        "attachment": attachment,
    }, "File uploaded successfully")
}

func (h *AttachmentHandler) GetAttachments(c *gin.Context) {
    defectID := c.Param("defect_id")
    
    var attachments []models.Attachment
    if err := h.DB.
        Where("defect_id = ?", defectID).
        Find(&attachments).Error; err != nil {
        h.internalError(c, "Failed to fetch attachments")
        return
    }
    
    h.success(c, gin.H{
        "attachments": attachments,
    }, "Attachments retrieved successfully")
}

func (h *AttachmentHandler) DownloadAttachment(c *gin.Context) {
    attachmentID := c.Param("id")
    
    var attachment models.Attachment
    if err := h.DB.First(&attachment, attachmentID).Error; err != nil {
        h.notFound(c, "Attachment not found")
        return
    }
    
    filePath := filepath.Join(h.UploadPath, attachment.Filepath)
    
    // Проверяем существование файла
    if !h.FileStorage.FileExists(filePath) {
        h.notFound(c, "File not found")
        return
    }
    
    // Устанавливаем заголовки для скачивания
    c.Header("Content-Disposition", "attachment; filename="+attachment.Filename)
    c.Header("Content-Type", attachment.MimeType)
    c.File(filePath)
}

func (h *AttachmentHandler) DeleteAttachment(c *gin.Context) {
    attachmentID := c.Param("id")
    
    var attachment models.Attachment
    if err := h.DB.First(&attachment, attachmentID).Error; err != nil {
        h.notFound(c, "Attachment not found")
        return
    }
    
    userID, userRole, err := h.GetUserFromContext(c)
    if err != nil {
        h.unauthorized(c, "User not authenticated")
        return
    }
    
    // Проверяем права (только автор загрузки или менеджер)
    if attachment.UploadedBy != userID && userRole != "manager" {
        h.error(c, http.StatusForbidden, "You can only delete your own attachments")
        return
    }
    
    filePath := filepath.Join(h.UploadPath, attachment.Filepath)
    
    // Удаляем файл
    if err := h.FileStorage.DeleteFile(filePath); err != nil {
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