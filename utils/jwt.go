package utils

import (
	"assignmentday23/config"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(userId uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(config.GetJWTExpireTime()).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GetJWTSecretKey()))
}
