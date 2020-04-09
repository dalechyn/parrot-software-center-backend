package handlers

import (
	"encoding/json"
	"net/http"
	"parrot-software-center-backend/models"
	"parrot-software-center-backend/utils"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func GetRatings(w http.ResponseWriter, r *http.Request) {
	log.Debug("GetRatings attempt")

	packageName, exists := mux.Vars(r)["name"]
	if !exists {
		log.Debug("Bad request: ", r.URL.String())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db := utils.GetDB()

	lookedUpRating := &models.PackageRating{}
	row := db.QueryRow("select * from Ratings where package_name = $1", packageName)
	err := row.Scan(&lookedUpRating.Name, &lookedUpRating.Rating)
	if err != nil{
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resBytes, err := json.Marshal(&getResponse{lookedUpRating.Name, lookedUpRating.Rating})
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
