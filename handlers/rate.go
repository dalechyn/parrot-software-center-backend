package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"net/http"
	"parrot-software-center-backend/utils"
	"strings"

	log "github.com/sirupsen/logrus"
)

// PUT route to rate and comment a package
func Rate(w http.ResponseWriter, r *http.Request) {
	log.Debug("Rate attempt")

	// Decoding http request
	inRequest := &rateRequest{}
	err := json.NewDecoder(r.Body).Decode(inRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Connecting to Redis
	rdb := redis.NewFailoverClient(&redis.FailoverOptions{
		SentinelAddrs: []string{":26379", ":26380", ":26381"},
		MasterName: "mymaster",
		SentinelPassword: utils.GetSentinelPassword(),
		Password: utils.GetRedisPassword(),
	})

	userKey, err := utils.GetKeyFromToken(inRequest.Token)

	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	login := strings.Split(userKey, "-")[2]

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
