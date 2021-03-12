package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"net/http"
	"parrot-software-center-backend/utils"
	"strconv"
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
	log.Info("REPORT USERKEY:", userKey)

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

	includeReviewed := "0"
	if inRequest.ShowReviewed {
		includeReviewed = "1"
	}

	reportsKeys, err := rdb.ZRangeByScoreWithScores(ctx, "reports", &redis.ZRangeBy{
		Min:    "0",
		Max:    includeReviewed,
	}).Result()
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var lookedUpReports []reportResponse
	for _, z := range reportsKeys {
		reportKey := z.Member.(string)
		results, err := rdb.HMGet(ctx, reportKey, "commentary", "reviewed_by", "reviewed_date",
			"review", "date").Result()
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		splitted := strings.Split(reportKey, "_")

		reviewed := false
		if z.Score == 1 {
			reviewed = true
		}

		reviewedDate := 0
		if reviewed {
			reviewedDate, err = strconv.Atoi(results[2].(string))
			if err != nil {
				log.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		date, err := strconv.Atoi(results[4].(string))
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		lookedUpReports = append(lookedUpReports, reportResponse{
			PackageName: splitted[1],
			ReportedUser: splitted[2],
			ReportedBy: splitted[3],
			Commentary: results[0].(string),
			Reviewed: reviewed,
			ReviewedBy: results[1].(string),
			ReviewedDate: reviewedDate,
			Review: results[3].(string),
			Date: date,
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
