package handlers

import (
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

	log "github.com/sirupsen/logrus"
)

func Register(w http.ResponseWriter, r *http.Request) {
	log.Debug("Register attempt")
	db := utils.GetDB()

	// Decoding http request
	inRequest := &registerRequest{}
	err := json.NewDecoder(r.Body).Decode(inRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if user exists
	id := 0
	row := db.QueryRow("select id from users where username = $1 or email = $2", inRequest.Login,
		inRequest.Email)
	if err := row.Scan(&id); err == nil {
		log.Errorf("attempt to register existing user - username: %s, email: %s",
			inRequest.Login, inRequest.Email)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(inRequest.Password), 14)

	result, err := db.Exec("insert into users (email, username, password, confirmed) values ($1, $2, $3, 0)",
		inRequest.Email, inRequest.Login, string(bytes))
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newID, err := result.LastInsertId()
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		&models.ConfirmClaims{
			ID: newID,
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
