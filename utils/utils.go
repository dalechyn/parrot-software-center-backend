package utils

import (
	"errors"
	"fmt"
	"parrot-software-center-backend/models"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/mattn/go-sqlite3"
)

func GetKeyFromToken(tokenStr string) (string, error) {
	hmacSecret := []byte(GetSecret())
	token, err := jwt.ParseWithClaims(tokenStr, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return hmacSecret, nil
	})

	if err != nil {
		return "", errors.New("invalid token")
	}
	claims, ok := token.Claims.(*models.Claims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token")
	}

	return claims.Key, nil
}