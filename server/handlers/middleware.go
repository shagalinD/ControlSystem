package handlers

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware - middleware для проверки JWT токена
func (h *AuthHandler) AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            h.unauthorized(c, "Authorization header required")
            c.Abort()
            return
        }
        
        tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
        if tokenString == "" {
            h.unauthorized(c, "Bearer token required")
            c.Abort()
            return
        }
        
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, jwt.ErrSignatureInvalid
            }
            return []byte(h.JWTSecret), nil
        })
        
        if err != nil || !token.Valid {
            h.unauthorized(c, "Invalid or expired token")
            c.Abort()
            return
        }
        
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            h.unauthorized(c, "Invalid token claims")
            c.Abort()
            return
        }
        
        // Сохраняем данные пользователя в контекст
        c.Set("user_id", uint(claims["user_id"].(float64)))
        c.Set("user_role", claims["role"])
        c.Set("user_email", claims["email"])
        
        c.Next()
    }
}

// RoleMiddleware - middleware для проверки ролей
func (h *AuthHandler) RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole, exists := c.Get("user_role")
        if !exists {
            h.unauthorized(c, "User role not found")
            c.Abort()
            return
        }
        
        roleStr, ok := userRole.(string)
        if !ok {
            h.unauthorized(c, "Invalid user role")
            c.Abort()
            return
        }
        
        // Проверяем, есть ли роль пользователя в разрешенных ролях
        allowed := false
        for _, allowedRole := range allowedRoles {
            if roleStr == allowedRole {
                allowed = true
                break
            }
        }
        
        if !allowed {
            h.error(c, http.StatusForbidden, "Insufficient permissions")
            c.Abort()
            return
        }
        
        c.Next()
    }
}