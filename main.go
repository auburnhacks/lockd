package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/auburnhacks/lockd/config"
	"github.com/auburnhacks/lockd/handlers"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
)

var (
	cleanInterval *time.Duration
)

func init() {
	cleanInterval = flag.Duration("clean_duration", 20*time.Second, "interval when the cleanup should run")
	flag.Set("v", "0")
	flag.Set("logtostderr", "true")

	flag.Parse()
}

func main() {
	defer glog.Flush()
	glog.Infof("create new config")
	glog.Infof("running server on localhost:8000")

	config := config.New()
	go config.Cleanup(*cleanInterval)

	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.Handle("/", &handlers.IndexHandler{Config: config})
	r.Handle("/{service_name}/acquire/{ttl}", &handlers.AquireHandler{Config: config})
	r.Handle("/{service_name}/release/", &handlers.ReleaseHandler{Config: config})

	stopChan := make(chan os.Signal)
	go func() {
		if err := http.ListenAndServe("localhost:8000", r); err != nil {
			glog.Fatalf("could not start server: %v", err)
		}
	}()
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	<-stopChan
	glog.Infof("performing silent shutdown...")
	os.Exit(0)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		glog.Infof("%s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
