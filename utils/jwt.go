package utils

import (
    "os"
    "time"
    "github.com/golang-jwt/jwt/v5"
)
func GenerateToken(userID string, role string, fullName string, email string) (string, error) {
    secret := os.Getenv("JWT_SECRET")
    claims := jwt.MapClaims{
        "user_id":   userID,
        "role":      role,
        "full_name": fullName, 
        "email":     email,    
        "exp":       time.Now().Add(time.Hour * 72).Unix(),
        "iat":       time.Now().Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}