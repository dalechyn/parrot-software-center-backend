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
	"time"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

var ctx = context.Background()

// POST route to register a user
func Register(w http.ResponseWriter, r *http.Request) {
	log.Debug("Register attempt")

	// Decoding http request
	inRequest := &registerRequest{}
	err := json.NewDecoder(r.Body).Decode(inRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Connecting to Redis
	rdb := redis.NewFailoverClient(&redis.FailoverOptions{
		SentinelAddrs: []string{"sentinel:26379", "sentinel:26380", "sentinel:26381"},
		MasterName: "mymaster",
		SentinelPassword: utils.GetSentinelPassword(),
		Password: utils.GetRedisPassword(),
	})

	// Checking either user already exists or not
	userKey := fmt.Sprintf("user_%s", inRequest.Login)
	if exists, err := rdb.Exists(ctx, userKey).Result(); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		if exists == 1 {
			log.Errorf("attempt to register user with existing username - username: %s, email: %s",
				inRequest.Login, inRequest.Email)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	if exists, err := rdb.SIsMember(ctx, "emails", inRequest.Email).Result(); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		if exists {
			log.Errorf("attempt to register user with existing email - username: %s, email: %s",
				inRequest.Login, inRequest.Email)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(inRequest.Password), 14)

	if _, err := rdb.HSet(ctx,
		userKey,
		"email", inRequest.Email, "password", string(bytes), "confirmed", "0").Result(); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if _, err := rdb.Expire(ctx, userKey, time.Hour).Result(); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Managing keys in sets will be handy in future
	if _, err := rdb.SAdd(ctx, "users", userKey).Result(); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := rdb.SAdd(ctx, "emails", inRequest.Email).Result(); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Generating confirmation token that expires in an hour
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		&models.Claims{
			Key: userKey,
			StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Add(1 * time.Hour).Unix()}})

	emailSecret, emailSecretExists := os.LookupEnv("EMAIL_KEY")
	if !emailSecretExists {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	emailConfirmationJWT, _ := token.SignedString([]byte(emailSecret))

	// Connecting to smtp server and sending the confirmation email
	email, loginExists := os.LookupEnv("EMAIL_LOGIN")
	password, passwordExists := os.LookupEnv("EMAIL_PASSWORD")
	if !loginExists || !passwordExists {
		log.Error("Email credentials are not set")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	auth := smtp.PlainAuth(email, email, password, "smtp.parrotsec.org")
	to := []string{inRequest.Email}
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf(
		mime + `From: %s
To: %s
Subject: Parrot Software Center Email Confirmation

<p><img style="display: block; margin-left: auto; margin-right: auto;" src="https://parrotsec.org/images/logo.png" /></p>
<h1 style="color: #2196f3; text-align: center;">Parrot Software Center Email Confirmation</h1>
<p>Hi! To confirm your Parrot Software Center account, please follow this <a href="http://165.227.140.210:8000/confirm/%s">link</a>.</p>
<p>If that wasn't you, just ignore this letter.</p>
<p>&nbsp;</p>
<h4 style="text-align: center;">Copyright &copy; 2020 Parrot Security CIC</h4>`, email, to, emailConfirmationJWT)
	msg := []byte(body)

	if err := smtp.SendMail("smtp.parrotsec.org:587", auth, email, to, msg); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
