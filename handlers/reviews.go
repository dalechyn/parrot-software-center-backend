package handlers

import (
	"encoding/json"
	"net/http"
	"parrot-software-center-backend/utils"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func Reviews(w http.ResponseWriter, r *http.Request) {
	log.Debug("Reviews attempt")

	packageName, exists := mux.Vars(r)["name"]
	if !exists {
		log.Debug("Bad request: ", r.URL.String())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db := utils.GetDB()

	var lookedUpRatings []reviewResponse
	rows, err := db.Query("select author, rating, commentary from ratings where name = $1", packageName)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := reviewResponse{}
		err := rows.Scan(&r.Author, &r.Rating, &r.Commentary)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		lookedUpRatings = append(lookedUpRatings, r)
	}

	if len(lookedUpRatings) == 0 {
		log.Error(err)
		w.WriteHeader(http.StatusNoContent)
		return
	}


	resBytes, err := json.Marshal(&lookedUpRatings)
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
