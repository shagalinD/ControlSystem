package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Handler struct {
    DB           *gorm.DB
    Validate     *validator.Validate
    JWTSecret    string
    AuthServiceURL string
    ProjectDefectServiceURL string
}

func NewHandler(db *gorm.DB, jwtSecret, authServiceURL, projectDefectServiceURL string) *Handler {
    validate := validator.New()
    return &Handler{
        DB:            db,
        Validate:      validate,
        JWTSecret:     jwtSecret,
        AuthServiceURL: authServiceURL,
        ProjectDefectServiceURL: projectDefectServiceURL,
    }
}

type Response struct {
    Success bool   `json:"success"`
    Message string `json:"message,omitempty"`
    Data    any    `json:"data,omitempty"`
    Error   string `json:"error,omitempty"`
}

func (h *Handler) success(c *gin.Context, data interface{}, message string) {
    c.JSON(http.StatusOK, Response{
        Success: true,
        Message: message,
        Data:    data,
    })
}

func (h *Handler) error(c *gin.Context, status int, message string) {
    c.JSON(status, Response{
        Success: false,
        Error:   message,
    })
}

func (h *Handler) validateRequest(c *gin.Context, req interface{}) bool {
    if err := c.ShouldBindJSON(req); err != nil {
        h.error(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
        return false
    }
    
    if h.Validate == nil {
        h.internalError(c, "Validator not initialized")
        return false
    }
    
    if err := h.Validate.Struct(req); err != nil {
        h.error(c, http.StatusBadRequest, "Validation error: "+err.Error())
        return false
    }
    
    return true
}

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

func (h *Handler) GetUserFromContext(c *gin.Context) (uint, string, error) {
    userID, exists := c.Get("user_id")
    if !exists {
        return 0, "", fmt.Errorf("user not authenticated")
    }
    
    var userIDUint uint
    switch v := userID.(type) {
    case uint:
        userIDUint = v
    case float64:
        userIDUint = uint(v)
    case int:
        userIDUint = uint(v)
    default:
        return 0, "", fmt.Errorf("invalid user ID type: %T", userID)
    }
    
    userRole, exists := c.Get("user_role")
    if !exists {
        return 0, "", fmt.Errorf("user role not found")
    }
    
    roleStr, ok := userRole.(string)
    if !ok {
        return 0, "", fmt.Errorf("invalid user role type")
    }
    
    return userIDUint, roleStr, nil
}

func (h *Handler) getPaginationParams(c *gin.Context) (int, int) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
    
    if page < 1 {
        page = 1
    }
    if pageSize < 1 || pageSize > 100 {
        pageSize = 20
    }
    
    return page, pageSize
}