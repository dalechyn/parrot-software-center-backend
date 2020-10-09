package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/smtp"
	"os"
	"parrot-software-center-backend/models"
	"parrot-software-center-backend/utils"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

var ctx = context.Background()

func Register(w http.ResponseWriter, r *http.Request) {
	log.Debug("Register attempt")

	// Decoding http request
	inRequest := &registerRequest{}
	err := json.NewDecoder(r.Body).Decode(inRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if user exists
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:26379",
		Password: utils.GetRedisPassword(),
	})

	var cursor uint64
	var keys []string

	for {
		var err error
		var newKeys []string
		newKeys, cursor, err = rdb.SScan(ctx, "users", cursor, "user-*", 10).Result()

		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		for _, key := range newKeys {
			if split := strings.Split(key[4:], "-"); split[0] == inRequest.Email || split[1] == inRequest.Login {
				log.Errorf("attempt to register existing user - username: %s, email: %s",
					inRequest.Login, inRequest.Email)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		keys = append(keys, newKeys...)
		if cursor == 0 {
			break
		}
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(inRequest.Password), 14)

	userKey := fmt.Sprintf("user-%s-%s", inRequest.Email, inRequest.Login)
	if _, err := rdb.HSet(ctx,
		userKey,
		"email", inRequest.Email, "login", inRequest.Login, "password", string(bytes), "confirmed", "0").Result(); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := rdb.SAdd(ctx, "users", userKey).Result(); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		&models.ConfirmClaims{
			Key: userKey,
			StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Add(1 * time.Hour).Unix()}})

	emailSecret, emailSecretExists := os.LookupEnv("EMAIL_KEY")
	if !emailSecretExists {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	emailConfirmationJWT, _ := token.SignedString([]byte(emailSecret))

	email, loginExists := os.LookupEnv("EMAIL_LOGIN")
	password, passwordExists := os.LookupEnv("EMAIL_PASSWORD")
	if !loginExists || !passwordExists {
		log.Error("Email credentials are not set")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	auth := smtp.PlainAuth("Parrot Software Center", email, password, "smtp.gmail.com")
	to := []string{inRequest.Email}
	body := fmt.Sprintf(
		`From: noreply@parrot.sh
To: %s
Subject: Parrot Software Center Account Confirmation

Hi! To confirm your Parrot Software Center account, please follow the link: http://localhost:8000/confirm/%s`, to, emailConfirmationJWT)
	msg := []byte(body)

	if err := smtp.SendMail("smtp.gmail.com:587", auth, "vlad.dalechin@gmail.com", to, msg); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
