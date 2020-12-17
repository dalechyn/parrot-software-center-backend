package handlers

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"parrot-software-center-backend/models"
	"parrot-software-center-backend/utils"
)

// GET route to confirm registered account
func Confirm(w http.ResponseWriter, r *http.Request) {
	log.Debug("Confirm attempt")

	// Connecting to Redis
	rdb := redis.NewFailoverClient(&redis.FailoverOptions{
		SentinelAddrs: []string{"sentinel:26379", "sentinel:26380", "sentinel:26381"},
		MasterName: "mymaster",
		SentinelPassword: utils.GetSentinelPassword(),
		Password: utils.GetRedisPassword(),
	})

	tokenStr, exists := mux.Vars(r)["token"]
	if !exists {
		log.Debug("Bad request: ", r.URL.String())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	EMAIL_SECRET, exists := os.LookupEnv("EMAIL_KEY")
	if !exists {
		log.Error("EMAIL_SECRET is not set")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Parsing user key to confirm him
	token, err := jwt.ParseWithClaims(tokenStr, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(EMAIL_SECRET), nil
	})

	if err != nil {
		log.Error("invalid token")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	claims, ok := token.Claims.(*models.Claims)
	if !ok || !token.Valid {
		log.Error("invalid token")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := rdb.HSet(ctx, claims.Key, "confirm", "1").Result(); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := rdb.Persist(ctx, claims.Key).Result(); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	content, err := ioutil.ReadFile("static/success.html")
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(content)); err != nil {
		log.Error(err)
	}
}
