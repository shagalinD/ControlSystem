package handlers

import (
	"kopatel_online/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
    Handler
    JWTSecret string
}

func NewAuthHandler(db *gorm.DB, jwtSecret string) *AuthHandler {
    return &AuthHandler{
        Handler:   *NewHandler(db), // 
        JWTSecret: jwtSecret,
    }
}

// Register - регистрация нового пользователя
func (h *AuthHandler) Register(c *gin.Context) {
    var req models.UserCreateRequest
    if !h.validateRequest(c, &req) {
        return
    }
    
    // Проверка существования пользователя
    var existingUser models.User
    if err := h.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
        h.badRequest(c, "User with this email already exists")
        return
    }
    
    // Проверка существования роли
    var role models.Role
    if err := h.DB.First(&role, req.RoleID).Error; err != nil {
        h.badRequest(c, "Invalid role ID")
        return
    }
    
    // Создание пользователя
    user := models.User{
        Email:    req.Email,
        FullName: req.FullName,
        RoleID:   req.RoleID,
    }
    
    // Хеширование пароля
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        h.internalError(c, "Failed to hash password")
        return
    }
    user.PasswordHash = string(hashedPassword)
    
    if err := h.DB.Create(&user).Error; err != nil {
        h.internalError(c, "Failed to create user")
        return
    }
    
    // Загружаем роль для ответа
    h.DB.Preload("Role").First(&user, user.ID)
    
    h.success(c, gin.H{
        "user": user.ToResponse(),
    }, "User registered successfully")
}

// Login - аутентификация пользователя
func (h *AuthHandler) Login(c *gin.Context) {
    var req models.UserLoginRequest
    if !h.validateRequest(c, &req) {
        return
    }
    
    // Поиск пользователя
    var user models.User
    if err := h.DB.Preload("Role").Where("email = ?", req.Email).First(&user).Error; err != nil {
        h.unauthorized(c, "Invalid email or password")
        return
    }
    
    // Проверка пароля
    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
        h.unauthorized(c, "Invalid email or password")
        return
    }
    
    // Генерация JWT токена
    token, err := h.generateJWT(user)
    if err != nil {
        h.internalError(c, "Failed to generate token")
        return
    }
    
    h.success(c, gin.H{
        "token": token,
        "user":  user.ToResponse(),
    }, "Login successful")
}

// GenerateJWT - создание JWT токена
func (h *AuthHandler) generateJWT(user models.User) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "email":   user.Email,
        "role_id": user.RoleID,
        "role":    user.Role.RoleName,
        "exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 дней
    })
    
    return token.SignedString([]byte(h.JWTSecret))
}

// GetCurrentUser - получение текущего пользователя
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        h.unauthorized(c, "User not authenticated")
        return
    }
    
    var user models.User
    if err := h.DB.Preload("Role").First(&user, userID).Error; err != nil {
        h.unauthorized(c, "User not found")
        return
    }
    
    h.success(c, gin.H{
        "user": user.ToResponse(),
    }, "User data retrieved successfully")
}