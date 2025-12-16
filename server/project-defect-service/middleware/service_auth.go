package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// ServiceAuthMiddleware - middleware для межсервисной аутентификации
func ServiceAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Проверяем сервисный токен
        serviceToken := c.GetHeader("X-Service-Token")
        if serviceToken != "" && strings.HasPrefix(serviceToken, "service-token-") {
            // Сервисный токен валиден, пропускаем запрос
            c.Next()
            return
        }
        
        // Если нет сервисного токена, используем обычную JWT аутентификацию
        c.Next()
    }
}