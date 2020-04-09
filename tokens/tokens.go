package tokens

import (
	"github.com/dgrijalva/jwt-go"
	"parrot-software-center-backend/models"
	"parrot-software-center-backend/utils"
)



func Generate(username string) (string, error) {
	// JWT creation
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &models.Claims{Username:username})

	tokenString, _ := token.SignedString([]byte(utils.GetSecret()))
	return tokenString, nil
}
