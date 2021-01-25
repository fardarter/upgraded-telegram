package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"saul/eetest.git/pkg/config"
)

var (
	buildTime string
	gitHash   string
)

var log = logrus.WithField("ctx", "main")

func loggingMiddleware(next http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, next)
}

func main() {
	log.Info(buildTime)
	log.Info(gitHash)
	conf, configErr := config.NewConfig()
	if configErr != nil {
		log.Fatal("Could not load config")
	}
	wait := parseWait()
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	getServer(conf, wait, getMux(r))

}

func parseWait() time.Duration {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()
	return wait
}

func getServer(conf *config.Config, wait time.Duration, r *mux.Router) {
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		WriteTimeout: time.Second * time.Duration(conf.WriteTimeout),
		ReadTimeout:  time.Second * time.Duration(conf.ReadTimeout),
		IdleTimeout:  time.Second * time.Duration(conf.IdleTimeout),
		Handler:      r,
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	shutdownErr := srv.Shutdown(ctx)
	if shutdownErr != nil {
		log.Fatal("Could not shut down server")
		os.Exit(0)
	}
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}

func getMux(mux *mux.Router) *mux.Router {
	mux.HandleFunc("/", home).Methods(http.MethodGet)
	http.Handle("/", mux)
	return mux
}

func home(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}
	rw.WriteHeader(http.StatusOK)
	_, err := rw.Write([]byte("Hello World!"))
	if err != nil {
		log.Info("Something went wrong writing text to respose: %v", err)
	}
}
