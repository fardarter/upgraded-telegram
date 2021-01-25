package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/eetest.git/pkg/config"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var (
	buildTime string
	gitHash   string
)

var log = logrus.WithField("ctx", "main")

func main() {
	log.Info(buildTime)
	log.Info(gitHash)
	conf, configErr := config.NewConfig()
	if configErr != nil {
		log.Fatal("Could not load config")
	}
	wait := parseWait()
	getServer(conf, wait, getMux())

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
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}

func getMux() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", Home)
	http.Handle("/", r)
	return r
}

func Home(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	fmt.Fprintf(rw, "Hello World!")
}
