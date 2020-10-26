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

// PUT Delete to rate and comment a package
func Delete(w http.ResponseWriter, r *http.Request) {
	log.Debug("Delete attempt")

	// Decoding http request
	inRequest := &deleteRequest{}
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Validating moderator role
	if _, err := rdb.ZRank(ctx, "moderators", userKey).Result(); err != nil && err != redis.Nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		if err == redis.Nil {
			log.Error(fmt.Errorf("Unauthorized access: %s", userKey))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	ratingKey := fmt.Sprintf("rating_%s_%s", inRequest.Author, strings.Split(userKey, "_")[1])

	if _, err := rdb.ZRem(ctx, "ratings", ratingKey).Result(); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := rdb.Del(ctx, ratingKey).Result(); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
