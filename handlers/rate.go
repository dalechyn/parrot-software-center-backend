package handlers

import (
	"net/http"
	"parrot-software-center-backend/utils"

	log "github.com/sirupsen/logrus"
)

func Rate(w http.ResponseWriter, r *http.Request) {
	log.Debug("Rate attempt")

	q := r.URL.Query()
	token := q.Get("token")
	packageName := q.Get("name")
	packageRating := q.Get("mark")

	if token == "" || packageName == "" || packageRating == "" {
		log.Debug("Bad request: ", r.URL.String())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db := utils.GetDB()

	userId, err := utils.GetIDFromToken(token)

	_, err = db.Exec("insert into Ratings (user_id, package_name, package_rating) values ($1, $2, $3)",
		userId, packageName, packageRating)
	if err != nil{
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
