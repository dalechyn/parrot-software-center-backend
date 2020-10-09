package tokens

import (
	"github.com/dgrijalva/jwt-go"
	"parrot-software-center-backend/models"
	"parrot-software-center-backend/utils"
)



func Generate(key string) (string, error) {
	// JWT creation
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &models.Claims{Key: key})

	tokenString, _ := token.SignedString([]byte(utils.GetSecret()))
	return tokenString, nil
}
