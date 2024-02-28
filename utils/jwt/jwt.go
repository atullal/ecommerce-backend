package utils

import (
	"github.com/golang-jwt/jwt/v5"
    "os"
    "time"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
    UserID uint `json:"user_id"`
    jwt.StandardClaims
}

// GenerateToken generates a JWT token for a given user ID
func GenerateToken(userID uint) (string, error) {
    expirationTime := time.Now().Add(30 * time.Minute)
    claims := &Claims{
        UserID: userID,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}

// ValidateToken validates a given JWT token and returns user ID
func ValidateToken(tokenString string) (uint, error) {
    claims := &Claims{}

    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })

    if err != nil {
        return 0, err
    }

    if !token.Valid {
        return 0, err
    }

    return claims.UserID, nil
}
