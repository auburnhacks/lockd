package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang/glog"

	"github.com/gorilla/mux"

	"github.com/auburnhacks/lockd/config"
	"github.com/auburnhacks/lockd/lock"
)

type AquireHandler struct {
	Config *config.Config
}

func (h *AquireHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceName := vars["service_name"]
	ttl, err := time.ParseDuration(vars["ttl"])
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %v", err), http.StatusInternalServerError)
		return
	}

	// see if a lock exists
	lck, err := h.Config.GetLockWithServiceName(serviceName)
	if err == config.NotAvailable {
		// TODO: do something here
	} else {
		glog.Infof("lock for service %v found", lck.ServiceName)
		if lckErr := lck.Lock(); lckErr != nil {
			http.Error(w, fmt.Sprintf("error: %v", lckErr), http.StatusInternalServerError)
			return
		}
	}

	// create a new lock
	glog.Infof("creating a new lock")
	l := lock.NewLock(serviceName, ttl)
	h.Config.SetLock(l)
	// aquire a lock
	l.Lock()

	// sending the metadata associated with the lock
	bb, err := json.Marshal(l)
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %v", err), http.StatusInternalServerError)
		return
	}
	w.Write(bb)
}

type ReleaseHandler struct {
	Config *config.Config
}

func (h *ReleaseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	glog.Infof("releasing a lock")
	vars := mux.Vars(r)
	serviceName := vars["service_name"]

	lck, err := h.Config.GetLockWithServiceName(serviceName)
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %v", err), http.StatusInternalServerError)
		return
	}
	lck.Unlock()
	h.Config.DeleteLock(lck)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
