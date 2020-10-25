package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"net/http"
	"parrot-software-center-backend/utils"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// GET route to get average rating of a package
func Ratings(w http.ResponseWriter, r *http.Request) {
	log.Debug("Ratings attempt")

	packageName, exists := mux.Vars(r)["name"]
	if !exists {
		log.Debug("Bad request: ", r.URL.String())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Connecting to Redis
	rdb := redis.NewFailoverClient(&redis.FailoverOptions{
		SentinelAddrs: []string{"sentinel:26379", "sentinel:26380", "sentinel:26381"},
		MasterName: "mymaster",
		SentinelPassword: utils.GetSentinelPassword(),
		Password: utils.GetRedisPassword(),
	})

	// Scanning for keys related to given package name
	rating := 0
	quantity := 0
	var cursor uint64
	for {
		var newKeys []string
		var err error
		newKeys, cursor, err = rdb.ZScan(ctx, "ratings", cursor, fmt.Sprintf("rating_%s_*", packageName), 10).Result()
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		quantity += len(newKeys) / 2
		for i := 0; i < len(newKeys) / 2; i += 2 {
			zRating, err := strconv.Atoi(newKeys[i + 1])
			if err != nil {
				log.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			rating += zRating
		}
		if cursor == 0 {
			break
		}
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
