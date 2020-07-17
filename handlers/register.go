package handlers

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"parrot-software-center-backend/utils"

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
	row := db.QueryRow("select id from Users where username = $1 or email = $2", inRequest.Login,
		inRequest.Email)
	if err := row.Scan(&id); err == nil {
		log.Errorf("attempt to register existing user - username: %s, email: %s",
			inRequest.Login, inRequest.Email)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(inRequest.Password), 14)

	_, err = db.Exec("insert into Users (email, username, password) values ($1, $2, $3)",
		inRequest.Email, inRequest.Login, string(bytes))
	if err != nil{
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
}
