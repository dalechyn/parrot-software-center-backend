package main

import (
	. "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"parrot-software-center-backend/handlers"
)

func Router() http.Handler  {
	r := mux.NewRouter()

	r.HandleFunc("/ratings/{name}", handlers.Ratings).Methods("GET")
	r.HandleFunc("/reviews/{name}", handlers.Reviews).Methods("GET")
	r.HandleFunc("/confirm/{token}", handlers.Confirm).Methods("GET")
	r.HandleFunc("/rate", handlers.Rate).Methods("PUT")
	r.HandleFunc("/register", handlers.Register).Methods("POST")
	r.HandleFunc("/login", handlers.Login).Methods("POST")
	r.HandleFunc("/delete", handlers.Delete).Methods("POST")
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets")))).
		Methods("GET")

	loggedHandler := LoggingHandler(log.New().Writer(), r)

	return CORS(
		AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		AllowedMethods([]string{"GET", "POST", "PUT", "OPTIONS"}))(loggedHandler)
}
