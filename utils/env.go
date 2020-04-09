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