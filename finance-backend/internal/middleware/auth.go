package middleware

import (
    "finance-backend/internal/auth"
    "net/http"
    "strings"
    "github.com/gin-gonic/gin"
)

// AuthMiddleware проверяет JWT токен
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Берём токен из заголовка Authorization
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }

        // Проверяем формат: Bearer <token>
        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
            c.Abort()
            return
        }

        // Проверяем токен
        claims, err := auth.ValidateToken(parts[1])
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
            c.Abort()
            return
        }

        // Сохраняем ID пользователя в контексте
        c.Set("userID", claims.UserID)
        c.Next()
    }
}

// GetUserID возвращает ID пользователя из контекста
func GetUserID(c *gin.Context) int64 {
    return c.GetInt64("userID")
}
