package handlers

import (
	"encoding/json"
	"net/http"
	"parrot-software-center-backend/utils"

	log "github.com/sirupsen/logrus"
)

func Rate(w http.ResponseWriter, r *http.Request) {
	log.Debug("Rate attempt")

	// Decoding http request
	inRequest := &rateRequest{}
	err := json.NewDecoder(r.Body).Decode(inRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db := utils.GetDB()

	userId, err := utils.GetIDFromToken(inRequest.Token)

	username := ""
	row := db.QueryRow("select username from users where id = $1", userId)
	if err := row.Scan(&username); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = db.Exec("replace into ratings (user_id, name, author, rating, commentary) values ($1, $2, $3, $4, $5)",
		userId, inRequest.Name, username, inRequest.Rating, inRequest.Comment)
	if err != nil{
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
