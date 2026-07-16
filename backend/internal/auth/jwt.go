package auth

import (
	"os"
	"strings"
	"time"

	"errors"

	"github.com/golang-jwt/jwt/v5"
)

func getEnv(k string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return ""
}

var jwtSecret = []byte(getEnv("JWT_SECRET"))

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID, email, role string) (string, error) {
	exp := time.Now().Add(72 * time.Hour)
	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "blog_projectss",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateJWT(tokenStr string) (*Claims, error) {
	tokenStr = strings.TrimSpace(tokenStr)
	if tokenStr == "" {
		return nil, errors.New("empty token")
	}
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	return claims, nil
}
