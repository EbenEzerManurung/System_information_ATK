package utils

import (
    "time"
    "errors"
    "github.com/golang-jwt/jwt/v5"
    "atk-backend/config"
)

type Claims struct {
    UserID   uint   `json:"user_id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}

func GenerateToken(userID uint, username, email, role string) (string, error) {
    cfg := config.LoadConfig()
    
    claims := Claims{
        UserID:   userID,
        Username: username,
        Email:    email,
        Role:     role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
            Issuer:    "atk-enterprise",
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(cfg.JWTSecret))
}

func ValidateToken(tokenString string) (*Claims, error) {
    cfg := config.LoadConfig()
    
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(cfg.JWTSecret), nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, errors.New("invalid token")
}