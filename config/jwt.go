package config

import (
	"os"
	"time"
)

func GetJWTSecretKey() []byte {
	return []byte(os.Getenv("JWT_SECRET_KEY"))
}

func GetJWTExpireTime() time.Duration {
	duration, err := time.ParseDuration(os.Getenv("JWT_EXPIRES_IN"))
	if err != nil {
		panic(err)
	}
	return duration
}
