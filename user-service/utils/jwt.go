package jwt

import (
	"github.com/golang-jwt/jwt/v5"
    "os"
    "time"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
    UserID uint `json:"user_id"`
    Role int `json:"role"`
    jwt.RegisteredClaims
}

// GenerateToken generates a JWT token for a given user ID
func GenerateToken(userID uint, role int) (string, error) {
    expirationTime := time.Now().Add(30 * time.Minute)
    claims := &Claims{
        UserID: userID,
        Role: role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}

// ValidateToken validates a given JWT token and returns user ID
func ValidateToken(tokenString string) (uint, int, error) {
    claims := &Claims{}

    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })

    if err != nil {
        return 0, 0, err
    }

    if !token.Valid {
        return 0, 0, jwt.ErrTokenUnverifiable
    }

    return claims.UserID, claims.Role, nil
}
