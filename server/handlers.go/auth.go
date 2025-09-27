package handlers

import (
	"time"

	"kopatel_online/models"

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
        Handler:   Handler{DB: db},
        JWTSecret: jwtSecret,
    }
}

func (h *AuthHandler) Register(c *gin.Context) {
    var req models.UserCreateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        h.badRequest(c, "Invalid request data")
        return
    }
    
    // Проверка существования пользователя
    var existingUser models.User
    if err := h.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
        h.badRequest(c, "User already exists")
        return
    }
    
    // Хеширование пароля
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        h.internalError(c, "Failed to hash password")
        return
    }
    
    user := models.User{
        Email:        req.Email,
        PasswordHash: string(hashedPassword),
        FullName:     req.FullName,
        RoleID:       req.RoleID,
    }
    
    if err := h.DB.Create(&user).Error; err != nil {
        h.internalError(c, "Failed to create user")
        return
    }
    
    h.success(c, user.ToResponse(), "User registered successfully")
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req models.UserLoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        h.badRequest(c, "Invalid request data")
        return
    }
    
    var user models.User
    if err := h.DB.Preload("Role").Where("email = ?", req.Email).First(&user).Error; err != nil {
        h.badRequest(c, "Invalid credentials")
        return
    }
    
    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
        h.badRequest(c, "Invalid credentials")
        return
    }
    
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

func (h *AuthHandler) generateJWT(user models.User) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "email":   user.Email,
        "role_id": user.RoleID,
        "exp":     time.Now().Add(time.Hour * 24).Unix(),
    })
    
    return token.SignedString([]byte(h.JWTSecret))
}