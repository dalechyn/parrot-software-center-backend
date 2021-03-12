package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"net/http"
	"parrot-software-center-backend/utils"
	"strings"
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
		SentinelAddrs: []string{"sentinel:26379", "sentinel:26380", "sentinel:26381"},
		MasterName: "mymaster",
		SentinelPassword: utils.GetSentinelPassword(),
		Password: utils.GetRedisPassword(),
	})

	userKey, err := utils.GetKeyFromToken(inRequest.Token)

	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	ratingKey := fmt.Sprintf("rating_%s_%s", inRequest.Name, strings.Split(userKey, "_")[1])

	_, err = rdb.Set(ctx, ratingKey, inRequest.Comment, 0).Result()
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = rdb.ZAdd(ctx, "ratings", &redis.Z{Score: inRequest.Rating, Member: ratingKey}).Result()
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
