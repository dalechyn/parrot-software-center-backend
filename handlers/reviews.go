package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"net/http"
	"parrot-software-center-backend/utils"
	"strings"

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

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:26379",
		Password: utils.GetRedisPassword(),
	})

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

	var lookedUpRatings []reviewResponse
	for _, key := range keys {
		res, err := rdb.HMGet(ctx, key, "rating", "commentary").Result()
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		lookedUpRatings = append(lookedUpRatings, reviewResponse{
			Author:     strings.Split(key, "-")[2],
			Rating:     res[0].(int),
			Commentary: res[1].(string),
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
