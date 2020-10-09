package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"net/http"
	"parrot-software-center-backend/utils"

	log "github.com/sirupsen/logrus"
)

func Rate(w http.ResponseWriter, r *http.Request) {
	log.Debug("Rate attempt")

	// Decoding http request
	inRequest := &rateRequest{}
	err := json.NewDecoder(r.Body).Decode(inRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:26379",
		Password: utils.GetRedisPassword(),
	})

	userKey, err := utils.GetKeyFromToken(inRequest.Token)

	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	login, err := rdb.HGet(ctx, userKey, "login").Result()

	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ratingKey := fmt.Sprintf("rating-%s-%s", inRequest.Name, login)

	_, err = rdb.HSet(ctx, ratingKey, "rating", inRequest.Rating, "commentary", inRequest.Comment).Result()
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = rdb.SAdd(ctx, "ratings", ratingKey).Result()
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
