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
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Connecting to Redis
	rdb := redis.NewFailoverClient(&redis.FailoverOptions{
		SentinelAddrs:    []string{"sentinel:26379", "sentinel:26380", "sentinel:26381"},
		MasterName:       "mymaster",
		SentinelPassword: utils.GetSentinelPassword(),
		Password:         utils.GetRedisPassword(),
	})

	// Checking if user exists
	userKey := fmt.Sprintf("user_%s", inRequest.Username)
	if exists, err := rdb.Exists(ctx, userKey).Result(); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		if exists == 0 {
			log.Error("attempt to login a user which does not exist")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	resMap, err := rdb.HMGet(ctx, userKey, "password", "confirmed").Result()
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

	token := ""
	moderator := false

	if exists, err := rdb.SIsMember(ctx, "moderators", userKey).Result(); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		if exists {
			moderator = true
		}
	}

	// Encode to JSON and send him
	var resBytes []byte
	if moderator {
		token, err = tokens.Generate(userKey, RoleModerator)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resBytes, err = json.Marshal(&loginResponse{token, RoleModerator})
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		token, err = tokens.Generate(userKey, RoleUser)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resBytes, err = json.Marshal(&loginResponse{token, RoleUser})
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
}