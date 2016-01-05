package web

import (
	"net/http"
	"github.com/gorilla/mux"
)

// TODO: move to a different file
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!\n"))
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK\n"))
}

// registerRoutes registers all the HTTP routes for the web server.
func registerRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler)
	r.HandleFunc("/v1/status", StatusHandler)
	return r
}
