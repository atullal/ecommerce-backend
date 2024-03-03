package middleware

import (
	"strings"
    "net/http"
    "github.com/gin-gonic/gin"
    "rest-service/utils" // Replace with your actual jwt package path
)

// AdminMiddleware checks if the authenticated user is an admin
func AdminMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        tokenString = strings.TrimPrefix(tokenString, "Bearer ")

        // Validate and parse the token
        _, role, err := jwt.ValidateToken(tokenString)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            return
        }

        // Extract role from token and check if it's an admin
        if err != nil || role != 3 {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access forbidden"})
            return
        }

        c.Next()
    }
}
