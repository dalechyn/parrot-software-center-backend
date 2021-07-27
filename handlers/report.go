package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"net/http"
	"parrot-software-center-backend/tokens"
	"parrot-software-center-backend/utils"
	"strings"
	"time"
)

const (
	ReportUnreviewed = 0
	ReportReviewed = 1
)

// POST report to report
func Report(w http.ResponseWriter, r *http.Request) {
	log.Debug("Report attempt")

	// Decoding http request
	inRequest := &reportRequest{}
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

	authorKey, err := tokens.GetKeyFromToken(inRequest.Token)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// report_{packageName}_{reportedUser}_{whoreported}
	reportKey := fmt.Sprintf("report_%s_%s_%s", inRequest.PackageName, inRequest.ReportedUser, strings.Split(authorKey, "_")[1])
	if _, err := rdb.HSet(ctx, reportKey, "commentary", inRequest.Commentary,
		"date", time.Now().Unix(), "reviewed_by", "", "reviewed_date", "", "review", "").Result(); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := rdb.ZAdd(ctx, "reports", &redis.Z{
		Score:  ReportUnreviewed,
		Member: reportKey,
	}).Result(); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
