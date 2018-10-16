package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/auburnhacks/lockd/config"
	"github.com/auburnhacks/lockd/handlers"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
)

func init() {
	flag.Set("v", "0")
	flag.Set("logtostderr", "true")

	flag.Parse()
}

func main() {
	defer glog.Flush()
	glog.Infof("create new config")
	glog.Infof("running server on localhost:8000")
	config := config.New()
	r := mux.NewRouter()
	r.Handle("/", &handlers.IndexHandler{Config: config})
	r.Handle("/{service_name}/aquire/{ttl}", &handlers.AquireHandler{Config: config})
	r.Handle("/{service_name}/release/", &handlers.ReleaseHandler{Config: config})
	log.Fatal(http.ListenAndServe("localhost:8000", r))
}
