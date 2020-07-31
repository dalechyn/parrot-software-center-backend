package handlers

import (
	"encoding/json"
	"net/http"
	"parrot-software-center-backend/utils"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func Ratings(w http.ResponseWriter, r *http.Request) {
	log.Debug("Ratings attempt")

	packageName, exists := mux.Vars(r)["name"]
	if !exists {
		log.Debug("Bad request: ", r.URL.String())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db := utils.GetDB()

	rows, err := db.Query("select rating from ratings where name = $1", packageName)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer rows.Close()
	rating := 0
	quantity := 0
	for rows.Next() {
		quantity++
		rowRating := 0
		err := rows.Scan(&rowRating)
		if err != nil{
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		rating += rowRating
	}

	if quantity == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	resBytes, err := json.Marshal(&getResponse{
		float64(rating) / float64(quantity)})
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
