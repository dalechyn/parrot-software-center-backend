package models

import "github.com/dgrijalva/jwt-go"

type User struct {
	ID       int
	Username string
	Password string
	Confirmed bool
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type ConfirmClaims struct {
	ID int64
	jwt.StandardClaims
}
