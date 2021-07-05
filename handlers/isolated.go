package handlers

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"net/http"
	"parrot-software-center-backend/utils"
)

// GET route to get a list of isolated packages
func Isolated(w http.ResponseWriter, r *http.Request) {
	log.Debug("Isolated attempt")

	// Connecting to Redis
	rdb := redis.NewFailoverClient(&redis.FailoverOptions{
		SentinelAddrs: []string{"sentinel:26379", "sentinel:26380", "sentinel:26381"},
		MasterName: "mymaster",
		SentinelPassword: utils.GetSentinelPassword(),
		Password: utils.GetRedisPassword(),
	})
	packages, err := rdb.SMembers(ctx, "isolated").Result()
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resBytes, err := json.Marshal(&isolatedResponse{
		packages})
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
