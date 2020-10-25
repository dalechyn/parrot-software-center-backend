package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"net/http"
	"parrot-software-center-backend/utils"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// GET route to get whole ratings/reviews information
func Reviews(w http.ResponseWriter, r *http.Request) {
	log.Debug("Reviews attempt")

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

	// Scanning keysAndScores related to given package name
	var cursor uint64
	var keysAndScores []string
	for {
		var newKeys []string
		var err error
		newKeys, cursor, err = rdb.ZScan(ctx, "ratings", cursor, fmt.Sprintf("rating_%s_*", packageName), 10).Result()
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		keysAndScores = append(keysAndScores, newKeys...)
		if cursor == 0 {
			break
		}
	}

	// Filling up data
	var lookedUpRatings []reviewResponse
	fmt.Print("HEY", keysAndScores)
	for i := 0; i < len(keysAndScores); i += 2 {
		res, err := rdb.Get(ctx, keysAndScores[i]).Result()
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		rating, err := strconv.Atoi(keysAndScores[i + 1])
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		lookedUpRatings = append(lookedUpRatings, reviewResponse{
			Author:     strings.Split(keysAndScores[i], "_")[2],
			Rating:     rating,
			Commentary: res,
		})
	}

	if len(lookedUpRatings) == 0 {
		log.Info(fmt.Sprintf("no reviews on %s", packageName))
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
