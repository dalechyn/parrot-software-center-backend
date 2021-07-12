package main

import (
	"context"
	"flag"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"parrot-software-center-backend/utils"
	"time"
)

var config struct {
	gracefulExitWait time.Duration
	debug            bool
}

func init() {
	flag.DurationVar(&config.gracefulExitWait, "graceful-timeout", time.Second*15,
		"the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.BoolVar(&config.debug, "debug", false, "debug mode")
	flag.Parse()
	utils.Load()
}

func main() {
	var serverAddr string
	if !config.debug {
		serverAddr = ":80"
	} else {
		serverAddr = ":8000"
	}

	srv := &http.Server{
		Addr: serverAddr,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      Router(), // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	log.WithFields(log.Fields{
		"address": srv.Addr,
	}).Info("Server successfully started")

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), config.gracefulExitWait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("Shutting down")
	os.Exit(0)
}
