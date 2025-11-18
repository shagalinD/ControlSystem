package handlers

import (
	"kopatel_online/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
    Handler
}

func NewUserHandler(db *gorm.DB) *UserHandler {
    return &UserHandler{Handler: *NewHandler(db)}
}

// GetEngineers - получение списка инженеров
func (h *UserHandler) GetEngineers(c *gin.Context) {
    var engineers []models.User
    
    // Ищем пользователей с ролью "engineer"
    if err := h.DB.
        Preload("Role").
        Joins("JOIN roles ON users.role_id = roles.id").
        Where("roles.role_name = ?", "engineer").
        Find(&engineers).Error; err != nil {
        h.internalError(c, "Failed to fetch engineers")
        return
    }
    
    // Преобразуем в response format
    var engineerResponses []models.UserResponse
    for _, engineer := range engineers {
        engineerResponses = append(engineerResponses, engineer.ToResponse())
    }
    
    h.success(c, gin.H{
        "engineers": engineerResponses,
    }, "Engineers retrieved successfully")
}

// GetManagers - получение списка менеджеров (на будущее)
func (h *UserHandler) GetManagers(c *gin.Context) {
    var managers []models.User
    
    if err := h.DB.
        Preload("Role").
        Joins("JOIN roles ON users.role_id = roles.id").
        Where("roles.role_name = ?", "manager").
        Find(&managers).Error; err != nil {
        h.internalError(c, "Failed to fetch managers")
        return
    }
    
    var managerResponses []models.UserResponse
    for _, manager := range managers {
        managerResponses = append(managerResponses, manager.ToResponse())
    }
    
    h.success(c, gin.H{
        "managers": managerResponses,
    }, "Managers retrieved successfully")
}

// GetAllUsers - получение всех пользователей (для админки)
func (h *UserHandler) GetAllUsers(c *gin.Context) {
    var users []models.User
    
    if err := h.DB.
        Preload("Role").
        Find(&users).Error; err != nil {
        h.internalError(c, "Failed to fetch users")
        return
    }
    
    var userResponses []models.UserResponse
    for _, user := range users {
        userResponses = append(userResponses, user.ToResponse())
    }
    
    h.success(c, gin.H{
        "users": userResponses,
    }, "Users retrieved successfully")
}

// UpdateUserData - обновляем данные пользователя
func (h *UserHandler) UpdateUserData(c *gin.Context) {
    var req models.UserUpdateRequest
    if !h.validateRequest(c, &req) {
        return
    }
    
    user, err := h.GetUserFromContext(c)
    if err != nil {
        h.unauthorized(c, "user is not authorized")
    }

    user.FullName = req.FullName
    user.Email = req.Email 
    
    if err := h.DB.Save(&user).Error; err != nil {
        h.internalError(c, "error on updating user")
        return 
    }
    
    h.success(c, gin.H{
        "user": user,
    }, "Update user successful")
}