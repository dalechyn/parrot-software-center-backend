package main

import (
	"net/http"
	"parrot-software-center-backend/handlers"

	. "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func Router() http.Handler  {
	r := mux.NewRouter()

	r.HandleFunc("/ratings/{name}", handlers.GetRatings).Methods("GET")
	r.HandleFunc("/rate", handlers.Rate).Methods("PUT")
	r.HandleFunc("/register", handlers.Register).Methods("POST")
	r.HandleFunc("/login", handlers.Login).Methods("POST")

	loggedHandler := LoggingHandler(log.New().Writer(), r)

	return CORS(
		AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		AllowedMethods([]string{"GET", "HEAD", "POST", "OPTIONS"}))(loggedHandler)
}
