package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"net/http"
	"parrot-software-center-backend/utils"
)

// GET route to get a list of isolated packages
func RemoveIsolated(w http.ResponseWriter, r *http.Request) {
	log.Debug("Isolated attempt")

	// Connecting to Redis
	rdb := redis.NewFailoverClient(&redis.FailoverOptions{
		SentinelAddrs: []string{"sentinel:26379", "sentinel:26380", "sentinel:26381"},
		MasterName: "mymaster",
		SentinelPassword: utils.GetSentinelPassword(),
		Password: utils.GetRedisPassword(),
	})

	// Decoding http request
	inRequest := &isolatedRequest{}
	err := json.NewDecoder(r.Body).Decode(inRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userKey, err := utils.GetKeyFromToken(inRequest.Token)

	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Validating moderator role
	if exists, err := rdb.SIsMember(ctx, "moderators", userKey).Result(); err != nil && err != redis.Nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		if !exists {
			log.Error(fmt.Errorf("Unauthorized access: %s", userKey))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	if _, err := rdb.SRem(ctx, "isolated", inRequest.Name).Result(); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}