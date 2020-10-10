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
		SentinelAddrs: []string{":26379", ":26380", ":26381"},
		MasterName: "mymaster",
		SentinelPassword: utils.GetSentinelPassword(),
		Password: utils.GetRedisPassword(),
	})

	// Scanning for keys related to given package name
	var cursor uint64
	var keys []string
	for {
		var newKeys []string
		var err error
		newKeys, cursor, err = rdb.SScan(ctx, "ratings", cursor, fmt.Sprintf("rating-%s-*", packageName), 10).Result()
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		keys = append(keys, newKeys...)
		if cursor == 0 {
			break
		}
	}

	// Summing up all ratings from users
	rating := 0
	quantity := 0
	for _, key := range keys {
		ratingStr, err := rdb.HGet(ctx, key, "rating").Result()
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		r, err := strconv.Atoi(ratingStr)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		rating += r
		quantity++
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
