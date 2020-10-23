package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"parrot-software-center-backend/tokens"
	"parrot-software-center-backend/utils"
)

// POST route to handle user login
func Login(w http.ResponseWriter, r *http.Request) {
	log.Debug("Login attempt")

	// Decoding http request
	inRequest := &loginRequest{}
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

	// Checking if user exists
	newKeys, _, err := rdb.SScan(ctx, "users", 0, fmt.Sprintf("user-*-%s", inRequest.Username), 1).Result()

	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(newKeys) == 0 {
		log.Error("attempt to login a user which does not exist")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resMap, err := rdb.HMGet(ctx, newKeys[0], "password", "confirmed").Result()
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	byteHash := []byte(resMap[0].(string))
	err = bcrypt.CompareHashAndPassword(byteHash, []byte(inRequest.Password))
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Restrict user from logging in if account is not confirmed
	if resMap[1].(string) == "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Generate some tokens for him
	token, err := tokens.Generate(newKeys[0])
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Encode to JSON and send him
	resBytes, err := json.Marshal(&loginResponse{token})
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
