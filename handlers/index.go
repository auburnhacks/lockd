package handlers

import (
	"net/http"

	"github.com/auburnhacks/lockd/config"
)

type IndexHandler struct {
	Config *config.Config
}

func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello\n"))
}
