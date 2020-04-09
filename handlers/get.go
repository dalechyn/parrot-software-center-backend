package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"parrot-software-center-backend/models"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func GetPackage(w http.ResponseWriter, r *http.Request) {
	log.Debug("GetPackage package attempt")

	packageId, exists := mux.Vars(r)["id"]
	if !exists {
		log.Debug("Bad request: ", r.URL.String())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	lookedUpPackage := &models.PackageSchema{}
	if err := packageCollection.
		FindOne(context.TODO(), bson.D{{"id", packageId}}).
		Decode(&lookedUpPackage); err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
	}

	resBytes, err := json.Marshal(&getResponse{lookedUpPackage.ID, lookedUpPackage.Name,
		lookedUpPackage.Description})
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
