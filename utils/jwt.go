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

func ValidateToken(tokenString string) (uint, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetJWTSecretKey()), nil
	})
	if err != nil {
		return 0, err
	}
	if claims, ok := token.Claims.(*jwt.MapClaims); ok && token.Valid {
		return uint((*claims)["user_id"].(float64)), nil
	}
	return 0, err
}
