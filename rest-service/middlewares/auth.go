package middleware

import (
	"strings"
    "net/http"
    "github.com/gin-gonic/gin"
    "rest-service/utils" // Replace with your actual jwt package path
)

// AuthMiddleware checks for a valid JWT token in the request
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        // Assuming the token is in the format "Bearer <token>"
        tokenString = strings.TrimPrefix(tokenString, "Bearer ")

        if tokenString == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization token not provided"})
            return
        }

        userID, role, err := jwt.ValidateToken(tokenString)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            return
        }

        // Store user ID in context for further use in handlers
        c.Set("userID", userID)
        c.Set("role", role)
        c.Next()
    }
}
