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

	r.HandleFunc("/packages/{id}", handlers.GetPackage).
		Methods("GET")

	loggedHandler := LoggingHandler(log.New().Writer(), r)

	return CORS(
		AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		AllowedOrigins([]string{"http://localhost:3000", "http://192.168.1.107:3000"}),
		AllowedMethods([]string{"GET", "HEAD", "POST", "OPTIONS"}))(loggedHandler)
}
