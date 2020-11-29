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

// POST route to get all reports for moderators
func Reports(w http.ResponseWriter, r *http.Request) {
	log.Debug("Reports attempt")

	// Connecting to Redis
	rdb := redis.NewFailoverClient(&redis.FailoverOptions{
		SentinelAddrs: []string{"sentinel:26379", "sentinel:26380", "sentinel:26381"},
		MasterName: "mymaster",
		SentinelPassword: utils.GetSentinelPassword(),
		Password: utils.GetRedisPassword(),
	})


	// Decoding http request
	inRequest := &reportsRequest{}
	err := json.NewDecoder(r.Body).Decode(inRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userKey, err := utils.GetKeyFromToken(inRequest.Token)

	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
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

	reportsKeys, err := rdb.SMembers(ctx, "reports").Result()
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var lookedUpReports []reportResponse
	for _, reportKey := range reportsKeys {
		results, err := rdb.HMGet(ctx, reportKey, "commentary", "reviewed", "reviewed_by", "reviewed_date", "review").Result()
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		splitted := strings.Split(reportKey, "_")

		reviewed := false
		if results[1] == "1" {
			reviewed = true
		}

		lookedUpReports = append(lookedUpReports, reportResponse{
			PackageName: splitted[1],
			ReportedUser: splitted[2],
			ReportedBy: splitted[4],
			Commentary: results[0].(string),
			Reviewed: reviewed,
			ReviewedBy: results[2].(string),
			ReviewedDate: results[3].(string),
			Review: results[4].(string),
		})
	}

	resBytes, err := json.Marshal(&lookedUpReports)
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
