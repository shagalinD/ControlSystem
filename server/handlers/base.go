package handlers

import (
	"fmt"
	"kopatel_online/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Handler struct {
    DB       *gorm.DB
    Validate *validator.Validate
}

func NewHandler(db *gorm.DB) *Handler {
    return &Handler{
        DB:       db,
        Validate: validator.New(),
    }
}

// Базовая структура ответа
type Response struct {
    Success bool        `json:"success"`
    Message string      `json:"message,omitempty"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

// Успешный ответ
func (h *Handler) success(c *gin.Context, data interface{}, message string) {
    c.JSON(http.StatusOK, Response{
        Success: true,
        Message: message,
        Data:    data,
    })
}


// Ошибка
func (h *Handler) error(c *gin.Context, status int, message string) {
    c.JSON(status, Response{
        Success: false,
        Error:   message,
    })
}

// Валидация входных данных
func (h *Handler) validateRequest(c *gin.Context, req interface{}) bool {
    if err := c.ShouldBindJSON(req); err != nil {
        h.error(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
        return false
    }
    
    if err := h.Validate.Struct(req); err != nil {
        h.error(c, http.StatusBadRequest, "Validation error: "+err.Error())
        return false
    }
    
    return true
}

// Вспомогательные методы для ошибок
func (h *Handler) badRequest(c *gin.Context, message string) {
    h.error(c, http.StatusBadRequest, message)
}

func (h *Handler) notFound(c *gin.Context, message string) {
    h.error(c, http.StatusNotFound, message)
}

func (h *Handler) internalError(c *gin.Context, message string) {
    h.error(c, http.StatusInternalServerError, message)
}

func (h *Handler) unauthorized(c *gin.Context, message string) {
    h.error(c, http.StatusUnauthorized, message)
}

// GetUserFromContext - получение пользователя из контекста (общий метод)
func (h *Handler) GetUserFromContext(c *gin.Context) (*models.User, error) {
    userID, exists := c.Get("user_id")
    if !exists {
        return nil, fmt.Errorf("user not authenticated")
    }
    
    userIDUint, ok := userID.(uint)
    if !ok {
        // Если тип uint не подходит, пробуем float64 (из JWT)
        if userIDFloat, ok := userID.(float64); ok {
            userIDUint = uint(userIDFloat)
        } else {
            return nil, fmt.Errorf("invalid user ID type")
        }
    }
    
    var user models.User
    if err := h.DB.Preload("Role").First(&user, userIDUint).Error; err != nil {
        return nil, err
    }
    
    return &user, nil
}