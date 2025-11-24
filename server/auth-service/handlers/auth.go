package handlers

import (
	"auth-service/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
    Handler
}

func NewAuthHandler(db *gorm.DB, jwtSecret string) *AuthHandler {
    return &AuthHandler{
        Handler: *NewHandler(db, jwtSecret),
    }
}

func (h *AuthHandler) Register(c *gin.Context) {
    var req models.UserCreateRequest
    if !h.validateRequest(c, &req) {
        return
    }
    
    var existingUser models.User
    if err := h.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
        h.badRequest(c, "User with this email already exists")
        return
    }
    
    var role models.Role
    if err := h.DB.First(&role, req.RoleID).Error; err != nil {
        h.badRequest(c, "Invalid role ID")
        return
    }
    
    user := models.User{
        Email:    req.Email,
        FullName: req.FullName,
        RoleID:   req.RoleID,
    }
    
    if err := user.SetPassword(req.Password); err != nil {
        h.internalError(c, "Failed to hash password")
        return
    }
    
    if err := h.DB.Create(&user).Error; err != nil {
        h.internalError(c, "Failed to create user")
        return
    }
    
    h.DB.Preload("Role").First(&user, user.ID)
    
    h.success(c, gin.H{
        "user": user.ToResponse(),
    }, "User registered successfully")
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req models.UserLoginRequest
    if !h.validateRequest(c, &req) {
        return
    }
    
    var user models.User
    if err := h.DB.Preload("Role").Where("email = ?", req.Email).First(&user).Error; err != nil {
        h.unauthorized(c, "Invalid email or password")
        return
    }
    
    if !user.CheckPassword(req.Password) {
        h.unauthorized(c, "Invalid email or password")
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

func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
    user, err := h.GetUserFromContext(c)
    if err != nil {
        h.unauthorized(c, "User not authenticated")
        return
    }
    
    h.success(c, gin.H{
        "user": user.ToResponse(),
    }, "User data retrieved successfully")
}

func (h *AuthHandler) generateJWT(user models.User) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "email":   user.Email,
        "role_id": user.RoleID,
        "role":    user.Role.RoleName,
        "exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
    })
    
    return token.SignedString([]byte(h.JWTSecret))
}