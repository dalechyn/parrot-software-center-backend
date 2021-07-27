package tokens

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"parrot-software-center-backend/utils"
	"time"
)

const (
	accessTokenExpiry = 24 * time.Second
	accessTokenType   = "access"

	refreshTokenExpiry = 5 * 24 * time.Hour
	refreshTokenType   = "refresh"
)

type AccessTokenClaims struct {
	UserKey   string
	Role      string
	TokenType string
	jwt.StandardClaims
}

type refreshTokenClaims struct {
	UserKey   string
	Role      string
	TokenType string
	jwt.StandardClaims
}

// Generates Access token which has the UserKey inside
func generateAccessToken(userKey string, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, AccessTokenClaims{
		userKey,
		role,
		accessTokenType,
		jwt.StandardClaims{ExpiresAt: time.Now().Add(accessTokenExpiry).Unix()}})

	secretKey := utils.GetSecret()
	return token.SignedString([]byte(secretKey))
}

// Generates Refresh token
func generateRefreshToken(userID string, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims{
		userID,
		role,
		refreshTokenType,
		jwt.StandardClaims{ExpiresAt: time.Now().Add(refreshTokenExpiry).Unix()}})

	secretKey := utils.GetSecret()

	res, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Error(err)
		return "", err
	}

	return res, nil
}

// Replaces old Refresh token ID with a new one given

// Validates refresh token by type, algorithm and expiry
func parseRefreshToken(refreshToken string) (string, string, error) {
	// Pulling out secret key to validate and decode refresh JWT
	secretKey:= utils.GetSecret()

	token, err := jwt.ParseWithClaims(refreshToken, &refreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Checking encryption algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		log.Error(err)
		return "", "", err
	}

	claims, ok := token.Claims.(*refreshTokenClaims)

	// Checking if token is valid
	if !ok || !token.Valid {
		log.Error(err)
		return "", "", errors.New("invalid token")
	}

	// RefreshToken type must be a refresh token
	if claims.TokenType != refreshTokenType {
		err := errors.New("bad token type")
		log.Error(err)
		return "", "", err
	}

	return claims.UserKey, claims.Role, nil
}
/*
func ValidateAccessToken(accessToken string) error {
	// Pulling out secret key to validate and decode access JWT
	secretKey, err := getSecretKey()
	if err != nil {
		log.Error(err)
		return err
	}
	token, err := jwt.ParseWithClaims(accessToken, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Checking encryption algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secretKey), nil
	})

	return nil
}
*/
/*func Generate(key string, Role string) (string, error) {
	// JWT creation
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &models.Claims{Key: key, Role: Role})

	tokenString, _ := token.SignedString([]byte(utils.GetSecret()))
	return tokenString, nil
}*/

func GenerateTokens(userKey string, role string) (string, string, error) {
	accessToken, err := generateAccessToken(userKey, role)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := generateRefreshToken(userKey, role)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func UpdateTokens(oldRefreshToken string) (string, string, error) {
	userKey, role, err := parseRefreshToken(oldRefreshToken)
	if err != nil {
		return "", "", err
	}

	accessToken, err := generateAccessToken(userKey, role)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := generateRefreshToken(userKey, role)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func GetKeyFromToken(tokenStr string) (string, error) {
	hmacSecret := []byte(utils.GetSecret())
	token, err := jwt.ParseWithClaims(tokenStr, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return hmacSecret, nil
	})

	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token")
	}

	return claims.UserKey, nil
}