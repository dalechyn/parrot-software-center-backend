package handlers

import (
	"database/sql"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"parrot-software-center-backend/models"
	"parrot-software-center-backend/tokens"
	"parrot-software-center-backend/utils"

	log "github.com/sirupsen/logrus"
)

func Login(w http.ResponseWriter, r *http.Request) {
	log.Debug("Login attempt")
	db := utils.GetDB()

	// Decoding http request
	inRequest := &loginRequest{}
	err := json.NewDecoder(r.Body).Decode(inRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if user exists
	lookedUpUser := &models.User{}
	row := db.QueryRow("select id, username, password from Users where username = $1", inRequest.Username)
	if err := row.Scan(&lookedUpUser.ID, &lookedUpUser.Username, &lookedUpUser.Password); err != nil && err != sql.ErrNoRows {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	byteHash := []byte(lookedUpUser.Password)
	err = bcrypt.CompareHashAndPassword(byteHash, []byte(inRequest.Password))
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// Generate some tokens for him
	token, err := tokens.Generate(lookedUpUser.Username)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Encode to JSON and send him
	resBytes, err := json.Marshal(&loginResponse{token})
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resBytes); err != nil {
		log.Error(err)
	}
}
