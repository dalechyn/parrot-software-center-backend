package models

import "github.com/dgrijalva/jwt-go"

type User struct {
	ID       int
	Username string
	Password string
	Confirmed bool
}

type Claims struct {
	Key string
	jwt.StandardClaims
}
