package utils

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func GetSecret() string {
	secretKey, exists := os.LookupEnv("SECRET_KEY")
	if !exists {
		log.Fatal("env variable SECRET_KEY is empty")
	}
	return secretKey
}

func GetRedisPassword() string {
	password, exists := os.LookupEnv("REDIS_PASSWORD")
	if !exists {
		log.Fatal("env variable REDIS_PASSWORD is empty")
	}
	return password
}