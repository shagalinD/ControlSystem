package handlers

import (
	"net/http"
	"strconv"

	"auth-service/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
    Handler
}

func NewUserHandler(db *gorm.DB, jwtSecret string) *UserHandler {
    return &UserHandler{
        Handler: *NewHandler(db, jwtSecret),
    }
}

// GetEngineers - получение списка инженеров
func (h *UserHandler) GetEngineers(c *gin.Context) {
    var engineers []models.User
    
    if err := h.DB.
        Preload("Role").
        Joins("JOIN roles ON users.role_id = roles.id").
        Where("roles.role_name = ?", "engineer").
        Find(&engineers).Error; err != nil {
        h.internalError(c, "Failed to fetch engineers")
        return
    }
    
    var engineerResponses []models.UserResponse
    for _, engineer := range engineers {
        engineerResponses = append(engineerResponses, engineer.ToResponse())
    }
    
    h.success(c, gin.H{
        "engineers": engineerResponses,
    }, "Engineers retrieved successfully")
}

// GetManagers - получение списка менеджеров
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

// GetAllUsers - получение всех пользователей
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

// UpdateUserData - обновление данных пользователя
func (h *UserHandler) UpdateUserData(c *gin.Context) {
    userIDStr := c.Param("id")
    
    userID, err := strconv.ParseUint(userIDStr, 10, 32)
    if err != nil {
        h.badRequest(c, "Invalid user ID")
        return
    }
    
    var req models.UserUpdateRequest
    if !h.validateRequest(c, &req) {
        return
    }
    
    // Получаем текущего пользователя из контекста
    currentUser, err := h.GetUserFromContext(c)
    if err != nil {
        h.unauthorized(c, "User not authenticated")
        return
    }
    
    // Проверяем что пользователь обновляет свои данные или имеет права
    if currentUser.ID != uint(userID) && currentUser.Role.RoleName != "manager" {
        h.error(c, http.StatusForbidden, "You can only update your own profile")
        return
    }
    
    var user models.User
    if err := h.DB.Preload("Role").First(&user, userID).Error; err != nil {
        h.notFound(c, "User not found")
        return
    }
    
    // Проверяем email на уникальность если он меняется
    if user.Email != req.Email {
        var existingUser models.User
        if err := h.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
            h.badRequest(c, "User with this email already exists")
            return
        }
    }
    
    user.Email = req.Email
    user.FullName = req.FullName
    
    if err := h.DB.Save(&user).Error; err != nil {
        h.internalError(c, "Failed to update user")
        return
    }
    
    h.success(c, gin.H{
        "user": user.ToResponse(),
    }, "User updated successfully")
}

// GetUserByID - получение пользователя по ID
func (h *UserHandler) GetUserByID(c *gin.Context) {
    userIDStr := c.Param("id")
    
    userID, err := strconv.ParseUint(userIDStr, 10, 32)
    if err != nil {
        h.badRequest(c, "Invalid user ID")
        return
    }
    
    var user models.User
    if err := h.DB.Preload("Role").First(&user, userID).Error; err != nil {
        h.notFound(c, "User not found")
        return
    }
    
    h.success(c, gin.H{
        "user": user.ToResponse(),
    }, "User retrieved successfully")
}