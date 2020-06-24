package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
	"parrot-software-center-backend/utils"

	log "github.com/sirupsen/logrus"
)

func Rate(w http.ResponseWriter, r *http.Request) {
	log.Debug("Rate attempt")

	vars := mux.Vars(r)
	token, tokenExists := vars["token"]
	packageName, nameExists := vars["name"]
	packageRating, ratingExists := vars["rating"]

	if !tokenExists || !nameExists || !ratingExists {
		log.Debug("Bad request: ", r.URL.String())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db := utils.GetDB()

	userId, err := utils.GetIDFromToken(token)

	_, err = db.Exec("replace into Ratings (user_id, package_name, package_rating) values ($1, $2, $3)",
		userId, packageName, packageRating)
	if err != nil{
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
